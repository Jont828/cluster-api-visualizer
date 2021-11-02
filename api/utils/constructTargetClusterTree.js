const k8s = require('@kubernetes/client-node');
const { default: cluster } = require('cluster');
const { assert } = require('console');

const HttpStatus = require('http-status-codes')

const resourceMap = {
  clusterresourcesetbindings: { group: "addons.cluster.x-k8s.io", category: "clusterInfra" },
  clusterresourcesets: { group: "addons.cluster.x-k8s.io", category: "clusterInfra" },
  // clusterclasses: { group: "cluster.x-k8s.io", category: "clusterInfra" },
  clusters: { group: "cluster.x-k8s.io", category: null },
  machinedeployments: { group: "cluster.x-k8s.io", category: "workers" },
  // machinehealthchecks: { group: "cluster.x-k8s.io", category: "clusterInfra" },
  machinepools: { group: "cluster.x-k8s.io", category: "workers" },
  machinesets: { group: "cluster.x-k8s.io", category: "workers" },
  machines: { group: "cluster.x-k8s.io", category: null },
  azureclusteridentities: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  azureclusters: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  dockerclusters: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  azuremachinepoolmachines: { group: "infrastructure.cluster.x-k8s.io", category: "workers" },
  azuremachinepools: { group: "infrastructure.cluster.x-k8s.io", category: "workers" },
  dockermachinepools: { group: "infrastructure.cluster.x-k8s.io", category: "workers" },
  azuremachines: { group: "infrastructure.cluster.x-k8s.io", category: null },
  dockermachines: { group: "infrastructure.cluster.x-k8s.io", category: null },
  azuremachinetemplates: { group: "infrastructure.cluster.x-k8s.io", category: null },
  dockermachinetemplates: { group: "infrastructure.cluster.x-k8s.io", category: null },
  // azuremanagedclusters: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  // azuremanagedcontrolplanes: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  // azuremanagedmachinepools: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  // azureserviceprincipals: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  // azuresystemassignedidentites: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  // azureuserassignedidentites: { group: "infrastructure.cluster.x-k8s.io", category: "clusterInfra" },
  kubeadmconfigs: { group: "bootstrap.cluster.x-k8s.io", category: "clusterInfra" },
  kubeadmconfigtemplates: { group: "bootstrap.cluster.x-k8s.io", category: "clusterInfra" },
  // kubeadmconfigs: { group: "bootstrap.cluster.x-k8s.io", category: null },
  // kubeadmconfigtemplates: { group: "bootstrap.cluster.x-k8s.io", category: null },
  kubeadmcontrolplanes: { group: "controlplane.cluster.x-k8s.io", category: "controlPlane" },
  // kubeadmcontrolplanetemplates: { group: "controlplane.cluster.x-k8s.io", category: "controlPlane" },
};

function resolveCategory(crd) {
  if (crd.kind == 'Cluster')
    return null
  else if (crd.labels !== undefined)
    return ('cluster.x-k8s.io/control-plane' in crd.labels) ? 'controlPlane' : 'workers'
  else // TODO: Don't rely on names of CRDs
    return ((crd.name.indexOf('control-plane') > -1) || (crd.name.indexOf('controlplane') > -1)) ? 'controlPlane' : 'workers'
}

const multipleOwners = {
  // Format = Kind: { ExpectedOwner, RedundantOwners }
  'AzureMachine': { expectedOwner: 'Machine', redundantOwners: ['KubeadmControlPlane'] },
  'DockerMachine': { expectedOwner: 'Machine', redundantOwners: ['KubeadmControlPlane'] },
  'KubeadmConfig': { expectedOwner: 'Machine', redundantOwners: ['KubeadmControlPlane'] },
  'ClusterResourceSetBinding': { expectedOwner: 'ClusterResourceSet', redundantOwners: ['Cluster'] },
}

function resolveOwners(crd) {
  let owners = crd.ownerRefs;

  if (owners.length > 1) { // If multiple owners 

    // If kind in lookup table
    if (crd.kind in multipleOwners) {
      let expectedOwner = multipleOwners[crd.kind].expectedOwner;
      let allOwners = new Set(multipleOwners[crd.kind].redundantOwners);
      allOwners.add(expectedOwner);

      // If owners match owners in lookup table for kind
      if (owners.length == allOwners.size) {
        let match = true;
        owners.forEach((o, i) => {
          match = match && allOwners.has(o.kind);
        });

        if (match)
          return owners.filter(o => o.kind == expectedOwner)[0].uid; // Return ID of expected owner type in owner refs if matched
      }
      console.log('Cannot resolve multiple owners for', crd.kind);
      console.log(owners);
      throw 'Failed to resolve multiple owners!';
    }

  } else { // If only one owner, easy case
    return owners[0].uid;
  }
}

async function getCRDInstances(group, plural, initCategory, clusterName, clusterUid) {
  const kc = new k8s.KubeConfig();
  kc.loadFromDefault();
  const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

  try {
    const res = await k8sCrd.listClusterCustomObject(group, 'v1beta1', plural);
    let crds = [];
    res.body.items.forEach((e, i) => {
      // 1. Init easy fields
      let crd = {
        id: e.metadata.uid,
        name: e.metadata.name,
        kind: e.kind,
        group: group,
        plural: plural,
        provider: group.substr(0, group.indexOf('.')),
        ownerRefs: e.metadata.ownerReferences,
        labels: e.metadata.labels,
        spec: e.spec
      }

      // 2. If the category depends on context, i.e. Machine, then resolve it now
      crd.category = initCategory ? initCategory : resolveCategory(crd)

      // 3. If there are resources left without owners, bind them to the root
      let owner;
      if (crd.kind == 'Cluster') { // Root node has no owner
        owner = null;
      } if (e.metadata.ownerReferences === undefined) { // If no owners and not the root, i.e. bucket/category nodes
        owner = clusterUid;
      } else {
        owner = resolveOwners(crd);
      }

      // Lastly, take all the parents that point to the root and bind them to their respective category node
      if (owner == clusterUid)
        owner = crd.category;

      crd.parent = owner;
      crds.push(crd)
    })

    return crds;
  } catch (error) {
    if (error.statusCode == HttpStatus.NOT_FOUND)
      return [];
    console.log(error);
    throw 'Error fetching for ' + plural + ' in ' + clusterName
  }
}

module.exports = async function constructTargetClusterTree(clusterName) {
  // Hack since getClusterCustomObject is getting a 404
  const kc = new k8s.KubeConfig();
  kc.loadFromDefault();
  const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

  const response = await k8sCrd.listClusterCustomObject('cluster.x-k8s.io', 'v1beta1', 'clusters');
  let clusters = response.body.items.filter(e => e.metadata.name == clusterName);
  assert(clusters.length == 1);
  let clusterUid = clusters[0].metadata.uid;
  let clusterLabels = clusters[0].metadata.labels;
  // End hack

  let allCrds = [];

  for (const [plural, value] of Object.entries(resourceMap)) {
    const instances = await getCRDInstances(value.group, plural, value.category, clusterName, clusterUid);
    allCrds = allCrds.concat(instances);
  }

  const whitelistKinds = ['ClusterResourceSet', 'ClusterResourceSetBinding'];

  let crds = allCrds.filter(crd => (
    (crd.labels !== undefined && crd.labels['cluster.x-k8s.io/cluster-name'] == clusterName) ||
    (crd.ownerRefs !== undefined && crd.ownerRefs.find(o => o.uid == clusterUid)) ||
    crd.name.indexOf(clusterName) == 0 ||
    whitelistKinds.includes(crd.kind)
  ));
  // console.log(crds);

  // TODO: Filter types with cluster-name label instead
  // let crds = allCrds.filter(crd => (crd.labels['cluster.x-k8s.io/cluster-name'] == clusterName || whitelistKinds.includes(crd.kind)));

  let binding = allCrds.find(crd => (
    crd.kind == 'ClusterResourceSetBinding' &&
    crd.ownerRefs.find(e => e.kind == 'Cluster').name == clusterName
  ));

  if (binding) {
    let resourceSetNames = new Set();
    binding.spec.bindings.forEach(e => {
      resourceSetNames.add(e.clusterResourceSetName);
    });
    crds = crds.filter(crd => (
      !whitelistKinds.includes(crd.kind) || // Keep non binding or resource set
      (crd.kind == 'ClusterResourceSet' && resourceSetNames.has(crd.name)) ||
      (crd.kind == 'ClusterResourceSetBinding' && crd.name == binding.name)
    ));
  } else {
    crds = crds.filter(crd => (!whitelistKinds.includes(crd.kind)));
  }


  // Add dummy nodes with CRDs
  let dummyNodes = [
    {
      name: '',
      kind: 'ClusterInfrastructure',
      id: 'clusterInfra',
      provider: '',
      collapsable: true,
      parent: clusterUid,
    },
    {
      name: '',
      kind: 'ControlPlane',
      id: 'controlPlane',
      provider: '',
      collapsable: true,
      parent: clusterUid,
    },
    {
      name: '',
      kind: 'Workers',
      id: 'workers',
      provider: '',
      collapsable: true,
      parent: clusterUid,
    },
  ];

  crds = crds.concat(dummyNodes);

  // Create mapping to prepare to construct tree
  const idMapping = crds.reduce((acc, e, i) => {
    acc[e.id] = i;
    return acc;
  }, {});

  console.log(idMapping);

  let root;
  // console.log(crds);
  crds.forEach(e => {
    // Handle the root element
    if (e.parent == null) {
      root = e;
      return;
    }
    // Use our mapping to locate the parent element in our data array
    let parentNode = crds[idMapping[e.parent]];

    // console.log('Parent', parentNode);
    // console.log('Child', e);
    // console.log('\n');

    // Add our current e to its parent's `children` array
    if (parentNode.children === undefined)
      parentNode.children = [];

    parentNode.children.push(e)


  });

  console.log('Final tree:');
  console.log(root);
  return root;

}


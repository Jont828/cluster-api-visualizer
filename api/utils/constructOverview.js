const k8s = require('@kubernetes/client-node');
const HttpStatus = require('http-status-codes')

function getIcon(cluster) {
  let clusterType = cluster.spec.infrastructureRef.kind;
  switch (clusterType) {
    case 'AzureCluster':
      return 'microsoft-azure'
    case 'DockerCluster':
      return 'docker'
    case 'GCPCluster':
      return 'google-cloud'
    case 'AWSCluster':
      return 'aws'
    default:
      return 'kubernetes'
  }
}

module.exports = async function constructOverview() {

  const kc = new k8s.KubeConfig();
  let k8sCrd;
  try {
    kc.loadFromDefault();
    k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);
  } catch (error) {
    return null;
  }

  const context = kc.currentContext;
  const cluster = kc.clusters.find(ctx => ctx.name == context);
  if (cluster === undefined)
    return null;

  let root = {
    name: cluster.name,
    isRoot: true,
    icon: "kubernetes",
    children: [],
  }

  // TODO: Make this work recursively, i.e. if a child is another management cluster
  try {
    console.log('Looking for target clusters');
    const response = await k8sCrd.listClusterCustomObject('cluster.x-k8s.io', 'v1beta1', 'clusters');
    // console.log(response.body);
    response.body.items.forEach((e, i) => {
      let clusterName = e.metadata.name;
      console.log('Found cluster', clusterName);
      root.children.push({
        name: clusterName,
        icon: getIcon(e),
        children: []
      })
    });
  } catch (error) {
    console.log('Error is');
    console.log(error.message);
    // if (error.statusCode == HttpStatus.NOT_FOUND)
  }

  return root;
}

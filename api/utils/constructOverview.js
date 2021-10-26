const k8s = require('@kubernetes/client-node');
const { default: cluster } = require('cluster');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

module.exports = async function constructOverview() {
  let tree = {
    name: "kind-capz",
    isRoot: true,
    icon: "kubernetes",
    children: [],
  }

  // TODO: Make this work recursively, i.e. if a child is another management cluster
  try {
    const response = await k8sCrd.listClusterCustomObject('cluster.x-k8s.io', 'v1beta1', 'clusters');
    // console.log(response.body);
    response.body.items.forEach((e, i) => {
      let clusterName = e.metadata.name;
      console.log('Found cluster', clusterName);
      tree.children.push({
        name: clusterName,
        icon: 'microsoft-azure',
        children: []
      })
    });
  } catch (error) {
    console.log('Error fetching cluster overview');
    console.log(error);
  }

  return tree;
}

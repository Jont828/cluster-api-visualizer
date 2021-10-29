const k8s = require('@kubernetes/client-node');

module.exports = async function constructOverview() {

  const kc = new k8s.KubeConfig();
  kc.loadFromDefault();
  const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

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
    const response = await k8sCrd.listClusterCustomObject('cluster.x-k8s.io', 'v1beta1', 'clusters');
    // console.log(response.body);
    response.body.items.forEach((e, i) => {
      let clusterName = e.metadata.name;
      console.log('Found cluster', clusterName);
      root.children.push({
        name: clusterName,
        icon: 'microsoft-azure',
        children: []
      })
    });
  } catch (error) {
    console.log(error);
    throw 'Error fetching target clusters';
  }

  return root;
}

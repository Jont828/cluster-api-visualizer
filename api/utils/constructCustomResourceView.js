const k8s = require('@kubernetes/client-node');
const { default: cluster } = require('cluster');
const { assert } = require('console');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

module.exports = async function constructCustomResourceView(group, plural, name) {
  // Hack since getClusterCustomObject is getting a 404
  const response = await k8sCrd.listClusterCustomObject(group, 'v1beta1', plural);
  let items = response.body.items.filter(e => e.metadata.name == name);
  assert(items.length == 1);
  // End hack

  return formatToTreeview(items[0]);
}

function formatToTreeview(resource, id = 0) {
  let result = [];
  if (typeof (resource) == 'string') {
    return [{ name: resource }]
  } else if (Array.isArray(resource)) {
    let children = [];
    resource.forEach((e, i) => {
      result.push({
        id: id++,
        name: i.toString(),
        children: formatToTreeview(e, id)
      });
    });

  } else { // isObject
    Object.entries(resource).forEach(([key, value]) => {
      let name = '';
      let children = [];
      if (typeof (value) == 'string') {
        name = key + ': ' + value
      } else {
        name = key;
        children = formatToTreeview(value, id);
      }
      result.push({
        id: id++,
        name: name,
        children: children
      });
    });

  }

  return result;
}
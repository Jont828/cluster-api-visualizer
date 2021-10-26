const yaml = require('js-yaml');
const fs = require('fs');

module.exports = function constructCustomResourceView() {
  const file = yaml.load(fs.readFileSync('./temp-assets/azureclusters.infrastructure.cluster.x-k8s.io-default-1495.yaml', 'utf8'));
  console.log(JSON.stringify(file, null, 2));
  return formatToTreeview(file);
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
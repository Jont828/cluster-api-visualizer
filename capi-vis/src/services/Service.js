const axios = require('axios');

const http = axios.create({
  baseURL: 'http://localhost:3080/',
})

export async function getClusterOverview() {
  console.log('Getting cluster overview');
  const response = await axios.get(`/api/cluster-overview/`);
  return response.data;
}


export async function getCluster(clusterId) {
  console.log('Getting Cluster ' + clusterId);
  const response = await axios.get(`/api/cluster/`, { params: { ID: clusterId } });
  return response.data;
}

export async function getClusterResource(clusterId, resourceKind, resourceName) {
  console.log('Getting CRD ' + clusterId + ' ' + resourceKind + ' ' + resourceName);
  const response = await axios.get(`/api/cluster-resource/`, {
    params: { ID: clusterId, resourceKind: resourceKind, resourceName: resourceName }
  });
  return response.data;
}

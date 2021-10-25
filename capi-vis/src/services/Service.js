const axios = require('axios');

const http = axios.create({
  baseURL: 'http://localhost:3080/',
})

export async function getAllClusters() {

  const response = await axios.get('/api/cluster');
  return response.data;
}

export async function getCluster(clusterId) {
  console.log("Getting Cluster " + clusterId);
  const response = await axios.get(`/api/cluster/`, { id: clusterId } );
  return response.data;
}

export async function postCluster(data) {
  console.log("Posting getCluster");
  const response = await axios.post(`/api/cluster`, { cluster: data });
  return response.data;
}
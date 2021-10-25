const express = require('express');
const path = require('path');
const app = express(),
  bodyParser = require("body-parser");
port = 3080;

// place holder for the data
const clusters = [];

app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../capi-vis/build')));

app.get('/api/cluster', (req, res) => {
  console.log('api/clusters called!')
  let id = req.query.id;
  console.log('Got cluster ID' + id);
  res.json({cluster: 'myCluster'});
});

app.post('/api/cluster', (req, res) => {
  const cluster = req.body.cluster;
  console.log('Adding cluster:::::', cluster);
  clusters.push(cluster);
  res.json("cluster added");
});

app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, '../capi-vis/build/index.html'));
});

app.listen(port, () => {
  console.log(`Server listening on the port::${port}`);
});
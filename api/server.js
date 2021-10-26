const express = require('express');
const path = require('path');
const app = express(),
  bodyParser = require("body-parser");
port = 3080;

const constructOverview = require('./utils/constructOverview.js');
const constructTargetClusterTree = require('./utils/constructTargetClusterTree.js');
const constructCustomResourceView = require('./utils/constructCustomResourceView.js');

app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../capi-vis/build')));

app.get('/api/cluster-overview', async (req, res) => {
  console.log('api/clusters-overview called!')
  const tree = await constructOverview();
  res.json(tree);
});

app.get('/api/cluster', (req, res) => {
  console.log('api/clusters called!')
  let id = req.query.ID;
  res.json(constructTargetClusterTree(id));
});

app.get('/api/cluster-resource', (req, res) => {
  console.log('api/cluster-resource called!')
  console.log(req.query);
  let id = req.query.ID;

  try {
    let result = constructCustomResourceView();
    res.json(result);
  } catch (e) {
    console.log(e);
    res.status(500);
  }
});

// app.get('/', (req, res) => {
//   res.sendFile(path.join(__dirname, '../capi-vis/build/index.html'));
// });

app.listen(port, () => {
  console.log(`Server listening on the port::${port}`);
});


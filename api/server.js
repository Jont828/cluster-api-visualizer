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
  try {
    const tree = await constructOverview();
    res.json(tree);
  } catch (e) {
    console.log(e);
    res.status(500);
  }
});

app.get('/api/cluster', async (req, res) => {
  console.log('api/clusters called!')
  try {
    let id = req.query.ID;
    const tree = await constructTargetClusterTree(id);
    res.json(tree);
  } catch (e) {
    console.log(e);
    res.status(500);
  }
});

app.get('/api/cluster-resource', async (req, res) => {
  console.log('api/cluster-resource called!')

  try {
    let group = req.query.group;
    let plural = req.query.plural;
    let name = req.query.name;
    const result = await constructCustomResourceView(group, plural, name);
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


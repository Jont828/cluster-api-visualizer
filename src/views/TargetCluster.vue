<template>
  <div class="container">
    <h1>Target Cluster Resources: {{ this.$route.params.id }}</h1>
    <vue-tree
      id="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <template v-slot:node="{ node, collapsed }">
        <div
          class="node"
          :style="{ border: collapsed ? '2px solid grey' : '' }"
        >
          <router-link :to="'/'">
            <span>{{ node.name }}</span>
          </router-link>
        </div>
      </template>
    </vue-tree>
  </div>
</template>

<script>
import VueTree from "@ssthouse/vue-tree-chart";

export default {
  name: "TargetCluster",
  components: {
    VueTree,
  },
  data() {
    return {
      treeData: {
        name: "default-1",
        children: [
          {
            name: "Cluster Infra",
            children: [
              {
                name: "AzureCluster",
              },
              {
                name: "ClusterResource SetBindings",
              },
              {
                name: "AzureClusterIdentity",
              },
              {
                name: "ClusterResourceSet",
              },
            ],
          },
          {
            name: "Control Plane",
            children: [
              {
                name: "Machine",
              },
              {
                name: "Azure Machine",
              },
              {
                name: "AzureMachine Template",
              },
              {
                name: "KubeadmConfig",
              },
              {
                name: "KubeadmCtrlPlane",
              },
            ],
          },
          {
            name: "Workers",
            children: [
              {
                name: "Machines",
                children: [
                  {
                    name: "default-1-asdf",
                    children: [
                      {
                        name: "Machine",
                      },
                      {
                        name: "AzureMachine",
                      },
                      {
                        name: "Kubeadmconfig",
                      },
                    ],
                  },
                  {
                    name: "default-1-jkl;",
                    children: [
                      {
                        name: "Machine",
                      },
                      {
                        name: "AzureMachine",
                      },
                      {
                        name: "Kubeadmconfig",
                      },
                    ],
                  },
                ],
              },
              {
                name: "Kubeadmconfig template",
              },
              {
                name: "Machine Deployment",
              },
              {
                name: "AzureMachine Template",
              },
              {
                name: "MachineSet",
              },
            ],
          },
        ],
      },
      treeConfig: { nodeWidth: 120, nodeHeight: 50, levelHeight: 100 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style scoped>
#tree {
  width: 100%;
  height: 1000px;
  border: 1px solid black;
}

.container {
  height: 100%;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 100px;
  height: 50px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #2c3e50;
  background-color: #dae8fc;
  border-radius: 4px;
}

.node span {
  font-size: 10px;
  margin: 0px;
}
</style>

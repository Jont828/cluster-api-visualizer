package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	"sigs.k8s.io/cluster-api/controllers/external"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed frontend/dist
var frontend embed.FS

func test() {
	fmt.Println("Testing")
	_, err := getCustomResource("AzureCluster", "infrastructure.cluster.x-k8s.io/v1beta1", "default", "default-7370")
	if err != nil {
		fmt.Println(err)
	}
	// _, err = getCustomResource("Cluster", "cluster.x-k8s.io/v1beta1", "default", "default-7370")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "The port to listen on")
	flag.Parse()

	http.Handle("/api/v1/cluster-resources/", http.HandlerFunc(handleObjectTree))
	http.Handle("/api/v1/multicluster/", http.HandlerFunc(handleMultiClusterTree))
	http.Handle("/api/v1/custom-resource/", http.HandlerFunc(handleCustomResourceTree))

	stripped, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalln(err)
	}

	frontendFS := http.FileServer(http.FS(stripped))
	http.Handle("/", frontendFS)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type CustomResourceRequestBody struct {
	APIVersion string `json:"version"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

func handleCustomResourceTree(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting custom resource tree\n")
	fmt.Println("GET params were:", r.URL.Query())
	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")

	object, err := getCustomResource(kind, apiVersion, "default", name)
	if err != nil {
		fmt.Println("Failed to get CRD:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := object.MarshalJSON()
	if err != nil {
		fmt.Println("Couldn't unmarshal CRD", err)
		http.Error(w, "couldn't marshall CRD json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(data))
}

func handleMultiClusterTree(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting multicluster tree\n")
	tree, err := client.MultiDiscovery(dc.kubeconfig)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	exportTree := constructMultiClusterTree(tree)
	if exportTree != nil {
		marshalled, err := json.MarshalIndent(*exportTree, "", "  ")
		if err != nil {
			http.Error(w, "couldn't retrieve multi cluster", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, bytes.NewReader(marshalled))
	}
}

func handleObjectTree(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Path[len("/api/v1/cluster-resources/"):]
	fmt.Printf("Getting object tree for %s\n", clusterName)
	objTree, err := getObjectTree(clusterName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// cmd.PrintObjectTree(objTree)
	exportTree := constructClusterResourceTree(objTree, objTree.GetRoot())
	if exportTree != nil {
		marshalled, err := json.MarshalIndent(*exportTree, "", "  ")
		if err != nil {
			http.Error(w, "couldn't retrieve cluster", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, bytes.NewReader(marshalled))
	}
}

type describeClusterOptions struct {
	kubeconfig          string
	kubeconfigContext   string
	namespace           string
	showOtherConditions string
	showAllResources    bool
	showMachineSets     bool
	// echo                bool
	// disableNoEcho       bool
	// grouping            bool
	// disableGrouping     bool
}

var dc = &describeClusterOptions{
	kubeconfig:       "/home/jonathan/.kube/config",
	showAllResources: true,
}

func getObjectTree(name string) (*tree.ObjectTree, error) {
	c, err := client.New("")
	if err != nil {
		return nil, err
	}

	return c.DescribeCluster(client.DescribeClusterOptions{
		Kubeconfig:          client.Kubeconfig{Path: dc.kubeconfig, Context: dc.kubeconfigContext},
		Namespace:           dc.namespace,
		ClusterName:         name,
		ShowOtherConditions: dc.showOtherConditions,
		ShowAllResources:    dc.showAllResources,
		ShowMachineSets:     dc.showMachineSets,
	})
}

type ClusterResourceNode struct {
	Name      string                 `json:"name"`
	Kind      string                 `json:"kind"`
	Group     string                 `json:"group"`
	Version   string                 `json:"version"`
	Provider  string                 `json:"provider"`
	UID       string                 `json:"uid"`
	IsVirtual bool                   `json:"isVirtual"`
	Children  []*ClusterResourceNode `json:"children"`
}

func constructClusterResourceTree(objTree *tree.ObjectTree, object ctrlclient.Object) *ClusterResourceNode {
	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version
	fmt.Printf("Name: %s/%s\n", kind, object.GetName())
	fmt.Printf("Group: %s\n", group)
	fmt.Printf("Version: %s\n\n", version)
	node := &ClusterResourceNode{
		Name:      object.GetName(),
		Kind:      kind,
		Group:     group,
		Version:   version,
		Provider:  group[:strings.IndexByte(group, '.')],
		IsVirtual: tree.IsVirtualObject(object),
		Children:  []*ClusterResourceNode{},
		UID:       string(object.GetUID()),
	}
	for _, child := range objTree.GetObjectsByParent(object.GetUID()) {
		node.Children = append(node.Children, constructClusterResourceTree(objTree, child))
	}
	return node
}

type MultiClusterTreeNode struct {
	Name                   string `json:"name"`
	Namespace              string
	Icon                   string `json:"icon"`
	InfrastructureProvider string
	IsVirtual              bool
	Children               []*MultiClusterTreeNode `json:"children"`
	Kubeconfig             string
}

func constructMultiClusterTree(tree *client.MultiClusterTree) *MultiClusterTreeNode {
	if tree == nil {
		return nil
	}

	node := &MultiClusterTreeNode{
		Name:      tree.Name,
		Namespace: tree.Namespace,
		Icon:      getIcon(to.String(tree.InfrastructureProvider)),
		Children:  []*MultiClusterTreeNode{},
		IsVirtual: false,
	}
	if tree.Kubeconfig != nil {
		node.Kubeconfig = *tree.Kubeconfig
	} else {
		node.Kubeconfig = ""
	}
	if tree.InfrastructureProvider != nil {
		node.InfrastructureProvider = *tree.InfrastructureProvider
	} else {
		node.InfrastructureProvider = ""
	}

	for _, child := range tree.WorkloadClusters {
		node.Children = append(node.Children, constructMultiClusterTree(child))
	}
	return node
}

func getIcon(provider string) string {
	switch provider {
	case "AzureCluster":
		return "microsoft-azure"
	case "DockerCluster":
		return "docker"
	case "GCPCluster":
		return "google-cloud"
	case "AWSCluster":
		return "aws"
	default:
		return "kubernetes"
	}
}

func getCustomResource(kind string, apiVersion string, namespace string, name string) (*unstructured.Unstructured, error) {
	cfgFile := ""
	configClient, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}

	clusterClient := cluster.New(cluster.Kubeconfig{Path: dc.kubeconfig, Context: ""}, configClient)

	// Fetch the Cluster client.
	client, err := clusterClient.Proxy().NewClient()
	if err != nil {
		return nil, err
	}
	objectRef := corev1.ObjectReference{
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
		APIVersion: apiVersion,
	}
	object, err := external.Get(context.TODO(), client, &objectRef, namespace)
	if err != nil {
		return nil, err
	}

	data, err := object.MarshalJSON()
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(string(data), "", "\t")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(out))

	return object, nil
}

// type CustomResourceTreeNode struct {
// 	ID      string
// 	Name    string
// 	Chidren []*CustomResourceTreeNode
// }

// func constructCustomResourceTree(object *unstructured.Unstructured) []*ClusterResourceNode {

// }

package main

import (
	"bytes"
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
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed frontend/dist
var frontend embed.FS

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "The port to listen on")
	flag.Parse()

	http.Handle("/api/v1/cluster/", http.HandlerFunc(handleObjectTree))
	http.Handle("/api/v1/multicluster/", http.HandlerFunc(handleMultiTree))

	stripped, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalln(err)
	}

	frontendFS := http.FileServer(http.FS(stripped))
	http.Handle("/", frontendFS)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleMultiTree(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting multicluster tree\n")
	tree, err := client.MultiDiscovery(dc.kubeconfig)
	if err != nil {
		fmt.Println(err)
	}
	exportTree := constructMultiExportTree(tree)
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

func handleObjectTree(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Path[len("/api/v1/cluster/"):]
	fmt.Printf("Getting object tree for %s\n", clusterName)
	objTree, err := getObjectTree(clusterName)
	if err != nil {
		fmt.Println(err)
	}
	// cmd.PrintObjectTree(objTree)
	exportTree := constructObjectExportTree(objTree, objTree.GetRoot())
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

// type ObjectTree struct {
// 	root      client.Object
// 	options   ObjectTreeOptions
// 	items     map[types.UID]client.Object
// 	ownership map[types.UID]map[types.UID]bool
// }

type TreeNode struct {
	Name      string      `json:"name"`
	Kind      string      `json:"kind"`
	Group     string      `json:"group"`
	Version   string      `json:"version"`
	Provider  string      `json:"provider"`
	UID       string      `json:"uid"`
	IsVirtual bool        `json:"isVirtual"`
	Children  []*TreeNode `json:"children"`
}

func constructObjectExportTree(objTree *tree.ObjectTree, object ctrlclient.Object) *TreeNode {
	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version
	fmt.Printf("Name: %s/%s\n", kind, object.GetName())
	fmt.Printf("Group: %s\n", group)
	fmt.Printf("Version: %s\n\n", version)
	node := &TreeNode{
		Name:      object.GetName(),
		Kind:      kind,
		Group:     group,
		Version:   version,
		Provider:  group[:strings.IndexByte(group, '.')],
		IsVirtual: tree.IsVirtualObject(object),
		Children:  []*TreeNode{},
		UID:       string(object.GetUID()),
	}
	for _, child := range objTree.GetObjectsByParent(object.GetUID()) {
		node.Children = append(node.Children, constructObjectExportTree(objTree, child))
	}
	return node
}

type MultiTreeNode struct {
	Name                   string `json:"name"`
	Namespace              string
	Icon                   string `json:"icon"`
	InfrastructureProvider string
	IsVirtual              bool
	Children               []*MultiTreeNode `json:"children"`
	Kubeconfig             string
}

func constructMultiExportTree(tree *client.MultiClusterTree) *MultiTreeNode {
	if tree == nil {
		return nil
	}

	node := &MultiTreeNode{
		Name:      tree.Name,
		Namespace: tree.Namespace,
		Icon:      getIcon(to.String(tree.InfrastructureProvider)),
		Children:  []*MultiTreeNode{},
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
		node.Children = append(node.Children, constructMultiExportTree(child))
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

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

	"github.com/Jont828/capi-visualization/internal"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed web/dist
var frontend embed.FS

var kubeconfigPath = ""
var kubeContext = ""
var namespace string

var defaultClient client.Client
var clusterClient cluster.Client
var runtimeClient ctrlClient.Client
var k8sConfigClient *api.Config

func init() {
	var err error
	defaultClient, err = client.New("")
	if err != nil {
		panic(err)
	}

	k8sConfigClient, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		panic(err)
	} else if k8sConfigClient == nil {
		panic("api config client is nil")
	}
	kubeContext = k8sConfigClient.CurrentContext

	configClient, err := config.New("")
	if err != nil {
		panic(err)
	}

	clusterClient = cluster.New(cluster.Kubeconfig{Path: kubeconfigPath, Context: kubeContext}, configClient)

	namespace, err = clusterClient.Proxy().CurrentNamespace()
	if err != nil {
		panic(err)
	}

	runtimeClient, err = clusterClient.Proxy().NewClient()
	if err != nil {
		panic(err)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "The port to listen on")
	flag.Parse()

	http.Handle("/api/v1/multicluster/", http.HandlerFunc(handleMultiClusterTree))
	http.Handle("/api/v1/custom-resource/", http.HandlerFunc(handleCustomResourceTree))
	http.Handle("/api/v1/cluster-resources/", http.HandlerFunc(handleClusterResourceTree))

	stripped, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalln(err)
	}

	frontendFS := http.FileServer(http.FS(stripped))
	http.Handle("/", frontendFS)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleMultiClusterTree(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Getting multicluster tree\n")
	tree, err := internal.ConstructMultiClusterTree(clusterClient, k8sConfigClient)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if tree != nil {
		marshalled, err := json.MarshalIndent(*tree, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, bytes.NewReader(marshalled))
	}
}

func handleClusterResourceTree(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Path[len("/api/v1/cluster-resources/"):]
	// fmt.Printf("Getting object tree for %s\n", clusterName)

	dcOptions := client.DescribeClusterOptions{
		Kubeconfig:              client.Kubeconfig{Path: kubeconfigPath, Context: kubeContext},
		Namespace:               "",
		ClusterName:             clusterName,
		ShowOtherConditions:     "",
		ShowMachineSets:         true,
		ShowClusterResourceSets: true,
		ShowTemplates:           true,
		Echo:                    true,
	}

	tree, err := internal.ConstructClusterResourceTree(defaultClient, dcOptions)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if tree != nil {
		marshalled, err := json.MarshalIndent(*tree, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, bytes.NewReader(marshalled))
	}
}

func handleCustomResourceTree(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Getting custom resource tree\n")
	// fmt.Println("GET params were:", r.URL.Query())
	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")

	object, err := internal.GetCustomResource(runtimeClient, kind, apiVersion, namespace, name)
	if err != nil {
		fmt.Println("Failed to get CRD:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := object.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(data))
}

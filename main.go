package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
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
	err := initClients()
	if err != nil {
		log.Println(err) // Allow app to start even if initClients return an error
	}
}

func initClients() error {
	var err error

	if defaultClient == nil {
		defaultClient, err = client.New("")
		if err != nil {
			return err
		}
	}

	if clusterClient == nil {
		configClient, err := config.New("")
		if err != nil {
			return err
		} else if k8sConfigClient == nil {
			return err
		}

		clusterClient = cluster.New(cluster.Kubeconfig{Path: kubeconfigPath, Context: kubeContext}, configClient)
		// log.Println("Using kubeconfig context:", clusterClient.Kubeconfig().Context)
		// log.Println("Using kubeconfig path:", clusterClient.Kubeconfig().Path)

		contexts, err := clusterClient.Proxy().GetContexts("")
		if err != nil {
			log.Println("Error fetching contexts:", err)
			return errors.New("unable to find kubecontexts, is the management cluster running?")
		}
		log.Println("Contexts:", contexts)
		if len(contexts) == 0 {
			return errors.New("no kubecontexts available, is the management cluster running?")
		}
	}

	if runtimeClient == nil {
		runtimeClient, err = clusterClient.Proxy().NewClient()
		if err != nil {
			return err
		}
	}

	namespace, err = clusterClient.Proxy().CurrentNamespace()
	if err != nil {
		return err
	}

	if k8sConfigClient == nil {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		rules.ExplicitPath = clusterClient.Kubeconfig().Path
		k8sConfigClient, err = rules.Load()
		if err != nil {
			return err
		} else if k8sConfigClient == nil {
			return err
		}
	}

	return nil
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

	log.Printf("Listening on port %d\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleMultiClusterTree(w http.ResponseWriter, r *http.Request) {
	log.Println("GET call to " + r.URL.Path)

	// Attempt to initialize any clients that are nil
	err := initClients()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tree, err := internal.ConstructMultiClusterTree(clusterClient, k8sConfigClient)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	log.Println("GET call to " + r.URL.Path)
	clusterName := r.URL.Path[len("/api/v1/cluster-resources/"):]

	// Uncomment these fields when changes merge to CAPI main
	dcOptions := client.DescribeClusterOptions{
		Kubeconfig:          client.Kubeconfig{Path: kubeconfigPath, Context: kubeContext},
		Namespace:           "",
		ClusterName:         clusterName,
		ShowOtherConditions: "",
		ShowMachineSets:     true,
		Echo:                true,
		Grouping:            false,
		// ShowClusterResourceSets: true,
		// ShowTemplates:           true,
	}

	tree, err := internal.ConstructClusterResourceTree(defaultClient, dcOptions)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	log.Println("GET call to " + r.URL.Path)
	log.Println("GET params are: ", r.URL.Query())
	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")

	object, err := internal.GetCustomResource(runtimeClient, kind, apiVersion, namespace, name)
	if err != nil {
		log.Println("Failed to get CRD:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := object.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(data))
}

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
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed web/dist
var frontend embed.FS

type Client struct {
	DefaultClient    client.Client
	ClusterClient    cluster.Client
	RuntimeClient    ctrlclient.Client
	K8sConfigClient  *api.Config
	CurrentNamespace string
}

var c *Client

var kubeconfigPath = ""
var kubeContext = ""

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var httpErr *internal.HTTPError
	c, httpErr = newClient()
	if httpErr != nil {
		log.Println(httpErr) // Try to initialize client but allow GUI to start anyway even if it fails
	}
}

func newClient() (*Client, *internal.HTTPError) {
	c := &Client{}
	var err error

	c.DefaultClient, err = client.New("")
	if err != nil {
		return nil, internal.NewInternalError(err)
	}

	configClient, err := config.New("")
	if err != nil {
		return nil, internal.NewInternalError(err)
	}

	c.ClusterClient = cluster.New(cluster.Kubeconfig{Path: kubeconfigPath, Context: kubeContext}, configClient)
	log.Println("Using kubeconfig context:", c.ClusterClient.Kubeconfig().Context)
	log.Println("Using kubeconfig path:", c.ClusterClient.Kubeconfig().Path)

	err = c.ClusterClient.Proxy().CheckClusterAvailable()
	if err != nil {
		log.Println("Cluster unavailable:", err)
		return nil, &internal.HTTPError{Status: http.StatusNotFound, Message: err.Error()}
	}

	c.RuntimeClient, err = c.ClusterClient.Proxy().NewClient()
	if err != nil {
		return nil, internal.NewInternalError(err)
	}

	c.CurrentNamespace, err = c.ClusterClient.Proxy().CurrentNamespace()
	if err != nil {
		return nil, internal.NewInternalError(err)
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = c.ClusterClient.Kubeconfig().Path
	c.K8sConfigClient, err = rules.Load()
	if err != nil {
		return nil, internal.NewInternalError(err)
	} else if c.K8sConfigClient == nil {
		return nil, internal.NewInternalError(err)
	}

	return c, nil
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

	// Attempt to initialize clients
	c, httpErr := newClient()
	if httpErr != nil {
		log.Println(httpErr)
		http.Error(w, httpErr.Error(), httpErr.Status)
		return
	}

	// TODO: should we pass in the runtimeClient here or regenerate it in the function?
	tree, httpErr := internal.ConstructMultiClusterTree(c.ClusterClient, c.K8sConfigClient)
	if httpErr != nil {
		log.Println(httpErr)
		http.Error(w, httpErr.Error(), httpErr.Status)
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

	tree, httpErr := internal.ConstructClusterResourceTree(c.DefaultClient, dcOptions)
	if httpErr != nil {
		log.Println(httpErr)
		http.Error(w, httpErr.Error(), httpErr.Status)
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

	// TODO: should the runtimeClient be regenerated here?
	object, httpErr := internal.GetCustomResource(c.RuntimeClient, kind, apiVersion, c.CurrentNamespace, name)
	if httpErr != nil {
		log.Println(httpErr)
		http.Error(w, httpErr.Error(), httpErr.Status)
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

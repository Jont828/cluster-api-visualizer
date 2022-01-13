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

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

//go:embed web/dist
var frontend embed.FS

var kubeconfig string = ""

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
	fmt.Printf("Getting multicluster tree\n")
	tree, err := client.MultiDiscovery(kubeconfig)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	exportTree := internal.ConstructMultiClusterTree(tree)
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

func handleClusterResourceTree(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Path[len("/api/v1/cluster-resources/"):]
	fmt.Printf("Getting object tree for %s\n", clusterName)

	tree, err := internal.ConstructClusterResourceTree(clusterName)
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
	fmt.Printf("Getting custom resource tree\n")
	fmt.Println("GET params were:", r.URL.Query())
	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")

	object, err := internal.GetCustomResource(kind, apiVersion, "default", name)
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

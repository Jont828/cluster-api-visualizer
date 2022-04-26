package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"

	// "log"
	"net/http"
	"os"
	"strings"

	"github.com/Jont828/cluster-api-visualizer/internal"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

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
	log := klogr.New()

	var httpErr *internal.HTTPError
	c, httpErr = newClient()
	if httpErr != nil {
		log.Error(httpErr, "failed to initialize client, will allow frontend to start") // Try to initialize client but allow GUI to start anyway even if it fails
	}
}

func newClient() (*Client, *internal.HTTPError) {
	log := klogr.New()

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

	err = c.ClusterClient.Proxy().CheckClusterAvailable()
	if err != nil {
		log.Error(err, "failed to check cluster availability for cluster client")
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
	var host string
	var port int

	flag.StringVar(&host, "host", "localhost", "Host to listen on")
	flag.IntVar(&port, "port", 8081, "The port to listen on")

	klog.InitFlags(nil)
	flag.Set("v", "2")

	flag.Parse()

	log := klogr.New()

	http.Handle("/api/v1/multicluster/", http.HandlerFunc(handleMultiClusterTree))
	http.Handle("/api/v1/custom-resource/", http.HandlerFunc(handleCustomResourceTree))
	http.Handle("/api/v1/cluster-resources/", http.HandlerFunc(handleClusterResourceTree))

	var frontend fs.FS = os.DirFS("web/dist")
	httpFS := http.FS(frontend)
	fileServer := http.FileServer(httpFS)
	serveIndex := serveFileContents("index.html", httpFS)

	http.Handle("/", intercept404(fileServer, serveIndex))

	uri := fmt.Sprintf("%s:%d", host, port)
	log.V(2).Info(fmt.Sprintf("Listening at http://%s", uri))
	if host == "0.0.0.0" {
		log.V(2).Info(fmt.Sprintf("View at http://localhost:%d in browser", port))
	}

	http.ListenAndServe(uri, nil)
}

type hookedResponseWriter struct {
	http.ResponseWriter
	got404 bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	log := klogr.New()

	log.V(4).Info("Writing header", "status", status)
	if status == http.StatusNotFound {
		// Don't actually write the 404 header, just set a flag.
		hrw.got404 = true
	} else {
		hrw.ResponseWriter.WriteHeader(status)
	}
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
	log := klogr.New()

	log.V(4).Info("Writing content", "content", string(p))
	if hrw.got404 {
		// No-op, but pretend that we wrote len(p) bytes to the writer.
		return len(p), nil
	}

	return hrw.ResponseWriter.Write(p)
}

func intercept404(handler, on404 http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hookedWriter := &hookedResponseWriter{ResponseWriter: w}
		handler.ServeHTTP(hookedWriter, r)

		if hookedWriter.got404 {
			on404.ServeHTTP(w, r)
		}
	})
}

func serveFileContents(file string, files http.FileSystem) http.HandlerFunc {
	log := klogr.New()

	log.V(4).Info("Serving file", "filename", file)
	return func(w http.ResponseWriter, r *http.Request) {
		// Restrict only to instances where the browser is looking for an HTML file
		if !strings.Contains(r.Header.Get("Accept"), "text/html") {
			log.V(4).Info("404 file not found", "filename", file)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 not found")
			return
		}

		// Open the file and return its contents using http.ServeContent
		index, err := files.Open(file)
		if err != nil {
			log.Error(err, "open file error")

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "`%s` not found", file)

			return
		}

		fi, err := index.Stat()
		if err != nil {
			log.Error(err, "stat file error")

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "`%s` not found", file)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), index)
	}
}

func handleMultiClusterTree(w http.ResponseWriter, r *http.Request) {
	log := klogr.New()

	log.V(2).Info("GET call to url", "url", r.URL.Path)

	// Attempt to initialize clients
	c, httpErr := newClient()
	if httpErr != nil {
		log.Error(httpErr, "failed to initialize clients")
		http.Error(w, httpErr.Error(), httpErr.Status)
		return
	}

	// TODO: should we pass in the runtimeClient here or regenerate it in the function?
	tree, httpErr := internal.ConstructMultiClusterTree(c.ClusterClient, c.K8sConfigClient)
	if httpErr != nil {
		log.Error(httpErr, "failed to construct management cluster tree view")
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
	log := klogr.New()

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	clusterName := r.URL.Path[len("/api/v1/cluster-resources/"):]

	// Uncomment these fields when changes merge to CAPI main
	dcOptions := client.DescribeClusterOptions{
		Kubeconfig:              client.Kubeconfig{Path: kubeconfigPath, Context: kubeContext},
		Namespace:               "",
		ClusterName:             clusterName,
		ShowOtherConditions:     "",
		ShowMachineSets:         true,
		Echo:                    true,
		Grouping:                false,
		ShowClusterResourceSets: true,
		ShowTemplates:           true,
	}

	tree, httpErr := internal.ConstructClusterResourceTree(c.DefaultClient, dcOptions)
	if httpErr != nil {
		log.Error(httpErr, "failed to construct resource tree for target cluster", "clusterName", clusterName)
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
	log := klogr.New()

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	log.V(2).Info("GET call params are", "params", r.URL.Query())

	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")

	// TODO: should the runtimeClient be regenerated here?
	object, httpErr := internal.GetCustomResource(c.RuntimeClient, kind, apiVersion, c.CurrentNamespace, name)
	if httpErr != nil {
		log.Error(httpErr, "failed to construct tree for custom resource", "kind", kind, "name", name)
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

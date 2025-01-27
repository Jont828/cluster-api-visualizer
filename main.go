package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/Jont828/cluster-api-visualizer/internal"
	"github.com/Jont828/cluster-api-visualizer/version"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	configclient "sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Client struct {
	ClusterctlClient        client.Client     // The clusterctl client needed to run clusterctl operations like `c.DescribeCluster()`
	ClusterClient           cluster.Client    // The client used by clusterctl to interact with the management cluster
	ControllerRuntimeClient ctrlclient.Client // The Kubernetes controller-runtime client needed to run `external.Get()` to fetch any CRD as a JSON object
	K8sConfigClient         *api.Config       // This is the Kubernetes config client needed to access information from the kubeconfig like the namespace and context
	CurrentNamespace        string
}

var c *Client

var kubeconfigPath = ""
var kubeContext = ""
var clusterctlConfigPath = ""

func newClient(ctx context.Context) (*Client, *internal.HTTPError) {
	log := ctrl.LoggerFrom(ctx)

	c := &Client{}
	var err error

	clusterKubeconfig := cluster.Kubeconfig{Path: kubeconfigPath, Context: kubeContext}

	c.ClusterctlClient, err = client.New(ctx, clusterctlConfigPath)
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	}

	configClient, err := config.New(ctx, clusterctlConfigPath)
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	}

	clusterClient := cluster.New(clusterKubeconfig, configClient)
	c.ClusterClient = clusterClient

	err = clusterClient.Proxy().CheckClusterAvailable(ctx)
	if err != nil {
		log.Error(err, "failed to check cluster availability for cluster client")
		return nil, &internal.HTTPError{Status: http.StatusNotFound, Message: err.Error()}
	}

	c.ControllerRuntimeClient, err = clusterClient.Proxy().NewClient(ctx)
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	}

	c.CurrentNamespace, err = clusterClient.Proxy().CurrentNamespace()
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = clusterClient.Kubeconfig().Path
	c.K8sConfigClient, err = rules.Load()
	if err != nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	} else if c.K8sConfigClient == nil {
		log.Error(err, "failed to create client")
		return nil, internal.NewInternalError(err)
	}

	return c, nil
}

func main() {
	var host string
	var port int
	var generateConfig bool

	flag.StringVar(&host, "host", "localhost", "Host to listen on")
	flag.IntVar(&port, "port", 8081, "The port to listen on")
	flag.BoolVar(&generateConfig, "generate-config", false, "Generate a kubeconfig file and write it to disk. Useful for running inside a pod within a cluster.")

	klog.InitFlags(nil)
	flag.Set("v", "2")

	flag.Parse()

	ctrl.SetLogger(klog.Background())
	ctx := ctrl.SetupSignalHandler()

	log := ctrl.LoggerFrom(ctx)

	log.Info("Starting app with version", "version", version.Get().String())

	if generateConfig {
		log.V(2).Info("Generating kubeconfig file")
		restConfig := configclient.GetConfigOrDie()

		apiConfig, err := internal.ConstructInClusterKubeconfig(ctx, restConfig, "")
		if err != nil {
			log.Error(err, "error constructing in-cluster kubeconfig")
			return
		}
		filePath := "tmp/management.kubeconfig"
		if err = internal.WriteKubeconfigToFile(ctx, filePath, *apiConfig); err != nil {
			log.Error(err, "error writing kubeconfig to file")
			return
		}
		kubeconfigPath = filePath
		kubeContext = apiConfig.CurrentContext
	}

	var httpErr *internal.HTTPError
	c, httpErr = newClient(ctx)
	if httpErr != nil {
		log.Error(httpErr, "failed to initialize client, will allow frontend to start") // Try to initialize client but allow GUI to start anyway even if it fails
	}

	http.Handle("/api/v1/management-cluster/", http.HandlerFunc(handleManagementClusterTree))
	http.Handle("/api/v1/custom-resource-definition/", http.HandlerFunc(handleCustomResourceDefinitionTree))
	http.Handle("/api/v1/resource-logs/", http.HandlerFunc(handleGetResourceLogs))
	http.Handle("/api/v1/describe-cluster/", http.HandlerFunc(handleDescribeClusterTree))
	http.Handle("/api/v1/cluster-kubeconfig/", http.HandlerFunc(handleClusterKubeConfig))
	http.Handle("/api/v1/version/", http.HandlerFunc(handleGetVersion))

	var frontend fs.FS = os.DirFS("web/dist")
	httpFS := http.FS(frontend)
	fileServer := http.FileServer(httpFS)
	serveIndex := serveFileContents(ctx, "index.html", httpFS)

	http.Handle("/", intercept404(fileServer, serveIndex))

	uri := fmt.Sprintf("%s:%d", host, port)
	log.V(2).Info(fmt.Sprintf("KYLE Listening at http://%s", uri))
	if host == "0.0.0.0" {
		log.V(2).Info(fmt.Sprintf("View at http://localhost:%d in browser", port))
	}

	srv := &http.Server{
		Addr: uri,
		// Pass root context to the server so it gets propagated to all requests.
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	// srv.Handler is nil so it uses default serve mux, which http.Handle configures by default.

	srv.ListenAndServe()
}

type hookedResponseWriter struct {
	http.ResponseWriter
	got404 bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	if status == http.StatusNotFound {
		// Don't actually write the 404 header, just set a flag.
		hrw.got404 = true
	} else {
		hrw.ResponseWriter.WriteHeader(status)
	}
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
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

func serveFileContents(ctx context.Context, file string, files http.FileSystem) http.HandlerFunc {
	log := ctrl.LoggerFrom(ctx)

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

		r = r.WithContext(ctx)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), index)
	}
}

func handleManagementClusterTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)

	// Attempt to initialize clients
	c, httpErr := newClient(ctx)
	if httpErr != nil {
		log.Error(httpErr, "failed to initialize clients")
		http.Error(w, httpErr.Error(), httpErr.Status)
		return
	}

	tree, httpErr := internal.ConstructMultiClusterTree(ctx, c.ControllerRuntimeClient, c.K8sConfigClient)
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

func handleClusterKubeConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	log.V(2).Info("GET call params are", "params", r.URL.Query())
	name := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")

	found, secret, err := internal.GetKubeConfigSecret(ctx, name, namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if found {
		w.Header().Set("Content-Type", "text/vnd.yaml")
		io.Copy(w, bytes.NewReader(secret.Data["value"]))
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func handleDescribeClusterTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	log.V(2).Info("GET call params are", "params", r.URL.Query())
	name := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")

	dcOptions := client.DescribeClusterOptions{
		Kubeconfig:              client.Kubeconfig{Path: kubeconfigPath, Context: kubeContext},
		Namespace:               namespace,
		ClusterName:             name,
		ShowOtherConditions:     "",
		ShowMachineSets:         true,
		Echo:                    true,
		Grouping:                false,
		AddTemplateVirtualNode:  true,
		ShowClusterResourceSets: true,
		ShowTemplates:           true,
	}

	tree, httpErr := internal.ConstructClusterResourceTree(ctx, c.ClusterctlClient, c.ControllerRuntimeClient, dcOptions)
	if httpErr != nil {
		log.Error(httpErr, "failed to construct resource tree for target cluster", "clusterName", name)
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

func handleCustomResourceDefinitionTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	log.V(2).Info("GET call params are", "params", r.URL.Query())

	kind := r.URL.Query().Get("kind")
	apiVersion := r.URL.Query().Get("apiVersion")
	name := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")

	// TODO: should the runtimeClient be regenerated here?
	object, httpErr := internal.GetCustomResource(ctx, c.ControllerRuntimeClient, kind, apiVersion, namespace, name)
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

func handleGetResourceLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)
	log.V(2).Info("GET call params are", "params", r.URL.Query())

	kind := r.URL.Query().Get("kind")
	name := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")

	config, err := c.ClusterClient.Proxy().GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	logs, err := internal.GetPodLogsForResource(ctx, c.ControllerRuntimeClient, config, kind, namespace, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := json.Marshal(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(data))
}

func handleGetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctrl.LoggerFrom(ctx)

	log.V(2).Info("GET call to url", "url", r.URL.Path)

	versionInfo := version.Get()

	data, err := json.MarshalIndent(versionInfo, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(data))
}

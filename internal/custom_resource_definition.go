package internal

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	"sigs.k8s.io/cluster-api/controllers/external"
)

func GetCustomResource(kind string, apiVersion string, namespace string, name string) (*unstructured.Unstructured, error) {
	kubeconfig := "/home/jonathan/.kube/config"

	cfgFile := ""
	configClient, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}

	clusterClient := cluster.New(cluster.Kubeconfig{Path: kubeconfig, Context: ""}, configClient)

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

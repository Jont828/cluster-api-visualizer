package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/controllers/external"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetCustomResource(runtimeClient ctrlclient.Client, kind string, apiVersion string, namespace string, name string) (*unstructured.Unstructured, error) {
	objectRef := corev1.ObjectReference{
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
		APIVersion: apiVersion,
	}
	object, err := external.Get(context.TODO(), runtimeClient, &objectRef, namespace)
	if err != nil {
		return nil, err
	}

	// data, err := object.MarshalJSON()
	// if err != nil {
	// 	return nil, err
	// }

	// out, err := json.MarshalIndent(string(data), "", "\t")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(out))

	return object, nil
}

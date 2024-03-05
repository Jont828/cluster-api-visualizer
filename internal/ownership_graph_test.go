package internal

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRemoveTransitiveOwners(t *testing.T) {
	objects := map[types.UID]ctrlclient.Object{
		"1": &unstructured.Unstructured{Object: map[string]interface{}{"kind": "Cluster"}},
		"2": &unstructured.Unstructured{Object: map[string]interface{}{"kind": "Machine"}},
		"3": &unstructured.Unstructured{Object: map[string]interface{}{"kind": "MachineSet"}},
		"4": &unstructured.Unstructured{Object: map[string]interface{}{"kind": "MachineDeployment"}},
	}
	startNode := "1"

	testcases := []struct {
		name           string
		ownerRefs      map[types.UID]map[types.UID]struct{}
		expectedOwners map[types.UID]struct{} // expected owners of just the start node
	}{
		{
			name: "no connected edges",
			ownerRefs: map[types.UID]map[types.UID]struct{}{
				"1": {"2": {}},
				"3": {"4": {}},
			},
			expectedOwners: map[types.UID]struct{}{"2": {}},
		},
		{
			name: "owned by one owner and a transitive owner",
			ownerRefs: map[types.UID]map[types.UID]struct{}{
				"1": {"2": {}, "5": {}},
				"2": {"3": {}},
				"3": {"4": {}},
				"4": {"5": {}},
			},
			expectedOwners: map[types.UID]struct{}{"2": {}},
		},
		{
			name: "multiple owners left after transitive removal",
			ownerRefs: map[types.UID]map[types.UID]struct{}{
				"1": {"2": {}, "3": {}, "4": {}, "5": {}},
				"2": {"3": {}},
				"4": {"5": {}},
			},
			expectedOwners: map[types.UID]struct{}{"2": {}, "4": {}},
		},
	}

	for _, testcase := range testcases {
		test := testcase
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			objectGraph := &OwnershipGraph{
				Objects:   objects,
				OwnerRefs: test.ownerRefs,
			}

			RemoveTransitiveOwners(types.UID(startNode), objectGraph)
			g.Expect(objectGraph.OwnerRefs[types.UID(startNode)]).To(BeEquivalentTo(test.expectedOwners))
		})
	}
}

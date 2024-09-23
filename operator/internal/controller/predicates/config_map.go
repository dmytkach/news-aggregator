package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ConfigMapNamePredicate creates a predicate that filters Kubernetes events based on the name of a ConfigMap.
//
// The predicate checks whether the name of the ConfigMap involved in the event matches the provided name.
// This predicate is applied to various types of events such as Create, Update, Delete, and Generic events.
func ConfigMapNamePredicate(configMapName string) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return e.Object.GetName() == configMapName
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetName() == configMapName
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Object.GetName() == configMapName
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return e.Object.GetName() == configMapName
		},
	}
}

package predicates_test

import (
	"com.teamdev/news-aggregator/internal/controller/predicates"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"testing"
)

func TestPredicates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Predicates Suite")
}

var _ = BeforeSuite(func() {
	_ = v1.AddToScheme(scheme.Scheme)
})
var _ = Describe("ConfigMapNamePredicate", func() {
	var (
		predicateFn   predicate.Predicate
		configMapName string
	)

	BeforeEach(func() {
		configMapName = "test-configmap"
		predicateFn = predicates.ConfigMapNamePredicate(configMapName)
	})

	Context("when processing CreateEvent", func() {
		It("should return true when ConfigMap name matches", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName(configMapName)

			createEvent := event.CreateEvent{
				Object: obj,
			}

			Expect(predicateFn.Create(createEvent)).To(BeTrue())
		})
		It("should return false when ConfigMap name does not match", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName("different-configmap")

			createEvent := event.CreateEvent{
				Object: obj,
			}

			Expect(predicateFn.Create(createEvent)).To(BeFalse())
		})
	})

	Context("when processing UpdateEvent", func() {
		It("should return true when the new ConfigMap name matches", func() {
			objNew := &unstructured.Unstructured{}
			objNew.SetName(configMapName)

			updateEvent := event.UpdateEvent{
				ObjectNew: objNew,
			}

			Expect(predicateFn.Update(updateEvent)).To(BeTrue())
		})
		It("should return false when the new ConfigMap name does not match", func() {
			objNew := &unstructured.Unstructured{}
			objNew.SetName("different-configmap")

			updateEvent := event.UpdateEvent{
				ObjectNew: objNew,
			}

			Expect(predicateFn.Update(updateEvent)).To(BeFalse())
		})
	})

	Context("when processing DeleteEvent", func() {
		It("should return true when ConfigMap name matches", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName(configMapName)

			deleteEvent := event.DeleteEvent{
				Object: obj,
			}

			Expect(predicateFn.Delete(deleteEvent)).To(BeTrue())
		})
		It("should return false when ConfigMap name does not match", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName("different-configmap")

			deleteEvent := event.DeleteEvent{
				Object: obj,
			}

			Expect(predicateFn.Delete(deleteEvent)).To(BeFalse())
		})
	})

	Context("when processing GenericEvent", func() {
		It("should return true when ConfigMap name matches", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName(configMapName)

			genericEvent := event.GenericEvent{
				Object: obj,
			}

			Expect(predicateFn.Generic(genericEvent)).To(BeTrue())
		})
		It("should return false when ConfigMap name does not match", func() {
			obj := &unstructured.Unstructured{}
			obj.SetName("different-configmap")

			genericEvent := event.GenericEvent{
				Object: obj,
			}

			Expect(predicateFn.Generic(genericEvent)).To(BeFalse())
		})
	})
})

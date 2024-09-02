package v1_test

import (
	v12 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	v13 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ = Describe("ConfigMapWebHook", func() {
	var (
		webhook v12.ConfigMapWebHook
		client  client.Client
	)

	BeforeEach(func() {
		client = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		webhook = v12.ConfigMapWebHook{
			Client:  client,
			Decoder: admission.NewDecoder(scheme.Scheme),
		}
	})

	Context("when handling a ConfigMap admission request", func() {
		It("should return an error if a referenced feed does not exist", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"feeds": "non-existent-feed",
				},
			}
			rawConfigMap, err := json.Marshal(configMap)
			Expect(err).NotTo(HaveOccurred())

			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Raw: rawConfigMap,
					},
					Operation: v13.Create,
					Namespace: "default",
				},
			}

			response := webhook.Handle(context.TODO(), req)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Code).To(Equal(int32(http.StatusNotFound)))
			Expect(response.Result.Message).To(ContainSubstring("feed 'non-existent-feed' not found"))
		})

		It("should succeed if feed exist", func() {
			feed := &v12.Feed{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "feed1",
					Namespace: "default",
				},
			}
			err := client.Create(context.Background(), feed)
			Expect(err).NotTo(HaveOccurred())

			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"feeds": "feed1",
				},
			}

			rawConfigMap, err := json.Marshal(configMap)
			Expect(err).NotTo(HaveOccurred())

			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Raw: rawConfigMap,
					},
					Operation: v13.Create,
					Namespace: "default",
				},
			}

			response := webhook.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeTrue())
		})

		It("should return an error if the ConfigMap cannot be decoded", func() {
			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Object: &v1.ConfigMap{},
					},
					Operation: v13.Create,
					Namespace: "default",
				},
			}

			response := webhook.Handle(context.Background(), req)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Code).To(Equal(int32(http.StatusBadRequest)))
			Expect(response.Result.Message).To(ContainSubstring("there is no content to decode"))
		})
	})
})

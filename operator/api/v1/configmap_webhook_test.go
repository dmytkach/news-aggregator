package v1_test

import (
	v12 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"errors"
	v13 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

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
	)

	BeforeEach(func() {
		client := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		webhook = v12.ConfigMapWebHook{
			Client:  client,
			Decoder: admission.NewDecoder(scheme.Scheme),
		}
	})

	Context("when handling a CREATE/UPDATE admission request", func() {
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
			err := webhook.Client.Create(context.Background(), feed)
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

	Context("when handling a DELETE admission request", func() {
		It("should allow deletion if no FeedGroups are used as keys in the ConfigMap", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			}

			rawConfigMap, err := json.Marshal(configMap)
			Expect(err).NotTo(HaveOccurred())

			hotNews := &v12.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews1",
					Namespace: "default",
				},
				Spec: v12.HotNewsSpec{
					FeedGroups: []string{"non-existent-feedgroup"},
				},
			}
			err = webhook.Client.Create(context.TODO(), hotNews)
			Expect(err).NotTo(HaveOccurred())

			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Raw: rawConfigMap,
					},
					Namespace: "default",
					Operation: v13.Delete,
				},
			}

			response := webhook.Handle(context.TODO(), req)
			Expect(response.Allowed).To(BeTrue())
			Expect(response.Result.Code).To(Equal(int32(http.StatusOK)))
		})
		It("should deny deletion if a FeedGroup exists as a key in the ConfigMap", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"feedgroup1": "value",
				},
			}

			rawConfigMap, err := json.Marshal(configMap)
			Expect(err).NotTo(HaveOccurred())

			hotNews := &v12.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews1",
					Namespace: "default",
				},
				Spec: v12.HotNewsSpec{
					FeedGroups: []string{"feedgroup1"},
				},
			}

			err = webhook.Client.Create(context.TODO(), hotNews)
			Expect(err).NotTo(HaveOccurred())

			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Raw: rawConfigMap,
					},
					Namespace: "default",
					Operation: v13.Delete,
				},
			}

			response := webhook.Handle(context.TODO(), req)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Code).To(Equal(int32(http.StatusForbidden)))
			Expect(response.Result.Message).To(ContainSubstring("ConfigMap 'test-configmap' contains feed group 'feedgroup1', deletion is not allowed"))
		})
		It("should deny deletion if unable to retrieve HotNews", func() {
			fakeClient := fake.NewClientBuilder().WithScheme(scheme.Scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("error with getting feeds")
					},
				}).Build()
			webhook.Client = fakeClient
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"feedgroup1": "value",
				},
			}

			rawConfigMap, err := json.Marshal(configMap)
			Expect(err).NotTo(HaveOccurred())

			req := admission.Request{
				AdmissionRequest: v13.AdmissionRequest{
					Object: runtime.RawExtension{
						Raw: rawConfigMap,
					},
					Namespace: "default",
					Operation: v13.Delete,
				},
			}

			response := webhook.Handle(context.TODO(), req)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Code).To(Equal(int32(http.StatusInternalServerError)))
		})
	})
})

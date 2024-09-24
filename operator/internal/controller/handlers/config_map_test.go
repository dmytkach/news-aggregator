package handlers_test

import (
	"com.teamdev/news-aggregator/internal/controller/handlers"
	"context"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("ConfigMapHandler", func() {
	var (
		fakeClient    client.Client
		handler       *handlers.ConfigMapHandler
		ctx           context.Context
		cancelFunc    context.CancelFunc
		configMapName string
		configMap     *v1.ConfigMap
		feed          *aggregatorv1.Feed
	)

	BeforeEach(func() {
		configMapName = "test-configmap"
		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: "default",
			},
			Data: map[string]string{
				"group": "test-feed",
			},
		}
		feed = &aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
		}
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		ctx, cancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
		handler = &handlers.ConfigMapHandler{
			Client: fakeClient,
		}
	})

	AfterEach(func() {
		cancelFunc()
	})

	Context("when the object is not a ConfigMap", func() {
		It("should return an empty list of requests", func() {
			otherObj := &aggregatorv1.HotNews{}

			requests := handler.Handle(ctx, otherObj)

			Expect(requests).To(BeEmpty())
		})
	})

	Context("when the ConfigMap name does not match", func() {
		It("should return an empty list of requests", func() {
			nonMatchingConfigMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-configmap",
					Namespace: "default",
				},
			}

			Expect(fakeClient.Create(context.Background(), nonMatchingConfigMap)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), configMap)).To(Succeed())

			requests := handler.Handle(ctx, nonMatchingConfigMap)

			Expect(requests).To(BeEmpty())
		})
	})

	Context("when listing HotNews fails", func() {
		It("should return an empty list of requests", func() {
			handler.Client = fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()

			Expect(fakeClient.Create(context.Background(), configMap)).To(Succeed())

			requests := handler.Handle(ctx, configMap)

			Expect(requests).To(BeEmpty())
		})
	})

	Context("when HotNews list is successfully retrieved", func() {
		var hotNews *aggregatorv1.HotNews

		BeforeEach(func() {
			hotNews = &aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hot-news-1",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					FeedGroups: []string{"group"},
				},
			}
		})

		It("should return reconcile requests if HotNews is associated with ConfigMap feeds", func() {
			Expect(fakeClient.Create(context.Background(), configMap)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), hotNews)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), feed)).To(Succeed())

			requests := handler.Handle(ctx, configMap)

			Expect(requests).To(HaveLen(1))
			Expect(requests[0].NamespacedName).To(Equal(types.NamespacedName{
				Namespace: "default",
				Name:      "hot-news-1",
			}))
		})

		It("should return an empty list if HotNews is not associated with ConfigMap feeds", func() {
			hotNews.Namespace = "different-namespace"
			Expect(fakeClient.Create(context.Background(), hotNews)).To(Succeed())

			requests := handler.Handle(ctx, configMap)

			Expect(requests).To(BeEmpty())
		})
	})
})

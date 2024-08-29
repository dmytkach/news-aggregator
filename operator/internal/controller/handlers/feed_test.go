package handlers_test

import (
	"com.teamdev/news-aggregator/internal/controller/handlers"
	"context"
	. "github.com/onsi/ginkgo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("FeedHandler", func() {

	var (
		fakeClient client.Client
		handler    *handlers.FeedHandler
		ctx        context.Context
		cancelFunc context.CancelFunc
		feed       *aggregatorv1.Feed
	)

	BeforeEach(func() {
		feed = &aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
		}
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		ctx, cancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
		handler = &handlers.FeedHandler{
			Client: fakeClient,
		}
	})

	AfterEach(func() {
		cancelFunc()
	})

	Context("when the object is not a Feed", func() {
		It("should return an empty list of requests", func() {
			otherObj := &aggregatorv1.HotNews{}

			requests := handler.Handle(ctx, otherObj)

			Expect(requests).To(BeEmpty())
		})
	})

	Context("when listing HotNews fails", func() {
		It("should return an empty list of requests", func() {
			handler.Client = fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()

			requests := handler.Handle(ctx, feed)

			Expect(requests).To(BeEmpty())
		})
	})

	Context("when HotNews list is successfully retrieved", func() {
		It("should return reconcile requests for each matching HotNews", func() {

			hotNews := &aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hot-news-1",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Feeds: []string{"test-feed"},
				},
			}

			Expect(fakeClient.Create(context.Background(), feed)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), hotNews)).To(Succeed())

			requests := handler.Handle(ctx, feed)

			Expect(requests).To(HaveLen(1))
			Expect(requests[0].NamespacedName).To(Equal(types.NamespacedName{
				Namespace: "default",
				Name:      "hot-news-1",
			}))
		})
	})

	Context("when Feed and HotNews are in different namespaces", func() {
		It("should return an empty list of requests", func() {
			hotNews := &aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hot-news-1",
					Namespace: "other",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Feeds: []string{"test-feed"},
				},
			}

			Expect(fakeClient.Create(context.Background(), feed)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), hotNews)).To(Succeed())

			requests := handler.Handle(ctx, feed)

			Expect(requests).To(BeEmpty())
		})
	})
	Context("when Feed doesn't exist in HotNews", func() {
		It("should return an empty list of requests", func() {
			hotNews := &aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hot-news-1",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Feeds: []string{"other-feed"},
				},
			}

			Expect(fakeClient.Create(context.Background(), feed)).To(Succeed())
			Expect(fakeClient.Create(context.Background(), hotNews)).To(Succeed())

			requests := handler.Handle(ctx, feed)

			Expect(requests).To(BeEmpty())
		})
	})
})

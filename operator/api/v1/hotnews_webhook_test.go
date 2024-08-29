package v1_test

import (
	v1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"time"
)

var _ = Describe("HotNews Webhook", func() {
	var (
		scheme   *runtime.Scheme
		ctx      context.Context
		testFeed *v1.Feed
		hotNews  *v1.HotNews
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		_ = v1.AddToScheme(scheme)
		v1.Client = fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = context.TODO()
		testFeed = &v1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testFeed",
				Namespace: "default",
			},
			Spec: v1.FeedSpec{
				Name: "valid-name",
				Link: "http://valid-url.com",
			},
		}
		hotNews = &v1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-hotnews",
				Namespace: "default",
			},
			Spec: v1.HotNewsSpec{
				Keywords:  []string{"news"},
				DateStart: "2024-01-01",
				DateEnd:   "2024-01-10",
				Feeds:     []string{"testFeed"},
			},
		}
	})

	Context("Default function", func() {
		It("should set default values correctly", func() {
			hotNews.Spec.Feeds = []string{}
			hotNews.Default()

			Expect(hotNews.Spec.SummaryConfig.TitlesCount).To(Equal(10), "Default TitlesCount should be set to 10")
			Expect(hotNews.Spec.Feeds).To(BeEmpty(), "Feeds should remain empty if no feeds or feed groups specified")
		})
		It("should set the Feeds field with existing feeds", func() {
			Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
			hotNews.Spec.Feeds = []string{}
			hotNews.Default()

			Expect(hotNews.Spec.Feeds).To(ContainElement("testFeed"), "Feeds should include 'feed1' as it exists in namespace")
		})
	})

	Context("ValidateCreate", func() {
		It("should pass with a valid configuration", func() {
			Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
			_, err := hotNews.ValidateCreate()
			Expect(err).NotTo(HaveOccurred(), "Valid HotNews configuration should pass validation")
		})

		It("should fail when keywords are missing", func() {
			Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
			hotNews.Spec.Keywords = []string{}
			_, err := hotNews.ValidateCreate()
			Expect(err).To(HaveOccurred(), "HotNews without keywords should fail validation")
		})
		Context("Date problems", func() {

			It("should fail when dates are in incorrect format", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = "01-01-2023"
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with incorrect date format should fail validation")
			})

			It("should fail when dateStart is in the future", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = time.Now().AddDate(0, 0, 1).Format(v1.DateFormat)
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with dateStart in the future should fail validation")
			})

			It("should fail when dateEnd is in the future", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateEnd = time.Now().AddDate(0, 0, 1).Format(v1.DateFormat)
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with dateEnd in the future should fail validation")
			})

			It("should fail when dateEnd is before dateStart", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = "2023-01-02"
				hotNews.Spec.DateEnd = "2023-01-01"
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with dateEnd before dateStart should fail validation")
			})

			It("should fail when dateStart is too old", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = "1800-01-01"
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with dateStart too old should fail validation")
			})

			It("should fail when dateEnd is too old", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateEnd = "1800-01-01"
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "HotNews with dateEnd too old should fail validation")
			})

			It("should pass when dateStart and dateEnd are valid", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = "2023-01-01"
				hotNews.Spec.DateEnd = "2023-12-31"
				_, err := hotNews.ValidateCreate()
				Expect(err).NotTo(HaveOccurred(), "HotNews with valid dateStart and dateEnd should pass validation")
			})

			It("should pass when dateStart and dateEnd are empty", func() {
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				hotNews.Spec.DateStart = ""
				hotNews.Spec.DateEnd = ""
				_, err := hotNews.ValidateCreate()
				Expect(err).NotTo(HaveOccurred(), "HotNews with empty dateStart and dateEnd should pass validation")
			})
		})
		Context("Feed problems", func() {

			It("feeds or configmap not found", func() {
				hotNews.Spec.Feeds = []string{}
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "at least one feed must be specified")
			})
			It("feeds or configmap not found", func() {
				hotNews.Spec.Feeds = []string{"feed12"}
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "feed feed12 does not exist in namespace default")
			})
			It("error with getting feeds", func() {
				v1.Client = fake.NewClientBuilder().WithScheme(scheme).
					WithInterceptorFuncs(interceptor.Funcs{
						List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
							return errors.New("error with getting feeds")
						},
						Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
							return client.Create(ctx, obj)
						}}).Build()
				Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
				_, err := hotNews.ValidateCreate()
				Expect(err).To(HaveOccurred(), "error with getting feeds")
			})
		})
	})

	Context("ValidateUpdate", func() {
		It("should pass with a valid update", func() {
			Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
			_, err := hotNews.ValidateUpdate(&v1.HotNews{})
			Expect(err).NotTo(HaveOccurred(), "Valid HotNews update should pass validation")
		})

		It("should fail when updating to have no keywords", func() {
			Expect(v1.Client.Create(ctx, testFeed)).Should(Succeed())
			hotNews.Spec.Keywords = []string{}
			_, err := hotNews.ValidateUpdate(&v1.HotNews{})
			Expect(err).To(HaveOccurred(), "Updating HotNews to have no keywords should fail validation")
		})
	})
})

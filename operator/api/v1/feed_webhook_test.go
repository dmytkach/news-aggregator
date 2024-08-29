package v1_test

import (
	v1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

var _ = Describe("FeedWebhook", func() {
	var (
		testFeed *v1.Feed
		ctx      context.Context
		c        client.Client
		old      runtime.Object
	)

	BeforeEach(func() {
		ctx = context.TODO()
		testFeed = &v1.Feed{
			ObjectMeta: v12.ObjectMeta{
				Name:      "testFeed",
				Namespace: "testFeedNamespace",
			},
			Spec: v1.FeedSpec{
				Name: "valid-name",
				Link: "http://valid-url.com",
			},
		}
		c = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		v1.Client = c
	})
	It("ValidateCreate", func() {
		warnings, err := testFeed.ValidateCreate()
		Expect(err).To(BeNil())
		Expect(warnings).To(BeNil())
	})
	It("ValidateUpdate", func() {
		warnings, err := testFeed.ValidateUpdate(old)
		Expect(err).To(BeNil())
		Expect(warnings).To(BeNil())
	})

	It("ValidateDelete", func() {
		warnings, err := testFeed.ValidateDelete()
		Expect(err).To(BeNil())
		Expect(warnings).To(BeNil())
	})

	Context("feed has an invalid name", func() {
		It("empty name", func() {
			testFeed.Spec.Name = ""
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
		})
		It("long name", func() {
			testFeed.Spec.Name = "verysuperlongnameformytestfeed"
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name must not exceed 20 characters"))
		})
	})

	Context("feed has an invalid link", func() {
		It("empty link", func() {
			testFeed.Spec.Link = ""
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("link cannot be empty"))
		})
		It("invelid link", func() {
			testFeed.Spec.Link = "invalid-url"
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("link must be a valid"))
		})

	})
	Context("feed is not unique", func() {
		It("try to create existing feed", func() {
			Expect(c.Create(ctx, testFeed)).Should(Succeed())
			testFeed.UID = "12345"
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already exists in namespace"))
		})
		It("try to create existing feed", func() {
			c = fake.NewClientBuilder().WithScheme(scheme.Scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("error with get feeds")
					},
				}).Build()
			v1.Client = c
			Expect(c.Create(ctx, testFeed)).Should(Succeed())
			testFeed.UID = "12345"
			_, err := testFeed.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error with get feeds"))
		})
	})

})

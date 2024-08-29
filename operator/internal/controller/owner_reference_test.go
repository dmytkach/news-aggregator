package controller_test

import (
	"com.teamdev/news-aggregator/internal/controller"
	"context"
	"errors"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	_ "testing"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("HotNewsOwner", func() {
	var (
		fakeClient client.Client
		ctx        context.Context
		hotNews    *aggregatorv1.HotNews
		owner      *controller.HotNewsOwner
	)

	BeforeEach(func() {
		ctx = context.TODO()
		hotNews = &aggregatorv1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-hotnews",
				Namespace: "default",
				UID:       "12345",
			},
			Spec: aggregatorv1.HotNewsSpec{
				Feeds: []string{"feed1", "feed3"},
			},
		}
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		owner = &controller.HotNewsOwner{
			Client:  fakeClient,
			Ctx:     ctx,
			HotNews: hotNews,
		}
	})
	Describe("CleanupOwnerReferences", func() {
		Context("when feeds not found", func() {
			It("should return an error", func() {
				c := fake.NewClientBuilder().
					WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
					Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
						return fakeClient.Get(ctx, key, obj, opts...)
					},
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("failed to list feeds")
					},
				}).Build()
				owner.Client = c
				err := owner.CleanupOwnerReferences()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to list feeds"))
			})
		})

		Context("when feed was found", func() {
			var feed1 aggregatorv1.Feed

			BeforeEach(func() {
				feed1 = aggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "feed1",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  hotNews.UID,
								Name: hotNews.Name,
							},
						},
					},
				}
			})

			It("should remove OwnerReferences and update feeds", func() {
				Expect(fakeClient.Create(ctx, &feed1)).NotTo(HaveOccurred())
				err := owner.CleanupOwnerReferences()
				Expect(err).NotTo(HaveOccurred())
				var updatedFeed aggregatorv1.Feed
				Expect(fakeClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "feed1"}, &updatedFeed)).NotTo(HaveOccurred())

				Expect(updatedFeed.OwnerReferences).NotTo(ContainElement(
					metav1.OwnerReference{
						UID:  hotNews.UID,
						Name: hotNews.Name,
					},
				))
			})

			It("should return an error if updating feed fails", func() {
				c := fake.NewClientBuilder().
					WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
					Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
						return errors.New("failed to update feeds")
					},
					Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
						return client.Create(ctx, obj)
					},
				}).Build()
				Expect(c.Create(ctx, &feed1)).NotTo(HaveOccurred())
				owner.Client = c
				err := owner.CleanupOwnerReferences()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to update feeds"))
			})
		})
	})

	Describe("updateOwnerReferences", func() {
		Context("when feeds not found", func() {
			It("should return an error", func() {
				c := fake.NewClientBuilder().
					WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
					Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
						return fakeClient.Get(ctx, key, obj, opts...)
					},
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("failed to list feeds")
					},
				}).Build()
				owner.Client = c
				err := owner.UpdateOwnerReferences()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to list feeds"))
			})
		})

		Context("when feed was found", func() {
			var feed1 aggregatorv1.Feed
			var feed2 aggregatorv1.Feed
			var feed3 aggregatorv1.Feed
			BeforeEach(func() {
				feed1 = aggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "feed1",
					},
				}
				feed2 = aggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "feed2",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  hotNews.UID,
								Name: hotNews.Name,
							},
						},
					},
				}
				feed3 = aggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "feed3",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  hotNews.UID,
								Name: hotNews.Name,
							},
						},
					},
				}
			})

			It("should add OwnerReferences", func() {
				Expect(fakeClient.Create(ctx, &feed1)).NotTo(HaveOccurred())
				Expect(fakeClient.Create(ctx, &feed2)).NotTo(HaveOccurred())
				Expect(fakeClient.Create(ctx, &feed3)).NotTo(HaveOccurred())
				err := owner.UpdateOwnerReferences()
				Expect(err).NotTo(HaveOccurred())
				var updatedFeed aggregatorv1.Feed
				Expect(fakeClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "feed1"}, &updatedFeed)).NotTo(HaveOccurred())

				Expect(updatedFeed.OwnerReferences).To(ContainElement(
					metav1.OwnerReference{
						UID:  hotNews.UID,
						Name: hotNews.Name,
					},
				))
				Expect(fakeClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "feed2"}, &updatedFeed)).NotTo(HaveOccurred())

				Expect(updatedFeed.OwnerReferences).ToNot(ContainElement(
					metav1.OwnerReference{
						UID:  hotNews.UID,
						Name: hotNews.Name,
					},
				))
				Expect(fakeClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "feed3"}, &updatedFeed)).NotTo(HaveOccurred())

				Expect(updatedFeed.OwnerReferences).To(ContainElement(
					metav1.OwnerReference{
						UID:  hotNews.UID,
						Name: hotNews.Name,
					},
				))
			})

			It("should return an error if updating feed fails", func() {
				c := fake.NewClientBuilder().
					WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
					Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
						return errors.New("failed to update feeds")
					},
					Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
						return client.Create(ctx, obj)
					},
				}).Build()
				Expect(c.Create(ctx, &feed1)).NotTo(HaveOccurred())
				owner.Client = c
				err := owner.UpdateOwnerReferences()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to update feeds"))
			})
		})
	})
})

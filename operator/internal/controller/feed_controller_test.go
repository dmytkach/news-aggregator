package controller_test

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	mockaggregator "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("FeedReconciler", func() {
	var (
		fakeClient     client.Client
		mockHTTPClient *mockaggregator.MockHttpClient
		reconciler     *controller.FeedReconciler
		ctx            context.Context
		feed           *aggregatorv1.Feed
		ctrl           *gomock.Controller
		err            error
		req            reconcile.Request
		mgr            manager.Manager
	)

	BeforeEach(func() {
		ctx = context.TODO()
		ctrl = gomock.NewController(GinkgoT())
		mockHTTPClient = mockaggregator.NewMockHttpClient(ctrl)
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		mgr, err = manager.New(config.GetConfigOrDie(), manager.Options{
			Scheme: scheme.Scheme,
		})
		req = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "non-existent-feed",
			},
		}
		reconciler = &controller.FeedReconciler{
			Client:        fakeClient,
			Scheme:        scheme.Scheme,
			HttpClient:    mockHTTPClient,
			ServiceURL:    "http://test-service",
			FeedFinalizer: "test-finalizer",
		}

		feed = &aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "test-feed",
				Link: "http://test-example",
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when Feed resource is not found", func() {
		It("should return no error and not attempt to reconcile with a not found error", func() {
			_, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(BeNil())
		})

		It("should return an error when there is a find error", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("find Error")
				},
			}).Build()

			reconciler.Client = fakeClient

			_, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("find Error"))
		})
	})

	Context("when creating a Feed resource", func() {
		It("should succeed and add a finalizer", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Post(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())

			feed := &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.ObjectMeta.Finalizers).To(ContainElement("test-finalizer"))
			Expect(feed.Status.Conditions[0].Status).To(Equal(true))
			Expect(feed.Status.Conditions[0].Type).To(Equal(aggregatorv1.ConditionAdded))
			Expect(feed.Status.Conditions[0].Message).To(Equal("Feed added successfully"))
		})
		It("should handle POST success but fail to update status", func() {
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Post(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())
		})
		It("should handle error when adding finalizer", func() {
			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithInterceptorFuncs(
					interceptor.Funcs{
						Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
							return client.Get(ctx, key, obj)
						},
						Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
							return client.Create(ctx, obj)
						},
						Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
							return errors.New("error with Update feed")
						},
					}).Build()
			reconciler.Client = fakeClient

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(MatchError("error with Update feed"))
			Expect(res.Requeue).To(BeFalse())
		})
		It("should handle POST fails", func() {
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Post(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())
		})
		It("should handle POST request failure", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Post("http://test-service?name=test-feed&url=http://test-example", "application/json", nil).
				Return(nil, errors.New("error with Post request"))

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(MatchError("error with Post request"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.ObjectMeta.Finalizers).To(ContainElement("test-finalizer"))
			Expect(feed.Status.Conditions[0].Status).To(Equal(false))
			Expect(feed.Status.Conditions[0].Type).To(Equal(aggregatorv1.ConditionAdded))
			Expect(feed.Status.Conditions[0].Message).To(Equal("Feed didn't add successfully"))

		})
		Context("when POST request returns an error status", func() {
			var req reconcile.Request
			BeforeEach(func() {
				mockHTTPClient.EXPECT().
					Post(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
					}, nil)

				req = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      "test-feed",
						Namespace: "default",
					},
				}
			})
			It("should handle POST request failure and update the feed status", func() {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
				reconciler.Client = fakeClient
				Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

				res, err := reconciler.Reconcile(ctx, req)
				Expect(err).To(MatchError("failed to create source, status code: 500"))
				Expect(res.Requeue).To(BeFalse())

				feed = &aggregatorv1.Feed{}
				err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
				Expect(err).ToNot(HaveOccurred())
				Expect(feed.ObjectMeta.Finalizers).To(ContainElement("test-finalizer"))
				Expect(feed.Status.Conditions[0].Status).To(Equal(false))
				Expect(feed.Status.Conditions[0].Type).To(Equal(aggregatorv1.ConditionAdded))
				Expect(feed.Status.Conditions[0].Message).To(Equal("Feed didn't add successfully"))

			})
			It("should handle POST request failure but feed not updated", func() {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
				reconciler.Client = fakeClient
				Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

				res, err := reconciler.Reconcile(ctx, req)
				Expect(err.Error()).To(ContainSubstring("not found"))
				Expect(res.Requeue).To(BeFalse())

			})
		})
	})

	Context("when updating a Feed resource", func() {
		BeforeEach(func() {
			feed.Status.Conditions = []aggregatorv1.Condition{
				{
					Type:   aggregatorv1.ConditionAdded,
					Status: true,
				},
			}
		})
		It("should successfully update the feed", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())

			feed := &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.Status.Conditions[1].Status).To(Equal(true))
			Expect(feed.Status.Conditions[1].Type).To(Equal(aggregatorv1.ConditionUpdated))
			Expect(feed.Status.Conditions[1].Message).To(Equal("Feed updated successfully"))
		})
		It("should handle a successful update but fail to update the status", func() {
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())
		})
		It("should handle errors when sending the Put request", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(nil, errors.New("error with PUT request"))

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(MatchError("error with PUT request"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.Status.Conditions[1].Status).To(Equal(false))
			Expect(feed.Status.Conditions[1].Type).To(Equal(aggregatorv1.ConditionUpdated))
			Expect(feed.Status.Conditions[1].Message).To(Equal("Feed didn't update successfully"))
		})
		It("should handle issues with the URL scheme in update request", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			reconciler.ServiceURL = "://bad-url"

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err.Error()).To(ContainSubstring("missing protocol scheme"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.Status.Conditions[1].Status).To(Equal(false))
			Expect(feed.Status.Conditions[1].Type).To(Equal(aggregatorv1.ConditionUpdated))
			Expect(feed.Status.Conditions[1].Message).To(Equal("Feed didn't update successfully"))
		})
		Context("PUT Request Error Status Handling", func() {
			var req reconcile.Request
			BeforeEach(func() {
				mockHTTPClient.EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
					}, nil)

				req = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      "test-feed",
						Namespace: "default",
					},
				}
			})

			It("should handle PUT request failure but feed updated", func() {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
				reconciler.Client = fakeClient
				Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

				res, err := reconciler.Reconcile(ctx, req)
				Expect(err).To(MatchError("failed to update source, status code: 500"))
				Expect(res.Requeue).To(BeFalse())

				feed = &aggregatorv1.Feed{}
				err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
				Expect(err).ToNot(HaveOccurred())
				Expect(feed.Status.Conditions[1].Status).To(Equal(false))
				Expect(feed.Status.Conditions[1].Type).To(Equal(aggregatorv1.ConditionUpdated))
				Expect(feed.Status.Conditions[1].Message).To(Equal("Feed didn't update successfully"))
			})

			It("should handle PUT request failure but feed not updated", func() {
				fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
				reconciler.Client = fakeClient
				Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())

				res, err := reconciler.Reconcile(ctx, req)
				Expect(err.Error()).To(ContainSubstring("not found"))
				Expect(res.Requeue).To(BeFalse())
			})
		})
	})

	Context("when deleting a Feed resource", func() {
		It("should successfully delete the feed", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			feed.Finalizers = []string{"test-finalizer"}

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			Expect(reconciler.Client.Delete(ctx, feed)).To(Succeed())

			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).To(HaveOccurred())
		})
		It("should handle successful deletion but fail to update the status", func() {
			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithInterceptorFuncs(
					interceptor.Funcs{
						Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
							return client.Get(ctx, key, obj)
						},
						Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
							return client.Create(ctx, obj)
						},
						Delete: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.DeleteOption) error {
							return client.Delete(ctx, obj)
						},
						Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
							return errors.New("error with Update feed")
						},
					}).Build()
			reconciler.Client = fakeClient
			feed.Finalizers = []string{"test-finalizer"}

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			Expect(reconciler.Client.Delete(ctx, feed)).To(Succeed())

			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(MatchError("error with Update feed"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should handle errors during the DELETE request ", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			feed.Finalizers = []string{"test-finalizer"}

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			Expect(reconciler.Client.Delete(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(nil, errors.New("error with DELETE request"))

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(MatchError("error with DELETE request"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.Status.Conditions[0].Status).To(Equal(false))
			Expect(feed.Status.Conditions[0].Type).To(Equal(aggregatorv1.ConditionDeleted))
			Expect(feed.Status.Conditions[0].Message).To(Equal("Failed to delete feed"))
		})
		It("should handle issues with an invalid URL scheme during the DELETE request", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()
			reconciler.Client = fakeClient
			reconciler.ServiceURL = "://bad-url"
			feed.Finalizers = []string{"test-finalizer"}

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			Expect(reconciler.Client.Delete(ctx, feed)).To(Succeed())

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err.Error()).To(ContainSubstring("missing protocol scheme"))
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
			Expect(feed.Status.Conditions[0].Status).To(Equal(false))
			Expect(feed.Status.Conditions[0].Type).To(Equal(aggregatorv1.ConditionDeleted))
			Expect(feed.Status.Conditions[0].Message).To(Equal("Failed to delete feed"))
		})
		It("should handle issues during the DELETE request without updating feed status", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
			reconciler.Client = fakeClient
			feed.Finalizers = []string{"test-finalizer"}

			Expect(reconciler.Client.Create(ctx, feed)).To(Succeed())
			Expect(reconciler.Client.Delete(ctx, feed)).To(Succeed())
			mockHTTPClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
				}, nil)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-feed",
					Namespace: "default",
				},
			}

			res, err := reconciler.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(res.Requeue).To(BeFalse())

			feed = &aggregatorv1.Feed{}
			err = reconciler.Client.Get(ctx, req.NamespacedName, feed)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("SetupWithManager", func() {
		It("should setup the controller without errors", func() {
			err = reconciler.SetupWithManager(mgr)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

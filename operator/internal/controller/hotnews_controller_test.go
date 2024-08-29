package controller_test

import (
	"bytes"
	v1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	mockaggregator "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var _ = Describe("HotNewsReconciler", func() {
	var (
		fakeClient     client.Client
		mockHTTPClient *mockaggregator.MockHttpClient
		reconciler     *controller.HotNewsReconciler
		ctx            context.Context
		feed           *v1.Feed
		ctrl           *gomock.Controller
		mgr            manager.Manager
		err            error
	)

	BeforeEach(func() {
		ctx = context.TODO()
		ctrl = gomock.NewController(GinkgoT())
		mockHTTPClient = mockaggregator.NewMockHttpClient(ctrl)
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		mgr, err = manager.New(config.GetConfigOrDie(), manager.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).ToNot(HaveOccurred())
		reconciler = &controller.HotNewsReconciler{
			Client:     fakeClient,
			Scheme:     scheme.Scheme,
			HttpClient: mockHTTPClient,
			ServiceURL: "http://test-service",
			ConfigMap:  "test-configmap",
			Finalizer:  "test-finalizer",
		}

		feed = &v1.Feed{
			ObjectMeta: v12.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
			Spec: v1.FeedSpec{
				Name: "test-feed",
				Link: "http://test-example",
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when HotNews resource is not found", func() {
		It("should return no error and not attempt to reconcile with a not found error", func() {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "non-existent-hotnews",
				},
			}

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

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "non-existent-hotnews",
				},
			}

			_, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("find Error"))
		})
	})
	Context("when HotNews uses finalizers", func() {
		var (
			hotNews *v1.HotNews
		)
		BeforeEach(func() {
			hotNews = &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: v1.HotNewsSpec{
					Feeds:     []string{"test-feed"},
					Keywords:  []string{"test-keyword"},
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-02",
					SummaryConfig: v1.SummaryConfig{
						TitlesCount: 3,
					},
				},
			}
		})
		It("Finalizer not found with wrong status update", func() {
			c := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).WithStatusSubresource(&v1.HotNews{}).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return fakeClient.Get(ctx, key, obj, opts...)
				},
				Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					return errors.New("fail to add finalizer")
				},
			}).Build()
			reconciler.Client = c

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("fail to add finalizer"))
			Expect(result.Requeue).To(BeFalse())

		})
	})
	Context("when HotNews uses feeds", func() {
		var (
			hotNews *v1.HotNews
		)

		BeforeEach(func() {
			hotNews = &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:       "test-hotnews",
					Namespace:  "default",
					Finalizers: []string{"test-finalizer"},
				},
				Spec: v1.HotNewsSpec{
					Feeds:     []string{"test-feed"},
					Keywords:  []string{"test-keyword"},
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-02",
					SummaryConfig: v1.SummaryConfig{
						TitlesCount: 3,
					},
				},
			}
		})
		It("should update the HotNews status with fetched news", func() {
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			mockHTTPClient.EXPECT().Do(gomock.Any()).
				DoAndReturn(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet {
						return nil, fmt.Errorf("expected GET method, got %s", req.Method)
					}

					response := `[
					{"Title": "News 1"},
					{"Title": "News 2"},
					{"Title": "News 3"}
				]`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(response)),
					}, nil
				})

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(BeNil())
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())

			Expect(updatedHotNews.Status.ArticlesCount).To(Equal(3))
			Expect(updatedHotNews.Status.ArticlesTitles).To(ConsistOf("News 1", "News 2", "News 3"))
			Expect(updatedHotNews.Status.NewsLink).
				To(Equal(fmt.Sprintf("http://test-service?date-end=%s&date-start=%s&keywords=test-keyword&sort-order=asc&sources=test-feed",
					hotNews.Spec.DateEnd, hotNews.Spec.DateStart)))
		})
		It("should be error Feed not found with wrong status update", func() {
			c := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).WithStatusSubresource(&v1.HotNews{}).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return fakeClient.Get(ctx, key, obj, opts...)
				},
				List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
					return errors.New("list error")
				},
				Create: func(ctx context.Context, client client.WithWatch, object client.Object, opts ...client.CreateOption) error {
					return fakeClient.Create(ctx, object, opts...)
				},
			}).Build()
			reconciler.Client = c

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("error updating Feed default/test-hotnews"))
			Expect(result.Requeue).To(BeFalse())

		})
		It("should be error Feed not found with status update", func() {
			c := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&v1.HotNews{}).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return client.Get(ctx, key, obj, opts...)
				},
				List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
					return errors.New("get feeds error")
				},
				Create: func(ctx context.Context, client client.WithWatch, object client.Object, opts ...client.CreateOption) error {
					return client.Create(ctx, object, opts...)
				},
				Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					return client.Update(ctx, obj, opts...)
				},
			}).Build()
			reconciler.Client = c

			Expect(c.Create(ctx, hotNews)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("get feeds error"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(c.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("get feeds error"))
			Expect(updatedHotNews.Status.Condition.Status).To(Equal(false))
		})
	})
	Context("when buildRequestURL fails", func() {
		It("should return an error due to an invalid url", func() {
			reconciler.ServiceURL = "http://example.com/%ZZ"
			fakeClient := fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient

			hotNews := &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: v1.HotNewsSpec{
					Feeds:    []string{"test-feed"},
					Keywords: []string{"test-keyword"},
				},
			}

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError(fmt.Sprintf("invalid service URL: %s", reconciler.ServiceURL)))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal(fmt.Sprintf("invalid service URL: %s", reconciler.ServiceURL)))
		})

		It("should return an error for missing keyword", func() {
			fakeClient := fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient

			hotNews := &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: v1.HotNewsSpec{
					Feeds:    []string{"test-feed"},
					Keywords: nil,
				},
			}

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("keywords not found"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("keywords not found"))
		})
	})
	Context("when makeRequest fails", func() {
		var (
			req        reconcile.Request
			hotNews    *v1.HotNews
			fakeClient client.Client
		)

		BeforeEach(func() {
			fakeClient = fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient

			hotNews = &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: v1.HotNewsSpec{
					Feeds:     []string{"test-feed"},
					Keywords:  []string{"test-keyword"},
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-02",
					SummaryConfig: v1.SummaryConfig{
						TitlesCount: 3,
					},
				},
			}

			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())

			req = reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
		})
		It("should return an error with HTTP request", func() {
			mockHTTPClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("request error"))

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("request error"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("request error"))
		})
		It("should return an error with HTTP response decoding", func() {

			mockHTTPClient.EXPECT().Do(gomock.Any()).
				DoAndReturn(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("not-a-json-response")),
					}, nil
				})

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("invalid character 'o' in literal null (expecting 'u')"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("invalid character 'o' in literal null (expecting 'u')"))
		})
		It("should return an error with HTTP response status", func() {
			mockHTTPClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			})

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("failed to create source, status code: 500"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("failed to create source, status code: 500"))
		})
		It("should return an error with closing response body", func() {
			mockHTTPClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&errorReader{}),
				}, nil
			})

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("failed to close response body: mock error"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("failed to close response body: mock error"))
		})
	})
	Context("when HotNews uses FeedGroups", func() {
		var (
			configMap *v13.ConfigMap
			hotNews   *v1.HotNews
		)

		BeforeEach(func() {
			configMap = &v13.ConfigMap{
				ObjectMeta: v12.ObjectMeta{
					Name:      reconciler.ConfigMap,
					Namespace: "default",
				},
				Data: map[string]string{
					"group": "test-feed",
				},
			}
			hotNews = &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:       "test-hotnews",
					Namespace:  "default",
					Finalizers: []string{"test-finalizer"},
				},
				Spec: v1.HotNewsSpec{
					FeedGroups: []string{"group"},
					Keywords:   []string{"test-keyword"},
					DateStart:  "2024-01-01",
					DateEnd:    "2024-01-02",
					SummaryConfig: v1.SummaryConfig{
						TitlesCount: 3,
					},
				},
			}
		})

		It("should update the HotNews status with fetched news", func() {
			fakeClient := fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient
			Expect(fakeClient.Create(ctx, configMap)).To(Succeed())
			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			mockHTTPClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				response := `[
                {"Title": "News 1"},
                {"Title": "News 2"},
                {"Title": "News 3"}
            ]`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(response)),
				}, nil
			})

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(BeNil())
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())

			Expect(updatedHotNews.Status.ArticlesCount).To(Equal(3))
			Expect(updatedHotNews.Status.ArticlesTitles).To(ConsistOf("News 1", "News 2", "News 3"))
			Expect(updatedHotNews.Status.NewsLink).
				To(Equal(fmt.Sprintf("http://test-service?date-end=%s&date-start=%s&keywords=test-keyword&sort-order=asc&sources=test-feed",
					hotNews.Spec.DateEnd, hotNews.Spec.DateStart)))
		})
		It("should return an error status update if ConfigMap is not found", func() {

			c := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
						if key.Name == "test-configmap" {
							return errors.New("ConfigMap not found")
						}
						return fakeClient.Get(ctx, key, obj, opts...)
					},
					Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
						return fakeClient.Create(ctx, obj, opts...)
					},
				}).Build()

			reconciler.Client = c

			Expect(c.Create(ctx, hotNews)).To(Succeed())

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("error updating Feed default/test-hotnews"))
			Expect(result.Requeue).To(BeFalse())
		})
		It("should successful status update if ConfigMap is not found", func() {
			c := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&v1.HotNews{}).
				WithInterceptorFuncs(interceptor.Funcs{
					Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
						if key.Name == "test-configmap" {
							return errors.New("ConfigMap not found")
						}
						return client.Get(ctx, key, obj, opts...)
					},
					Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
						return client.Create(ctx, obj, opts...)
					},
					Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
						return client.Update(ctx, obj, opts...)
					},
				}).Build()

			reconciler.Client = c

			Expect(c.Create(ctx, hotNews)).To(Succeed())

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("ConfigMap not found"))
			Expect(result.Requeue).To(BeFalse())

			var updatedHotNews v1.HotNews
			Expect(c.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("ConfigMap not found"))
			Expect(updatedHotNews.Status.Condition.Status).To(Equal(false))
		})
		It("should return an error not found if ConfigMap is not found", func() {
			fakeClient := fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())
			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}

			_, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(BeNil())
		})
	})
	Context("when no feeds or feed groups are provided", func() {
		var (
			hotNews *v1.HotNews
		)

		BeforeEach(func() {
			hotNews = &v1.HotNews{
				ObjectMeta: v12.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: v1.HotNewsSpec{
					Keywords:  []string{"test-keyword"},
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-02",
					SummaryConfig: v1.SummaryConfig{
						TitlesCount: 3,
					},
				},
			}
		})
		It("should return an error and update status", func() {
			fakeClient := fake.NewClientBuilder().WithStatusSubresource(&v1.HotNews{}).Build()
			reconciler.Client = fakeClient
			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("no feeds or feed groups provided in HotNews spec"))
			Expect(result.Requeue).To(BeFalse())
			var updatedHotNews v1.HotNews
			Expect(fakeClient.Get(ctx, req.NamespacedName, &updatedHotNews)).To(Succeed())
			Expect(updatedHotNews.Status.Condition.Reason).To(Equal("no feeds or feed groups provided in HotNews spec"))
			Expect(updatedHotNews.Status.Condition.Status).To(Equal(false))
		})
		It("should return an error and not update status", func() {
			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-hotnews",
				},
			}
			result, err := reconciler.Reconcile(ctx, req)

			Expect(err).To(MatchError("error updating Feed default/test-hotnews"))
			Expect(result.Requeue).To(BeFalse())

		})
	})
	Context("SetupWithManager", func() {
		It("should setup the controller without errors", func() {
			err = reconciler.SetupWithManager(mgr)
			Expect(err).ToNot(HaveOccurred())
		})
	})

})

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to close response body: mock error")
}

func (e *errorReader) Close() error {
	return fmt.Errorf("failed to close response body: mock error")
}

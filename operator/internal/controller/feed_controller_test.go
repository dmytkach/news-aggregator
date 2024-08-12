package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestFeedReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	initialFeed := &aggregatorv1.Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-feed",
			Namespace: "default",
		},
		Spec: aggregatorv1.FeedSpec{
			Name: "TestFeed",
			Link: "https://example.com/feed",
		},
		Status: aggregatorv1.FeedStatus{
			Conditions: []aggregatorv1.Condition{},
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initialFeed).Build()

	feed := &aggregatorv1.Feed{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      "test-feed",
		Namespace: "default",
	}, feed)
	assert.NoError(t, err, "initial Feed object should be found")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPClient := controller.NewMockHttpClient(ctrl)

	mockHTTPClient.EXPECT().
		Post("https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources?name=TestFeed&url=https://example.com/feed",
			"application/json", gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil)

	r := &FeedReconciler{
		Client:        client,
		Scheme:        scheme,
		HttpClient:    mockHTTPClient,
		FeedFinalizer: "feed.finalizers.news.teamdev.com",
		ServiceURL:    "https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources",
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-feed",
			Namespace: "default",
		},
	}

	res, err := r.Reconcile(context.Background(), req)
	assert.False(t, res.Requeue)

	err = client.Get(context.Background(), req.NamespacedName, feed)
	assert.NoError(t, err, "Feed object should be found after reconciliation")

	if !assert.Contains(t, feed.Finalizers, "feed.finalizers.news.teamdev.com", "Finalizer should be added") {
		t.Logf("Finalizers found: %v", feed.Finalizers)
	}
}
func TestFeedReconciler_addFeed(t *testing.T) {
	tests := []struct {
		name             string
		feed             aggregatorv1.Feed
		mockPostResponse *http.Response
		mockPostError    error
		expectedError    bool
	}{
		{
			name: "Success request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPostResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			mockPostError: nil,
			expectedError: false,
		},
		{
			name: "HTTP error from POST request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPostResponse: nil,
			mockPostError:    fmt.Errorf("network error"),
			expectedError:    true,
		},
		{
			name: "Non-200 status code",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPostResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			mockPostError: nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Post(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.mockPostResponse, tt.mockPostError)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				ServiceURL: "http://mock-server/create-feed",
			}

			err := reconciler.createFeed(tt.feed)

			if (err != nil) != tt.expectedError {
				t.Errorf("createFeed() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestFeedReconciler_deleteFeed(t *testing.T) {
	tests := []struct {
		name               string
		feed               aggregatorv1.Feed
		mockDeleteResponse *http.Response
		mockDeleteError    error
		expectedError      bool
	}{
		{
			name: "Success delete request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockDeleteResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			mockDeleteError: nil,
			expectedError:   false,
		},
		{
			name: "Failed delete request with error",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockDeleteResponse: nil,
			mockDeleteError:    fmt.Errorf("delete request failed"),
			expectedError:      true,
		},
		{
			name: "Failed delete request with non-OK status code",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockDeleteResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			mockDeleteError: nil,
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Do(gomock.Any()).
				Return(tt.mockDeleteResponse, tt.mockDeleteError).
				Times(1)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				ServiceURL: "http://mock-server/delete-feed",
			}

			err := reconciler.deleteFeed(tt.feed)
			if (err != nil) != tt.expectedError {
				t.Errorf("deleteFeed() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
func TestFeedReconciler_updateFeed(t *testing.T) {
	tests := []struct {
		name            string
		feed            aggregatorv1.Feed
		mockPutResponse *http.Response
		mockPutError    error
		expectedError   bool
	}{
		{
			name: "Success update request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPutResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			mockPutError:  nil,
			expectedError: false,
		},
		{
			name: "Failed update request with error",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPutResponse: nil,
			mockPutError:    fmt.Errorf("update request failed"),
			expectedError:   true,
		},
		{
			name: "Failed update request with non-OK status code",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "TestFeed",
					Link: "http://example.com/feed",
				},
			},
			mockPutResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			mockPutError:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Do(gomock.Any()).
				Return(tt.mockPutResponse, tt.mockPutError).
				Times(1)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				ServiceURL: "http://mock-server/update-feed",
			}

			err := reconciler.updateFeed(tt.feed)
			if (err != nil) != tt.expectedError {
				t.Errorf("updateFeed() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

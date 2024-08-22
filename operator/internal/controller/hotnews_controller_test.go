package controller

import (
	"bytes"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

func TestHotNewsReconciler_Reconcile(t *testing.T) {
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
	hotNews := &aggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hotnews",
			Namespace: "default",
		},
		Spec: aggregatorv1.HotNewsSpec{
			Keywords: []string{"test"},
			Feeds:    []string{"test-feed"},
			SummaryConfig: aggregatorv1.SummaryConfig{
				TitlesCount: 5,
			},
		},
		Status: aggregatorv1.HotNewsStatus{},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).
		WithObjects(hotNews).WithObjects(initialFeed).Build()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := controller.NewMockHttpClient(ctrl)

	mockHTTPClient.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBufferString(`[
				{"Title": "News1"},
				{"Title": "News2"},
				{"Title": "News3"}
			]`)),
		}, nil).
		Times(1)

	r := &HotNewsReconciler{
		Client:     client,
		Scheme:     scheme,
		HttpClient: mockHTTPClient,
		ServiceURL: "http://test-service",
		Finalizer:  "example.com.finalizer",
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-hotnews",
			Namespace: "default",
		},
	}

	res, err := r.Reconcile(context.Background(), req)
	assert.False(t, res.Requeue)

	err = client.Get(context.Background(), req.NamespacedName, hotNews)
	assert.NoError(t, err)

	if !assert.Contains(t, hotNews.Finalizers, "example.com.finalizer", "Finalizer should be added") {
		t.Logf("Finalizers found: %v", hotNews.Finalizers)
	}
}

func TestHotNewsReconciler_Reconcile_ErrorGettingHotNews(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := &HotNewsReconciler{
		Client:    client,
		Scheme:    scheme,
		Finalizer: "example.com.finalizer",
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "non-existent-hotnews",
			Namespace: "default",
		},
	}

	res, err := r.Reconcile(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.Requeue, "Reconcile should not requeue when HotNews is not found")
}

func TestFetchNewsData_ErrorFetching(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := controller.NewMockHttpClient(ctrl)
	reconciler := &HotNewsReconciler{
		HttpClient: mockHTTPClient,
		ServiceURL: "http://test-service",
	}

	mockHTTPClient.EXPECT().
		Do(gomock.Any()).
		Return(nil, fmt.Errorf("network error")).
		Times(1)

	status, err := reconciler.fetchNewsData([]string{"source1"}, aggregatorv1.HotNewsSpec{
		Keywords:      []string{"keyword"},
		DateStart:     "2024-01-01",
		DateEnd:       "2024-01-31",
		SummaryConfig: aggregatorv1.SummaryConfig{TitlesCount: 2},
	})

	assert.Error(t, err)
	assert.Empty(t, status)
}

func TestFetchNewsData_InvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := controller.NewMockHttpClient(ctrl)
	reconciler := &HotNewsReconciler{
		HttpClient: mockHTTPClient,
		ServiceURL: "http://test-service",
	}

	mockHTTPClient.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`invalid json`)),
		}, nil).
		Times(1)

	status, err := reconciler.fetchNewsData([]string{"source1"}, aggregatorv1.HotNewsSpec{
		Keywords:      []string{"keyword"},
		DateStart:     "2024-01-01",
		DateEnd:       "2024-01-31",
		SummaryConfig: aggregatorv1.SummaryConfig{TitlesCount: 2},
	})

	assert.Error(t, err)
	assert.Empty(t, status)
}

func TestFetchNewsData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := controller.NewMockHttpClient(ctrl)
	reconciler := &HotNewsReconciler{
		HttpClient: mockHTTPClient,
		ServiceURL: "http://test-service",
	}

	mockHTTPClient.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`[
				{"Title": "Test News 1"},
				{"Title": "Test News 2"}
			]`)),
		}, nil).
		Times(1)

	status, err := reconciler.fetchNewsData([]string{"source1"}, aggregatorv1.HotNewsSpec{
		Keywords:      []string{"keyword"},
		DateStart:     "2024-01-01",
		DateEnd:       "2024-01-31",
		SummaryConfig: aggregatorv1.SummaryConfig{TitlesCount: 2},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, status.ArticlesCount)
	assert.ElementsMatch(t, []string{"Test News 1", "Test News 2"}, status.ArticlesTitles)
}
func TestGetFeedNamesFromFeedGroups(t *testing.T) {
	r := &HotNewsReconciler{}

	configMap := v1.ConfigMap{
		Data: map[string]string{
			"group1": "feed1,feed2",
			"group2": "feed3",
		},
	}

	sources := r.extractFeedsFromGroups([]string{"group1", "group2"}, configMap)

	expectedSources := []string{"feed1", "feed2", "feed3"}
	assert.ElementsMatch(t, expectedSources, sources)
}
func TestGetNewsSourcesFromFeeds(t *testing.T) {
	r := &HotNewsReconciler{}
	feed1 := aggregatorv1.Feed{Spec: aggregatorv1.FeedSpec{Name: "news1"}}
	feed1.Name = "feed1"
	feed2 := aggregatorv1.Feed{Spec: aggregatorv1.FeedSpec{Name: "news2"}}
	feed2.Name = "feed2"
	feeds := aggregatorv1.FeedList{
		Items: []aggregatorv1.Feed{feed1, feed2},
	}
	newsSources := r.getNewsSourcesFromFeeds([]string{"feed1", "feed2"}, feeds)

	expectedNewsSources := []string{"news1", "news2"}
	assert.ElementsMatch(t, expectedNewsSources, newsSources)
}
func TestBuildRequestURL(t *testing.T) {
	url, err := buildRequestURL(
		"https://example.com/api",
		[]string{"feed1", "feed2"},
		[]string{"keyword1", "keyword2"},
		"2024-01-01",
		"2024-01-31",
	)

	assert.NoError(t, err)
	expectedURL := "https://example.com/api?date-end=2024-01-31&date-start=2024-01-01&keywords=keyword1%2Ckeyword2&sort-order=asc&sources=feed1%2Cfeed2"
	assert.Equal(t, expectedURL, url)
}

func TestMakeRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := controller.NewMockHttpClient(ctrl)
	reconciler := &HotNewsReconciler{
		HttpClient: mockHTTPClient,
		ServiceURL: "http://example.com/news",
	}

	mockHTTPClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`[{"Title": "Title1"}, {"Title": "Title2"}, {"Title": "Title3"}]`)),
	}, nil).Times(1)

	status, err := reconciler.fetchNewsData([]string{"source1"}, aggregatorv1.HotNewsSpec{
		Keywords:      []string{"keyword"},
		DateStart:     "2024-01-01",
		DateEnd:       "2024-01-31",
		SummaryConfig: aggregatorv1.SummaryConfig{TitlesCount: 3},
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, status.ArticlesCount)
	assert.Equal(t, "http://example.com/news?date-end=2024-01-31&date-start=2024-01-01&keywords=keyword&sort-order=asc&sources=source1", status.NewsLink)
	assert.ElementsMatch(t, []string{"Title1", "Title2", "Title3"}, status.ArticlesTitles)
}

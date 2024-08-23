package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HotNewsReconciler manages HotNews resources within a k8s cluster.
// It interacts with the Kubernetes API and external news services to synchronize
// and update the state of HotNews resources.
type HotNewsReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	HttpClient HttpClient
	ServiceURL string
	ConfigMap  string
}

// NewsTitle represents a single news article title retrieved from
// the external news service.
type NewsTitle struct {
	Title string `json:"Title"`
}

// NewsResponse represents a collection of NewsTitle elements the external news service returns.
type NewsResponse []NewsTitle

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile performs the reconciliation logic for HotNews resources.
// It manages the lifecycle of the HotNews resource,
// retrieving relevant Feeds and ConfigMaps, and updating the HotNews status with fetched news data.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for HotNews %s/%s", req.Namespace, req.Name)

	var hotNews aggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			log.Print("Reconcile: HotNews was not found. Error ignored")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	var feeds aggregatorv1.FeedList
	if err := r.List(ctx, &feeds, client.InNamespace(req.Namespace)); err != nil {
		log.Printf("Error listing Feed resources: %v", err)
		hotNews.Status = aggregatorv1.HotNewsStatus{
			Condition: aggregatorv1.HotNewsCondition{
				Status: false,
				Reason: err.Error(),
			},
		}
		if err := r.Status().Update(ctx, &hotNews); err != nil {
			log.Printf("Error updating Feed %s/%s: %v", req.Namespace, req.Name, err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	var feedNames []string
	if len(hotNews.Spec.Feeds) > 0 {
		feedNames = hotNews.Spec.Feeds
	} else if len(hotNews.Spec.FeedGroups) > 0 {
		var feedGroupConfigMap v1.ConfigMap
		if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: r.ConfigMap}, &feedGroupConfigMap); err != nil {
			if errors.IsNotFound(err) {
				log.Print("ConfigMap not found, retrying later")
				return ctrl.Result{}, err
			}
			log.Printf("Error retrieving ConfigMap %s from k8s Cluster: %v", "feed-group-source", err)
			hotNews.Status = aggregatorv1.HotNewsStatus{
				Condition: aggregatorv1.HotNewsCondition{
					Status: false,
					Reason: err.Error(),
				},
			}
			if err := r.Status().Update(ctx, &hotNews); err != nil {
				log.Printf("Error updating Feed %s/%s: %v", req.Namespace, req.Name, err)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
		feedNames = r.extractFeedsFromGroups(hotNews.Spec.FeedGroups, feedGroupConfigMap)
	} else {
		feedNames = feeds.GetAllFeedNames()
		err := fmt.Errorf("no feeds or feed groups provided in HotNews spec")
		hotNews.Status = aggregatorv1.HotNewsStatus{
			Condition: aggregatorv1.HotNewsCondition{
				Status: false,
				Reason: err.Error(),
			},
		}
		if err := r.Status().Update(ctx, &hotNews); err != nil {
			log.Printf("Error updating Feed %s/%s: %v", req.Namespace, req.Name, err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	sources := r.getNewsSourcesFromFeeds(feedNames, feeds)

	log.Printf("Final list of news names: %v", sources)
	status, err := r.fetchNewsData(sources, hotNews.Spec)
	if err != nil {
		hotNews.Status = aggregatorv1.HotNewsStatus{
			Condition: aggregatorv1.HotNewsCondition{
				Status: false,
				Reason: err.Error(),
			},
		}
		if err := r.Status().Update(ctx, &hotNews); err != nil {
			log.Printf("Error updating Feed %s/%s: %v", req.Namespace, req.Name, err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	hotNews.Status = aggregatorv1.HotNewsStatus{
		ArticlesCount:  status.ArticlesCount,
		NewsLink:       status.NewsLink,
		ArticlesTitles: status.ArticlesTitles,
		Condition: aggregatorv1.HotNewsCondition{
			Status: true,
		},
	}

	if err := r.Status().Update(ctx, &hotNews); err != nil {
		log.Printf("Error updating status of HotNews %s/%s: %v", req.Namespace, req.Name, err)
		return reconcile.Result{}, err
	}

	log.Printf("Successfully updated HotNews %s/%s", req.Namespace, req.Name)
	return reconcile.Result{}, nil
}

// extractFeedsFromGroups retrieves the feed names associated
// with the specified FeedGroups by referencing a ConfigMap.
func (r *HotNewsReconciler) extractFeedsFromGroups(feedGroups []string, configMap v1.ConfigMap) []string {
	var feedNames []string

	for _, feedGroup := range feedGroups {
		log.Printf("Processing FeedGroup: %s", feedGroup)
		if value, ok := configMap.Data[feedGroup]; ok {
			feedNames = append(feedNames, strings.Split(value, ",")...)
			log.Printf("Matched FeedGroup '%s' in ConfigMap, added values: %v", feedGroup, feedNames)
		}
	}

	return feedNames
}

// getNewsSourcesFromFeeds retrieves the news sources based on the specified feed names.
func (r *HotNewsReconciler) getNewsSourcesFromFeeds(feedNames []string, feeds aggregatorv1.FeedList) []string {
	var sources []string

	for _, feed := range feeds.Items {
		if slices.Contains(feedNames, feed.Name) && !slices.Contains(sources, feed.Spec.Name) {
			sources = append(sources, feed.Spec.Name)
			log.Printf("Matched Feed: %s, adding to sources: %s", feed.Name, feed.Spec.Name)
		}
	}

	return sources
}

// fetchNewsData constructs a request URL using the provided sources and HotNews specifications,
// then sends a request to the news service to retrieve news data. It returns the status of the HotNews
// with the fetched articles or an error if the process fails.
func (r *HotNewsReconciler) fetchNewsData(sources []string, hotNews aggregatorv1.HotNewsSpec) (aggregatorv1.HotNewsStatus, error) {
	log.Printf("Starting fetchNewsData with sources: %v, keywords: %v, dateStart: %s, dateEnd: %s", sources, hotNews.Keywords, hotNews.DateStart, hotNews.DateEnd)
	reqURL, err := buildRequestURL(r.ServiceURL, sources, hotNews.Keywords, hotNews.DateStart, hotNews.DateEnd)
	if err != nil {
		log.Printf("Error building request: %v", err)
		return aggregatorv1.HotNewsStatus{Condition: aggregatorv1.HotNewsCondition{
			Status: false,
			Reason: err.Error(),
		}}, err
	}

	status, err := r.makeRequest(reqURL, hotNews.SummaryConfig.TitlesCount)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return aggregatorv1.HotNewsStatus{Condition: aggregatorv1.HotNewsCondition{
			Status: false,
			Reason: err.Error(),
		}}, err
	}
	log.Printf("Request completed successfully, received status: %+v", status)

	return status, nil
}

// buildRequestURL constructs the URL for the news service request based on the provided sources,
// keywords, and date range. It returns the formatted URL or an error if the URL cannot be constructed
func buildRequestURL(serviceURL string, sources, keywords []string, dateStart, dateEnd string) (string, error) {
	baseURL, err := url.Parse(serviceURL)
	if err != nil {
		return "", fmt.Errorf("invalid service URL: %v", err)
	}

	params := url.Values{}

	if len(sources) > 0 {
		params.Add("sources", strings.Join(sources, ","))
	}

	if len(keywords) > 0 {
		params.Add("keywords", strings.Join(keywords, ","))
	} else {
		return "", fmt.Errorf("keywords not found")
	}

	if dateStart != "" {
		params.Add("date-start", dateStart)
	}
	if dateEnd != "" {
		params.Add("date-end", dateEnd)
	}

	params.Add("sort-order", "asc")
	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}

// makeRequest performs the HTTP request to the news service and processes the response.
// It returns the HotNewsStatus populated with the articles and titles or an error if the request fails.
func (r *HotNewsReconciler) makeRequest(reqURL string, titleCount int) (aggregatorv1.HotNewsStatus, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("Failed to create Get request: %v", err)
		return aggregatorv1.HotNewsStatus{}, err
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make Get request: %v", err)
		return aggregatorv1.HotNewsStatus{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to create source, status code: %d", resp.StatusCode)
		return aggregatorv1.HotNewsStatus{}, fmt.Errorf("failed to create source, status code: %d", resp.StatusCode)
	}

	var newsResponse NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResponse); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return aggregatorv1.HotNewsStatus{}, err
	}

	var titles []string
	for i := range newsResponse {
		titles = append(titles, newsResponse[i].Title)
		if i >= titleCount-1 {
			break
		}
	}

	return aggregatorv1.HotNewsStatus{
		ArticlesCount:  len(newsResponse),
		NewsLink:       reqURL,
		ArticlesTitles: titles,
		Condition:      aggregatorv1.HotNewsCondition{Status: true},
	}, nil
}

// SetupWithManager configures the HotNewsReconciler to manage resources and adds the necessary event predicates
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.handlerConfigMap),
		).
		Watches(
			&aggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(r.handlerFeed)).
		Complete(r)
}

// handlerConfigMap processes changes to ConfigMap objects and generates reconcile requests for HotNews resources.
func (r *HotNewsReconciler) handlerConfigMap(ctx context.Context, obj client.Object) []ctrl.Request {
	var requests []ctrl.Request

	log.Printf("Starting handlerConfigMap for object: %v", obj)

	configMap, ok := obj.(*v1.ConfigMap)
	if !ok {
		log.Printf("Object is not a ConfigMap: %v", obj)
		return requests
	}

	if configMap.Name == r.ConfigMap {

		namespace := configMap.Namespace

		hotNewsList := &aggregatorv1.HotNewsList{}
		err := r.Client.List(ctx, hotNewsList, client.InNamespace(namespace))
		if err != nil {
			log.Printf("Error listing HotNews in namespace %s: %v", namespace, err)
			return requests
		}

		log.Printf("Found %d HotNews in namespace %s", len(hotNewsList.Items), namespace)

		for _, hotNews := range hotNewsList.Items {
			if len(r.extractFeedsFromGroups(hotNews.Spec.FeedGroups, *configMap)) > 0 {
				requests = append(requests, ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: hotNews.Namespace,
						Name:      hotNews.Name,
					},
				})

				log.Printf("Enqueued request for HotNews: Name=%s, Namespace=%s", hotNews.Name, hotNews.Namespace)
			}
		}
	}

	log.Printf("Completed handlerConfigMap with %d requests", len(requests))

	return requests
}

// handlerFeed processes changes to Feed objects and generates reconcile requests for HotNews resources.
func (r *HotNewsReconciler) handlerFeed(context.Context, client.Object) []ctrl.Request {
	var requests []ctrl.Request
	//get feeds
	feedList := &aggregatorv1.FeedList{}
	err := r.Client.List(context.Background(), feedList, &client.ListOptions{})
	if err != nil {
		return requests
	}

	for _, feed := range feedList.Items {
		namespace := feed.Namespace

		hotNewsList := &aggregatorv1.HotNewsList{}
		// get hotNews in the feed namespace
		err = r.Client.List(context.Background(), hotNewsList, client.InNamespace(namespace))
		if err != nil {
			return requests
		}
		// fing hotnews and save to request
		for _, hotNews := range hotNewsList.Items {
			if slices.Contains(hotNews.Spec.Feeds, feed.Name) {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: hotNews.Namespace,
						Name:      hotNews.Name,
					},
				})
			}
		}
	}
	return requests
}

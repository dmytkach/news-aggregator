package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	HttpClient HttpClient
	ServiceURL string
	Finalizer  string
}
type NewsItem struct {
	Title string `json:"Title"`
}

type NewsResponse []NewsItem

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile reads that state of the cluster for a HotNews object and makes changes based on the state read
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for HotNews %s/%s", req.Namespace, req.Name)

	// Fetch the HotNews resource
	var hotNews aggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			log.Print("Reconcile: HotNews was not found. Error ignored")
			return ctrl.Result{}, nil
		}
		log.Printf("Error retrieving HotNews %s/%s from k8s Cluster: %v", req.Namespace, req.Name, err)
		return ctrl.Result{}, err
	}
	if !slices.Contains(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
		hotNews.ObjectMeta.Finalizers = append(hotNews.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &hotNews); err != nil {
			log.Printf("Error adding finalizer to Feed %s/%s: %v", req.Namespace, req.Name, err)
			return ctrl.Result{}, err
		}
	}

	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if slices.Contains(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
			log.Printf("Handling deletion of Feed %s/%s", req.Namespace, req.Name)
			hotNews.ObjectMeta.Finalizers = removeString(hotNews.ObjectMeta.Finalizers, r.Finalizer)
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				log.Printf("Error removing finalizer from Feed %s/%s: %v", req.Namespace, req.Name, err)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	var feeds aggregatorv1.FeedList
	if err := r.List(ctx, &feeds, client.InNamespace(req.Namespace)); err != nil {
		log.Printf("Error listing Feed resources: %v", err)
		return ctrl.Result{}, err
	}
	var feedNames []string
	if len(hotNews.Spec.Feeds) > 0 {
		feedNames = hotNews.Spec.Feeds
	} else if len(hotNews.Spec.FeedGroups) > 0 {
		var feedGroupConfigMap v1.ConfigMap
		if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: "feed-group-source"}, &feedGroupConfigMap); err != nil {
			if errors.IsNotFound(err) {
				log.Print("ConfigMap not found, retrying later")
				return ctrl.Result{}, err
			}
			log.Printf("Error retrieving ConfigMap %s from k8s Cluster: %v", "feed-group-source", err)
			return ctrl.Result{}, err
		}
		feedNames = r.getFeedNamesFromFeedGroups(hotNews.Spec.FeedGroups, feedGroupConfigMap)
	} else {
		feedNames = feeds.GetAllFeedNames()
		return ctrl.Result{}, fmt.Errorf("no feeds or feed groups provided in HotNews spec")
	}

	sources := r.getNewsSourcesFromFeeds(feedNames, feeds)

	log.Printf("Final list of news names: %v", sources)
	log.Printf("Final list of feedNames: %v", feedNames)
	// Process feeds and update HotNews
	status, err := r.fetchNewsData(sources, hotNews.Spec)
	if err != nil {
		if err := r.Client.Update(ctx, &hotNews); err != nil {
			log.Printf("Error updating Feed %s/%s: %v", req.Namespace, req.Name, err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Update HotNews status
	hotNews.Status = status

	if err := r.Status().Update(ctx, &hotNews); err != nil {
		log.Printf("Error updating status of HotNews %s/%s: %v", req.Namespace, req.Name, err)
		return reconcile.Result{}, err
	}

	log.Printf("Successfully updated HotNews %s/%s", req.Namespace, req.Name)
	return reconcile.Result{}, nil
}

// fetchNewsData constructs a request URL using the provided sources and HotNews specifications,
// then sends a request to the news service to retrieve news data. It returns the status of the HotNews
// with the fetched articles or an error if the process fails.
func (r *HotNewsReconciler) fetchNewsData(sources []string, hotNews aggregatorv1.HotNewsSpec) (aggregatorv1.HotNewsStatus, error) {
	log.Printf("Starting fetchNewsData with sources: %v, keywords: %v, dateStart: %s, dateEnd: %s", sources, hotNews.Keywords, hotNews.DateStart, hotNews.DateEnd)
	reqURL, err := buildRequestURL(r.ServiceURL, sources, hotNews.Keywords, hotNews.DateStart, hotNews.DateEnd)
	if err != nil {
		log.Print()
		return aggregatorv1.HotNewsStatus{}, err
	}
	status, err := r.makeRequest(reqURL, hotNews.SummaryConfig.TitlesCount)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return aggregatorv1.HotNewsStatus{}, err
	}
	log.Printf("Request completed successfully, received status: %+v", status)

	return status, nil
}

// getFeedNamesFromFeedGroups retrieves sources from FeedGroups based on the ConfigMap
func (r *HotNewsReconciler) getFeedNamesFromFeedGroups(feedGroups []string, configMap v1.ConfigMap) []string {
	var sources []string

	for _, feedGroup := range feedGroups {
		log.Printf("Processing FeedGroup: %s", feedGroup)
		if value, ok := configMap.Data[feedGroup]; ok {
			sources = append(sources, strings.Split(value, ",")...)
			log.Printf("Matched FeedGroup '%s' in ConfigMap, added values: %v", feedGroup, sources)
		}
	}

	return sources
}

// getNewsSourcesFromFeeds retrieves news names from Feed resources based on the sources
func (r *HotNewsReconciler) getNewsSourcesFromFeeds(sources []string, feeds aggregatorv1.FeedList) []string {
	var newsNames []string

	for _, feed := range feeds.Items {
		if slices.Contains(sources, feed.Name) {
			newsNames = append(newsNames, feed.Spec.Name)
			log.Printf("Matched Feed: %s, adding to newsNames: %s", feed.Name, feed.Spec.Name)
		}
	}

	return newsNames
}

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

// makeRequest performs the HTTP request and processes the response
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
	}, nil
}

// SetupWithManager configures the HotNewsReconciler to manage resources and adds the necessary event predicates
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		Watches(
			&aggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(r.updateHotNews),
		).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.updateHotNews),
		).
		Complete(r)
}

// updateHotNews is a handler function that is triggered when relevant changes
// occur to resources that the controller watches.
func (r *HotNewsReconciler) updateHotNews(context.Context, client.Object) []reconcile.Request {
	var hotNewsList aggregatorv1.HotNewsList
	if err := r.List(context.TODO(), &hotNewsList); err != nil {
		log.Printf("Failed to list HotNews resources %v", err)
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {
		requests = append(requests, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Name:      hotNews.Name,
				Namespace: hotNews.Namespace,
			},
		})
	}

	return requests
}

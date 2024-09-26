package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/handlers"
	"com.teamdev/news-aggregator/internal/controller/predicates"
	"context"
	"encoding/json"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"
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
	Finalizer  string
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
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile performs the reconciliation logic for HotNews resources.
// It manages the lifecycle of the HotNews resource,
// retrieving relevant Feeds and ConfigMaps, and updating the HotNews status with fetched news data.
//
// Upon invocation, the Reconcile function utilizes the Kubernetes client to retrieve the HotNews
// resource. It ensures proper cleanup during deletion by managing HotNews finalizers.
// If the resource does not already possess a finalizer, one is added to facilitate this process.
//
// The function further processes the specified Feeds and FeedGroups.
// When Feeds are explicitly provided in the HotNews specification, they are utilized directly.
// In cases where FeedGroups are specified, the function attempts to retrieve the corresponding
// ConfigMap to extract the associated feeds.
//
// Additionally, the function updates the owner references of the associated Feed resources.
// This step ensures correct tracking and ownership, thereby maintaining the relationship
// between HotNews and its Feed resources.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for HotNews %s/%s", req.Namespace, req.Name)

	var hotNews aggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	owner := HotNewsOwner{
		Client:  r.Client,
		Ctx:     ctx,
		HotNews: &hotNews,
	}
	if !slices.Contains(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
		hotNews.ObjectMeta.Finalizers = append(hotNews.ObjectMeta.Finalizers, r.Finalizer)
		log.Printf("Add finalizer for HotNews %s/%s", req.Namespace, req.Name)
		if err := r.Client.Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
	}
	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if slices.Contains(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
			if cleanupErr := owner.CleanupOwnerReferences(); cleanupErr != nil {
				log.Print(cleanupErr, "Failed to clean up OwnerReferences after HotNews deletion")
				return ctrl.Result{}, cleanupErr
			}
			hotNews.ObjectMeta.Finalizers = removeString(hotNews.ObjectMeta.Finalizers, r.Finalizer)
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	var feedNames []string
	if len(hotNews.Spec.Feeds) > 0 {
		feedNames = hotNews.Spec.Feeds
	} else {
		var feedGroupConfigMap v1.ConfigMap
		if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: r.ConfigMap}, &feedGroupConfigMap); err != nil {
			if errors.IsNotFound(err) {
				return ctrl.Result{}, nil
			}
			hotNews.Status = aggregatorv1.SetHotNewsErrorStatus(err.Error())
			if err := r.Status().Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating Feed %s/%s", req.Namespace, req.Name)
			}
			return ctrl.Result{}, err
		}
		feedNames = hotNews.ExtractFeedsFromGroups(feedGroupConfigMap)
	}

	var feeds aggregatorv1.FeedList
	if err := r.List(ctx, &feeds, client.InNamespace(req.Namespace)); err != nil {
		log.Printf("Error listing Feed resources: %v", err)
		hotNews.Status = aggregatorv1.SetHotNewsErrorStatus(err.Error())
		if err := r.Status().Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating Feed %s/%s", req.Namespace, req.Name)
		}
		return ctrl.Result{}, err
	}

	sources := feeds.GetNewsSources(feedNames)
	log.Printf("Final list of news names: %v", sources)

	status, err := r.fetchNewsData(sources, hotNews.Spec)
	if err != nil {
		hotNews.Status = aggregatorv1.SetHotNewsErrorStatus(err.Error())
		if err := r.Status().Update(ctx, &hotNews); err != nil {
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
	if err := owner.UpdateOwnerReferences(); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Status().Update(ctx, &hotNews); err != nil {
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
	reqURL, err := r.buildRequestURL(sources, hotNews.Keywords, hotNews.DateStart, hotNews.DateEnd)
	if err != nil {
		log.Printf("Error building request: %v", err)
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

// buildRequestURL constructs the URL for the news service request based on the provided sources,
// keywords, and date range. It returns the formatted URL or an error if the URL cannot be constructed
func (r *HotNewsReconciler) buildRequestURL(sources, keywords []string, dateStart, dateEnd string) (string, error) {
	baseURL, err := url.Parse(r.ServiceURL)
	if err != nil {
		return "", fmt.Errorf("invalid service URL: %s", r.ServiceURL)
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
	cm := &handlers.ConfigMapHandler{
		Client: r.Client,
	}
	f := &handlers.FeedHandler{
		Client: r.Client,
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(cm.Handle),
			builder.WithPredicates(predicates.ConfigMapNamePredicate(r.ConfigMap)),
		).
		Watches(
			&aggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(f.Handle)).
		Complete(r)
}

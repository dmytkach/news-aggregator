package controller

import (
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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
}
type NewsItem struct {
	Title string `json:"title"`
}

type NewsResponse struct {
	Articles []NewsItem `json:"articles"`
}

var k8sClient client.Client

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update

// Reconcile reads that state of the cluster for a HotNews object and makes changes based on the state read
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for HotNews %s/%s", req.Namespace, req.Name)

	// Fetch the HotNews resource
	var hotNews aggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			log.Print("Reconcile: HotNews was not found. Error ignored")
			return reconcile.Result{}, nil
		}
		log.Printf("Error retrieving HotNews %s/%s from k8s Cluster: %v", req.Namespace, req.Name, err)
		return reconcile.Result{}, err
	}

	// Fetch ConfigMap
	var feedGroupConfigMap v1.ConfigMap
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: "feed-group-source"}, &feedGroupConfigMap); err != nil {
		if errors.IsNotFound(err) {
			log.Print("ConfigMap not found, retrying later")
			return reconcile.Result{}, err
		}
		log.Printf("Error retrieving ConfigMap %s from k8s Cluster: %v", "feed-group-source", err)
		return reconcile.Result{}, err
	}

	// Fetch Feed resources
	var feeds aggregatorv1.FeedList
	if err := r.List(ctx, &feeds, client.InNamespace(req.Namespace)); err != nil {
		log.Printf("Error listing Feed resources: %v", err)
		return reconcile.Result{}, err
	}

	// Process feeds and update HotNews
	articlesCount, articlesTitles, newsLink, err := r.processFeeds(feeds, feedGroupConfigMap, hotNews)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update HotNews status
	hotNews.Status = aggregatorv1.HotNewsStatus{
		ArticlesCount:  articlesCount,
		NewsLink:       newsLink,
		ArticlesTitles: articlesTitles,
	}

	if err := r.Status().Update(ctx, &hotNews); err != nil {
		log.Printf("Error updating status of HotNews %s/%s: %v", req.Namespace, req.Name, err)
		return reconcile.Result{}, err
	}

	log.Printf("Successfully updated HotNews %s/%s", req.Namespace, req.Name)
	return reconcile.Result{}, nil
}

// processFeeds processes the feeds based on the provided HotNews and ConfigMap
func (r *HotNewsReconciler) processFeeds(feeds aggregatorv1.FeedList, configMap v1.ConfigMap, hotNews aggregatorv1.HotNews) (int, []string, string, error) {
	sources := make([]string, 0)
	feedsNames := hotNews.Spec.Feeds
	if len(feedsNames) != 0 {
		sources = append(sources, feedsNames...)

	}
	for _, news := range hotNews.Spec.FeedGroups {
		for _, key := range configMap.Data {
			if key == news {
				val := strings.Split(configMap.Data[key], ",")
				sources = append(sources, val...)
			}
		}
	}
	newsName := make([]string, 0)
	for _, i := range feeds.Items {
		if slices.Contains(sources, i.Name) {
			newsName = append(newsName, i.Spec.Name)
		}
	}
	reqURL, err := buildRequestURL(r.ServiceURL, sources, hotNews.Spec.Keywords, hotNews.Spec.DateStart, hotNews.Spec.DateEnd)
	if err != nil {
		return 0, nil, "", err
	}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("Failed to create PUT request: %v", err)
		return 0, nil, "", err
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make PUT request: %v", err)
		return 0, nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to create source, status code: %d, response: %s", resp.StatusCode, string(body))
		return 0, nil, "", fmt.Errorf("failed to create source, status code: %d", resp.StatusCode)
	}

	return 0, nil, "", err
}

// SetupWithManager configures the HotNewsReconciler to manage resources and adds the necessary event predicates
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	err := ctrl.NewControllerManagedBy(mgr).
		For(&v1.ConfigMap{}).
		Complete(r)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return !e.DeleteStateUnknown
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).Complete(r)
}
func (r *HotNewsReconciler) getExistingSources(feeds []string, currentNs string) ([]string, error) {
	feedList := &aggregatorv1.FeedList{}
	listOpts := client.ListOptions{Namespace: currentNs}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("checkNameUniqueness: failed to list feeds: %v", err)
	}
	sources := make([]string, 0)
	for _, existingFeed := range feedList.Items {
		if slices.Contains(feeds, existingFeed.Name) {
			sources = append(sources, existingFeed.Spec.Name)
		}
	}
	return sources, nil
}

func buildRequestURL(serviceURL string, sources, keywords []string, dateStart, dateEnd string) (string, error) {
	baseURL, err := url.Parse(serviceURL)
	if err != nil {
		return "", fmt.Errorf("invalid service URL: %v", err)
	}

	params := url.Values{}
	for _, source := range sources {
		params.Add("sources", source)
	}
	for _, keyword := range keywords {
		params.Add("keywords", keyword)
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

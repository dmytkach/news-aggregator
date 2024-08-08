package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"k8s.io/client-go/util/retry"
	"log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

// FeedReconciler is a k8s controller that manages Feed resources.
type FeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	newsAggregatorServiceUrl = "https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources"
	feedFinalizer            = "feeds.finalizers.teamdev.com"
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

// Reconcile reconciles a Feed object.
// This function is called when a change is made to a Feed resource.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for Feed %s/%s", req.Namespace, req.Name)

	var feed aggregatorv1.Feed

	if err := r.Client.Get(ctx, req.NamespacedName, &feed); err != nil {
		if errors.IsNotFound(err) {
			log.Printf("Feed %s/%s not found, returning", req.Namespace, req.Name)
			return ctrl.Result{}, nil
		}
		log.Printf("Failed to get Feed %s/%s: %v", req.Namespace, req.Name, err)
		return ctrl.Result{}, err
	}

	if !feed.ObjectMeta.DeletionTimestamp.IsZero() {
		if slices.Contains(feed.ObjectMeta.Finalizers, feedFinalizer) {
			log.Printf("Finalizer 'feeds.finalizers.teamdev.com' present for Feed %s/%s, handling deletion", req.Namespace, req.Name)

			if _, err := r.handleDelete(feed); err != nil {
				log.Printf("Failed to handle deletion for Feed %s/%s: %v", req.Namespace, req.Name, err)
				return ctrl.Result{}, err
			}

			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				if err := r.Client.Get(ctx, req.NamespacedName, &feed); err != nil {
					return err
				}

				feed.ObjectMeta.Finalizers = removeString(feed.ObjectMeta.Finalizers, feedFinalizer)
				return r.Client.Update(ctx, &feed)
			})

			if retryErr != nil {
				log.Printf("Failed to remove finalizer for Feed %s/%s: %v", req.Namespace, req.Name, retryErr)
				return ctrl.Result{}, retryErr
			}
			log.Printf("Finalizer removed from Feed %s/%s", req.Namespace, req.Name)
		}
		return ctrl.Result{}, nil
	}
	if _, err := r.handleCreate(feed); err != nil {
		log.Printf("Failed to сreate Feed %s/%s: %v. Try to update feed.", feed.Namespace, feed.Name, err)

		if _, err := r.handleUpdate(feed); err != nil {
			log.Printf("Failed to handle update for Feed %s/%s: %v", feed.Namespace, feed.Name, err)
			return ctrl.Result{}, err
		}

		log.Printf("Successfully handled update for Feed %s/%s after failed creation", feed.Namespace, feed.Name)
		return ctrl.Result{}, nil
	}

	log.Printf("Successfully handled creation for Feed %s/%s", feed.Namespace, feed.Name)
	return ctrl.Result{}, nil
}

// handleCreate handles the creation of a new Feed.
// It sends a POST request to the news aggregator service to add the new news source.
func (r *FeedReconciler) handleCreate(feed aggregatorv1.Feed) (ctrl.Result, error) {
	log.Printf("Handling create for Feed %s with URLs %s", feed.Name, feed.Spec.NewUrl)
	if feed.Spec.PreviousURL != "" {
		return ctrl.Result{}, fmt.Errorf("error creating feed %s, try to update", feed.Spec.Name)
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	log.Printf("Start creating resource with url %s:", feed.Spec.NewUrl)
	reqURL := fmt.Sprintf("%s?url=%s", newsAggregatorServiceUrl, feed.Spec.NewUrl)

	resp, err := httpClient.Post(reqURL, "application/json", nil)
	if err != nil {
		log.Printf("Failed to make POST request: %v", err)
		return ctrl.Result{}, err
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
		return ctrl.Result{}, fmt.Errorf("failed to create source, status code: %d", resp.StatusCode)
	}
	for _, condition := range feed.Status.Conditions {
		if condition.Type == aggregatorv1.ConditionAdded {
			err = r.updateFeedStatus(&feed, aggregatorv1.ConditionUpdated, metav1.ConditionTrue, "Link for Source successfully created", "")
			if err != nil {
				log.Print("Failed to add link for existing Source")
				return ctrl.Result{}, err
			}
			log.Printf("Link for existing source %s added successfully", feed.Name)
			return ctrl.Result{}, nil
		}
	}
	err = r.updateFeedStatus(&feed, aggregatorv1.ConditionAdded, metav1.ConditionTrue, "Source successfully created", "")
	if err != nil {
		log.Print("Failed to  create successfully")
		return ctrl.Result{}, err
	}
	log.Print("Source create successfully")
	feed.ObjectMeta.Finalizers = append(feed.ObjectMeta.Finalizers, feedFinalizer)
	if err := r.Client.Update(context.Background(), &feed); err != nil {
		log.Printf("Failed to add finalizer to Feed %s: %v", feed.Name, err)
		return ctrl.Result{}, err
	}
	log.Printf("Finalizer added to Feed %s", feed.Name)
	return ctrl.Result{}, nil
}

// handleUpdate handles the update of an existing Feed.
// It sends a PUT request to the news aggregator service to update the source.
func (r *FeedReconciler) handleUpdate(feed aggregatorv1.Feed) (ctrl.Result, error) {
	log.Printf("Handling update for Feed %s with URL%s", feed.Name, feed.Spec.PreviousURL)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqURL := fmt.Sprintf("%s?oldUrl=%s&newUrl=%s", newsAggregatorServiceUrl, feed.Spec.PreviousURL, feed.Spec.NewUrl)

	req, err := http.NewRequest(http.MethodPut, reqURL, nil)
	if err != nil {
		log.Printf("Failed to create PUT request: %v", err)
		return ctrl.Result{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make PUT request: %v", err)
		return ctrl.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to update source, status code: %d, response: %s", resp.StatusCode, string(body))
		return ctrl.Result{}, fmt.Errorf("failed to update source, status code: %d", resp.StatusCode)
	}

	err = r.updateFeedStatus(&feed, aggregatorv1.ConditionUpdated, metav1.ConditionTrue, "Source successfully updated", "")
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Printf("Source %s updated successfully", feed.Name)
	return ctrl.Result{}, nil
}

// handleDelete handles the deletion of a Feed.
// It sends a DELETE request to the news aggregator service to delete the source.
func (r *FeedReconciler) handleDelete(feed aggregatorv1.Feed) (ctrl.Result, error) {
	log.Printf("Handling deletion for Feed %s", feed.Name)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s?name=%s", newsAggregatorServiceUrl, feed.Spec.Name), nil)
	if err != nil {
		log.Printf("Failed to create DELETE request: %v", err)
		return ctrl.Result{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make DELETE request: %v", err)
		return ctrl.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to delete source, status code: %d, response: %s", resp.StatusCode, string(body))
		return ctrl.Result{}, fmt.Errorf("failed to delete source, status code: %d", resp.StatusCode)
	}

	err = r.updateFeedStatus(&feed, aggregatorv1.ConditionDeleted, metav1.ConditionTrue, "Source successfully deleted", "")
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Printf("Source %s deleted successfully", feed.Spec.Name)
	return ctrl.Result{}, nil
}

// updateFeedStatus updates the status of a Feed resource.
// It modifies the status conditions and updates the resource.
func (r *FeedReconciler) updateFeedStatus(feed *aggregatorv1.Feed, conditionType aggregatorv1.ConditionType, status metav1.ConditionStatus, reason, message string) error {
	log.Printf("Updating status of Feed %s to conditionType %s with status %s", feed.Name, conditionType, status)

	var conditionUpdated bool
	for i, condition := range feed.Status.Conditions {
		if condition.Type == conditionType {
			feed.Status.Conditions[i] = aggregatorv1.Condition{
				Type:           conditionType,
				Status:         status,
				Reason:         reason,
				Message:        message,
				LastUpdateTime: metav1.Now(),
			}
			conditionUpdated = true
			break
		}
	}

	if !conditionUpdated {
		newCondition := aggregatorv1.Condition{
			Type:           conditionType,
			Status:         status,
			Reason:         reason,
			Message:        message,
			LastUpdateTime: metav1.Now(),
		}
		feed.Status.Conditions = append(feed.Status.Conditions, newCondition)
	}

	if err := r.Client.Status().Update(context.Background(), feed); err != nil {
		log.Printf("Failed to update status for Feed %s: %v", feed.Name, err)
		return err
	}
	log.Printf("Status for Feed %s updated successfully", feed.Name)
	return nil
}

// SetupWithManager configures the FeedReconciler to manage resources and
// adds the necessary event predicates to filter events based on changes.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log.Printf("Setting up FeedReconciler with manager")
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

// removeString removes a string from a slice
func removeString(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

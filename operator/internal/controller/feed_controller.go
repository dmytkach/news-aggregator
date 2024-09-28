package controller

import (
	"context"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"slices"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

// HttpClient defines methods for executing HTTP requests.
// It is used by FeedReconciler to communicate with a news aggregator service.
//
//go:generate mockgen -source=feed_controller.go -destination=mock_aggregator/mock_http_client.go -package=controller  news-aggregator/operator/internal/controller HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// FeedReconciler is a k8s controller that manages Feed resources.
// It uses the Client to interact with the Kubernetes API
// and HttpClient to communicate with an external news aggregator service.
type FeedReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	HttpClient    HttpClient
	ServiceURL    string
	FeedFinalizer string
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

// Reconcile is triggered whenever a Feed resource changes.
// It manages the resource's state by handling creation, updates, and deletion processes.
// Additionally, it manages finalizers to ensure that any necessary
// cleanup tasks are performed before the resource is deleted.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Starting reconciliation for Feed %s/%s", req.Namespace, req.Name)

	var feed aggregatorv1.Feed
	if err := r.Client.Get(ctx, req.NamespacedName, &feed); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !slices.Contains(feed.ObjectMeta.Finalizers, r.FeedFinalizer) {
		feed.ObjectMeta.Finalizers = append(feed.ObjectMeta.Finalizers, r.FeedFinalizer)
		if err := r.Client.Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !feed.ObjectMeta.DeletionTimestamp.IsZero() {
		if slices.Contains(feed.ObjectMeta.Finalizers, r.FeedFinalizer) {
			log.Printf("Handling deletion of Feed %s/%s", req.Namespace, req.Name)
			if err := r.deleteFeed(feed); err != nil {
				feed.Status.AddCondition(aggregatorv1.Condition{
					Type:    aggregatorv1.ConditionDeleted,
					Status:  false,
					Message: "Failed to delete feed",
					Reason:  err.Error(),
				})
				if err := r.Client.Status().Update(ctx, &feed); err != nil {
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, err
			}
			feed.Status.AddCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionDeleted,
				Status:  true,
				Message: "Feed deleted successfully",
			})
			feed.ObjectMeta.Finalizers = removeString(feed.ObjectMeta.Finalizers, r.FeedFinalizer)
			if err := r.Client.Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	log.Printf("Current Feed Status Conditions: %v", feed.Status.Conditions)
	if feed.Status.Contains(aggregatorv1.ConditionAdded, true) {
		if err := r.updateFeed(feed); err != nil {
			feed.Status.AddCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionUpdated,
				Status:  false,
				Reason:  err.Error(),
				Message: "Feed didn't update successfully",
			})
			if err := r.Client.Status().Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
		feed.Status.AddCondition(aggregatorv1.Condition{
			Type:    aggregatorv1.ConditionUpdated,
			Status:  true,
			Message: "Feed updated successfully",
		})
	} else {
		if err := r.createFeed(feed); err != nil {
			feed.Status.AddCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionAdded,
				Status:  false,
				Reason:  err.Error(),
				Message: "Feed didn't add successfully",
			})
			if err := r.Client.Status().Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
		feed.Status.AddCondition(aggregatorv1.Condition{
			Type:    aggregatorv1.ConditionAdded,
			Status:  true,
			Message: "Feed added successfully",
		})
	}

	if err := r.Client.Status().Update(ctx, &feed); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// deleteFeed handles the deletion of a Feed.
// It sends a DELETE request to the news aggregator service to delete the source.
func (r *FeedReconciler) deleteFeed(feed aggregatorv1.Feed) error {
	log.Printf("Handling deletion for Feed %s", feed.Name)

	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s?name=%s", r.ServiceURL, feed.Spec.Name), nil)
	if err != nil {
		log.Printf("Failed to create DELETE request: %v", err)
		return err
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make DELETE request: %v", err)
		return err
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
		return fmt.Errorf("failed to delete source, status code: %d", resp.StatusCode)
	}

	log.Printf("Source %s deleted successfully", feed.Spec.Name)
	return nil
}

// createFeed handles the creation of a new Feed.
// It sends a POST request to the news aggregator service to add the new news source.
func (r *FeedReconciler) createFeed(feed aggregatorv1.Feed) error {
	log.Printf("Create Feed %s with URLs %s", feed.Spec.Name, feed.Spec.Link)

	reqURL := fmt.Sprintf("%s?name=%s&url=%s", r.ServiceURL, feed.Spec.Name, feed.Spec.Link)

	resp, err := r.HttpClient.Post(reqURL, "application/json", nil)
	if err != nil {
		log.Printf("Failed to make POST request: %v", err)
		return err
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
		return fmt.Errorf("failed to create source, status code: %d", resp.StatusCode)
	}
	log.Print("Successfully created feed")
	return nil
}

// updateFeed handles the updating of an existing Feed.
// It sends a PUT request to the news aggregator service to update the news source.
func (r *FeedReconciler) updateFeed(feed aggregatorv1.Feed) error {
	log.Printf("Updating Feed %s with URLs %s", feed.Spec.Name, feed.Spec.Link)

	reqURL := fmt.Sprintf("%s?newUrl=%s&name=%s", r.ServiceURL, feed.Spec.Link, feed.Spec.Name)

	req, err := http.NewRequest(http.MethodPut, reqURL, nil)
	if err != nil {
		log.Printf("Failed to create PUT request: %v", err)
		return err
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make PUT request: %v", err)
		return err
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
		return fmt.Errorf("failed to update source, status code: %d", resp.StatusCode)
	}
	log.Print("Successfully updated feed")
	return nil
}

// SetupWithManager configures the FeedReconciler to manage resources and
// adds the necessary event predicates to filter events based on changes.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log.Printf("Setting up FeedReconciler with manager")
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
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

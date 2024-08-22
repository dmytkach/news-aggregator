package v1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"time"
)

const DateFormat = "2006-01-02"

// SetupWebhookWithManager configures the webhook for the HotNews resource with the provided manager.
func (h *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(h).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default implements webhook.Defaulter and sets default values for the HotNews resource.
func (h *HotNews) Default() {
	if h.Spec.SummaryConfig.TitlesCount == 0 {
		h.Spec.SummaryConfig.TitlesCount = 10
	}

	if len(h.Spec.Feeds) == 0 && len(h.Spec.FeedGroups) == 0 {
		feedList := &FeedList{}
		listOpts := client.ListOptions{Namespace: h.Namespace}
		err := k8sClient.List(context.Background(), feedList, &listOpts)
		if err != nil {
			log.Printf("validateFeeds: failed to list feeds: %v", err)
		}
		h.Spec.Feeds = feedList.GetAllFeedNames()
	}

	log.Print("default ", "name ", h.Name)

}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate validates the HotNews resource during creation.
func (h *HotNews) ValidateCreate() (admission.Warnings, error) {
	log.Print("validate create ", "name ", h.Name)
	return h.validate()
}

// ValidateUpdate validates the HotNews resource during updates.
func (h *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	log.Print("validate update", "name", h.Name)

	return h.validate()
}

// ValidateDelete validates the HotNews resource during deletion.
func (h *HotNews) ValidateDelete() (admission.Warnings, error) {
	log.Print("validate delete ", "name ", h.Name)

	return nil, nil
}

// validate the HotNews resource.
// It checks the validity of keywords, dates, and feeds and returns any errors found.
func (h *HotNews) validate() (admission.Warnings, error) {
	var errorsList field.ErrorList
	specPath := field.NewPath("spec")

	err := h.validateKeywords()
	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("keywords"), err.Error()))
	}
	err = h.validateDate()

	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("dateStart"), err.Error()),
			field.Required(specPath.Child("dateEnd"), err.Error()))
	}
	err = h.validateFeeds()
	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("feeds"), err.Error()))
	}

	log.Print("Error list length: ", len(errorsList))
	log.Print("Errors from error list: ", errorsList.ToAggregate())

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

// validateKeywords ensures that at least one keyword is specified in the HotNews resource.
// Returns an error if no keywords are provided.
func (h *HotNews) validateKeywords() error {
	if len(h.Spec.Keywords) == 0 {
		return fmt.Errorf("keywords is required")
	}
	return nil
}

// validateFeeds verifies that the feeds listed in the HotNews resource exist in the namespace.
// Returns an error if any of the specified feeds do not exist.
func (h *HotNews) validateFeeds() error {
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: h.Namespace}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
	if err != nil {
		return fmt.Errorf("validateFeeds: failed to list feeds: %v", err)
	}

	existingFeeds := make(map[string]bool)
	for _, feed := range feedList.Items {
		existingFeeds[feed.Name] = true
	}
	for _, feedName := range h.Spec.Feeds {
		if !existingFeeds[feedName] {
			return fmt.Errorf("validateFeeds: feed %s does not exist in namespace %s", feedName, h.Namespace)
		}
	}
	return nil
}

// validateDate ensures that the start and end dates in the HotNews resource are valid.
// Checks include proper formatting, non-future dates, and that the end date is after the start date.
func (h *HotNews) validateDate() error {
	if h.Spec.DateStart != "" {
		startDate, err := time.Parse(DateFormat, h.Spec.DateStart)
		if err != nil {
			return fmt.Errorf("invalid start date format. Please use YYYY-MM-DD")
		}
		if startDate.After(time.Now()) {
			return fmt.Errorf("start date cannot be in the future")
		}
		if startDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return fmt.Errorf("start date is too old. Please use a more recent date")
		}
	}

	if h.Spec.DateEnd != "" {
		endDate, err := time.Parse(DateFormat, h.Spec.DateEnd)
		if err != nil {
			return fmt.Errorf("invalid end date format. Please use YYYY-MM-DD")
		}
		if endDate.After(time.Now()) {
			return fmt.Errorf("end date cannot be in the future")
		}
		if endDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return fmt.Errorf("end date is too old. Please use a more recent date")
		}
		if h.Spec.DateStart != "" {
			startDate, _ := time.Parse(DateFormat, h.Spec.DateStart)
			if endDate.Before(startDate) {
				return fmt.Errorf("end date must be after start date")
			}
		}
	}
	return nil
}

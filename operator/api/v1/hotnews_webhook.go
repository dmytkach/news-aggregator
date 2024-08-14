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

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HotNews) Default() {
	if r.Spec.SummaryConfig.TitlesCount == 0 {
		r.Spec.SummaryConfig.TitlesCount = 10
	}

	if len(r.Spec.Feeds) == 0 && len(r.Spec.FeedGroups) == 0 {
		feedList := &FeedList{}
		listOpts := client.ListOptions{Namespace: r.Namespace}
		err := k8sClient.List(context.Background(), feedList, &listOpts)
		if err != nil {
			log.Printf("validateFeeds: failed to list feeds: %v", err)
		}
		var feedNameList []string
		for _, feed := range feedList.Items {
			feedNameList = append(feedNameList, feed.Name)
		}
		r.Spec.Feeds = feedNameList
	}

	log.Print("default ", "name ", r.Name)

}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	log.Print("validate create ", "name ", r.Name)
	return r.validateHotNews()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	log.Print("validate update", "name", r.Name)

	return r.validateHotNews()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	log.Print("validate delete ", "name ", r.Name)

	return nil, nil
}

func (r *HotNews) validateHotNews() (admission.Warnings, error) {
	var errorsList field.ErrorList
	specPath := field.NewPath("spec")

	err := r.validateKeywords()
	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("keywords"), err.Error()))
	}
	err = r.validateDate()

	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("dateStart"), err.Error()),
			field.Required(specPath.Child("dateEnd"), err.Error()))
	}
	err = r.validateFeeds()
	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("feeds"), err.Error()))
	}

	log.Print("Error list lenght: ", len(errorsList))
	log.Print("Errors from error list: ", errorsList.ToAggregate())

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

func (r *HotNews) validateKeywords() error {
	if len(r.Spec.Keywords) == 0 {
		return fmt.Errorf("keywords is required")
	}
	return nil
}

func (r *HotNews) validateFeeds() error {
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: r.Namespace}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
	if err != nil {
		return fmt.Errorf("validateFeeds: failed to list feeds: %v", err)
	}

	existingFeeds := make(map[string]bool)
	for _, feed := range feedList.Items {
		existingFeeds[feed.Name] = true
	}
	for _, feedName := range r.Spec.Feeds {
		if !existingFeeds[feedName] {
			return fmt.Errorf("validateFeeds: feed %s does not exist in namespace %s", feedName, r.Namespace)
		}
	}
	return nil
}

func (r *HotNews) validateDate() error {
	if r.Spec.DateStart != "" {
		startDate, err := time.Parse(DateFormat, r.Spec.DateStart)
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

	if r.Spec.DateEnd != "" {
		endDate, err := time.Parse(DateFormat, r.Spec.DateEnd)
		if err != nil {
			return fmt.Errorf("invalid end date format. Please use YYYY-MM-DD")
		}
		if endDate.After(time.Now()) {
			return fmt.Errorf("end date cannot be in the future")
		}
		if endDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return fmt.Errorf("end date is too old. Please use a more recent date")
		}
		if r.Spec.DateStart != "" {
			startDate, _ := time.Parse(DateFormat, r.Spec.DateStart)
			if endDate.Before(startDate) {
				return fmt.Errorf("end date must be after start date")
			}
		}
	}
	return nil
}

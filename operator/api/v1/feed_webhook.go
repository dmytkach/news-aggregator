package v1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"log"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var Client client.Client

// SetupWebhookWithManager configures the manager to handle webhooks for the Feed type.
func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	Client = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Feed{}

// ValidateCreate validates the Feed object during creation.
func (r *Feed) ValidateCreate() (admission.Warnings, error) {
	log.Print("validate create ", "name", r.Name)

	return r.validateFeed()

}

// ValidateUpdate validates the Feed object during updates.
func (r *Feed) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	log.Print("validate update ", "name", r.Name)

	return r.validateFeed()

}

// ValidateDelete validates the Feed object during deletion. No validation is performed in this case.
func (r *Feed) ValidateDelete() (admission.Warnings, error) {
	log.Print("validate delete ", "name", r.Name)

	return nil, nil
}

// validateFeed performs common validation logic for Feed creation and updates.
func (r *Feed) validateFeed() (admission.Warnings, error) {
	var errorsList field.ErrorList
	specPath := field.NewPath("spec")

	if r.Spec.Name == "" {
		errorsList = append(errorsList, field.Required(specPath.Child("name"), "name cannot be empty"))
	} else if len(r.Spec.Name) > 20 {
		errorsList = append(errorsList, field.Invalid(specPath.Child("name"), r.Spec.Name, "name must not exceed 20 characters"))
	}

	if r.Spec.Link == "" {
		errorsList = append(errorsList, field.Required(specPath.Child("link"), "link cannot be empty"))
	} else if !isValidLink(r.Spec.Link) {
		errorsList = append(errorsList, field.Invalid(specPath.Child("link"), r.Spec.Link, "link must be a valid"))
	}

	if err := checkNameUniqueness(r); err != nil {
		errorsList = append(errorsList, field.Invalid(specPath.Child("name"), r.Spec.Name, err.Error()))
	}
	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

// isValidLink checks if the provided string is a valid URL.
func isValidLink(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// checkNameUniqueness ensures that no other Feed with the same name exists in the same namespace.
func checkNameUniqueness(feed *Feed) error {
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: feed.Namespace}
	if err := Client.List(context.Background(), feedList, &listOpts); err != nil {
		return fmt.Errorf("checkNameUniqueness: failed to list feeds: %v", err)

	}

	for _, existingFeed := range feedList.Items {
		if existingFeed.Spec.Name == feed.Spec.Name && existingFeed.Namespace == feed.Namespace && existingFeed.UID != feed.UID {
			return fmt.Errorf("checkNameUniqueness: a Feed with name '%s' already exists in namespace '%s'", feed.Spec.Name, feed.Namespace)
		}
	}
	return nil
}

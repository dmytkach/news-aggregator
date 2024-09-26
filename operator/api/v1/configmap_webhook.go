package v1

import (
	"context"
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// +kubebuilder:webhook:path=/mutate-v1-configmap,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update;delete,versions=v1,name=mconfigmap.kb.io,admissionReviewVersions=v1

// ConfigMapWebHook handles admission requests for ConfigMap resources.
type ConfigMapWebHook struct {
	Client  client.Client
	Decoder admission.Decoder
}

// Handle processes admission requests and applies different logic based on the operation.
func (m *ConfigMapWebHook) Handle(ctx context.Context, req admission.Request) admission.Response {
	configMap := &corev1.ConfigMap{}

	if err := m.Decoder.Decode(req, configMap); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if req.Operation == admissionv1.Delete {
		return m.handleDelete(ctx, req, configMap)
	}

	return m.handleCreateOrUpdate(ctx, req, configMap)
}

// handleCreateOrUpdate processes the creation or update of ConfigMap.
func (m *ConfigMapWebHook) handleCreateOrUpdate(ctx context.Context, req admission.Request, configMap *corev1.ConfigMap) admission.Response {
	for _, value := range configMap.Data {
		feeds := strings.Split(value, ",")
		if res := m.validateFeeds(ctx, req.Namespace, feeds); !res.Allowed {
			return res
		}
	}
	return m.marshalAndPatch(req, configMap)
}

// validateFeeds checks if all feed names specified in ConfigMap exist.
func (m *ConfigMapWebHook) validateFeeds(ctx context.Context, namespace string, feeds []string) admission.Response {
	for _, feedName := range feeds {
		feedName = strings.TrimSpace(feedName)
		feed := &Feed{}

		if err := m.Client.Get(ctx, client.ObjectKey{Name: feedName, Namespace: namespace}, feed); err != nil {
			return admission.Errored(http.StatusNotFound, fmt.Errorf("feed '%s' not found in namespace '%s'", feedName, namespace))
		}
	}
	return admission.Allowed("")
}

// marshalAndPatch marshals the ConfigMap and returns a patch response.
func (m *ConfigMapWebHook) marshalAndPatch(req admission.Request, configMap *corev1.ConfigMap) admission.Response {
	marshaledConfigMap, err := json.Marshal(configMap)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledConfigMap)
}

// handleDelete contains the logic that blocks deletion if any FeedGroup exists in the ConfigMap.
func (m *ConfigMapWebHook) handleDelete(ctx context.Context, req admission.Request, configMap *corev1.ConfigMap) admission.Response {
	hotNewsList := &HotNewsList{}
	listOpts := client.ListOptions{Namespace: req.Namespace}

	if err := m.Client.List(ctx, hotNewsList, &listOpts); err != nil {
		return admission.Errored(http.StatusInternalServerError, fmt.Errorf("failed to list HotNews objects: %v", err))
	}

	return m.checkFeedGroupsInConfigMap(hotNewsList, configMap)
}

// checkFeedGroupsInConfigMap checks if any FeedGroup in HotNews exists as a key in ConfigMap.
func (m *ConfigMapWebHook) checkFeedGroupsInConfigMap(hotNewsList *HotNewsList, configMap *corev1.ConfigMap) admission.Response {
	for _, hotNews := range hotNewsList.Items {
		for _, feedGroup := range hotNews.Spec.FeedGroups {
			if _, exists := configMap.Data[feedGroup]; exists {
				return admission.Denied(fmt.Sprintf("ConfigMap '%s' contains feed group '%s', deletion is not allowed", configMap.Name, feedGroup))
			}
		}
	}

	return admission.Allowed(fmt.Sprintf("ConfigMap '%s' is not used, deletion is allowed", configMap.Name))
}

//// Handle checks whether all feed names specified in the ConfigMap data exist in the cluster.
//func (m *ConfigMapWebHook) Handle(ctx context.Context, req admission.Request) admission.Response {
//	configMap := &corev1.ConfigMap{}
//
//	err := (m.Decoder).Decode(req, configMap)
//	if err != nil {
//		return admission.Errored(http.StatusBadRequest, err)
//	}
//	if req.Operation == admissionv1.Delete {
//		return m.handleDelete(ctx, req)
//	}
//	for _, value := range configMap.Data {
//		feeds := strings.Split(value, ",")
//
//		for _, feedName := range feeds {
//			feedName = strings.TrimSpace(feedName)
//
//			feed := &Feed{}
//			err = m.Client.Get(ctx, client.ObjectKey{Name: feedName, Namespace: req.Namespace}, feed)
//			if err != nil {
//				return admission.Errored(http.StatusNotFound, fmt.Errorf("feed '%s' not found in namespace '%s'", feedName, req.Namespace))
//			}
//		}
//	}
//
//	marshaledConfigMap, err := json.Marshal(configMap)
//	if err != nil {
//		return admission.Errored(http.StatusInternalServerError, err)
//	}
//
//	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledConfigMap)
//}
//
//// handleDelete содержит логику, которая блокирует удаление ConfigMap, если FeedGroup используется
//func (m *ConfigMapWebHook) handleDelete(ctx context.Context, req admission.Request) admission.Response {
//	configMap := &corev1.ConfigMap{}
//	err := m.Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, configMap)
//	if err != nil {
//		return admission.Errored(http.StatusBadRequest, fmt.Errorf("failed to get ConfigMap: %v", err))
//	}
//
//	hotNewsList := &HotNewsList{}
//	listOpts := client.ListOptions{Namespace: req.Namespace}
//	err = m.Client.List(ctx, hotNewsList, &listOpts)
//	if err != nil {
//		return admission.Errored(http.StatusInternalServerError, fmt.Errorf("failed to list HotNews objects: %v", err))
//	}
//
//	for _, hotNews := range hotNewsList.Items {
//		for _, feedGroup := range hotNews.Spec.FeedGroups {
//			if _, exists := configMap.Data[feedGroup]; exists {
//				return admission.Denied(fmt.Sprintf("ConfigMap '%s' содержит feed group '%s', удаление не разрешено", configMap.Name, feedGroup))
//			}
//		}
//	}
//
//	return admission.Allowed(fmt.Sprintf("ConfigMap '%s' не используется, удаление разрешено", configMap.Name))
//}

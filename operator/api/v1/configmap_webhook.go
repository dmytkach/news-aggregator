package v1

import (
	"context"
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// +kubebuilder:webhook:path=/mutate-v1-configmap,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update,versions=v1,name=mconfigmap.kb.io,admissionReviewVersions=v1

// ConfigMapMutator implements admission.Handler
type ConfigMapMutator struct {
	Client  client.Client
	Decoder admission.Decoder
}

// Handle checks the ConfigMap and mutates it if necessary
func (m *ConfigMapMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	configMap := &corev1.ConfigMap{}

	err := (m.Decoder).Decode(req, configMap)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for _, value := range configMap.Data {
		feeds := strings.Split(value, ",")

		for _, feedName := range feeds {
			feedName = strings.TrimSpace(feedName)

			feed := &Feed{}
			err = m.Client.Get(ctx, client.ObjectKey{Name: feedName, Namespace: req.Namespace}, feed)
			if err != nil {
				return admission.Errored(http.StatusNotFound, fmt.Errorf("feed '%s' not found in namespace '%s'", feedName, req.Namespace))
			}
		}
	}

	marshaledConfigMap, err := json.Marshal(configMap)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledConfigMap)
}

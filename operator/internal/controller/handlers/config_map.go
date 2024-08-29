package handlers

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	v1 "k8s.io/api/core/v1"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type ConfigMapHandler struct {
	Client        client.Client
	ConfigMapName string
}

// Handle processes changes to ConfigMap objects and generates reconcile requests for HotNews resources.
func (h *ConfigMapHandler) Handle(ctx context.Context, obj client.Object) []ctrl.Request {
	var requests []ctrl.Request

	log.Print("Starting HandleConfigMap")

	configMap, ok := obj.(*v1.ConfigMap)
	if !ok {
		log.Printf("Object is not a ConfigMap: %v", obj)
		return requests
	}

	if configMap.Name == h.ConfigMapName {
		namespace := configMap.Namespace
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		hotNewsList := &aggregatorv1.HotNewsList{}
		err := h.Client.List(timeoutCtx, hotNewsList, client.InNamespace(namespace))
		if err != nil {
			log.Printf("Error listing HotNews in namespace %s: %v", namespace, err)
			return requests
		}

		log.Printf("Found %d HotNews in namespace %s", len(hotNewsList.Items), namespace)

		for _, hotNews := range hotNewsList.Items {

			if len(hotNews.ExtractFeedsFromGroups(*configMap)) > 0 {
				requests = append(requests, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: hotNews.Namespace,
						Name:      hotNews.Name,
					},
				})

				log.Printf("Enqueued request for HotNews: Name=%s, Namespace=%s", hotNews.Name, hotNews.Namespace)
			}
		}
	}

	log.Printf("Completed HandleConfigMap with %d requests", len(requests))

	return requests
}

package handlers

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	"time"
)

type FeedHandler struct {
	Client client.Client
}

// Handle processes changes to Feed objects and generates reconcile requests for HotNews resources.
func (h *FeedHandler) Handle(ctx context.Context, obj client.Object) []ctrl.Request {
	var requests []ctrl.Request

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	feed, ok := obj.(*aggregatorv1.Feed)
	if !ok {
		log.Printf("Object is not a Feed: %v", obj)
		return requests
	}

	namespace := feed.Namespace
	hotNewsList := &aggregatorv1.HotNewsList{}
	err := h.Client.List(timeoutCtx, hotNewsList, client.InNamespace(namespace))
	if err != nil {
		log.Printf("Error listing HotNews in namespace %s: %v", namespace, err)
		return requests
	}

	for _, hotNews := range hotNewsList.Items {
		if slices.Contains(hotNews.Spec.Feeds, feed.Name) {
			requests = append(requests, ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: hotNews.Namespace,
					Name:      hotNews.Name,
				},
			})
		}
	}

	return requests
}

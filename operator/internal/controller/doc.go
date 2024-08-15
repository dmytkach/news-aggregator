// Package controller provides an API to manage the lifecycle of feed and HotNews resources in a k8s cluster.
// This package includes reconciler implementations that ensure the desired state of these resources is maintained
// within the cluster by interacting with external services and handling the necessary updates.
//
// FeedReconciler is responsible for managing the state of feed resources.
// It ensures that the news sources associated with feed resources are correctly created,
// updated, or deleted in accordance with the cluster's desired state.
//
// HotNewsReconciler manages the lifecycle of HotNews resources. It synchronizes the latest news data with
// the HotNews resources in the cluster by fetching articles from external news services
// based on the specified feeds and feed groups.
package controller

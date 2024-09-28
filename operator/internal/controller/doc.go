// Package controller provides an API to maintain the lifecycle of feed resources in a k8s cluster.
// FeedReconciler is responsible for bringing the state of feed resources in the cluster to a desired state.
// This includes interacting with a news aggregator service to create, update, and delete
// news sources based on the state of the feed resources.
package controller

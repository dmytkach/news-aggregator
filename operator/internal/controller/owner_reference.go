package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
)

type HotNewsOwner struct {
	client.Client
	Ctx     context.Context
	HotNews *aggregatorv1.HotNews
}

// CleanupOwnerReferences removes OwnerReferences to the deleted HotNews from all related Feeds.
func (o *HotNewsOwner) CleanupOwnerReferences() error {
	feeds, err := o.listFeeds()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		o.removeOwnerReferenceFromFeed(&feed)
		if err := o.Client.Update(o.Ctx, &feed); err != nil {
			return fmt.Errorf("failed to remove OwnerReference to Feed %s: %w", feed.Name, err)
		}
	}
	return nil
}

// UpdateOwnerReferences manages owner references for feeds based on their usage in HotNews.
func (o *HotNewsOwner) UpdateOwnerReferences() error {
	feeds, err := o.listFeeds()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		if slices.Contains(o.HotNews.Spec.Feeds, feed.Name) {
			o.addOwnerReferenceToFeed(&feed)
			log.Printf("Successfully added OwnerReference to Feed %s", feed.Name)
		} else {
			o.removeOwnerReferenceFromFeed(&feed)
			log.Printf("Successfully removed OwnerReference from Feed %s", feed.Name)
		}
		if err := o.Client.Update(o.Ctx, &feed); err != nil {
			return fmt.Errorf("failed to update OwnerReference for Feed %s: %w", feed.Name, err)
		}
	}

	return nil
}

// listFeeds retrieves a list of Feeds from the specified namespace.
func (o *HotNewsOwner) listFeeds() ([]aggregatorv1.Feed, error) {
	var feedList aggregatorv1.FeedList
	listOpts := client.ListOptions{Namespace: o.HotNews.Namespace}

	if err := o.Client.List(o.Ctx, &feedList, &listOpts); err != nil {
		return nil, fmt.Errorf("failed to list Feeds in namespace %s: %w", o.HotNews.Namespace, err)
	}

	return feedList.Items, nil
}

// addOwnerReferenceToFeed adds an OwnerReference to a Feed.
func (o *HotNewsOwner) addOwnerReferenceToFeed(feed *aggregatorv1.Feed) {
	for _, ref := range feed.ObjectMeta.OwnerReferences {
		if ref.UID == o.HotNews.UID {
			return
		}
	}

	feed.ObjectMeta.OwnerReferences = append(feed.ObjectMeta.OwnerReferences, metav1.OwnerReference{
		APIVersion: o.HotNews.APIVersion,
		Kind:       o.HotNews.Kind,
		Name:       o.HotNews.Name,
		UID:        o.HotNews.UID,
	})
}

// removeOwnerReferenceFromFeed removes an OwnerReference from a Feed.
func (o *HotNewsOwner) removeOwnerReferenceFromFeed(feed *aggregatorv1.Feed) {
	var updatedReferences []metav1.OwnerReference

	for _, ref := range feed.ObjectMeta.OwnerReferences {
		if ref.Name != o.HotNews.Name {
			updatedReferences = append(updatedReferences, ref)
		}
	}
	feed.ObjectMeta.OwnerReferences = updatedReferences
}

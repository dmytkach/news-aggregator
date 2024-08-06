/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

type FeedReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var feed aggregatorv1.Feed

	err := r.Client.Get(ctx, req.NamespacedName, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Printf("Feed %s : %s", feed.Spec.Name, feed.Spec.Link)

	reqBody, err := json.Marshal(feed.Spec.Link)
	if err != nil {
		log.Printf("Failed to marshal source request: %s", err)
		return ctrl.Result{}, err
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := httpClient.Post("https://news-aggregator-service.news-aggregator.svc.cluster.local:8443/sources", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print("Failed to make POST request: ", err)
		return ctrl.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Print("Failed to create source, status code: ", resp.StatusCode)
		return ctrl.Result{}, err
	}

	feed.Status.Status = "Source is added"
	err = r.Client.Status().Update(ctx, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Print("Status updated.")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
		Complete(r)
}

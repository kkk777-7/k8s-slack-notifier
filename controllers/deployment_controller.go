/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/kkk777-7/k8s-slack-notifier/pkg/notify"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crlog "sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Notifier notify.Notifier
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := crlog.FromContext(ctx)

	// Setup Event Notify
	if r.Notifier == nil {
		n, err := notify.NewNotifier("slack", "/etc/notifier/slack.yaml")
		if err != nil {
			return ctrl.Result{}, err
		}
		r.Notifier = n
	}

	// Get Pod Information
	var pod corev1.Pod
	err := r.Get(ctx, req.NamespacedName, &pod)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		logger.Error(err, "unable to get pod", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	// Deleting Pod
	if !pod.ObjectMeta.DeletionTimestamp.IsZero() {
		title := "Pod DELETING"
		message := fmt.Sprintf("[%v] NameSpace: %s, Name: %s", pod.DeletionTimestamp.In(jst), pod.Namespace, pod.Name)
		logger.Info("Deleting!!! Pod", "podName", pod.Name)
		err := r.Notifier.SendFailEvent(title, message)
		if err != nil {
			logger.Error(err, "unable to send fail event")
		}
	}

	// Created Pod
	if pod.Status.Phase == "Running" && IsCreatePod(pod.ObjectMeta.CreationTimestamp) {
		title := "Pod CREATE"
		message := fmt.Sprintf("[%v] NameSpace: %s, Name: %s", pod.CreationTimestamp.In(jst), pod.Namespace, pod.Name)
		logger.Info("Create!!! Pod", "podName", pod.Name)
		err := r.Notifier.SendSuccessEvent(title, message)
		if err != nil {
			logger.Error(err, "unable to send success event")
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}

// Calculate the difference between the time the pod was created and the current time
// Positive if less than 5 seconds
func IsCreatePod(t v1.Time) bool {
	now := time.Now()
	return now.Sub(t.Time) < 5*time.Second
}

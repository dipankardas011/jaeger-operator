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
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	jaegertracingiov2alpha1 "github.com/jaegertracing/jaeger-operator/v2/operator/api/v2alpha1"
	"github.com/jaegertracing/jaeger-operator/v2/operator/pkg/helpers"
	"github.com/jaegertracing/jaeger-operator/v2/operator/pkg/template"
)

// JaegerReconciler reconciles a Jaeger object
type JaegerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=jaegertracing.io,resources=jaegers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jaegertracing.io,resources=jaegers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jaegertracing.io,resources=jaegers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Jaeger object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *JaegerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.WithName("jaeger-v2-operator")

	jaegerResource := &jaegertracingiov2alpha1.Jaeger{}

	log.V(-1).Info(
		"Reconciling Jaeger",
		"Name", req.Name,
		"Namespace", req.Namespace,
	)

	err := r.Get(ctx, req.NamespacedName, jaegerResource)
	if err != nil {
		if errors.IsNotFound(err) {

			log.V(-1).Info("CleanUp resources")

			podSpec := template.PodSpec(helpers.DELETION_OPERATION, req.Name, req.Namespace)

			if err := r.Client.Delete(ctx, podSpec); err != nil {
				return ctrl.Result{}, err
			}
			svcSpec := template.ServiceSpec(helpers.DELETION_OPERATION, req.Name, req.Namespace)

			if err := r.Client.Delete(ctx, svcSpec); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to fetch <kind>")
		return ctrl.Result{}, nil
	}

	log.V(-1).Info("Creating resources", "configuration", jaegerResource.Spec.Config)
	podSpec := template.PodSpec(helpers.CREATION_OPERATION, jaegerResource.Name, jaegerResource.Namespace)

	if err := r.Client.Create(ctx, podSpec); err != nil {
		return ctrl.Result{}, err
	}

	svcSpec := template.ServiceSpec(helpers.CREATION_OPERATION, jaegerResource.Name, jaegerResource.Namespace)
	if err := r.Client.Create(ctx, svcSpec); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JaegerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jaegertracingiov2alpha1.Jaeger{}).
		Complete(r)
}

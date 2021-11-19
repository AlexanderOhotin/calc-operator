/*
Copyright 2021.

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
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"

	appsv1 "github.com/sd01dev/demo-operator/api/v1"
)

// CalculatorReconciler reconciles a Calculator object
type CalculatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.sd01dev.com,resources=calculators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.sd01dev.com,resources=calculators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.sd01dev.com,resources=calculators/finalizers,verbs=update

func (r *CalculatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	calc := &appsv1.Calculator{}

	err := r.Client.Get(ctx, req.NamespacedName, calc)
	if err != nil {
		r.Log.Info("Failed to get Calculator")
		return ctrl.Result{}, err
	}

	calc.Status.Processed = true
	calc.Status.Result = calc.Spec.X + calc.Spec.Z
	err = r.Status().Update(ctx, calc)
	if err != nil {
		r.Log.Info("Failed to set status")
		return ctrl.Result{}, err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      calc.Name,
			Namespace: calc.Namespace,
			Annotations: map[string]string{
				"managed-by": "calc-operator",
			},
		},
		Immutable: nil,
		Data:      nil,
		StringData: map[string]string{
			"result": strconv.FormatInt(int64(calc.Status.Result), 10),
		},
		Type: v1.SecretTypeOpaque,
	}

	err = r.Client.Create(ctx, secret)
	if err != nil {
		r.Log.Info("Failed to create secret")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CalculatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Calculator{}).
		Complete(r)
}

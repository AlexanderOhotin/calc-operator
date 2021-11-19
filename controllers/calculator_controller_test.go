package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"

	appsv1 "github.com/sd01dev/demo-operator/api/v1"
)

func TestCalculatorReconciler_ReconcileSuccess(t *testing.T) {

	calc := &appsv1.Calculator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example one",
			Namespace: "default",
		},
		Spec: appsv1.CalculatorSpec{
			X: 42,
			Z: 34,
		},
		Status: appsv1.CalculatorStatus{
			Processed: false,
			Result:    0,
		},
	}

	secret := &v1.Secret{}

	objs := runtime.Object(calc)

	ch := scheme.Scheme
	ch.AddKnownTypes(v1.SchemeGroupVersion, calc)

	cl := fake.NewClientBuilder().WithRuntimeObjects(objs).Build()

	cr := CalculatorReconciler{
		Client: cl,
		Log:    log.NullLogger{},
		Scheme: ch,
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: calc.Namespace,
		Name:      calc.Name,
	}}

	_, err := cr.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = cl.Get(context.Background(), types.NamespacedName{
		Name:      calc.Name,
		Namespace: calc.Namespace,
	}, secret)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = cl.Get(context.Background(), types.NamespacedName{
		Namespace: calc.Namespace,
		Name:      calc.Name,
	}, calc)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// check status
	assert.Equal(t, int32(76), calc.Status.Result)
	assert.Equal(t, true, calc.Status.Processed)
	// check secret data
	assert.Equal(t, map[string]string{"managed-by": "calc-operator"}, secret.Annotations)
	assert.Equal(t, map[string]string{"result": "76"}, secret.StringData)

}

func TestCalculatorReconciler_ReconcileFailedToSetStatus(t *testing.T) {

	calc := appsv1.Calculator{}

	objs := runtime.Object(&calc)

	ch := scheme.Scheme
	ch.AddKnownTypes(v1.SchemeGroupVersion, &calc)

	cl := fake.NewClientBuilder().WithRuntimeObjects(objs).Build()

	cr := CalculatorReconciler{
		Client: cl,
		Log:    log.NullLogger{},
		Scheme: ch,
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: calc.Namespace,
		Name:      calc.Name,
	}}

	_, err := cr.Reconcile(context.Background(), req)
	if err != nil {
		assert.Error(t, errors.NewNotFound(schema.GroupResource{}, ""))
	}

}

func TestCalculatorReconciler_ReconcileObjNotFound(t *testing.T) {

	calc := appsv1.Calculator{}

	objs := runtime.Object(&calc)

	ch := scheme.Scheme
	ch.AddKnownTypes(v1.SchemeGroupVersion, &calc)

	cl := fake.NewClientBuilder().WithRuntimeObjects(objs).Build()

	cr := CalculatorReconciler{
		Client: cl,
		Log:    log.NullLogger{},
		Scheme: ch,
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: "none",
		Name:      calc.Name,
	}}

	_, err := cr.Reconcile(context.Background(), req)
	if err != nil {
		assert.Error(t, errors.NewNotFound(schema.GroupResource{}, "Failure"))
	}


}

func TestCalculatorReconciler_ReconcileGetObjFail(t *testing.T) {

	calc := appsv1.Calculator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example one",
			Namespace: "default",
		},
		Spec: appsv1.CalculatorSpec{
			X: 42,
			Z: 34,
		},
		Status: appsv1.CalculatorStatus{
			Processed: false,
			Result:    0,
		},
	}

	objs := runtime.Object(&calc)

	ch := scheme.Scheme
	ch.AddKnownTypes(v1.SchemeGroupVersion, &calc)

	cl := fake.NewClientBuilder().WithRuntimeObjects(objs).Build()

	cr := CalculatorReconciler{
		Client: cl,
		Log:    log.NullLogger{},
		Scheme: ch,
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: calc.Namespace,
		Name:      calc.Name,
	}}

	cr.Reconcile(context.Background(), req)

}

/*
Copyright 2023.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	traitv1alpha1 "github.com/sguruvar/oamvelaop/api/v1alpha1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// OamTraitReconciler reconciles a OamTrait object
type OamTraitReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=trait.oam.vela,resources=oamtraits,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=trait.oam.vela,resources=oamtraits/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=trait.oam.vela,resources=oamtraits/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OamTrait object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *OamTraitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	trait := &traitv1alpha1.OamTrait{}
	err := r.Get(ctx, req.NamespacedName, trait)
	if err != nil {
		l.Error(err, "unable to get OamTrait")
		return ctrl.Result{}, err
	}
	var ingressClassName = "aws-alb"
	hpaAvgUtilization := int32(trait.Spec.CpuTarget)
	img := trait.Spec.Image

	var traits traitv1alpha1.OamTraitList
	if err := r.List(ctx, &traits); err != nil {
		l.Error(err, "Unable to list OamTrait instances")
		return ctrl.Result{}, err
	}
	//for _, trait := range traits.Items {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trait.Name,
			Namespace: trait.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &trait.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": trait.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": trait.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  trait.Name,
							Image: img,
						},
					},
				},
			},
		},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trait.Name,
			Namespace: trait.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       trait.Spec.Port,
					TargetPort: intstr.FromInt(int(trait.Spec.Port)),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": trait.Name,
			},
		},
	}
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trait.Name,
			Namespace: trait.Namespace,
			Annotations: map[string]string{
				"alb.ingress.kubernetes.io/scheme":      "internet-facing",
				"alb.ingress.kubernetes.io/target-type": "ip",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: trait.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: trait.Spec.Port,
											},
										},
									},
									Path:     "/",
									PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
								},
							},
						},
					},
				},
			},
		},
	}
	hpa := &autoscalingv1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trait.Name,
			Namespace: trait.Namespace,
		},
		Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       trait.Name,
				APIVersion: "apps/v1",
			},
			MinReplicas:                    &trait.Spec.Replicas,
			MaxReplicas:                    trait.Spec.MaxReplicas,
			TargetCPUUtilizationPercentage: &hpaAvgUtilization,
		},
	}

	foundDeploy := &appsv1.Deployment{}
	foundService := &corev1.Service{}
	foundIngress := &networkingv1.Ingress{}
	foundHpa := &autoscalingv1.HorizontalPodAutoscaler{}
	err = r.Get(ctx, types.NamespacedName{Name: trait.Name, Namespace: trait.Namespace}, foundDeploy)
	errSvc := r.Get(ctx, types.NamespacedName{Name: trait.Name, Namespace: trait.Namespace}, foundService)
	errIngress := r.Get(ctx, req.NamespacedName, foundIngress)
	errHpa := r.Get(ctx, client.ObjectKey{Namespace: trait.Namespace, Name: trait.Name}, foundHpa)

	if err != nil && errors.IsNotFound(err) && errors.IsNotFound(errSvc) && errors.IsNotFound(errIngress) && errors.IsNotFound(errHpa) {
		//create deploy
		l.Info("creating deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
		if err := r.Create(ctx, deployment); err != nil {
			l.Error(err, "Failed to create deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		l.Info("creating service", "service.Namespace", service.Namespace, "service.Name", service.Name)
		if err := r.Create(ctx, service); err != nil {
			l.Error(err, "Failed to create service", "service.Namespace", service.Namespace, "service.Name", service.Name)
			return ctrl.Result{}, err
		}
		l.Info("creating ingress", "ingress.Namespace", ingress.Namespace, "ingress.Name", ingress.Name)
		if trait.Spec.IngressReq == "yes" {
			if err := r.Create(ctx, ingress); err != nil {
				l.Error(err, "Failed to create ingress", "ingress.Namespace", ingress.Namespace, "ingress.Name", ingress.Name)
				return ctrl.Result{}, err
			}
		}
		l.Info("creating hpa", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
		if err := r.Create(ctx, hpa); err != nil {
			l.Error(err, "Failed to create hpa", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		l.Error(err, "Failed to get deployment or service")
		return ctrl.Result{}, err
	}

	//deploy already exists
	foundDeploy.Spec = deployment.Spec
	l.Info("updating deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
	if err := r.Update(ctx, foundDeploy); err != nil {
		l.Error(err, "Failed to update deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
		return ctrl.Result{}, err
	}
	//service already exists
	foundService.Spec = service.Spec
	l.Info("updating service", "service.Namespace", service.Namespace, "service.Name", service.Name)
	if err := r.Update(ctx, foundService); err != nil {
		l.Error(err, "Failed to update service", "service.Namespace", service.Namespace, "service.Name", service.Name)
		return ctrl.Result{}, err
	}
	//ingress already exists
	foundIngress.Spec = ingress.Spec
	l.Info("updating ingress", "ingress.Namespace", ingress.Namespace, "ingress.Name", ingress.Name)
	if err := r.Update(ctx, foundIngress); err != nil {
		l.Error(err, "Failed to update ingress", "ingress.Namespace", ingress.Namespace, "ingress.Name", ingress.Name)
		return ctrl.Result{}, err
	}

	//hpa already exists
	foundHpa.Spec = hpa.Spec
	l.Info("updating hpa", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
	if err := r.Update(ctx, foundHpa); err != nil {
		l.Error(err, "Failed to update hpa", "hpa.Namespace", hpa.Namespace, "hpa.Name", hpa.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
	//}
	//return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OamTraitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traitv1alpha1.OamTrait{}).
		Complete(r)
}

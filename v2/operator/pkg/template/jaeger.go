package template

import (
	"github.com/jaegertracing/jaeger-operator/v2/operator/pkg/helpers"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			ContainerPort: 5775,
			Name:          "zk-compact-trft", // max 15 chars!
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 5778,
			Name:          "config-rest",
		},
		{
			ContainerPort: 6831,
			Name:          "jg-compact-trft",
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 6832,
			Name:          "jg-binary-trft",
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 9411,
			Name:          "zipkin",
		},
		{
			ContainerPort: 14267,
			Name:          "c-tchan-trft", // for collector
		},
		{
			ContainerPort: 14268,
			Name:          "c-binary-trft",
		},
		{
			ContainerPort: 16685,
			Name:          "grpc-query",
		},
		{
			ContainerPort: 16686,
			Name:          "query",
		},
		{
			ContainerPort: 14250,
			Name:          "grpc",
		},
	}
}

// TODO: add object reference as that of controller or deployment so that if the deployment goes the service goes as well

func ServiceSpec(op helpers.Operation, name, ns string) *corev1.Service {
	spec := &corev1.Service{
		TypeMeta: v1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: helpers.GetServiceName(name),
			Labels: map[string]string{
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": "jaeger-operator",
			},
			Namespace: ns,
		},
	}

	if op == helpers.CREATION_OPERATION {
		spec.Spec = corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": "jaeger-operator",
			},
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "ui",
					Port:       16686,
					TargetPort: intstr.FromInt(16686),
				},
				{
					Name:       "conn",
					Port:       14268,
					TargetPort: intstr.FromInt(14268),
				},
			},
		}
	}

	return spec
}

func PodSpec(op helpers.Operation, name, ns string) *corev1.Pod {
	spec := &corev1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      helpers.GetDeploymentName(name),
			Namespace: ns,
			Labels: map[string]string{
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": "jaeger-operator",
			},
		},
	}

	if op == helpers.CREATION_OPERATION {
		spec.Spec = corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:  "jaeger",
					Image: "jaegertracing/jaeger:latest",
					Ports: getPorts(),
				},
			},
		}
	}

	return spec
}

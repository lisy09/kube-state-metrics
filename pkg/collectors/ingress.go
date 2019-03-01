/*
Copyright 2019 The Kubernetes Authors All rights reserved.
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

package collectors

import (
	"k8s.io/kube-state-metrics/pkg/metrics"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	descIngressLabelsName          = "kube_ingress_labels"
	descIngressLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descIngressLabelsDefaultLabels = []string{"namespace", "ingress"}

	ingressMetricFamilies = []metrics.FamilyGenerator{
		metrics.FamilyGenerator{
			Name: descIngressLabelsName,
			Type: metrics.MetricTypeGauge,
			Help: descIngressLabelsHelp,
			GenerateFunc: wrapIngressFunc(func(s *v1beta1.Ingress) metrics.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(s.Labels)
				return metrics.Family{&metrics.Metric{
					Name:        descIngressLabelsName,
					LabelKeys:   labelKeys,
					LabelValues: labelValues,
					Value:       1,
				}}
			}),
		},
	}
)

func wrapIngressFunc(f func(*v1beta1.Ingress) metrics.Family) func(interface{}) metrics.Family {
	return func(obj interface{}) metrics.Family {
		ingress := obj.(*v1beta1.Ingress)
		metricFamily := f(ingress)

		for _, m := range metricFamily {
			m.LabelKeys = append(descIngressLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{ingress.Namespace, ingress.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createIngressListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return kubeClient.ExtensionsV1beta1().Ingresses(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return kubeClient.ExtensionsV1beta1().Ingresses(ns).Watch(opts)
		},
	}
}

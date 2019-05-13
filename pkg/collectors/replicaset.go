/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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
	descReplicaSetLabelsDefaultLabels = []string{"namespace", "replicaset"}

	replicaSetMetricFamilies = []metrics.FamilyGenerator{
		metrics.FamilyGenerator{
			Name: "kube_replicaset_created",
			Type: metrics.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapReplicaSetFunc(func(r *v1beta1.ReplicaSet) metrics.Family {
				f := metrics.Family{}

				if !r.CreationTimestamp.IsZero() {
					f = append(f, &metrics.Metric{
						Name:  "kube_replicaset_created",
						Value: float64(r.CreationTimestamp.Unix()),
					})
				}

				return f
			}),
		},
	}
)

func wrapReplicaSetFunc(f func(*v1beta1.ReplicaSet) metrics.Family) func(interface{}) metrics.Family {
	return func(obj interface{}) metrics.Family {
		replicaSet := obj.(*v1beta1.ReplicaSet)

		metricFamily := f(replicaSet)

		for _, m := range metricFamily {
			m.LabelKeys = append(descReplicaSetLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{replicaSet.Namespace, replicaSet.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createReplicaSetListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return kubeClient.ExtensionsV1beta1().ReplicaSets(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return kubeClient.ExtensionsV1beta1().ReplicaSets(ns).Watch(opts)
		},
	}
}

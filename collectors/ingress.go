/*
Copyright 2017 The Kubernetes Authors All rights reserved.

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
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

var (
	descIngressLabelsName          = "kube_ingress_labels"
	descIngressLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descIngressLabelsDefaultLabels = []string{"namespace", "ingress"}

	descIngressLabels = prometheus.NewDesc(
		descIngressLabelsName,
		descIngressLabelsHelp,
		descIngressLabelsDefaultLabels, nil,
	)
)

type IngressLister func() ([]v1beta1.Ingress, error)

func (l IngressLister) List() ([]v1beta1.Ingress, error) {
	return l()
}

func RegisterIngressCollector(registry prometheus.Registerer, kubeClient kubernetes.Interface, namespaces []string) {
	client := kubeClient.ExtensionsV1beta1().RESTClient()
	glog.Infof("collect ingress with %s", client.APIVersion())

	sinfs := NewSharedInformerList(client, "ingresses", namespaces, &v1beta1.Ingress{})

	ingressLister := IngressLister(func() (ingress []v1beta1.Ingress, err error) {
		for _, sinf := range *sinfs {
			for _, m := range sinf.GetStore().List() {
				ingress = append(ingress, *m.(*v1beta1.Ingress))
			}
		}
		return ingress, nil
	})

	registry.MustRegister(&ingressCollector{store: ingressLister})
	sinfs.Run(context.Background().Done())
}

type ingressStore interface {
	List() (Ingressvices []v1beta1.Ingress, err error)
}

// ingressCollector collects metrics about all ingresses in the cluster.
type ingressCollector struct {
	store ingressStore
}

// Describe implements the prometheus.Collector interface.
func (pc *ingressCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- descIngressLabels
}

// Collect implements the prometheus.Collector interface.
func (sc *ingressCollector) Collect(ch chan<- prometheus.Metric) {
	ingresses, err := sc.store.List()
	if err != nil {
		ScrapeErrorTotalMetric.With(prometheus.Labels{"resource": "ingress"}).Inc()
		glog.Errorf("listing ingresses failed: %s", err)
		return
	}
	ScrapeErrorTotalMetric.With(prometheus.Labels{"resource": "ingress"}).Add(0)

	ResourcesPerScrapeMetric.With(prometheus.Labels{"resource": "ingress"}).Observe(float64(len(ingresses)))
	for _, s := range ingresses {
		sc.collectIngress(ch, s)
	}
	glog.V(4).Infof("collected %d ingresses", len(ingresses))
}

func ingressLabelsDesc(labelKeys []string) *prometheus.Desc {
	return prometheus.NewDesc(
		descIngressLabelsName,
		descIngressLabelsHelp,
		append(descIngressLabelsDefaultLabels, labelKeys...),
		nil,
	)
}

func (sc *ingressCollector) collectIngress(ch chan<- prometheus.Metric, s v1beta1.Ingress) {
	addConstMetric := func(desc *prometheus.Desc, t prometheus.ValueType, v float64, lv ...string) {
		lv = append([]string{s.Namespace, s.Name}, lv...)
		ch <- prometheus.MustNewConstMetric(desc, t, v, lv...)
	}
	addGauge := func(desc *prometheus.Desc, v float64, lv ...string) {
		addConstMetric(desc, prometheus.GaugeValue, v, lv...)
	}

	labelKeys, labelValues := kubeLabelsToPrometheusLabels(s.Labels)
	addGauge(ingressLabelsDesc(labelKeys), 1, labelValues...)
}

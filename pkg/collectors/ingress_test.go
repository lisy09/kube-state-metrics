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
	"testing"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIngressCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.

	startTime := 1501569018
	metav1StartTime := metav1.Unix(int64(startTime), 0)

	const metadata = `
		# HELP kube_ingress_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_ingress_labels gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "ingress1",
					Namespace:       "ns1",
					ResourceVersion: "000000",
				},
			},
			Want: `
				kube_ingress_labels{namespace="ns1",ingress="ingress1"} 1
`,
			MetricNames: []string{"kube_ingress_labels"},
		},
		{
			Obj: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ingress2",
					Namespace:         "ns2",
					CreationTimestamp: metav1StartTime,
					ResourceVersion:   "123456",
				},
			},
			Want: `
				kube_ingress_labels{namespace="ns2",ingress="ingress2"} 1
				`,
			MetricNames: []string{"kube_ingress_labels"},
		},
		{
			Obj: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ingress3",
					Namespace:         "ns3",
					CreationTimestamp: metav1StartTime,
					Labels:            map[string]string{"test-3": "test-3"},
					ResourceVersion:   "abcdef",
				},
			},
			Want: `
				kube_ingress_labels{label_test_3="test-3",namespace="ns3",ingress="ingress3"} 1
`,
			MetricNames: []string{"kube_ingress_labels"},
		},
	}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(ingressMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}

	}
}

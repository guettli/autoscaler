/*
Copyright 2023 The Kubernetes Authors.

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

package safetoevict

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/simulator/drainability"
	"k8s.io/autoscaler/cluster-autoscaler/utils/drain"
)

func TestDrainable(t *testing.T) {
	for desc, tc := range map[string]struct {
		pod  *apiv1.Pod
		want drainability.Status
	}{
		"regular pod": {
			pod: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "ns",
				},
			},
			want: drainability.NewUndefinedStatus(),
		},
		"safe to evict pod": {
			pod: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
					Annotations: map[string]string{
						drain.PodSafeToEvictKey: "true",
					},
				},
			},
			want: drainability.NewDrainableStatus(),
		},
	} {
		t.Run(desc, func(t *testing.T) {
			got := New().Drainable(nil, tc.pod)
			if tc.want != got {
				t.Errorf("Rule.Drainable(%v) = %v, want %v", tc.pod.Name, got, tc.want)
			}
		})
	}
}

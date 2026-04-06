/*
Copyright 2024 The Kubermatic Kubernetes Platform contributors.

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

package headlamp

import (
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/reconciler/pkg/reconciling"

	rbacv1 "k8s.io/api/rbac/v1"
)

// ClusterRoleReconciler returns a ClusterRole for Headlamp with read-only access to cluster resources.
func ClusterRoleReconciler() reconciling.NamedClusterRoleReconcilerFactory {
	return func() (string, reconciling.ClusterRoleReconciler) {
		return resources.HeadlampClusterRoleName, func(cr *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
			cr.Labels = resources.BaseAppLabels(AppName, nil)

			cr.Rules = []rbacv1.PolicyRule{
				{
					APIGroups: []string{
						"",
						"apps",
						"batch",
						"networking.k8s.io",
						"rbac.authorization.k8s.io",
						"storage.k8s.io",
						"apiextensions.k8s.io",
						"policy",
						"autoscaling",
					},
					Resources: []string{"*"},
					Verbs:     []string{"get", "list", "watch"},
				},
				{
					APIGroups: []string{"metrics.k8s.io"},
					Resources: []string{"pods", "nodes"},
					Verbs:     []string{"get", "list", "watch"},
				},
			}
			return cr, nil
		}
	}
}

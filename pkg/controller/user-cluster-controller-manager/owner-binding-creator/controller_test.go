/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package ownerbindingcreator

import (
	"context"
	"testing"

	handlercommon "k8c.io/kubermatic/v2/pkg/handler/common"
	kubermaticlog "k8c.io/kubermatic/v2/pkg/log"
	"k8c.io/kubermatic/v2/pkg/test/diff"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	fakectrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcile(t *testing.T) {
	testCases := []struct {
		name            string
		clusterRole     *rbacv1.ClusterRole
		requestName     string
		ownerEmail      string
		expectedBinding rbacv1.ClusterRoleBinding
	}{
		{
			name: "cluster role not found, no error",
		},
		{
			name: "role binding created with cluster owner subject",
			clusterRole: &rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:   "cluster-admin",
				Labels: map[string]string{handlercommon.UserClusterComponentKey: handlercommon.UserClusterRoleComponentValue},
			}},
			requestName: "cluster-admin",
			ownerEmail:  "test@test.com",
			expectedBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Labels:          map[string]string{handlercommon.UserClusterComponentKey: handlercommon.UserClusterBindingComponentValue},
					ResourceVersion: "1",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "ClusterRole",
					Name:     "cluster-admin",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:     rbacv1.UserKind,
						APIGroup: rbacv1.GroupName,
						Name:     "test@test.com",
					},
				},
			},
		},
		{
			name: "role binding created for no admin role",
			clusterRole: &rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
				Name:   "view",
				Labels: map[string]string{handlercommon.UserClusterComponentKey: handlercommon.UserClusterRoleComponentValue},
			}},
			requestName: "view",
			ownerEmail:  "test@test.com",
			expectedBinding: rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Labels:          map[string]string{handlercommon.UserClusterComponentKey: handlercommon.UserClusterBindingComponentValue},
					ResourceVersion: "1",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.GroupName,
					Kind:     "ClusterRole",
					Name:     "view",
				},
			},
		},
	}

	for idx := range testCases {
		tc := testCases[idx]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientBuilder := fakectrlruntimeclient.NewClientBuilder().WithScheme(scheme.Scheme)
			if tc.clusterRole != nil {
				clientBuilder.WithObjects(tc.clusterRole)
			}

			client := clientBuilder.Build()
			r := &reconciler{
				log:        kubermaticlog.Logger,
				client:     client,
				recorder:   record.NewFakeRecorder(10),
				ownerEmail: tc.ownerEmail,
				clusterIsPaused: func(c context.Context) (bool, error) {
					return false, nil
				},
			}

			ctx := context.Background()
			request := reconcile.Request{NamespacedName: types.NamespacedName{Name: tc.requestName}}
			if _, err := r.Reconcile(ctx, request); err != nil {
				t.Fatalf("reconciling failed: %v", err)
			}

			if tc.clusterRole == nil {
				return
			}

			clusterRoleBindingList := &rbacv1.ClusterRoleBindingList{}
			if err := client.List(ctx, clusterRoleBindingList, ctrlruntimeclient.MatchingLabels{handlercommon.UserClusterComponentKey: handlercommon.UserClusterBindingComponentValue}); err != nil {
				t.Fatalf("failed to list cluster role bindigs: %v", err)
			}

			if len(clusterRoleBindingList.Items) != 1 {
				t.Fatalf("expected exactly one binding, got %d", len(clusterRoleBindingList.Items))
			}

			existingBinding := clusterRoleBindingList.Items[0]
			existingBinding.Name = tc.expectedBinding.Name

			if !diff.SemanticallyEqual(tc.expectedBinding, existingBinding) {
				t.Fatalf("bindings are not equal:\n%v", diff.ObjectDiff(tc.expectedBinding, existingBinding))
			}
		})
	}
}

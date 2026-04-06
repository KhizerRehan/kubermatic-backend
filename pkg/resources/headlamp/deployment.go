/*
Copyright 2025 The Kubermatic Kubernetes Platform contributors.

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
	"fmt"

	kubermaticv1 "k8c.io/kubermatic/sdk/v2/apis/kubermatic/v1"
	"k8c.io/kubermatic/sdk/v2/semver"
	"k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/resources/apiserver"
	"k8c.io/kubermatic/v2/pkg/resources/registry"
	"k8c.io/reconciler/pkg/reconciling"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

var (
	defaultResourceRequirements = map[string]*corev1.ResourceRequirements{
		resources.HeadlampDeploymentName: {
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("128Mi"),
				corev1.ResourceCPU:    resource.MustParse("100m"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("256Mi"),
				corev1.ResourceCPU:    resource.MustParse("250m"),
			},
		},
	}
)

const (
	name               = resources.HeadlampDeploymentName
	imageName          = "headlamp-k8s/headlamp"
	ContainerPort      = 4466
	AppLabel           = resources.AppLabelKey + "=" + name
	tmpVolumeName      = "tmp-volume"
	configVolumeName   = "headlamp-config"
	frontendVolumeName = "headlamp-frontend"

	headlampVersion = "v0.26.0"
)

// headlampData is the data needed to construct the Headlamp components.
type headlampData interface {
	Cluster() *kubermaticv1.Cluster
	RewriteImage(string) (string, error)
}

// DeploymentReconciler returns the function to create and update the Headlamp deployment.
func DeploymentReconciler(data headlampData) reconciling.NamedDeploymentReconcilerFactory {
	return func() (string, reconciling.DeploymentReconciler) {
		return name, func(dep *appsv1.Deployment) (*appsv1.Deployment, error) {
			baseLabels := resources.BaseAppLabels(name, nil)
			kubernetes.EnsureLabels(dep, baseLabels)

			dep.Spec.Replicas = resources.Int32(2)
			dep.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: baseLabels,
			}

			containers, err := getContainers(data, dep.Spec.Template.Spec.Containers)
			if err != nil {
				return nil, err
			}

			kubernetes.EnsureAnnotations(&dep.Spec.Template, map[string]string{
				resources.ClusterLastRestartAnnotation: data.Cluster().Annotations[resources.ClusterLastRestartAnnotation],
				// these volumes should not block the autoscaler from evicting the pod
				resources.ClusterAutoscalerSafeToEvictVolumesAnnotation: tmpVolumeName + "," + configVolumeName + "," + frontendVolumeName,
			})

			dep.Spec.Template.Spec.Volumes = getVolumes()
			dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: resources.ImagePullSecretName}}
			dep.Spec.Template.Spec.InitContainers = getInitContainers(containers[0].Image)
			dep.Spec.Template.Spec.Containers = containers
			err = resources.SetResourceRequirements(dep.Spec.Template.Spec.Containers, defaultResourceRequirements, nil, dep.Annotations)
			if err != nil {
				return nil, fmt.Errorf("failed to set resource requirements: %w", err)
			}
			dep.Spec.Template.Spec.Affinity = resources.HostnameAntiAffinity(name, kubermaticv1.AntiAffinityTypePreferred)

			dep.Spec.Template, err = apiserver.IsRunningWrapper(data, dep.Spec.Template, sets.New(name))
			if err != nil {
				return nil, fmt.Errorf("failed to add apiserver.IsRunningWrapper: %w", err)
			}

			return dep, nil
		}
	}
}

func getContainers(data headlampData, existingContainers []corev1.Container) ([]corev1.Container, error) {
	securityContext := &corev1.SecurityContext{}
	if len(existingContainers) == 1 && existingContainers[0].SecurityContext != nil {
		securityContext = existingContainers[0].SecurityContext
	}
	securityContext.RunAsUser = ptr.To[int64](100)
	securityContext.RunAsGroup = ptr.To[int64](101)
	securityContext.ReadOnlyRootFilesystem = ptr.To(true)
	securityContext.AllowPrivilegeEscalation = ptr.To(false)

	tag, err := HeadlampVersion(data.Cluster().Status.Versions.ControlPlane)
	if err != nil {
		return nil, err
	}

	return []corev1.Container{{
		Name:            name,
		Image:           registry.Must(data.RewriteImage(fmt.Sprintf("%s/%s:%s", "ghcr.io", imageName, tag))),
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         []string{"/headlamp/headlamp-server"},
		Args: []string{
			"-kubeconfig", "/etc/kubernetes/kubeconfig/kubeconfig",
			"-html-static-dir", "/headlamp/frontend",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      resources.HeadlampKubeconfigSecretName,
				MountPath: "/etc/kubernetes/kubeconfig",
				ReadOnly:  true,
			},
			{
				Name:      tmpVolumeName,
				MountPath: "/tmp",
			},
			{
				Name:      configVolumeName,
				MountPath: "/home/headlamp/.config",
			},
			{
				Name:      frontendVolumeName,
				MountPath: "/headlamp/frontend",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		SecurityContext: securityContext,
	}}, nil
}

func getVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: resources.HeadlampKubeconfigSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: resources.HeadlampKubeconfigSecretName,
				},
			},
		},
		{
			Name: tmpVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: configVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: frontendVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}

// getInitContainers returns an init container that copies the frontend static
// files from the image into the writable emptyDir volume. This is needed because
// the main container runs with readOnlyRootFilesystem and Headlamp writes
// index.baseUrl.html into the frontend directory at startup.
func getInitContainers(image string) []corev1.Container {
	return []corev1.Container{
		{
			Name:            "copy-frontend",
			Image:           image,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Command:         []string{"sh", "-c", "cp -r /headlamp/frontend/* /mnt/frontend/"},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      frontendVolumeName,
					MountPath: "/mnt/frontend",
				},
			},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:                ptr.To[int64](100),
				RunAsGroup:               ptr.To[int64](101),
				ReadOnlyRootFilesystem:   ptr.To(true),
				AllowPrivilegeEscalation: ptr.To(false),
			},
		},
	}
}

// HeadlampVersion returns the Headlamp image tag to use. Headlamp is not
// Kubernetes-version dependent, so a single fixed version is used regardless
// of the cluster's control plane version.
func HeadlampVersion(_ semver.Semver) (string, error) {
	return headlampVersion, nil
}

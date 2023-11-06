/*

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type SpillConfigSpec struct {
	// Enable the spill mechanism
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled,omitempty"`

	// the path of spilling
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=""
	Path string `json:"path,omitempty"`

	// the pvc name
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=""
	PVCName string `json:"pvcName,omitempty"`

	// low watermark of spilling memory
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="0.3"
	SpillLowerRate string `json:"spillLowerRate,omitempty"`

	// high watermark of triggering spilling
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="0.8"
	SpillUpperRate string `json:"spillUpperRate,omitempty"`
}

type MetricConfigSpec struct {
	// Enable metrics
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// represent the metric's image
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="vineyardcloudnative/vineyard-grok-exporter"
	Image string `json:"image,omitempty"`

	// the policy about pulling image
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type VineyarddTemplateSpec struct {
	// represent the vineyardd's image
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="vineyardcloudnative/vineyardd"
	Image string `json:"image,omitempty"`

	// the policy about pulling image
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// the tag of vineyardd's image
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="latest"
	ImageTag string `json:"imageTag,omitempty"`

	// the path of vineyard domain docket
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="/var/run"
	SocketPath string `json:"socketPath,omitempty"`

	// shared memory size for vineyardd
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=""
	Size string `json:"size,omitempty"`

	// reserve the shared memory for vineyardd
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	ReserveMemory bool `json:"reserveMemory,omitempty"`

	// Disable monitoring for Vineyard Runtime
	// Prometheus is enabled by default
	// +optional
	DisablePrometheus bool `json:"disablePrometheus,omitempty"`

	// vineyard environment configuration
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Env []corev1.EnvVar `json:"env,omitempty"`

	// the configuration of spilling
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Spill SpillConfigSpec `json:"spill,omitempty"`

	// the configuration of metric
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Metric MetricConfigSpec `json:"metric,omitempty"`
}

type FuseTemplateSpec struct {
	// the image of vineyard fuse (mount vineyard socket)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="vineyardcloudnative/mount-vineyard-socket"
	Image string `json:"image,omitempty"`

	// the image pull policy of vineyard fuse image (mount vineyard socket)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// the tag of vineyard fuse image (mount vineyard socket)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="latest"
	ImageTag string `json:"imageTag,omitempty"`

	// NodeSelector is a selector which must be true for the fuse client to fit on a node,
	// this option only effect when global is enabled
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// If the fuse client should be deployed in global mode,
	// otherwise the affinity should be considered
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Global bool `json:"global,omitempty"`

	// CleanPolicy decides when to clean Alluxio Fuse pods.
	// Currently Fluid supports two policies: OnDemand and OnRuntimeDeleted
	// OnDemand cleans fuse pod once the fuse pod on some node is not needed
	// OnRuntimeDeleted cleans fuse pod only when the cache runtime is deleted
	// Defaults to OnRuntimeDeleted
	// +optional
	CleanPolicy FuseCleanPolicy `json:"cleanPolicy,omitempty"`

	// livenessProbe of vineyard fuse pod
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// readinessProbe of vineyard fuse pod
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`

	// VolumeMounts specifies the volumes listed in ".spec.volumes" to mount into the vineyardruntime component's filesystem.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
}

// VineyardRuntimeSpec defines the desired state of VineyardRuntime
type VineyardRuntimeSpec struct {
	// Replicas is the number of vineyardd pods to deploy
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=3
	Replicas int32 `json:"replicas,omitempty"`

	// EtcdReplicas describe the etcd replicas
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	EtcdReplicas int32 `json:"etcdReplicas,omitempty"`

	// Vineyardd is the configuration of vineyardd
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Vineyardd VineyarddTemplateSpec `json:"vineyardd,omitempty"`

	// Fuse is the configuration of vineyard fuse
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Fuse FuseTemplateSpec `json:"fuse,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VineyardRuntime is the Schema for the vineyardruntimes API
type VineyardRuntime struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VineyardRuntimeSpec `json:"spec,omitempty"`
	Status RuntimeStatus       `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VineyardRuntimeList contains a list of VineyardRuntime
type VineyardRuntimeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VineyardRuntime `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VineyardRuntime{}, &VineyardRuntimeList{})
}

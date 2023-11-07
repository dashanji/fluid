/*
Copyright 2023 The Fluid Author.

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

package vineyard

import (
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/fluid-cloudnative/fluid/pkg/common"
	"github.com/fluid-cloudnative/fluid/pkg/ddc/base"
)

// The value yaml file
type Vineyard struct {
	FullnameOverride string `json:"fullnameOverride"`

	common.ImageInfo `json:",inline"`
	common.UserInfo  `json:",inline"`

	runtimeInfo     base.RuntimeInfoInterface
	Vineyardd       `json:"vineyardd,omitempty"`
	RuntimeIdentity common.RuntimeIdentity `json:"runtimeIdentity"`

	Owner *common.OwnerReference `json:"owner,omitempty"`

	Fuse `json:"fuse,omitempty"`
}

type HadoopConfig struct {
	ConfigMap       string `json:"configMap"`
	IncludeHdfsSite bool   `json:"includeHdfsSite"`
	IncludeCoreSite bool   `json:"includeCoreSite"`
}

type UFSPath struct {
	HostPath  string `json:"hostPath"`
	UFSVolume `json:",inline"`
}

type UFSVolume struct {
	Name          string `json:"name"`
	SubPath       string `json:"subPath,omitempty"`
	ContainerPath string `json:"containerPath"`
}

type Metastore struct {
	VolumeType string `json:"volumeType,omitempty"`
	Size       string `json:"size,omitempty"`
}

type Journal struct {
	VolumeType string `json:"volumeType,omitempty"`
	Size       string `json:"size,omitempty"`
}

type ShortCircuit struct {
	Enable     bool   `json:"enable,omitempty"`
	Policy     string `json:"policy,omitempty"`
	VolumeType string `json:"volumeType,omitempty"`
}

type Ports struct {
	Rpc      int `json:"rpc,omitempty"`
	Web      int `json:"web,omitempty"`
	Embedded int `json:"embedded,omitempty"`
	Data     int `json:"data,omitempty"`
	Rest     int `json:"rest,omitempty"`
}

type APIGateway struct {
	Enabled bool  `json:"enabled,omitempty"`
	Ports   Ports `json:"ports,omitempty"`
}

type JobMaster struct {
	Ports     Ports            `json:"ports,omitempty"`
	Resources common.Resources `json:"resources,omitempty"`
}

type JobWorker struct {
	Ports     Ports            `json:"ports,omitempty"`
	Resources common.Resources `json:"resources,omitempty"`
}

type Worker struct {
	JvmOptions   []string             `json:"jvmOptions,omitempty"`
	Env          map[string]string    `json:"env,omitempty"`
	NodeSelector map[string]string    `json:"nodeSelector,omitempty"`
	Properties   map[string]string    `json:"properties,omitempty"`
	HostNetwork  bool                 `json:"hostNetwork,omitempty"`
	Resources    common.Resources     `json:"resources,omitempty"`
	Ports        Ports                `json:"ports,omitempty"`
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	Volumes      []corev1.Volume      `json:"volumes,omitempty"`
	Labels       map[string]string    `json:"labels,omitempty"`
	Annotations  map[string]string    `json:"annotations,omitempty"`
}

type Spill struct {
	Enabled        bool   `json:"enabled,omitempty"`
	Path           string `json:"path,omitempty"`
	PvcName        string `json:"pvcName,omitempty"`
	SpillLowerRate string `json:"spillLowerRate,omitempty"`
	SpillUpperRate string `json:"spillUpperRate,omitempty"`
}

type Metric struct {
	Enabled         bool   `json:"enabled,omitempty"`
	Image           string `json:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type Vineyardd struct {
	Replicas    int32             `json:"replicas,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`

	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`

	Image           string `json:"image,omitempty"`
	ImageTag        string `json:"imageTag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	SocketPath string `json:"socketPath,omitempty"`

	Size          string `json:"size,omitempty"`
	ReserveMemory bool   `json:"reserveMemory,omitempty"`

	Spill  `json:"spill,omitempty"`
	Metric `json:"metric,omitempty"`
}

type Restore struct {
	Enabled bool   `json:"enabled,omitempty"`
	Path    string `json:"path,omitempty"`
	PVCName string `json:"pvcName,omitempty"`
}

type Fuse struct {
	Enabled         bool                   `json:"enabled,omitempty"`
	Image           string                 `json:"image,omitempty"`
	ImageTag        string                 `json:"imageTag,omitempty"`
	ImagePullPolicy string                 `json:"imagePullPolicy,omitempty"`
	Resources       common.Resources       `json:"resources,omitempty"`
	Ports           []corev1.ContainerPort `json:"ports,omitempty"`
	CriticalPod     bool                   `json:"criticalPod,omitempty"`
	HostNetwork     bool                   `json:"hostNetwork,omitempty"`
	TargetPath      string                 `json:"targetPath,omitempty"`
	NodeSelector    map[string]string      `json:"nodeSelector,omitempty"`
	Envs            []corev1.EnvVar        `json:"envs,omitempty"`
	Command         []string               `json:"command,omitempty"`
	Args            []string               `json:"args,omitempty"`
	Volumes         []corev1.Volume        `json:"volumes,omitempty"`
	VolumeMounts    []corev1.VolumeMount   `json:"volumeMounts,omitempty"`
	LivenessProbe   *corev1.Probe          `json:"livenessProbe,omitempty"`
	ReadinessProbe  *corev1.Probe          `json:"readinessProbe,omitempty"`
	CacheDir        string                 `json:"cacheDir,omitempty"`
	ConfigValue     string                 `json:"configValue"`
	ConfigStorage   string                 `json:"configStorage"`
}

type TieredStore struct {
	Levels []Level `json:"levels,omitempty"`
}

type Level struct {
	Alias      string `json:"alias,omitempty"`
	Level      int    `json:"level"`
	MediumType string `json:"mediumtype,omitempty"`
	Type       string `json:"type,omitempty"`
	Path       string `json:"path,omitempty"`
	Quota      string `json:"quota,omitempty"`
	High       string `json:"high,omitempty"`
	Low        string `json:"low,omitempty"`
}

type Affinity struct {
	NodeAffinity *NodeAffinity `json:"nodeAffinity"`
}

type cacheHitStates struct {
	cacheHitRatio  string
	localHitRatio  string
	remoteHitRatio string

	localThroughputRatio  string
	remoteThroughputRatio string
	cacheThroughputRatio  string

	bytesReadLocal  int64
	bytesReadRemote int64
	bytesReadUfsAll int64

	timestamp time.Time
}

type cacheStates struct {
	cacheCapacity string
	// cacheable        string
	// lowWaterMark     string
	// highWaterMark    string
	cached           string
	cachedPercentage string
	cacheHitStates   cacheHitStates
	// nonCacheable     string
}

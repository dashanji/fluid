/*
Copyright 2023 The Fluid Authors.

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
	"fmt"
	"time"

	datav1alpha1 "github.com/fluid-cloudnative/fluid/api/v1alpha1"
	"github.com/fluid-cloudnative/fluid/pkg/common"
	"github.com/fluid-cloudnative/fluid/pkg/utils"
	"github.com/fluid-cloudnative/fluid/pkg/utils/transfromer"
)

func (e *VineyardEngine) transform(runtime *datav1alpha1.VineyardRuntime) (value *Vineyard, err error) {
	if runtime == nil {
		err = fmt.Errorf("the vineyardRuntime is null")
		return
	}
	defer utils.TimeTrack(time.Now(), "VineyardRuntime.Transform", "name", runtime.Name)

	fmt.Println("transform vineyard runtime", *runtime)
	_, err = utils.GetDataset(e.Client, e.name, e.namespace)
	if err != nil {
		return value, err
	}

	spill := Spill{
		Enabled:        runtime.Spec.Vineyardd.Spill.Enabled,
		Path:           runtime.Spec.Vineyardd.Spill.Path,
		PvcName:        runtime.Spec.Vineyardd.Spill.PVCName,
		SpillLowerRate: runtime.Spec.Vineyardd.Spill.SpillLowerRate,
		SpillUpperRate: runtime.Spec.Vineyardd.Spill.SpillUpperRate,
	}

	metric := Metric{
		Enabled:         runtime.Spec.Vineyardd.Metric.Enabled,
		Image:           runtime.Spec.Vineyardd.Metric.Image,
		ImagePullPolicy: runtime.Spec.Vineyardd.Metric.ImagePullPolicy,
	}

	vineyardd := Vineyardd{
		Replicas:      runtime.Spec.Replicas,
		Image:         runtime.Spec.Vineyardd.Image,
		ImageTag:      runtime.Spec.Vineyardd.ImageTag,
		Size:          runtime.Spec.Vineyardd.Size,
		ReserveMemory: runtime.Spec.Vineyardd.ReserveMemory,
		SocketPath:    runtime.Spec.Vineyardd.SocketPath,
	}
	vineyardd.Spill = spill
	vineyardd.Metric = metric
	value = &Vineyard{
		RuntimeIdentity: common.RuntimeIdentity{
			Namespace: runtime.Namespace,
			Name:      runtime.Name,
		},
		Owner: transfromer.GenerateOwnerReferenceFromObject(runtime),
	}
	value.FullnameOverride = e.name

	value.Vineyardd = vineyardd

	etcdService := &EtcdService{
		Type:     runtime.Spec.Etcd.Service.Type,
		Protocol: runtime.Spec.Etcd.Service.Protocol,
		Ports: EtcdPorts{
			Client: int(runtime.Spec.Etcd.Service.ClientPort),
			Peer:   int(runtime.Spec.Etcd.Service.PeerPort),
		},
	}
	fmt.Println("service", *etcdService)

	auth := &EtcdAuth{
		Client: Transport{
			SecureTransport: runtime.Spec.Etcd.EnableSecureTransport,
		},
		Peer: Transport{
			SecureTransport: runtime.Spec.Etcd.EnableSecureTransport,
		},
	}
	fmt.Println("auth: ", *auth)
	persistent := &EtcdPersistent{
		Enabled:          runtime.Spec.Etcd.Persistent.Enabled,
		Annotaions:       runtime.Spec.Etcd.Persistent.Annotaions,
		Labels:           runtime.Spec.Etcd.Persistent.Labels,
		StorageClassName: runtime.Spec.Etcd.Persistent.StorageClassName,
		AccessModes:      runtime.Spec.Etcd.Persistent.AccessModes,
		Size:             runtime.Spec.Etcd.Persistent.Size,
		Selector:         runtime.Spec.Etcd.Persistent.Selector,
	}
	fmt.Println("persistent: ", *persistent)
	etcd := Etcd{
		Replicas: runtime.Spec.Etcd.Replicas,
		Persistent: EtcdPersistent{
			Enabled:          runtime.Spec.Etcd.Persistent.Enabled,
			Annotaions:       runtime.Spec.Etcd.Persistent.Annotaions,
			Labels:           runtime.Spec.Etcd.Persistent.Labels,
			StorageClassName: runtime.Spec.Etcd.Persistent.StorageClassName,
			AccessModes:      runtime.Spec.Etcd.Persistent.AccessModes,
			Size:             runtime.Spec.Etcd.Persistent.Size,
			Selector:         runtime.Spec.Etcd.Persistent.Selector,
		},
		Service: EtcdService{
			Type:     runtime.Spec.Etcd.Service.Type,
			Protocol: runtime.Spec.Etcd.Service.Protocol,
			Ports: EtcdPorts{
				Client: int(runtime.Spec.Etcd.Service.ClientPort),
				Peer:   int(runtime.Spec.Etcd.Service.PeerPort),
			},
		},
		Auth: EtcdAuth{
			Client: Transport{
				SecureTransport: runtime.Spec.Etcd.EnableSecureTransport,
			},
			Peer: Transport{
				SecureTransport: runtime.Spec.Etcd.EnableSecureTransport,
			},
		},
	}
	e.Log.Info("etcd value is", "value", etcd)
	fmt.Println("etcd: ", etcd)
	value.Etcd = etcd
	e.Log.Info("before value is", "value", value)
	fmt.Println("before value is: ", value)
	fuse := Fuse{
		Enabled:         true,
		Image:           runtime.Spec.Fuse.Image,
		ImageTag:        runtime.Spec.Fuse.ImageTag,
		ImagePullPolicy: runtime.Spec.Fuse.ImagePullPolicy,
	}
	value.Fuse = fuse
	value.Fuse.CriticalPod = common.CriticalFusePodEnabled()
	value.Fuse.TargetPath = e.getTargetPath()
	value.Fuse.LivenessProbe = runtime.Spec.Fuse.LivenessProbe
	value.Fuse.ReadinessProbe = runtime.Spec.Fuse.ReadinessProbe

	value.Fuse.NodeSelector = map[string]string{}
	if len(runtime.Spec.Fuse.NodeSelector) > 0 {
		value.Fuse.NodeSelector = runtime.Spec.Fuse.NodeSelector
	}
	value.Fuse.NodeSelector[e.getFuseLabelName()] = "true"

	value.ClusterDomain = runtime.Spec.ClusterDomain
	return value, nil
}

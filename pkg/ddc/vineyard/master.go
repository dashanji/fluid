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
	"context"
	"reflect"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"

	"github.com/fluid-cloudnative/fluid/pkg/utils/kubeclient"
)

// CheckMasterReady checks if the master is ready
func (e *VineyardEngine) CheckMasterReady() (ready bool, err error) {
	var (
		vineyarddName string = e.getVineyarddName()
		namespace     string = e.namespace
	)

	vineyardStatefulset, err := kubeclient.GetStatefulSet(e.Client, vineyarddName, namespace)
	if err != nil {
		return ready, err
	}

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if vineyardStatefulset.Status.ReadyReplicas == *vineyardStatefulset.Spec.Replicas {
			ready = true
			return nil
		}
		return err
	})

	return
}

// ShouldSetupMaster checks if we need setup the master
func (e *VineyardEngine) ShouldSetupMaster() (should bool, err error) {
	return true, nil
}

// SetupMaster setups the master and updates the status
// It will print the information in the Debug window according to the Master status
// It return any cache error encountered
func (e *VineyardEngine) SetupMaster() (err error) {
	vineyarddName := e.getVineyarddName()

	// 1. Setup the vineyardd
	master, err := kubeclient.GetStatefulSet(e.Client, vineyarddName, e.namespace)
	if err != nil && apierrs.IsNotFound(err) {
		//1. Is not found error
		e.Log.V(1).Info("SetupMaster", "vineyardd", vineyarddName)
		return e.setupMasterInternal()
	} else if err != nil {
		//2. Other errors
		return
	} else {
		//3.The master has been set up
		e.Log.V(1).Info("The master has been set.", "replicas", master.Status.ReadyReplicas)
	}

	// 2. Update the status of the runtime

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		runtime, err := e.getRuntime()
		if err != nil {
			return err
		}
		runtimeToUpdate := runtime.DeepCopy()

		/*
			runtimeToUpdate.Status.MasterPhase = datav1alpha1.RuntimePhaseNotReady
			replicas := runtimeToUpdate.Spec.Replicas
			if replicas == 0 {
				replicas = 1
			}

			// Init selector for worker
			runtimeToUpdate.Status.Selector = e.getWorkerSelectors()

			runtimeToUpdate.Status.DesiredMasterNumberScheduled = replicas
			runtimeToUpdate.Status.ValueFileConfigmap = e.getConfigmapName()

			if len(runtimeToUpdate.Status.Conditions) == 0 {
				runtimeToUpdate.Status.Conditions = []datav1alpha1.RuntimeCondition{}
			}
			cond := utils.NewRuntimeCondition(datav1alpha1.RuntimeMasterInitialized, datav1alpha1.RuntimeMasterInitializedReason,
				"The master is initialized.", corev1.ConditionTrue)
			runtimeToUpdate.Status.Conditions =
				utils.UpdateRuntimeCondition(runtimeToUpdate.Status.Conditions,
					cond)
		*/
		if !reflect.DeepEqual(runtime.Status, runtimeToUpdate.Status) {
			return e.Client.Status().Update(context.TODO(), runtimeToUpdate)
		}

		return nil
	})

	if err != nil {
		e.Log.Error(err, "Update runtime status")
		return err
	}

	return nil
}

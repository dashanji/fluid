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

import "github.com/fluid-cloudnative/fluid/pkg/utils/helm"

// shut down the Alluxio engine
func (e *VineyardEngine) Shutdown() (err error) {
	err = e.destroyMaster()
	if err != nil {
		return
	}
	return
}

// destroyMaster Destroies the master
func (e *VineyardEngine) destroyMaster() (err error) {
	var found bool
	found, err = helm.CheckRelease(e.name, e.namespace)
	if err != nil {
		return err
	}

	if found {
		err = helm.DeleteRelease(e.name, e.namespace)
		if err != nil {
			return
		}
	}
	return
}

/*
Copyright 2022 The Kubernetes Authors.

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

package cm

import (
    "fmt"
    "strings"

    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

// GetNodeAllocatableAbsolute returns the absolute value of Node Allocatable which is primarily useful for enforcement.
// Note that not all resources that are available on the node are included in the returned list of resources.
// Returns a ResourceList.
func (cm *containerManagerImpl) GetNodeAllocatableAbsolute() v1.ResourceList {
    return cm.getNodeAllocatableAbsoluteImpl(cm.capacity)
}

func (cm *containerManagerImpl) getNodeAllocatableAbsoluteImpl(capacity v1.ResourceList) v1.ResourceList {
    result := make(v1.ResourceList)
    for k, v := range capacity {
        value := v.DeepCopy()
        if cm.NodeConfig.SystemReserved != nil {
            value.Sub(cm.NodeConfig.SystemReserved[k])
        }
        if cm.NodeConfig.KubeReserved != nil {
            value.Sub(cm.NodeConfig.KubeReserved[k])
        }
        if value.Sign() < 0 {
            // Negative Allocatable resources don't make sense.
            value.Set(0)
        }
        result[k] = value
    }
    return result
}

// getNodeAllocatableInternalAbsolute is similar to getNodeAllocatableAbsolute except that
// it also includes internal resources (currently process IDs).  It is intended for setting
// up top level cgroups only.
func (cm *containerManagerImpl) getNodeAllocatableInternalAbsolute() v1.ResourceList {
    return cm.getNodeAllocatableAbsoluteImpl(cm.internalCapacity)
}

// GetNodeAllocatableReservation returns amount of compute or storage resource that have to be reserved on this node from scheduling.
func (cm *containerManagerImpl) GetNodeAllocatableReservation() v1.ResourceList {
    evictionReservation := hardEvictionReservation(cm.HardEvictionThresholds, cm.capacity)
    result := make(v1.ResourceList)
    for k := range cm.capacity {
        value := resource.NewQuantity(0, resource.DecimalSI)
        if cm.NodeConfig.SystemReserved != nil {
            value.Add(cm.NodeConfig.SystemReserved[k])
        }
        if cm.NodeConfig.KubeReserved != nil {
            value.Add(cm.NodeConfig.KubeReserved[k])
        }
        if evictionReservation != nil {
            value.Add(evictionReservation[k])
        }
        if !value.IsZero() {
            result[k] = *value
        }
    }
    return result
}

// validateNodeAllocatable ensures that the user specified Node Allocatable Configuration doesn't reserve more than the node capacity.
// Returns error if the configuration is invalid, nil otherwise.
func (cm *containerManagerImpl) validateNodeAllocatable() error {
    var errors []string
    nar := cm.GetNodeAllocatableReservation()
    for k, v := range nar {
        value := cm.capacity[k].DeepCopy()
        value.Sub(v)

        if value.Sign() < 0 {
            errors = append(errors, fmt.Sprintf("Resource %q has an allocatable of %v, capacity of %v", k, v, value))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("invalid Node Allocatable configuration. %s", strings.Join(errors, " "))
    }
    return nil
}

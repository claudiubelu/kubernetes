/*
Copyright 2021 The Kubernetes Authors.

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

package daemonset

import (
    "context"
    "fmt"
    "sort"

    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDaemonSet(dsName, image string, label map[string]string) *appsv1.DaemonSet {
    return &appsv1.DaemonSet{
        ObjectMeta: metav1.ObjectMeta{
            Name: dsName,
        },
        Spec: appsv1.DaemonSetSpec{
            Selector: &metav1.LabelSelector{
                MatchLabels: label,
            },
            Template: v1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: label,
                },
                Spec: v1.PodSpec{
                    Containers: []v1.Container{
                        {
                            Name:  "app",
                            Image: image,
                        },
                    },
                },
            },
        },
    }
}

func CheckRunningOnAllNodes(f *framework.Framework, ds *appsv1.DaemonSet) func() (bool, error) {
    return func() (bool, error) {
        nodeNames := schedulableNodes(f.ClientSet, ds)
        return checkDaemonPodOnNodes(f, ds, nodeNames)()
    }
}

func schedulableNodes(c clientset.Interface, ds *appsv1.DaemonSet) []string {
    nodeList, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    framework.ExpectNoError(err)
    nodeNames := make([]string, 0)
    for _, node := range nodeList.Items {
        if !canScheduleOnNode(node, ds) {
            framework.Logf("DaemonSet pods can't tolerate node %s with taints %+v, skip checking this node", node.Name, node.Spec.Taints)
            continue
        }
        nodeNames = append(nodeNames, node.Name)
    }
    return nodeNames
}

func checkDaemonPodOnNodes(f *framework.Framework, ds *appsv1.DaemonSet, nodeNames []string) func() (bool, error) {
    return func() (bool, error) {
        podList, err := f.ClientSet.CoreV1().Pods(f.Namespace.Name).List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            framework.Logf("could not get the pod list: %v", err)
            return false, nil
        }
        pods := podList.Items

        nodesToPodCount := make(map[string]int)
        for _, pod := range pods {
            if !metav1.IsControlledBy(&pod, ds) {
                continue
            }
            if pod.DeletionTimestamp != nil {
                continue
            }
            if podutil.IsPodAvailable(&pod, ds.Spec.MinReadySeconds, metav1.Now()) {
                nodesToPodCount[pod.Spec.NodeName]++
            }
        }
        framework.Logf("Number of nodes with available pods: %d", len(nodesToPodCount))

        // Ensure that exactly 1 pod is running on all nodes in nodeNames.
        for _, nodeName := range nodeNames {
            if nodesToPodCount[nodeName] != 1 {
                framework.Logf("Node %s is running more than one daemon pod", nodeName)
                return false, nil
            }
        }

        framework.Logf("Number of running nodes: %d, number of available pods: %d", len(nodeNames), len(nodesToPodCount))
        // Ensure that sizes of the lists are the same. We've verified that every element of nodeNames is in
        // nodesToPodCount, so verifying the lengths are equal ensures that there aren't pods running on any
        // other nodes.
        return len(nodesToPodCount) == len(nodeNames), nil
    }
}

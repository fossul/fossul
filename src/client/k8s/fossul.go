/*
Copyright 2019 The Fossul Authors.
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

package k8s

import (
	"context"

	"encoding/json"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
)

func CreateBackupCustomResource(accessWithinCluster, namespace, crName, profileName, configName, policyName string) error {
	client, err := getDynamicClient(accessWithinCluster)
	if err != nil {
		return err
	}

	// need to figure out how to get current namespace dynamically
	fossulNamespace := os.Getenv("FOSSUL_NAMESPACE")

	var backupCustomResourceGroup = schema.GroupVersionResource{Group: "fossul.io", Version: "v1", Resource: "backups"}

	backupCustomResource := &unstructured.Unstructured{}
	backupCustomResource.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "fossul.io/v1",
		"kind":       "Backup",
		"metadata": map[string]interface{}{
			"name":      crName,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"deployment_name":  configName,
			"policy":           policyName,
			"fossul_namespace": fossulNamespace,
		},
	})

	_, err = client.Resource(backupCustomResourceGroup).Namespace(namespace).Create(context.Background(), backupCustomResource, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func DeleteBackupCustomResource(accessWithinCluster, namespace, crName string) error {
	client, err := getDynamicClient(accessWithinCluster)
	if err != nil {
		return err
	}

	var backupCustomResourceGroup = schema.GroupVersionResource{Group: "fossul.io", Version: "v1", Resource: "backups"}

	err = client.Resource(backupCustomResourceGroup).Namespace(namespace).Delete(context.Background(), crName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func UpdateBackupCustomResource(accessWithinCluster, namespace, crName, op, specKey, specValue string) error {
	client, err := getDynamicClient(accessWithinCluster)
	if err != nil {
		return err
	}

	patch := []interface{}{
		map[string]interface{}{
			"op":    op,
			"path":  "/spec/" + specKey,
			"value": specValue,
		},
	}

	payload, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	var backupCustomResourceGroup = schema.GroupVersionResource{Group: "fossul.io", Version: "v1", Resource: "backups"}

	_, err = client.Resource(backupCustomResourceGroup).Namespace(namespace).Patch(context.Background(), crName, types.JSONPatchType, payload, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetBackupCustomResource(accessWithinCluster, namespace, crName string) (*unstructured.Unstructured, error) {
	var backupCustomResource *unstructured.Unstructured
	client, err := getDynamicClient(accessWithinCluster)
	if err != nil {
		return backupCustomResource, err
	}

	var backupCustomResourceGroup = schema.GroupVersionResource{Group: "fossul.io", Version: "v1", Resource: "backups"}

	backupCustomResource, err = client.Resource(backupCustomResourceGroup).Namespace(namespace).Get(context.Background(), crName, metav1.GetOptions{})
	if err != nil {
		return backupCustomResource, err
	}

	return backupCustomResource, nil
}

func ListBackupCustomResources(accessWithinCluster, namespace string) (*unstructured.UnstructuredList, error) {
	var backupCustomResourceList *unstructured.UnstructuredList
	client, err := getDynamicClient(accessWithinCluster)
	if err != nil {
		return backupCustomResourceList, err
	}

	var backupCustomResourceGroup = schema.GroupVersionResource{Group: "fossul.io", Version: "v1", Resource: "backups"}

	backupCustomResourceList, err = client.Resource(backupCustomResourceGroup).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return backupCustomResourceList, err
	}

	return backupCustomResourceList, nil
}

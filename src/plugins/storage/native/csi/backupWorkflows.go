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
package main

import (
	"fmt"

	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
)

func DeploymentConfigBackupWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	timeout := util.StringToInt(config.StoragePluginParameters["SnapshotTimeoutSeconds"])

	var backup util.Backup
	backup.Timestamp = timestampToString
	backup.Epoch = int(config.WorkflowTimestamp)
	backup.Policy = config.SelectedBackupPolicy
	backup.WorkflowId = config.WorkflowId

	var contents []util.Content
	var content util.Content

	deploymentConfig, err := k8s.GetDeploymentConfig(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	volumes := deploymentConfig.Spec.Template.Spec.Volumes
	for _, volume := range volumes {
		if volume.PersistentVolumeClaim == nil {
			msg := util.SetMessage("ERROR", "No persistent volume claim found for deployment ["+config.StoragePluginParameters["DeploymentName"]+"]")
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		pvcName := volume.PersistentVolumeClaim.ClaimName
		backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
		backup.Name = backupName
		snapshotName := backupName + "-" + pvcName

		pvc, err := k8s.GetPersistentVolumeClaim(config.StoragePluginParameters["Namespace"], pvcName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
		pvcSize := pvc.Status.Capacity.Storage().String()
		storageClassName := pvc.Spec.StorageClassName

		content.Size = pvcSize
		content.StorageClass = *storageClassName

		provisionerName, err := k8s.GetStorageClassProvisionerName(content.StorageClass, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		snapshotClassName, err := k8s.GetVolumeSnapshotClassName(provisionerName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.VolumeSnapshotClass = snapshotClassName

		msg := util.SetMessage("INFO", "Creating CSI snapshot ["+snapshotName+"] of pvc ["+pvcName+"] namespace ["+config.StoragePluginParameters["Namespace"]+"] snapshot class ["+snapshotClassName+" ] timeout ["+config.StoragePluginParameters["SnapshotTimeoutSeconds"]+"]")
		messages = append(messages, msg)

		err = k8s.CreateSnapshot(snapshotName, config.StoragePluginParameters["Namespace"], snapshotClassName, pvcName, config.AccessWithinCluster, timeout)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.Type = "volume"
		content.Source = pvcName
		content.Data = snapshotName
		contents = append(contents, content)

		msg = util.SetMessage("INFO", "CSI snapshot created successfully")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)

	backup.Contents = contents
	result.Backup = backup

	return result
}

func DeploymentBackupWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	timeout := util.StringToInt(config.StoragePluginParameters["SnapshotTimeoutSeconds"])

	var backup util.Backup
	backup.Timestamp = timestampToString
	backup.Epoch = int(config.WorkflowTimestamp)
	backup.Policy = config.SelectedBackupPolicy
	backup.WorkflowId = config.WorkflowId

	var contents []util.Content
	var content util.Content

	deployment, err := k8s.GetDeployment(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	volumes := deployment.Spec.Template.Spec.Volumes
	for _, volume := range volumes {
		pvcName := volume.PersistentVolumeClaim.ClaimName
		backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
		backup.Name = backupName
		snapshotName := backupName + "-" + pvcName

		pvc, err := k8s.GetPersistentVolumeClaim(config.StoragePluginParameters["Namespace"], pvcName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
		pvcSize := pvc.Status.Capacity.Storage().String()
		storageClassName := pvc.Spec.StorageClassName

		content.Size = pvcSize
		content.StorageClass = *storageClassName

		provisionerName, err := k8s.GetStorageClassProvisionerName(content.StorageClass, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		snapshotClassName, err := k8s.GetVolumeSnapshotClassName(provisionerName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.VolumeSnapshotClass = snapshotClassName

		msg := util.SetMessage("INFO", "Creating CSI snapshot ["+snapshotName+"] of pvc ["+pvcName+"] namespace ["+config.StoragePluginParameters["Namespace"]+"] snapshot class ["+snapshotClassName+" ] timeout ["+config.StoragePluginParameters["SnapshotTimeoutSeconds"]+"]")
		messages = append(messages, msg)

		err = k8s.CreateSnapshot(snapshotName, config.StoragePluginParameters["Namespace"], snapshotClassName, pvcName, config.AccessWithinCluster, timeout)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.Type = "volume"
		content.Source = pvcName
		content.Data = snapshotName
		contents = append(contents, content)

		msg = util.SetMessage("INFO", "CSI snapshot created successfully")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)

	backup.Contents = contents
	result.Backup = backup

	return result
}

func StatefulSetBackupWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	timeout := util.StringToInt(config.StoragePluginParameters["SnapshotTimeoutSeconds"])

	var backup util.Backup
	backup.Timestamp = timestampToString
	backup.Epoch = int(config.WorkflowTimestamp)
	backup.Policy = config.SelectedBackupPolicy
	backup.WorkflowId = config.WorkflowId

	var contents []util.Content
	var content util.Content

	statefulset, err := k8s.GetStatefulSet(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	volumes := statefulset.Spec.VolumeClaimTemplates
	volumeCount := 0
	for _, volume := range volumes {

		// unfortunately the pvs name is not in the statefulset you need to put a few pieces of info together, maybe there is better way
		instance := volume.ObjectMeta.Labels["app.kubernetes.io/instance"]
		dbName := volume.ObjectMeta.Labels["app.kubernetes.io/name"]
		volumeName := volume.ObjectMeta.Name
		pvcName := volumeName + "-" + instance + "-" + dbName + "-" + util.IntToString(volumeCount)
		volumeCount++
		backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
		backup.Name = backupName
		snapshotName := backupName + "-" + pvcName

		pvc, err := k8s.GetPersistentVolumeClaim(config.StoragePluginParameters["Namespace"], pvcName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
		pvcSize := pvc.Status.Capacity.Storage().String()
		storageClassName := pvc.Spec.StorageClassName

		content.Size = pvcSize
		content.StorageClass = *storageClassName

		provisionerName, err := k8s.GetStorageClassProvisionerName(content.StorageClass, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		snapshotClassName, err := k8s.GetVolumeSnapshotClassName(provisionerName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.VolumeSnapshotClass = snapshotClassName

		msg := util.SetMessage("INFO", "Creating CSI snapshot ["+snapshotName+"] of pvc ["+pvcName+"] namespace ["+config.StoragePluginParameters["Namespace"]+"] snapshot class ["+snapshotClassName+" ] timeout ["+config.StoragePluginParameters["SnapshotTimeoutSeconds"]+"]")
		messages = append(messages, msg)

		err = k8s.CreateSnapshot(snapshotName, config.StoragePluginParameters["Namespace"], snapshotClassName, pvcName, config.AccessWithinCluster, timeout)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.Type = "volume"
		content.Source = pvcName
		content.Data = snapshotName
		contents = append(contents, content)

		msg = util.SetMessage("INFO", "CSI snapshot created successfully")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)

	backup.Contents = contents
	result.Backup = backup

	return result
}

func VirtualMachineBackupWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	timeout := util.StringToInt(config.StoragePluginParameters["SnapshotTimeoutSeconds"])

	var backup util.Backup
	backup.Timestamp = timestampToString
	backup.Epoch = int(config.WorkflowTimestamp)
	backup.Policy = config.SelectedBackupPolicy
	backup.WorkflowId = config.WorkflowId

	var contents []util.Content
	var content util.Content

	vm, err := k8s.GetVirtualMachine(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, config.StoragePluginParameters["DeploymentName"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	volumes := vm.Spec.Template.Spec.Volumes
	for _, volume := range volumes {
		if volume.Name == "cloudinitdisk" {
			continue
		}

		var pvcName string
		if volume.PersistentVolumeClaim == nil {
			pvcName = volume.DataVolume.Name
		} else {
			pvcName = volume.PersistentVolumeClaim.ClaimName
		}

		backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
		backup.Name = backupName
		snapshotName := backupName + "-" + pvcName

		pvc, err := k8s.GetPersistentVolumeClaim(config.StoragePluginParameters["Namespace"], pvcName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
		pvcSize := pvc.Status.Capacity.Storage().String()
		storageClassName := pvc.Spec.StorageClassName

		content.Size = pvcSize
		content.StorageClass = *storageClassName

		provisionerName, err := k8s.GetStorageClassProvisionerName(content.StorageClass, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		snapshotClassName, err := k8s.GetVolumeSnapshotClassName(provisionerName, config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.VolumeSnapshotClass = snapshotClassName

		msg := util.SetMessage("INFO", "Creating CSI snapshot ["+snapshotName+"] of pvc ["+pvcName+"] namespace ["+config.StoragePluginParameters["Namespace"]+"] snapshot class ["+snapshotClassName+" ] timeout ["+config.StoragePluginParameters["SnapshotTimeoutSeconds"]+"]")
		messages = append(messages, msg)

		err = k8s.CreateSnapshot(snapshotName, config.StoragePluginParameters["Namespace"], snapshotClassName, pvcName, config.AccessWithinCluster, timeout)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		content.Type = "volume"
		content.Source = pvcName
		content.Metadata = volume.Name
		content.Data = snapshotName
		contents = append(contents, content)

		msg = util.SetMessage("INFO", "CSI snapshot created successfully")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)

	backup.Contents = contents
	result.Backup = backup

	return result
}

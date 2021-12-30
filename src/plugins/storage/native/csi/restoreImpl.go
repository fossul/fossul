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
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"

	//"github.com/fossul/fossul/src/plugins/pluginUtil"

	"time"
)

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Performing CSI snapshot restore")
	messages = append(messages, msg)

	msg = util.SetMessage("INFO", "Scaling down deployment ["+config.StoragePluginParameters["DeploymentName"]+"]")
	messages = append(messages, msg)

	var deploymentReplicasInt int32
	var scaleDownInt int32 = 0
	var err error
	if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
		deploymentReplicasInt, err = k8s.GetDeploymentConfigScaleInteger(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		err = k8s.ScaleDownDeploymentConfig(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster, scaleDownInt, 120)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
		deploymentReplicasIntRef, err := k8s.GetDeploymentScaleInteger(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
		deploymentReplicasInt = *deploymentReplicasIntRef
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		err = k8s.ScaleDownDeployment(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster, scaleDownInt, 120)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
		msg = util.SetMessage("INFO", "Stopping virtual machine ["+config.StoragePluginParameters["DeploymentName"]+"]")
		messages = append(messages, msg)

		err := k8s.StopVirtualMachine(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, config.StoragePluginParameters["DeploymentName"])
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

	} else {
		msg := util.SetMessage("ERROR", "Couldn't find Deployment or DeploymentConfig, check configuration")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	pvcList, err := k8s.ListPersistentVolumeClaims(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	for _, backedUpPvc := range config.Backup.Contents {
		var pvcRestoreName string
		if config.StoragePluginParameters["RestoreToNewPvc"] == "true" {
			pvcRestoreName = backedUpPvc.Source + "-restore"
		} else {
			pvcRestoreName = backedUpPvc.Source
		}

		var existsRestorePvc bool = false
		for _, existingPvc := range pvcList.Items {
			if existingPvc.Name == pvcRestoreName {
				existsRestorePvc = true
			}
		}

		var pvcDeleteTimeoutSeconds int
		pvcDeleteTimeoutSeconds = util.StringToInt(config.StoragePluginParameters["PvcDeletionTimeout"])

		if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {

			msg = util.SetMessage("INFO", "Updating volume source for virtual machine ["+config.StoragePluginParameters["DeploymentName"]+"]")
			messages = append(messages, msg)

			err := k8s.UpdateVirtualMachineVolumeSource(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, config.StoragePluginParameters["DeploymentName"], pvcRestoreName)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				return result
			}

			if config.StoragePluginParameters["OverwritePvcOnRestore"] == "true" && existsRestorePvc {
				msg := util.SetMessage("INFO", "Deleting existing pvc restore ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"]")
				messages = append(messages, msg)
				err = k8s.DeletePersistentVolumeClaim(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, pvcDeleteTimeoutSeconds)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}
			}
		}

		if !existsRestorePvc {
			msg := util.SetMessage("INFO", "The pvc ["+backedUpPvc.Source+"] from backup does not exist in namespace ["+config.StoragePluginParameters["Namespace"]+"]")
			messages = append(messages, msg)

			msg = util.SetMessage("INFO", "Restoring snapshot ["+backedUpPvc.Data+"] to pvc ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"] using storage class ["+backedUpPvc.StorageClass+"]")
			messages = append(messages, msg)

			err = k8s.CreatePersistentVolumeClaimFromSnapshot(pvcRestoreName, backedUpPvc.Size, backedUpPvc.Data, config.StoragePluginParameters["Namespace"], backedUpPvc.StorageClass, config.AccessWithinCluster)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				return result
			}

			msg = util.SetMessage("INFO", "Updating deployment type ["+config.StoragePluginParameters["DeploymentType"]+"] deployment name  ["+config.StoragePluginParameters["DeploymentName"]+"] to use restore pvc ["+pvcRestoreName+"]")
			messages = append(messages, msg)

			if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
				err := k8s.UpdateDeploymentConfigVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}

				time.Sleep(5 * time.Second)
			} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
				err := k8s.UpdateDeploymentVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}

				time.Sleep(5 * time.Second)
			} else if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
				time.Sleep(5 * time.Second)
			}
		} else {
			existingPvc, err := k8s.GetPersistentVolumeClaim(config.StoragePluginParameters["Namespace"], pvcRestoreName, config.AccessWithinCluster)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				return result
			}

			accessModes := existingPvc.Spec.AccessModes
			volumeMode := existingPvc.Spec.VolumeMode

			if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
				msg = util.SetMessage("INFO", "Restoring snapshot ["+backedUpPvc.Data+"] to new pvc ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"] using storage class ["+backedUpPvc.StorageClass+"]")
				messages = append(messages, msg)

				err = k8s.CreatePersistentVolumeClaimFromSnapshotWithModes(pvcRestoreName, backedUpPvc.Size, backedUpPvc.Data, config.StoragePluginParameters["Namespace"], backedUpPvc.StorageClass, config.AccessWithinCluster, accessModes, volumeMode)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}
			}

			if config.StoragePluginParameters["OverwritePvcOnRestore"] == "true" && existsRestorePvc {
				msg := util.SetMessage("INFO", "Deleting existing pvc ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"]")
				messages = append(messages, msg)
				err = k8s.DeletePersistentVolumeClaim(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, pvcDeleteTimeoutSeconds)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}
			}

			if config.StoragePluginParameters["DeploymentType"] == "Deployment" || config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
				msg = util.SetMessage("INFO", "Restoring snapshot ["+backedUpPvc.Data+"] to new pvc ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"] using storage class ["+backedUpPvc.StorageClass+"]")
				messages = append(messages, msg)

				err = k8s.CreatePersistentVolumeClaimFromSnapshot(pvcRestoreName, backedUpPvc.Size, backedUpPvc.Data, config.StoragePluginParameters["Namespace"], backedUpPvc.StorageClass, config.AccessWithinCluster)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}
			}

			msg = util.SetMessage("INFO", "Updating deployment type ["+config.StoragePluginParameters["DeploymentType"]+"] deployment name  ["+config.StoragePluginParameters["DeploymentName"]+"] to use restore pvc ["+pvcRestoreName+"]")
			messages = append(messages, msg)

			if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
				err := k8s.UpdateDeploymentConfigVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}

				time.Sleep(5 * time.Second)
			} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
				err := k8s.UpdateDeploymentVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
				if err != nil {
					msg := util.SetMessage("ERROR", err.Error())
					messages = append(messages, msg)

					result = util.SetResult(1, messages)
					return result
				}

				time.Sleep(5 * time.Second)
			} else if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
				time.Sleep(5 * time.Second)
			}
		}
	}

	if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
		msg = util.SetMessage("INFO", "Scaling up deployment config ["+config.StoragePluginParameters["DeploymentName"]+"]")
		messages = append(messages, msg)

		err = k8s.ScaleUpDeploymentConfig(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster, deploymentReplicasInt, 120)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

	} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
		msg = util.SetMessage("INFO", "Scaling up deployment ["+config.StoragePluginParameters["DeploymentName"]+"]")
		messages = append(messages, msg)

		err = k8s.ScaleUpDeployment(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster, deploymentReplicasInt, 120)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
		msg = util.SetMessage("INFO", "Starting virtual machine ["+config.StoragePluginParameters["DeploymentName"]+"]")
		messages = append(messages, msg)

		err = k8s.StartVirtualMachine(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, config.StoragePluginParameters["DeploymentName"])
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else {
		msg := util.SetMessage("ERROR", "Couldn't find Deployment or DeploymentConfig, check configuration")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Deployment ["+config.StoragePluginParameters["DeploymentName"]+"] scaled up")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

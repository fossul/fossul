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
	"fossul/src/engine/plugins/pluginUtil"
	"fossul/src/engine/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"path/filepath"
	"strings"
)

type archivePlugin string

var ArchivePlugin archivePlugin

func (r archivePlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (r archivePlugin) Archive(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	//timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	//backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
	backupPath := util.GetBackupPathFromConfig(config)
	bucketPrefix := config.ProfileName + "/" + config.ConfigName + "/"

	msg := util.SetMessage("INFO", "Archiving backup ["+backupPath+"] to AWS bucket ["+config.ArchivePluginParameters["BucketName"]+"]")
	messages = append(messages, msg)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.ArchivePluginParameters["AwsRegion"]),
	}))

	s3svc := s3.New(sess)

	bucketExists, err := BucketExists(s3svc, config.ArchivePluginParameters["BucketName"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	if !bucketExists {
		err := CreateBucket(s3svc, config.ArchivePluginParameters["BucketName"])
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)
			result = util.SetResult(1, messages)
			return result
		}
		msg := util.SetMessage("INFO", "Created S3 bucket ["+config.ArchivePluginParameters["BucketName"]+"] successfully")
		messages = append(messages, msg)

	}

	var s3FileList []string
	s3FileList, err = GetS3FileList(s3svc, config.ArchivePluginParameters["BucketName"], bucketPrefix, config.StoragePluginParameters["BackupDestPath"], backupPath)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	for _, file := range s3FileList {
		basePath := filepath.Dir(file)
		relativeBasePath := strings.Replace(basePath, config.StoragePluginParameters["BackupDestPath"]+"/", "", 1)

		uploadToS3(s3svc, config.ArchivePluginParameters["BucketName"], relativeBasePath+"/", file)
		msg := util.SetMessage("INFO", "Uploaded file ["+file+"] to bucket ["+config.ArchivePluginParameters["BucketName"]+relativeBasePath+" successfully")
		messages = append(messages, msg)
	}

	msg = util.SetMessage("INFO", "Archive backup ["+backupPath+"] to AWS bucket ["+config.ArchivePluginParameters["BucketName"]+"] completed successfully")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.ArchivePluginParameters["AwsRegion"]),
	}))

	s3svc := s3.New(sess)

	bucketPrefix := config.ProfileName + "/" + config.ConfigName + "/"
	objects, err := ListArchiveFolders(s3svc, config.ArchivePluginParameters["BucketName"], bucketPrefix)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	archiveList, err := pluginUtil.ListArchives(objects)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	archivesByPolicy := util.GetArchivesByPolicy(config.SelectedBackupPolicy, archiveList)
	archiveCount := len(archivesByPolicy)

	if archiveCount > config.SelectedArchiveRetention {
		count := 1
		for archive := range pluginUtil.ReverseArchiveList(archivesByPolicy) {
			if count > config.SelectedArchiveRetention {
				msg := util.SetMessage("INFO", fmt.Sprintf("Number of archives [%d] greater than archive retention [%d]", archiveCount, config.SelectedArchiveRetention))
				messages = append(messages, msg)
				archiveCount = archiveCount - 1

				archiveName := archive.Name + "_" + archive.Policy + "_" + archive.WorkflowId + "_" + util.IntToString(archive.Epoch)
				msg = util.SetMessage("INFO", "Deleting archive "+archiveName)
				messages = append(messages, msg)

				archivePrefix := bucketPrefix + archiveName
				err := DeleteFolder(s3svc, config.ArchivePluginParameters["BucketName"], archivePrefix)
				if err != nil {
					msg := util.SetMessage("ERROR", "Archive "+archiveName+" delete failed! "+err.Error())
					messages = append(messages, msg)
					result = util.SetResult(1, messages)
					return result
				}
				msg = util.SetMessage("INFO", "Archive "+archiveName+" deleted successfully")
				messages = append(messages, msg)
			}
			count = count + 1
		}
	} else {
		msg := util.SetMessage("INFO", fmt.Sprintf("Archive deletion skipped, there are [%d] archives but archive retention is [%d]", archiveCount, config.SelectedArchiveRetention))
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveList(config util.Config) util.Archives {
	var archives util.Archives
	var result util.Result
	var messages []util.Message

	bucketPrefix := config.ProfileName + "/" + config.ConfigName + "/"

	msg := util.SetMessage("INFO", "Archiving list for AWS bucket ["+config.ArchivePluginParameters["BucketName"]+"] path ["+bucketPrefix+"]")
	messages = append(messages, msg)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.ArchivePluginParameters["AwsRegion"]),
	}))

	s3svc := s3.New(sess)

	objects, err := ListArchiveFolders(s3svc, config.ArchivePluginParameters["BucketName"], bucketPrefix)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		archives.Result = result

		return archives
	}

	archiveList, err := pluginUtil.ListArchives(objects)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		archives.Result = result

		return archives
	}

	result = util.SetResult(0, messages)
	archives.Result = result
	archives.Archives = archiveList

	return archives
}

func (r archivePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample-archive"
	plugin.Description = "A sample archive plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "archive"

	var capabilities []util.Capability
	var archiveCap util.Capability
	archiveCap.Name = "archive"

	var archiveListCap util.Capability
	archiveListCap.Name = "archiveList"

	var archiveDeleteCap util.Capability
	archiveDeleteCap.Name = "archiveDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, archiveCap, archiveListCap, archiveDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}

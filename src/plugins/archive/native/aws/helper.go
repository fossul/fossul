package main

import (
	"github.com/fossul/fossul/src/engine/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var preserveDirStructureBool bool

func BucketExists(s3svc *s3.S3, bucketName string) (bool, error) {
	var bucketList []string
	s3ListBuckets, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return false, err
	}

	for _, bucket := range s3ListBuckets.Buckets {
		bucketList = append(bucketList, aws.StringValue(bucket.Name))
		if aws.StringValue(bucket.Name) == bucketName {
			return true, nil
		}
	}

	return false, nil
}

func CreateBucket(s3svc *s3.S3, bucketName string) error {
	_, err := s3svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	err = s3svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	return nil
}

func FolderExists(s3svc *s3.S3, bucketName, folderName string) (bool, error) {
	resp, err := s3svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: &folderName,
	})

	if err != nil {
		return false, err
	}

	if len(resp.Contents) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func CreateFolder(s3svc *s3.S3, bucketName, folder string) error {

	_, err := s3svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(folder),
	})

	if err != nil {
		return err
	}

	return nil
}

func GetS3FileList(s3svc *s3.S3, bucketName, bucketPrefix, backupDestPath string, dirPath string) ([]string, error) {
	var fileList []string
	treeList, err := util.DirectoryTreeList(dirPath)
	if err != nil {
		return fileList, err
	}

	for _, file := range treeList {
		isDir, err := util.IsDirectory(file)
		if err != nil {
			return fileList, err
		}

		if isDir {
			file = strings.Replace(file, backupDestPath+"/", "", 1)
			folderExists, err := FolderExists(s3svc, bucketName, file+"/")
			if err != nil {
				return fileList, err
			}

			if !folderExists {
				err = CreateFolder(s3svc, bucketName, file+"/")
				if err != nil {
					return fileList, err
				}
			}
		} else {
			fileList = append(fileList, file)
		}
	}

	return fileList, nil
}

func uploadToS3(s3svc *s3.S3, bucketName string, bucketPrefix string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var key string
	if preserveDirStructureBool {
		fileDirectory, _ := filepath.Abs(filePath)
		key = bucketPrefix + fileDirectory
	} else {
		key = bucketPrefix + path.Base(filePath)
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}
	_, err = s3svc.PutObject(params)
	if err != nil {
		return err
	}

	return nil
}

func ListArchiveFolders(s3svc *s3.S3, bucketName string, bucketPrefix string) ([]string, error) {
	var objects []string

	inputparams := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(bucketPrefix),
	}

	pageNum := 0
	err := s3svc.ListObjectsPages(inputparams, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++
		for _, value := range page.Contents {
			subfolderRegex := regexp.MustCompile(`^\S+\/\S+\/\S+_\S+_\d+_\d+\/\d+.*`)
			subFolderMatch := subfolderRegex.FindStringSubmatch(*value.Key)
			if subFolderMatch != nil {
				continue
			}

			archiveDirRegex := regexp.MustCompile(`^\S+\/\S+\/(\S+_\S+_\d+_\d+)\/`)
			archiveDirMatch := archiveDirRegex.FindStringSubmatch(*value.Key)

			objects = append(objects, archiveDirMatch[1])
		}

		return true
	})

	if err != nil {
		return objects, err
	}

	return objects, nil
}

func DeleteFolder(s3svc *s3.S3, bucketName, bucketPrefix string) error {
	iter := s3manager.NewDeleteListIterator(s3svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(bucketPrefix),
	})

	if err := s3manager.NewBatchDeleteWithClient(s3svc).Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	return nil
}

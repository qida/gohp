package oss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"gohp/logx"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// OSS
type OssMinIO struct {
	AccessKeyId     string
	AccessKeySecret string
	NameBucket      string
	Expires         int
	Client          *minio.Client
}

func NewOssMinIO(expires int, end_point, name_bucket, access_key_id, access_key_secret string) *OssMinIO {
	client, err := minio.New(end_point, &minio.Options{
		Creds:  credentials.NewStaticV4(access_key_id, access_key_secret, ""),
		Secure: false,
	})
	if err != nil {
		logx.Error(err.Error())
		return nil
	}
	return &OssMinIO{
		Expires:         expires,
		NameBucket:      name_bucket,
		AccessKeyId:     access_key_id,
		AccessKeySecret: access_key_secret,
		Client:          client,
	}
}

var options = minio.PutObjectOptions{
	ContentType: "application/octet-stream",
}

// 文件名不包括相对路径
func (t *OssMinIO) UploadFile(time_out time.Duration, file_name, prop string, file interface{}) (_file_key string, _err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time_out)
	defer cancel()
	switch file := file.(type) {
	case string:
		info, err := t.Client.FPutObject(ctx, t.NameBucket, file_name, file, options)
		if err != nil {
			_err = err
			logx.Error("UploadFile", zap.Any("func", GetCurrFuncName()), zap.Error(err))
			return
		}
		// _file_url = fmt.Sprintf("/%s/%s", info.Bucket, info.Key)
		_file_key = info.Key
	case []byte:
		info, err := t.Client.PutObject(ctx, t.NameBucket, file_name, bytes.NewReader(file), int64(len(file)), options)
		if err != nil {
			_err = err
			logx.Error("UploadFile", zap.Any("func", GetCurrFuncName()), zap.Error(err))
			return
		}
		// _file_url = fmt.Sprintf("/%s/%s", info.Bucket, info.Key)
		_file_key = info.Key
	default:
		_err = errors.New("file 参数错误")
	}
	return
}
func (t *OssMinIO) ListObjects() (_err error) {
	// List the content of the prefix:
	lopts := minio.ListObjectsOptions{
		Recursive: true,
		// Prefix:     "/",
		// Prefix:    string(prefix) + "/",
	}
	// List all objects from a bucket-name with a matching prefix.
	for object := range t.Client.ListObjects(context.Background(), t.NameBucket, lopts) {
		if object.Err != nil {
			_err = object.Err
			return
		}
		fmt.Println(object.Key, "Size:", object.Size)
	}
	return
}

func (t *OssMinIO) DeleteFile(file_name string) (_err error) {
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	objectsCh := make(chan minio.ObjectInfo)
	objectsCh <- minio.ObjectInfo{Key: file_name}
	err := <-t.Client.RemoveObjects(context.Background(), t.NameBucket, objectsCh, opts)
	_err = err.Err
	return
}

func (t *OssMinIO) DeleteAllFiles() (_err error) {
	defer func() {
		if _err != nil {
			logx.Error("DeleteAllFiles", zap.Any("func", GetCurrFuncName()), zap.Error(_err))
		}
	}()
	if t == nil {
		_err = errors.New("DeleteAllFiles 对没有完成初始化")
		return
	}
	logx.Info("限测试使用 删除所有图片")
	objectsCh := make(chan minio.ObjectInfo)
	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range t.Client.ListObjects(context.Background(), t.NameBucket, minio.ListObjectsOptions{Prefix: "/", Recursive: true}) {
			if object.Err != nil {
				logx.Error("删除", zap.Error(object.Err))
			}
			objectsCh <- object
		}
	}()
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	for rErr := range t.Client.RemoveObjects(context.Background(), t.NameBucket, objectsCh, opts) {
		fmt.Println("Error detected during deletion: ", rErr)
	}
	logx.Info("限测试使用 已删除所有图片")
	return nil
}

// func (t *OssMinIO) GetBucketStat() (_status map[string]interface{}, _err error) {
// 	stat, err := t.Client.GetBucketStat(t.NameBucket)
// 	if err != nil {
// 		_err = err
// 		return
// 	}
// 	_status = make(map[string]interface{}, 0)
// 	_status["Storage"] = stat.Storage
// 	_status["ObjectCount"] = stat.ObjectCount
// 	_status["MultipartUploadCount"] = stat.MultipartUploadCount
// 	_status["LiveChannelCount"] = stat.LiveChannelCount
// 	_status["LastModifiedTime"] = stat.LastModifiedTime
// 	_status["StandardStorage"] = stat.StandardStorage
// 	_status["StandardObjectCount"] = stat.StandardObjectCount
// 	_status["InfrequentAccessStorage"] = stat.InfrequentAccessStorage
// 	_status["InfrequentAccessRealStorage"] = stat.InfrequentAccessRealStorage
// 	_status["InfrequentAccessObjectCount"] = stat.InfrequentAccessObjectCount
// 	_status["ArchiveStorage"] = stat.ArchiveStorage
// 	_status["ArchiveObjectCount"] = stat.ArchiveObjectCount
// 	_status["ColdArchiveStorage"] = stat.ColdArchiveStorage
// 	_status["ColdArchiveRealStorage"] = stat.ColdArchiveRealStorage
// 	_status["ColdArchiveObjectCount"] = stat.ColdArchiveObjectCount
// 	return
// }

package oss

import (
	"bytes"
	"errors"
	"time"

	"github.com/qida/gohp/logx"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

// 阿里OSS
type OssAliyun struct {
	AccessKeyId     string
	AccessKeySecret string
	NameBucket      string
	Expires         int
	Client          *oss.Client
	Bucket          *oss.Bucket
}

func NewOssAliyun(expires int, end_point, name_bucket, access_key_id, access_key_secret string) *OssAliyun {
	client, err := oss.New(end_point, access_key_id, access_key_secret)
	if err != nil {
		logx.Error("错误", zap.Error(err))
		return nil
	}
	bucket, err := client.Bucket(name_bucket)
	if err != nil {
		logx.Error("错误", zap.Error(err))
		return nil
	}
	return &OssAliyun{
		Expires:         expires,
		NameBucket:      name_bucket,
		AccessKeyId:     access_key_id,
		AccessKeySecret: access_key_secret,
		Client:          client,
		Bucket:          bucket,
	}
}
func (t *OssAliyun) BucketList() (_buckets []string, _err error) {
	lsRes, err := t.Client.ListBuckets()
	if err != nil {
		_err = err
		return
	}
	for _, bucket := range lsRes.Buckets {
		logx.Debugf("Buckets: %s", bucket.Name)
		_buckets = append(_buckets, bucket.Name)
	}
	return
}

// file_path_name:应该为包含完整的相对路径 如: /22年10月/TPKEY_m_1665393213289024000_0001.jpg
func (t *OssAliyun) UploadFile(time_out time.Duration, file_path_name, prop string, file interface{}) (_file_url string, _err error) {
	options := []oss.Option{
		//TODO 正式 需修改
		// oss.Expires(time.Now().AddDate(t.Expires, 0, 0)), //3年
		oss.Expires(time.Now().Add(60 * time.Minute)), //1小时
		oss.ObjectACL(oss.ACLPublicRead),
		oss.Meta("prop", prop),
	}
	switch file := file.(type) {
	case string:
		err := t.Bucket.PutObjectFromFile(file_path_name, file, options...)
		if err != nil {
			_err = err
			return
		}
	case []byte:
		err := t.Bucket.PutObject(file_path_name, bytes.NewReader(file), options...)
		if err != nil {
			_err = err
			return
		}
	default:
		_err = errors.New("file 参数错误")
		return
	}
	_file_url = file_path_name
	return
}

func (t *OssAliyun) DeleteFile(file_name string) (_err error) {
	t.Bucket.DeleteObject(file_name)
	return
}

func (t *OssAliyun) DeleteAllFiles() (_err error) {
	if t == nil {
		return
	}
	defer func() {
		if _err != nil {
			logx.Error("DeleteAllFiles", zap.Any("func", GetCurrFuncName()), zap.Error(_err))
		}
	}()
	logx.Info("限测试使用 删除所有图片")
	marker := oss.Marker("")
	for {
		lor, err := t.Bucket.ListObjects(marker)
		if err != nil {
			return err
		}
		for _, object := range lor.Objects {
			err = t.Bucket.DeleteObject(object.Key)
			if err != nil {
				return err
			}
			logx.Debugf("删除:%s", object.Key)
		}
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	logx.Info("限测试使用 已删除所有图片")
	return nil
}

func (t *OssAliyun) GetBucketStat() (_status map[string]interface{}, _err error) {
	stat, err := t.Client.GetBucketStat(t.NameBucket)
	if err != nil {
		_err = err
		return
	}
	_status = make(map[string]interface{}, 0)
	_status["Storage"] = stat.Storage
	_status["ObjectCount"] = stat.ObjectCount
	_status["MultipartUploadCount"] = stat.MultipartUploadCount
	_status["LiveChannelCount"] = stat.LiveChannelCount
	_status["LastModifiedTime"] = stat.LastModifiedTime
	_status["StandardStorage"] = stat.StandardStorage
	_status["StandardObjectCount"] = stat.StandardObjectCount
	_status["InfrequentAccessStorage"] = stat.InfrequentAccessStorage
	_status["InfrequentAccessRealStorage"] = stat.InfrequentAccessRealStorage
	_status["InfrequentAccessObjectCount"] = stat.InfrequentAccessObjectCount
	_status["ArchiveStorage"] = stat.ArchiveStorage
	_status["ArchiveObjectCount"] = stat.ArchiveObjectCount
	_status["ColdArchiveStorage"] = stat.ColdArchiveStorage
	_status["ColdArchiveRealStorage"] = stat.ColdArchiveRealStorage
	_status["ColdArchiveObjectCount"] = stat.ColdArchiveObjectCount
	return
}

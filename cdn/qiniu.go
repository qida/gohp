package cdn

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/cdn"
	"github.com/qiniu/api.v7/v7/storage"
	"golang.org/x/net/context"
)

type QiNiu struct {
	Bucket        string
	Url           string
	Mac           *qbox.Mac
	CdnManager    *cdn.CdnManager
	Config        *storage.Config
	PutPolicy     *storage.PutPolicy
	BucketManager *storage.BucketManager
}

func NewQiNiu(bucket string, url string, accessKey, secretKey string) *QiNiu {
	zone, _ := storage.GetZone(accessKey, bucket)
	mac := qbox.NewMac(accessKey, secretKey)
	cdnManager := cdn.NewCdnManager(mac)
	cfg := &storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: true,
	}
	bucketManager := storage.NewBucketManager(mac, cfg)
	return &QiNiu{
		Bucket:        bucket,
		Url:           url,
		Mac:           mac,
		CdnManager:    cdnManager,
		BucketManager: bucketManager,
		Config:        cfg,
		PutPolicy: &storage.PutPolicy{
			Scope: bucket,
		},
	}
}

func (c *QiNiu) Upload(localFile io.Reader, size int64, file_name string, is_refresh bool) (url_file string, err error) {
	var isRefresh = false
	if is_refresh {
		_, err1 := c.BucketManager.Stat(c.Bucket, file_name)
		if err1 == nil { //存在
			isRefresh = true
		} else {
			log.Printf("%+v", err1)
		}
	}
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", c.Bucket, file_name), //覆盖上传
	}
	upToken := putPolicy.UploadToken(c.Mac)
	formUploader := storage.NewFormUploader(c.Config)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": file_name,
		},
	}
	err = formUploader.Put(context.Background(), &ret, upToken, file_name, localFile, size, &putExtra)
	if err != nil {
		return
	}
	url_file = ret.Key
	fmt.Printf("=====上传======\r\nKey:%s Hash:%s\r\n==============\r\n", ret.Key, ret.Hash)
	if isRefresh {
		err = c.Refresh(c.Url + url_file)
	}
	return
}

func (c *QiNiu) Refresh(url ...string) (err error) {
	for i := 0; i < len(url); i++ {
		if !strings.Contains(url[i], "http:") {
			url[i] = "http:" + url[i]
		}
	}
	fmt.Printf("=====刷新文件======\r\n%s\r\n==============\r\n", url)
	_, err = c.CdnManager.RefreshUrls(url)
	return
}

func (c *QiNiu) UploadFile(file_name string, file_data []byte) (url_file string, err error) {
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", c.Bucket, file_name),
	}
	upToken := putPolicy.UploadToken(c.Mac)
	formUploader := storage.NewFormUploader(c.Config)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": file_name,
		},
	}
	err = formUploader.Put(context.Background(), &ret, upToken, file_name, bytes.NewReader(file_data), int64(len(file_data)), &putExtra)
	if err != nil {
		return
	}
	url_file = ret.Key
	fmt.Printf("=====上传======\r\nKey:%s Hash:%s\r\n==============\r\n", ret.Key, ret.Hash)
	err = c.Refresh(c.Url + url_file)
	return
}

func (c *QiNiu) MoveFile(src_url string, dst_url string) (url_file string, err error) {
	bucketManager := storage.NewBucketManager(c.Mac, c.Config)
	//如果目标文件存在，是否强制覆盖，如果不覆盖，默认返回614 file exists
	err = bucketManager.Move(c.Bucket, src_url, c.Bucket, dst_url, true)
	if err != nil {
		log.Printf("移动文件 %+v", err1)
		return
	}
	url_file = dst_url
	return
}

// 返回私人有地址
func (c *QiNiu) GetPrivateMediaUrl(src_url string) (privateAccessURL string) {
	deadline := time.Now().Add(time.Minute * 60).Unix() //60分钟有效期
	if !strings.Contains(src_url, "http:") {
		src_url = "http:" + src_url
	}
	privateAccessURL = storage.MakePrivateURL(c.Mac, c.Url, src_url, deadline)
	return
}

func (c *QiNiu) Delete(url string) (err error) {
	if url == "" {
		// err = errors.New("文件url不能为空")
		return
	}
	bucketManager := storage.NewBucketManager(c.Mac, c.Config)
	key := strings.TrimPrefix(url, c.Url)
	err = bucketManager.Delete(c.Bucket, key)
	if err != nil {
		log.Printf("DeleteFile key:%s Err:%+v", key, err1)
	} else {
		fmt.Printf("成功删除：%s\r\n", key)
	}
	return
}

func (c *QiNiu) GetTokenUpload(region string, key string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	//ECN, SCN, NCN, NA, ASG
	// putPolicy.CallbackURL = "https://xxx/upload"
	m["Region"] = region
	m["UpTokenURL"] = c.Bucket
	if key == "" {
		m["Key"] = fmt.Sprintf("temp/%d", time.Now().Unix()) //不起作用
	} else {
		m["Key"] = key
	}
	m["UpToken"] = c.PutPolicy.UploadToken(c.Mac)
	m["Domain"] = strings.Replace(c.Url, "/", "", -1)
	return
}

func (c *QiNiu) RefreshDir(dir string) (err error) {
	// 刷新目录，刷新目录需要联系七牛技术支持开通权限
	// 单次请求链接不可以超过10个，如果超过，请分批发送请求
	dirsToRefresh := []string{dir}
	_, err = c.CdnManager.RefreshDirs(dirsToRefresh)
	fmt.Printf("刷新目录：%v\r\n", dirsToRefresh)
	return
}

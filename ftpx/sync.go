package ftpx

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"gohp/logx"

	"github.com/secsy/goftp"
)

type Sync struct {
	TimeOut int
	Client  *goftp.Client
}

func NewSync(host, port, user_name, pass_word string, time_out int) *Sync {
	config := goftp.Config{
		User:               user_name,
		Password:           pass_word,
		ConnectionsPerHost: 10,
		Timeout:            time.Duration(time_out) * time.Second,
		Logger:             os.Stderr,
	}
	client, err := goftp.DialConfig(config, net.JoinHostPort(host, port))
	if err != nil {
		log.Printf("初始化文件同步错误:%+v", err)
		return nil
	}
	return &Sync{
		TimeOut: time_out,
		Client:  client,
	}
}

// const (
// 	numberGoroutines = 1
// )

// var bufferWg sync.WaitGroup

func (t *Sync) UploadFile(file_paths ...string) (err error) {
	if len(file_paths) == 0 {
		return
	}
	tt := time.Now()
	// Upload a file from disk
	for i := 0; i < len(file_paths); i++ {
		bigFile, err := os.Open(file_paths[i])
		if err != nil {
			panic(err)
		}
		err = t.Client.Store(filepath.Base(file_paths[i]), bigFile)
		if err != nil {
			panic(err)
		}
	}

	// //创建了任务的缓冲通道
	// taskLoad := len(file_paths)
	// tasks := make(chan string, taskLoad)

	// bufferWg.Add(numberGoroutines)

	// //创建4个Goroutine
	// for gr := 1; gr <= numberGoroutines; gr++ {
	// 	go t.worker(tasks, gr)
	// }

	// //每张图214KB 10000个 2.04GB
	// // logx.Debugf("共%d个文件 %dMB 共耗时:%v", len(file_paths), countSize>>20, time.Since(tt))
	// // if err = t.Client.Quit(); err != nil {
	// // 	log.Printf("文件同步 退出 错误:%+v", err)
	// // }

	// //向缓冲通道中放入数据
	// for i := 0; i < taskLoad; i++ {
	// 	tasks <- file_paths[i]
	// }
	// close(tasks)
	// bufferWg.Wait()
	// // logx.Debugf("共%d个文件 %dMB 共耗时:%v", len(file_paths), countSize>>20, time.Since(tt))
	logx.Debugf("共耗时:%v", time.Since(tt))
	return
}

// func (t *Sync) worker(tasks chan string, work int) {
// 	defer bufferWg.Done()
// 	for {
// 		task, ok := <-tasks
// 		if !ok {
// 			logx.Debugf("Worker: %d : 结束工作 ", work)
// 			return
// 		}

// 		logx.Debugf("Worker: %d : 开始工作 %s", work, task)
// 		bytImage, err := os.ReadFile(task)
// 		if err != nil {
// 			log.Printf("文件同步 读取 错误:%+v", err)
// 			continue
// 		}
// 		t.Client, err = ftp.Dial(net.JoinHostPort(t.Host, t.Post), ftp.DialWithTimeout(time.Duration(t.Timeout)*time.Second))
// 		if err != nil {
// 			log.Printf("文件同步 连接 错误:%+v", err)
// 		}
// 		err = t.Client.Login(t.UserName, t.Password)
// 		if err != nil {
// 			log.Printf("文件同步 登陆 错误:%+v", err)
// 		}
// 		err = t.Client.Stor(filepath.Base(task), bytes.NewBuffer(bytImage))
// 		if err != nil {
// 			log.Printf("文件同步 上传 错误:%+v", err)
// 			continue
// 		}
// 		logx.Debugf("Worker: %d : 完成工作 %s", work, task)
// 		if err = t.Client.Quit(); err != nil {
// 			log.Printf("文件同步 退出 错误:%+v", err)
// 		}
// 	}
// }

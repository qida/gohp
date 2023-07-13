# 重装Linux手册

### 换源
``` sh
sudo sed -i.bak -e 's|^mirrorlist=|#mirrorlist=|' -e 's|^#baseurl=|baseurl=|' -e 's|http://mirror.centos.org|https://mirrors.aliyun.com|' /etc/yum.repos.d/CentOS-*.repo

sudo yum -y erase podman buildah  # 卸载Podman
sudo yum install git rsync -y     # 安装必要工具
sudo yum clean all                # 清空缓存
sudo yum makecache                # 建立新缓存
sudo yum update -y                # 升级

```
### 安装Docker
``` sh
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh --mirror Aliyun
rm get-docker.sh

sudo ln -s /usr/libexec/docker/cli-plugins/docker-compose /usr/bin/docker-compose
```

``` sh
sudo tee /etc/docker/daemon.json <<-'EOF'
{
   "registry-mirrors": [
       "https://mirror.ccs.tencentyun.com"
  ]
}
EOF

docker network create --driver=bridge --subnet=192.168.0.0/24 qida

systemctl daemon-reload
service docker restart

docker info
```
### SSH

``` sh
sudo tee -a /etc/ssh/sshd_config <<-'EOF'
Port 10022
ClientAliveInterval 60 #服务器端向客户端发送心跳以判断客户端是否存活（即客户端是否操作服务器）的时间间隔，单位为秒，默认是0。
ClientAliveCountMax 3
EOF

systemctl restart sshd.service
```

### 磁盘挂载
``` sh
sudo mkdir -p /mnt
sudo cp /etc/fstab /etc/fstab.bak
sudo echo '/dev/vdb1        /mnt ext4    defaults 0       0' >> /etc/fstab
sudo mount /dev/vdb1 /mnt
sudo df -h
```

### 环境变量

``` sh
cp /etc/profile /etc/profile.bak

sudo tee -a /etc/profile <<-'EOF'
export BASE=/mnt
export DOWNLOADS=$BASE/downloads
export GOPATH=$BASE/gopath
export GOROOT=$BASE/go
export GOBIN=$GOROOT/bin
export GO15VENDOREXPERIMENT=1
export GOROOT_BOOTSTRAP=$BASE/go1.4
export CGO_ENABLE=0
export GOARCH=amd64
export GOOS=linux
export GO111MODULE=auto
export GOPROXY=https://proxy.golang.com.cn,https://goproxy.io,direct
export GOPRIVATE=git.sunqida.cn,github.com/qida/go
export PATH=$PATH:$GOBIN
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$BASE/gopath/src/github.com/NICEXAI/WeWorkFinanceSDK/lib
EOF

source /etc/profile
#LD_LIBRARY_PATH
find / -name "*libxcb.so.1*"
```
### Go Package
``` sh
go version
go get github.com/nsf/gocode
go get github.com/rogpeppe/godef
go get github.com/golang/lint/golint
go get github.com/lukehoban/go-outline
go get sourcegraph.com/sqs/goreturns
go get golang.org/x/tools/cmd/gorename
go get github.com/tpng/gopkgs
go get github.com/newhook/go-symbols
go get golang.org/x/tools/cmd/guru
```

### 问题描述
/tmp/go-build3735770271/b001/exe/main: error while loading shared libraries: libWeWorkFinanceSdk_C.so: cannot open shared object file: No such file or directory exit status 127
解决方案 libxcb
``` sh
git config --global http.https://github.com.proxy socks5://10.10.10.3:1080
git config --global http.https://github.com.sslVerify false
git config --system http.sslVerify false

go get -u github.com/NICEXAI/WeWorkFinanceSDK
sudo tee -a /etc/ld.so.conf <<-'EOF'
/mnt/gopath/src/github.com/NICEXAI/WeWorkFinanceSDK/lib
EOF
sudo ldconfig
```

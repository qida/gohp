# netx

网络工具集：纯真 IP 库（qqwry.dat）归属地查询、本机/外网 IP 获取与 CIDR 计算、URL/Base64 编解码。

## 功能列表

### ip.go — 纯真 IP 库查询

- `QQwry`：纯真 IP 库查询器，导出字段含 Ip、Lng、Lat、Address（归属地/国家）、Ips（区域信息）。
- `NewQQwry(path_file string) *QQwry`：初始化查询器；本地库文件不存在时自动从 update.cz88.net 下载最新纯真 IP 库并保存到该路径。
- `(*QQwry) Find(ip string) string`：查询 IP 归属地（内部 GBK 自动转码）；内网/回环地址返回"内网地址"，格式非法返回"IP有误"，库未初始化返回"网络库，初始化失败"。
- `GetOnline() ([]byte, error)`：在线下载、解密并解压最新纯真 IP 库数据。
- `GetKey() (uint32, error)`：从 copywrite.rar 中获取库文件的解密密钥。
- `(*QQwry) ReadData(num int, offset ...int64)` / `SetOffset(offset)`：底层按偏移读取数据。

### ipx.go — 本机/外网 IP 与 CIDR

- `ExternalIP() (net.IP, error)`：遍历网卡，返回第一个处于 UP 状态的非回环 IPv4 地址。
- `GetIpFromAddr(addr net.Addr) net.IP`：从 `net.Addr` 提取 IPv4（排除回环与非 IPv4）。
- `GetIpOut() string`：并发请求多个公网 IP 服务（myexternalip.com、ip.cip.cc、4.ipw.cn、test.ipw.cn），返回最先拿到的结果，整体 10 秒超时。
- `GetIpIn() string` / `GetIpsIn() []string`：获取内网 IPv4（首个 / 全部非回环 IPv4）。
- `IpToUInt32(ip_str) uint32` / `IntToIp(ip_int) string`：IP 字符串与 uint32 整数互转。
- `GetCidrIpRangeIP / GetCidrIpRangeString / GetCidrIpRangeUInt32(cidr)`：根据 CIDR 计算网络地址（首个 IP）与广播地址（最后 IP），返回 net.IP / 字符串 / uint32 三种形式。
- `IsIPinCIDR(ipstr, netstr string) (bool, error)`：判断 IP 是否属于指定 CIDR 网段。

### url.go — URL 与 Base64 编解码

- `UrlEncode(str) string` / `UrlDecode(str) (string, error)`：URL 查询串编解码（空格编码为 "+"）。
- `Base64Encode(str) string` / `Base64Decode(str) (string, error)`：标准 Base64 编解码。

## 使用示例

查询 IP 归属地（提炼自 ip_test.go）：

```go
q := netx.NewQQwry("qqwry.dat") // 本地无库时自动联网下载
addr := q.Find("39.78.34.61")   // 返回归属地字符串
```

## 备注

- 外部依赖：`github.com/yinheli/mahonia`（纯真库 GBK 解码）。
- `NewQQwry` 在本地库缺失时需访问 `update.cz88.net`；`GetIpOut` 需访问公网 IP 查询服务，均依赖外网连通性。
- 纯真 IP 库解析逻辑来自 github.com/freshcn/qqwry 项目。

// from https://github.com/freshcn/qqwry
package netx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/yinheli/mahonia"
)

const (
	INDEX_LEN       = 7    // INDEX_LEN 索引长度
	REDIRECT_MODE_1 = 0x01 // REDIRECT_MODE_1 国家的类型, 指向另一个指向
	REDIRECT_MODE_2 = 0x02 // REDIRECT_MODE_2 国家的类型, 指向一个指向
)

type QQwry struct {
	Ip      string
	Lng     float64
	Lat     float64
	Address string
	Ips     string

	filepath string
	FileData []byte
	NumIp    int64
	Offset   int64
	Enc      mahonia.Decoder
}

func NewQQwry(path_file string) (qqwry *QQwry) {
	var err error
	qqwry = &QQwry{filepath: path_file, Enc: mahonia.NewDecoder("gbk")}
	var tmpData []byte
	// 判断文件是否存在
	_, err = os.Stat(qqwry.filepath)
	if err != nil && os.IsNotExist(err) {
		log.Printf("文件不存在: %s ，尝试从网络获取最新纯真 IP 库\r\n", path_file)
		tmpData, err = GetOnline()
		log.Printf("临时Data Len %d \r\n", len(tmpData))
		if err != nil {
			log.Println(err.Error())
			return
		} else {
			if err = ioutil.WriteFile(qqwry.filepath, tmpData, 0644); err == nil {
				log.Printf("已将最新的纯真 IP 库保存到本地 %s ", qqwry.filepath)
			} else {
				log.Printf("失败 最新的纯真 IP 库保存到本地 %s ", err.Error())
			}
		}
	} else {
		// 打开文件句柄
		// log.Printf("从本地数据库文件 %s 打开\n", f.FilePath)
		file, err := os.OpenFile(qqwry.filepath, os.O_RDONLY, 0400)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer file.Close()
		tmpData, err = ioutil.ReadAll(file)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	qqwry.FileData = tmpData
	// buf := qqwry.FileData[0:8]
	// start := binary.LittleEndian.Uint32(buf[:4])
	// end := binary.LittleEndian.Uint32(buf[4:])
	// qqwry.NumIp = int64((end-start)/INDEX_LEN + 1)
	return
}

// Find ip地址查询对应归属地信息
func (q *QQwry) Find(ip string) string {
	if strings.Count(ip, ".") != 3 {
		return "IP有误"
	}
	switch ip {
	case "127.0.0.1", "localhost":
		return "内网地址"
	default:
		log.Printf("IP:%s\r\n", ip)
	}
	if len(q.FileData) == 0 {
		return "网络库，初始化失败"
	}
	q.Ip = ip
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return ""
	}
	var country []byte
	var area []byte
	mode := q.readMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := q.readUInt24()
			country = q.readString(c)
			countryOffset += 4
		} else {
			country = q.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.readArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := q.readUInt24()
		country = q.readString(countryOffset)
		area = q.readArea(offset + 8)
	} else {
		country = q.readString(offset + 4)
		area = q.readArea(offset + uint32(5+len(country)))
	}
	q.Address = q.Enc.ConvertString(string(country))
	q.Ips = q.Enc.ConvertString(string(area))
	log.Println(q.Address)
	return q.Address
}

// @ref https://zhangzifan.com/update-qqwry-dat.html
func GetOnline() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://update.cz88.net/ip/qqwry.rar", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html, */*")
	req.Header.Set("User-Agent", "Mozilla/3.0 (compatible; Indy Library)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		if key, err := getKey(); err == nil {
			for i := 0; i < 0x200; i++ {
				key = key * 0x805
				key++
				key = key & 0xff
				body[i] = byte(uint32(body[i]) ^ key)
			}
			reader, err := zlib.NewReader(bytes.NewReader(body))
			if err != nil {
				return nil, err
			}
			return ioutil.ReadAll(reader)
		}
	}
	return nil, err
}

func getKey() (uint32, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://update.cz88.net/ip/copywrite.rar", strings.NewReader(""))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", `text/html, */*`)
	req.Header.Set("User-Agent", `Mozilla/3.0 (compatible; Indy Library)`)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Printf("body len1:%d\r\n", len(body))
		// @see https://stackoverflow.com/questions/34078427/how-to-read-packed-binary-data-in-go
		return binary.LittleEndian.Uint32(body[5*4:]), nil
	}
	return 0, err
}

// ReadData 从文件中读取数据
func (q *QQwry) ReadData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.SetOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.FileData))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.FileData[q.Offset:end]
	q.Offset = end
	return
}

// SetOffset 设置偏移量
func (q *QQwry) SetOffset(offset int64) {
	q.Offset = offset
}

// readMode 获取偏移值类型
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.ReadData(1, int64(offset))
	return mode[0]
}

// readArea 读取区域
func (q *QQwry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		}
		return q.readString(areaOffset)
	}
	return q.readString(offset)
}

// readString 获取字符串
func (q *QQwry) readString(offset uint32) []byte {
	q.SetOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.ReadData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// searchIndex 查找索引位置
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.ReadData(8, 0)
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, INDEX_LEN)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.ReadData(INDEX_LEN, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := byteToUInt32(buf[4:])
			buf = q.ReadData(INDEX_LEN)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		// 找到的比较大，向前移
		if _ip > ip {
			end = mid
		} else if _ip < ip { // 找到的比较小，向后移
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}
	}
}

// readUInt24
func (q *QQwry) readUInt24() uint32 {
	buf := q.ReadData(3)
	return byteToUInt32(buf)
}

// byteToUInt32 将 byte 转换为uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
func (this *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

package logx

import (
	"time"
)

const (
	//Msg            = "msg"
	//Time           = "time"
	//Level          = "level"
	//DateTime       = "dateTime"       //时间

	File           = "file"           //文件路径
	Title          = "title"          //名称
	Trans          = "Trans"          //事务  (列 : 整个视频的调用流程)
	Action         = "Action"         //操作  (列 : 具体的功能)
	Method         = "Method"         //方法  (列 : 具体的方法)
	AppID          = "AppID"          //
	ServiceCluster = "ServiceCluster" //集群ID
	ServiceID      = "ServiceID"      //服务ID
	ServiceName    = "ServiceName"    //服务名称
)

type Base struct {
	Code            string `thrift:"code,1" db:"code" json:"code"`
	LogType         string `thrift:"logType,2" db:"logType" json:"logType"`
	Level           string `thrift:"level,3" db:"level" json:"level"`
	Source          string `thrift:"source,4" db:"source" json:"source"`
	IP              string `thrift:"ip,5" db:"ip" json:"ip"`
	Trans           string `thrift:"trans,6" db:"trans" json:"trans"`
	TransId         string `thrift:"transId,7" db:"transId" json:"transId"`
	Action          string `thrift:"action,8" db:"action" json:"action"`
	Method          string `thrift:"method,9" db:"method" json:"method"`
	Description     string `thrift:"description,10" db:"description" json:"description"`
	ApiType         string `thrift:"apiType,11" db:"apiType" json:"apiType"`
	ApiUrl          string `thrift:"apiUrl,12" db:"apiUrl" json:"apiUrl"`
	ProtoData       string `thrift:"protoData,13" db:"protoData" json:"protoData"`
	CreateTimestamp string `thrift:"createTimestamp,14" db:"createTimestamp" json:"createTimestamp"`
}

type LogQuery struct {
	TransId     string `json:"TransId,omitempty"` //事务ID (雪花算法)
	Trans       string `json:"Trans,omitempty"`   //事务
	Level       string `json:"omitempty"`         //级别 (Debug, Info, Warn, Error, Panic, Fatal)
	Bucket      string `json:"Bucket,omitempty"`
	Measurement string `json:"Measurement,omitempty"` //测量维度
	StartTime   int64  `json:"StartTime,omitempty"`   //开始时间 (毫秒时间戳)
	EndTime     int64  `json:"EndTime,omitempty"`     //结束时间 (毫秒时间戳)
	LastTime    string `json:"LastTime,omitempty"`    //最近时间 (-2d/天;-5m/分钟)
	ServiceID   string `json:"ServiceID,omitempty"`   //服务ID
	ServiceName string `json:"ServiceName,omitempty"` //服务名
	AppID       string `json:"AppID,omitempty"`       //应用ID
	ClientIP    string `json:"ClientIP,omitempty"`    //客户端IP
}

type LogHeader struct {
	//ID (雪花算法)
	ID int64 `gorm:"type:bigint(20);primaryKey;comment: 自增 ID" json:"ID,omitempty"`
	//事务ID (雪花算法)
	TransId int64 `json:"TransId,omitempty"`
	//集群ID
	ServiceCluster string `json:"ServiceCluster,omitempty"`
	//服务ID
	ServiceID string `json:"ServiceID,omitempty"`
	//服务名
	ServiceName string `json:"ServiceName,omitempty"`

	//名称
	Title string `json:"title"`
	Level string `json:"level,omitempty"` //级别 (Debug, Info, Warn, Error, Panic, Fatal)
	Time  string `json:"time"`            //时间
	File  string `json:"file,omitempty"`  //
	Msg   string `json:"msg"`

	//事务  (列 : 整个视频的调用流程)
	Trans string `json:"Trans,omitempty"`
	//操作  (列 : 具体的功能)
	Action string `json:"Action,omitempty"`
	//方法  (列 : 具体的方法)
	Method string `json:"Method,omitempty"`

	//数据类型（QPS/服务/事件）
	DataType   string        `json:"DataType,omitempty"`
	CreateTime time.Duration `json:"CreateTime,omitempty"`
}

type LogBody struct {
	LogHeader
	//测量维度
	Measurement string `json:"Measurement,omitempty"`
	//应用ID
	AppID string `json:"AppID,omitempty"`
	////原始数据
	RawData string `json:"-"`
	//运行时间 (毫秒us) //耗时(毫秒us)
	ElapsedTime time.Duration `json:"ElapsedTime,omitempty"`
	Fields      []LogField    `gorm:"-"`
	//客户端IP
	ClientIP string `gorm:"comment:客户端IP" json:"ClientIP,omitempty"`
	//开始时间 (毫秒时间戳)
	StartTime string `gorm:"comment:开始时间" json:"StartTime,omitempty"`
	//结束时间 (毫秒时间戳)
	EndTime      string `gorm:"comment:结束时间" json:"EndTime,omitempty"`
	TimeUsed     int64
	TimeUsedType string
}

type LogField struct {
	Field string
	Value any `gorm:"-"` // influxdb _value
}

type LogStatistic struct {
	LogHeader
	FieldType string      // _field 统计值的类型, 如 在线设备 调阅耗时 网络延迟平均值 等等
	Field     string      // influxdb _field 四个类型: timeSince 时长 count 数量 num 计数 other flout64 其他
	Value     interface{} // influxdb _value
}

func LogAny(key string, val any) LogField {
	return LogField{
		Field: key,
		Value: val,
	}
}

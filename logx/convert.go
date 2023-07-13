package logx

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func toFile(data interface{}) ([]FileLogger, error) {
	maps, err := cast.ToSliceE(data)
	if err != nil {
		return nil, errors.Wrapf(err, "cast to map string string error data %v", data)
	}

	fls := make([]FileLogger, 0)
	for _, file := range maps {
		m := cast.ToStringMap(file)
		path, ok := m["path"]
		if !ok {
			continue
		}
		l, ok := m["level"]
		if !ok {
			continue
		}
		maxAge := 7
		if a, ok := m["maxAge"]; ok {
			maxAge = cast.ToInt(a)
		}
		size := 100
		if s, ok := m["size"]; ok {
			size = cast.ToInt(s)
		}
		backups := 0
		if c, ok := m["count"]; ok {
			backups = cast.ToInt(c)
		}
		isCompress := false
		if c, ok := m["isCompress"]; ok {
			isCompress = cast.ToBool(c)
		}
		level := InfoLevel
		switch l {
		case "debug":
			level = DebugLevel
		case "info":
			level = InfoLevel
		case "warn":
			level = WarnLevel
		case "error":
			level = ErrorLevel
		case "dpanic":
			level = DPanicLevel
		}
		fls = append(fls, FileLogger{
			Path:       cast.ToString(path),
			Level:      &level,
			MaxAge:     maxAge,  // 最大保留时间  天
			MaxSize:    size,    // 文件大小   mb
			MaxBackups: backups, // 分割周期   小时
			IsUTC:      false,
			IsCompress: isCompress,
		})
	}
	if len(fls) == 0 {
		return nil, errors.New("no file config")
	}
	return fls, nil
}

func toLogs(data interface{}) (FilesLogger, error) {
	m := cast.ToStringMap(data)
	path, ok := m["path"]
	if !ok {
		path = "logs/log.log"
	}
	maxAge := 7
	if a, ok := m["maxAge"]; ok {
		maxAge = cast.ToInt(a)
	}
	size := 100
	if s, ok := m["size"]; ok {
		size = cast.ToInt(s)
	}
	backups := 0
	if c, ok := m["count"]; ok {
		backups = cast.ToInt(c)
	}
	filesLogger := FilesLogger{
		Path:       cast.ToString(path),
		MaxAge:     maxAge,  // 最大保留时间  天
		MaxSize:    size,    // 文件大小   mb
		MaxBackups: backups, // 分割周期   小时
	}
	_, ok = m["isUTC"]
	if _, ok = m["isUTC"]; ok {
		filesLogger.IsUTC = false
	}

	if _, ok = m["isGzip"]; ok {
		filesLogger.IsCompress = true
	}

	return filesLogger, nil
}

func toLevel(data interface{}) (Level, error) {
	switch data.(string) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "dpanic":
		return DPanicLevel, nil
	}
	return InfoLevel, errors.New("not found level")
}

func toKafka(data interface{}) ([]LogKafka, error) {
	arrI, err := cast.ToSliceE(data)
	if err != nil {
		return nil, errors.New("type error")
	}

	lks := make([]LogKafka, 0)
	for _, a := range arrI {
		m, err := cast.ToStringMapE(a)
		if err != nil {
			continue
		}
		topic, ok := m["topic"]
		if !ok {
			continue
		}
		addressI, ok := m["address"]
		if !ok {
			continue
		}
		address, err := cast.ToStringSliceE(addressI)
		if err != nil {
			continue
		}
		level, err := toLevel(m["level"])
		if err != nil {
			continue
		}
		lks = append(lks, LogKafka{
			Topic:   cast.ToString(topic),
			Address: address,
			Level:   &level,
			fields:  getFields(m["fields"]),
		})
	}

	if len(lks) == 0 {
		return nil, errors.New("not found kafka config")
	}

	return lks, nil
}

func toRocketMQ(data interface{}) ([]LogRocketMQ, error) {
	arrI, err := cast.ToSliceE(data)
	if err != nil {
		return nil, errors.New("type error")
	}

	lrs := make([]LogRocketMQ, 0)
	for _, a := range arrI {
		m, err := cast.ToStringMapE(a)
		if err != nil {
			continue
		}
		topic, ok := m["topic"]
		if !ok {
			continue
		}
		group, ok := m["group"]
		if !ok {
			continue
		}
		addressI, ok := m["address"]
		if !ok {
			continue
		}
		address, err := cast.ToStringSliceE(addressI)
		if err != nil {
			continue
		}
		level, err := toLevel(m["level"])
		if err != nil {
			continue
		}
		lrs = append(lrs, LogRocketMQ{
			Topic:   cast.ToString(topic),
			Group:   cast.ToString(group),
			Address: address,
			Level:   &level,
			fields:  getFields(m["fields"]),
		})
	}

	if len(lrs) == 0 {
		return nil, errors.New("not found rocketmq config")
	}

	return lrs, nil
}

func toMail(data interface{}) ([]LogMail, error) {
	arrI, err := cast.ToSliceE(data)
	if err != nil {
		return nil, errors.New("type error")
	}

	lms := make([]LogMail, 0)
	for _, a := range arrI {
		m, err := cast.ToStringMapE(a)
		if err != nil {
			continue
		}
		level, err := toLevel(m["level"])
		if err != nil {
			continue
		}
		from, ok := m["from"]
		if !ok {
			continue
		}
		to, ok := m["to"]
		if !ok {
			continue
		}
		subject, ok := m["subject"]
		if !ok {
			continue
		}
		stmp, ok := m["stmp"]
		if !ok {
			continue
		}
		port := cast.ToInt(m["port"])
		if port == 0 {
			continue
		}
		password := cast.ToString(m["password"])
		lms = append(lms, LogMail{
			Level:    &level,
			From:     cast.ToString(from),
			To:       cast.ToString(to),
			Subject:  cast.ToString(subject),
			Stmp:     cast.ToString(stmp),
			Port:     port,
			Password: password,
			fields:   getFields(m["fields"]),
		})
	}

	if len(lms) == 0 {
		return nil, errors.New("not found mail config")
	}
	return lms, nil
}

func toDingding(data interface{}) (*LogDingding, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("not found dingding config")
	}

	level, err := toLevel(m["level"])
	if err != nil {
		level = InfoLevel
	}

	access_token, ok := m["access_token"]
	if !ok {
		return nil, errors.New("not found access_token")
	}

	secret, ok := m["secret"]
	if !ok {
		return nil, errors.New("not found secret")
	}

	fields := getFields(m["fields"])

	return &LogDingding{
		level:        &level,
		secret:       cast.ToString(secret),
		access_token: cast.ToString(access_token),
		fields:       fields,
	}, nil
}

func getFields(data interface{}) []zap.Field {
	fields := []zap.Field{}

	mfield := cast.ToStringMapString(data)
	for k, v := range mfield {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}

func formatFileMbSize(fileSize int64) (size int64) {
	return fileSize * 1024 * 1024
}

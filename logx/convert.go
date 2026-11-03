package logx

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func toLevel(data interface{}) (Level, error) {
	switch data.(string) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "panic":
		return LevelPanic, nil
	}
	return LevelInfo, errors.New("not found level")
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

func getFields(data interface{}) []zap.Field {
	fields := []zap.Field{}

	mfield := cast.ToStringMapString(data)
	for k, v := range mfield {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}

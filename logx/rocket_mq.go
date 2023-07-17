package logx

import (
	"github.com/qida/gohp/mqtt/rocketmq"
	"github.com/qida/gohp/slice"
	"go.uber.org/zap/zapcore"
)

type rocketCore struct {
	zapcore.LevelEnabler
	enc zapcore.Encoder
	//out zapcore.WriteSyncer
	MqProducer *rocketmq.MqProducer
	MsgTags
}

type MsgTags struct {
	topicName string
	key       string
	tag       string
}

// NewRocketMQCore creates a Core that writes logs to a WriteSyncer.
func NewRocketMQCore(enc zapcore.Encoder, topic string, MqProducer *rocketmq.MqProducer, enab zapcore.LevelEnabler) zapcore.Core {
	return &rocketCore{
		LevelEnabler: enab,
		enc:          enc,
		MqProducer:   MqProducer,
		MsgTags: MsgTags{
			topicName: topic,
			key:       "",
			tag:       "",
		},
	}
}

// Level 日志级别
func (c *rocketCore) Level() zapcore.Level {
	return zapcore.LevelOf(c.LevelEnabler)
}

func (c *rocketCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields, c)
	return clone
}

func (c *rocketCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *rocketCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	//_, err = c.out.Write(buf.Bytes())

	valTag, err := slice.SliceStructPop[zapcore.Field, string](fields, "Key", "tag")
	if err != nil {
		return err
	}

	if len(valTag) >= 1 {
		c.tag = valTag[0].Interface.(string)
	} else {
		c.tag = ""
		c.key = ""
	}

	MqMsg := rocketmq.NewMQMessage(c.topicName, c.tag, c.key, buf.Bytes(), 0)
	_, err = c.MqProducer.SendAsyncMessageBody(MqMsg)
	buf.Free()
	if err != nil {
		return err
	}

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		//c.Sync()
	}
	return nil
}

func (c *rocketCore) Sync() error {
	//return c.out.Sync()
	return nil
}

func (c *rocketCore) clone() *rocketCore {
	return &rocketCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		MqProducer:   c.MqProducer,
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field, c *rocketCore) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

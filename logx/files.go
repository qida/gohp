package logx

import (
	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type FilesLogger struct {
	Path       string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	IsUTC      bool
	IsCompress bool
}

func getFilesWriter(fl FilesLogger) (io.Writer, error) {
	if fl.MaxAge == 0 {
		fl.MaxAge = 7
	}

	log := lumberjack.Logger{
		Filename:   fl.Path,
		MaxSize:    fl.MaxSize,
		MaxAge:     fl.MaxAge,
		MaxBackups: fl.MaxBackups,
		LocalTime:  !fl.IsUTC,
		Compress:   fl.IsCompress,
	}
	c := cron.New()
	// c.AddFunc("* * * * * *", func() { lj_log.Rotate() })
	c.AddFunc("@daily", func() { log.Rotate() })
	c.Start()

	//hook, err := rotatelogs.New(
	//	strings.Replace(fl.Path, ".log", "", -1)+"-%Y%m%d.log",
	//	rotatelogs.WithLinkName(fl.Path),
	//	rotatelogs.WithMaxAge(fl.MaxAge),
	//	rotatelogs.WithRotationSize(fl.MaxSize),
	//	rotatelogs.WithRotationTime(fl.RotationTime),
	//)
	//
	//if err != nil {
	//	return nil, err
	//}
	//return hook, nil
	return &log, nil
}

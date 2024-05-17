package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	/*
		debugLevel := logger.DebugLevel
		infoLevel := logger.InfoLevel
		dpanicLevel := logger.DPanicLevel
	*/
	//logger.Debug("defaultLogger", zap.String("stringField", "hahahaha"))
	/*
		log, err := logger.New(
			logger.WithLevel(logger.DebugLevel),
			logger.WithCaller(),
			logger.WithFile([]logger.Filelogger{
				logger.Filelogger{
					Path:  "./logs/debug.log",
					Level: &debugLevel,
				},
				logger.Filelogger{
					Path:  "./logs/info.log",
					Level: &infoLevel,
				},
			}),
			logger.WithKafka([]logger.LogKafka{
				logger.LogKafka{
					Topic:   "optdev",
					Address: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
					Level:   &debugLevel,
				},
			}),
			logger.WithMail([]logger.LogMail{
				logger.LogMail{
					Level:    &dpanicLevel,
					From:     "zls3434@qq.com",
					To:       "zls3434@126.com",
					Subject:  "服务预警",
					Stmp:     "smtp.qq.com",
					Port:     465,
					Password: "1111111111",
				},
			}),
			logger.WithInitialFields("app", "mediaserver"),
		)
	*/
	//jsonStr := `{
	//    "caller":true,
	//    "dingding":{
	//        "access_token":"f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b",
	//        "fields":{
	//            "location":"gd-al-01",
	//            "service":"mediaserver"
	//        },
	//        "level":"dpanic",
	//        "location":"gd-al-01",
	//        "secret":"SECf6f2ea8f45d8a057c9566a33f99474da2e5c6a6604d736121650e2730c6fb0a3",
	//        "service":"mediaserver"
	//    },
	//    "file":[
	//        {
	//            "level":"debug",
	//            "path":"./logs/debug.log"
	//        },
	//        {
	//            "level":"info",
	//            "path":"./logs/info.log"
	//        }
	//    ],
	//    "kafka":[
	//        {
	//            "address":[
	//                "127.0.0.1:9092",
	//                "127.0.0.1:9093",
	//                "127.0.0.1:9094"
	//            ],
	//            "fields":{
	//                "location":"gd-al-01",
	//                "service":"mediaserver"
	//            },
	//            "level":"info",
	//            "topic":"optdev"
	//        }
	//    ],
	//    "level":"info",
	//    "mail":[
	//        {
	//            "from":"zls3434@qq.com",
	//            "level":"dpanic",
	//            "password":"1111111111",
	//            "port":"465",
	//            "stmp":"smtp.qq.com",
	//            "subject":"服务预警",
	//            "to":"zls3434@126.com",
	//            "fields":{
	//                "location":"gd-al-01",
	//                "service":"mediaserver"
	//            }
	//        }
	//    ]
	//}`

	// TODO: DPanic will send email to administrator
	//log.DPanic("this is panic", zap.Errors("error", []error{errors.New("error with div zero")}))
	//log.Info("hi i am from debug", zap.Any("白居易", baijuyi))
	//log.Info("hi i am from info")
	//log.Infof("hi this is infof with args %s  %d", "dididi", 123)
	//log.Info("this is string field", zap.String("stringField", "optdev"))
	//log.Info("this is any field", zap.Any("anyField", map[string]string{
	//	"name": "optdev",
	//	"age":  "30",
	//}))
	//for i := 0; i < 10000; i++ {
	//	log.Info("this is any field", zap.Any("anyField", map[string]string{
	//		"name": "optdev",
	//		"age":  "30",
	//	}))
	//}
	time.Sleep(5 * time.Second)
	//test1()
}

func test1() {
	addr, err := primitive.NewNamesrvAddr("192.168.114.50:9876")
	if err != nil {
		panic(err)
	}
	topic := "broker-a"
	p, err := rocketmq.NewProducer(
		producer.WithGroupName("my_service"),
		producer.WithNameServer(addr),
		producer.WithCreateTopicKey(topic),
		producer.WithRetry(1))
	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}

	// 发送异步消息
	res, err := p.SendSync(context.Background(), primitive.NewMessage(topic, []byte(baijuyi)))
	if err != nil {
		fmt.Printf("send sync message error:%s\n", err)
	} else {
		fmt.Printf("send sync message success. result=%s\n", res.String())
	}

	// 发送消息后回调
	//err = p.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
	//	if err != nil {
	//		fmt.Printf("receive message error:%v\n", err)
	//	} else {
	//		fmt.Printf("send message success. result=%s\n", result.String())
	//	}
	//}, primitive.NewMessage(topic, []byte("send async message")))
	//if err != nil {
	//	fmt.Printf("send async message error:%s\n", err)
	//}

	//// 批量发送消息
	//var msgs []*primitive.Message
	//for i := 0; i < 5; i++ {
	//	msgs = append(msgs, primitive.NewMessage(topic, []byte("batch send message. num:"+strconv.Itoa(i))))
	//}
	//res, err = p.SendSync(context.Background(), msgs...)
	//if err != nil {
	//	fmt.Printf("batch send sync message error:%s\n", err)
	//} else {
	//	fmt.Printf("batch send sync message success. result=%s\n", res.String())
	//}

	//// 发送延迟消息
	//msg := primitive.NewMessage(topic, []byte("delay send message"))
	//msg.WithDelayTimeLevel(3)
	//res, err = p.SendSync(context.Background(), msg)
	//if err != nil {
	//	fmt.Printf("delay send sync message error:%s\n", err)
	//} else {
	//	fmt.Printf("delay send sync message success. result=%s\n", res.String())
	//}

	//// 发送带有tag的消息
	//msg1 := primitive.NewMessage(topic, []byte("send tag message"))
	//msg1.WithTag("tagA")
	//res, err = p.SendSync(context.Background(), msg1)
	//if err != nil {
	//	fmt.Printf("send tag sync message error:%s\n", err)
	//} else {
	//	fmt.Printf("send tag sync message success. result=%s\n", res.String())
	//}
	//
	err = p.Shutdown()
	if err != nil {
		panic(err)
	}
}

var baijuyi = `白居易 〔唐代〕

汉皇重色思倾国，御宇多年求不得。
杨家有女初长成，养在深闺人未识。
天生丽质难自弃，一朝选在君王侧。
回眸一笑百媚生，六宫粉黛无颜色。
春寒赐浴华清池，温泉水滑洗凝脂。
侍儿扶起娇无力，始是新承恩泽时。
云鬓花颜金步摇，芙蓉帐暖度春宵。
春宵苦短日高起，从此君王不早朝。
承欢侍宴无闲暇，春从春游夜专夜。
后宫佳丽三千人，三千宠爱在一身。
金屋妆成娇侍夜，玉楼宴罢醉和春。
姊妹弟兄皆列土，可怜光彩生门户。
遂令天下父母心，不重生男重生女。
骊宫高处入青云，仙乐风飘处处闻。
缓歌慢舞凝丝竹，尽日君王看不足。

渔阳鼙鼓动地来，惊破霓裳羽衣曲。
九重城阙烟尘生，千乘万骑西南行。
翠华摇摇行复止，西出都门百余里。
六军不发无奈何，宛转蛾眉马前死。
花钿委地无人收，翠翘金雀玉搔头。
君王掩面救不得，回看血泪相和流。
黄埃散漫风萧索，云栈萦纡登剑阁。
峨嵋山下少人行，旌旗无光日色薄。
蜀江水碧蜀山青，圣主朝朝暮暮情。
行宫见月伤心色，夜雨闻铃肠断声。
天旋地转回龙驭，到此踌躇不能去。(地转 一作：日转)
马嵬坡下泥土中，不见玉颜空死处。
君臣相顾尽沾衣，东望都门信马归。
归来池苑皆依旧，太液芙蓉未央柳。
芙蓉如面柳如眉，对此如何不泪垂？
春风桃李花开日，秋雨梧桐叶落时。(花开日 一作：花开夜)
西宫南内多秋草，落叶满阶红不扫。(南内 一作：南苑)
梨园弟子白发新，椒房阿监青娥老。
夕殿萤飞思悄然，孤灯挑尽未成眠。
迟迟钟鼓初长夜，耿耿星河欲曙天。
鸳鸯瓦冷霜华重，翡翠衾寒谁与共？
悠悠生死别经年，魂魄不曾来入梦。

临邛道士鸿都客，能以精诚致魂魄。
为感君王辗转思，遂教方士殷勤觅。
排空驭气奔如电，升天入地求之遍。
上穷碧落下黄泉，两处茫茫皆不见。
忽闻海上有仙山，山在虚无缥缈间。
楼阁玲珑五云起，其中绰约多仙子。
中有一人字太真，雪肤花貌参差是。
金阙西厢叩玉扃，转教小玉报双成。
闻道汉家天子使，九华帐里梦魂惊。
揽衣推枕起徘徊，珠箔银屏迤逦开。
云鬓半偏新睡觉，花冠不整下堂来。

风吹仙袂飘飖举，犹似霓裳羽衣舞。(飘飖 一作：飘飘)
玉容寂寞泪阑干，梨花一枝春带雨。(阑 通：栏)
含情凝睇谢君王，一别音容两渺茫。
昭阳殿里恩爱绝，蓬莱宫中日月长。
回头下望人寰处，不见长安见尘雾。
惟将旧物表深情，钿合金钗寄将去。
钗留一股合一扇，钗擘黄金合分钿。
但令心似金钿坚，天上人间会相见。(但令 一作：但教)

临别殷勤重寄词，词中有誓两心知。
七月七日长生殿，夜半无人私语时。
在天愿作比翼鸟，在地愿为连理枝。
天长地久有时尽，此恨绵绵无绝期。
`

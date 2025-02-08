package idx

import (
	"sync"
	"time"
)

// var a= "{\"b\":1234567890123456789}";var b = JSON.parse(a);console.log(b);b.b=1234567890123456789;var c = JSON.stringify(b);console.log(c);
// var a= "{\"b\":999999999999999}";var b = JSON.parse(a);console.log(b);b.b=999999999999998;var c = JSON.stringify(b);console.log(c);
// 因为snowFlake目的是解决分布式下生成唯一id 所以ID中是包含集群和节点编号在内的
// JS 支持的最大数 999999999999999  15位
// 129511438879297536

// 定义一个woker工作节点所需要的基本参数
type SnowFlakeJS struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录时间戳
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// 实例化一个工作节点
// 0 < id_worker < 1024
func NewSnowFlakeJS() *SnowFlakeJS {
	return &SnowFlakeJS{
		timestamp: 0,
		number:    0,
	}
}

// 接下来我们开始生成id
// 生成方法一定要挂载在某个woker下，这样逻辑会比较清晰 指定某个节点生成id
func (t *SnowFlakeJS) GetId() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	// 获取生成时的时间戳
	now := time.Now().UnixMilli() // 毫秒
	if t.timestamp == now {
		t.number++
		// 这里要判断，当前工作节点是否在1毫秒内已经生成numberMax个ID
		if t.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= t.timestamp {
				now = time.Now().UnixMilli()
			}
			t.timestamp = now // 更新时间戳为新的毫秒数
			t.number = 0      // 重置序列号
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		t.number = 0
		t.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}
	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	// fmt.Printf("now: %d , %d %d ", now, now*100, t.number)
	ID := int64(now*100 + t.number)
	return ID
}

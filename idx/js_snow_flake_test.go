package idx

import (
	"fmt"
	"testing"
)

func TestSnowFlakeJS_GetId(t *testing.T) {
	snowFlake := NewSnowFlakeJS()
	ids := make(map[int64]struct{})
	for i := 0; i < 1000; i++ {
		id := snowFlake.GetId()
		if _, ok := ids[id]; ok {
			t.Errorf("SnowFlakeJS.GetId() 生成了重复ID: %d (第 %d 次)", id, i)
		}
		ids[id] = struct{}{}
	}
	if len(ids) != 1000 {
		t.Errorf("期望 1000 个唯一ID, 实际 %d 个", len(ids))
	}
}

func TestSnowFlakeJS_GetId_Concurrent(t *testing.T) {
	snowFlake := NewSnowFlakeJS()
	const goroutines = 10
	const perGoroutine = 100

	idCh := make(chan int64, goroutines*perGoroutine)
	for g := 0; g < goroutines; g++ {
		go func() {
			for i := 0; i < perGoroutine; i++ {
				idCh <- snowFlake.GetId()
			}
		}()
	}

	ids := make(map[int64]struct{})
	for i := 0; i < goroutines*perGoroutine; i++ {
		id := <-idCh
		if _, ok := ids[id]; ok {
			t.Errorf("并发场景下生成了重复ID: %d", id)
		}
		ids[id] = struct{}{}
	}
	if len(ids) != goroutines*perGoroutine {
		t.Errorf("期望 %d 个唯一ID, 实际 %d 个", goroutines*perGoroutine, len(ids))
	}
	fmt.Printf("SnowFlakeJS 并发生成 %d 个唯一ID\n", len(ids))
}

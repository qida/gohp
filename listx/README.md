# listx

链表队列工具，基于标准库 `container/list` 封装的线程安全（带互斥锁）队列实现。包名为 `lists`。

## 功能列表

### lists.go
- `Queue`：线程安全的双向链表队列（内部使用 `sync.Mutex` 保护）
- `NewQueue() *Queue`：创建队列
- `(*Queue) PushBack(v interface{})`：元素入队尾（nil 元素被忽略）
- `(*Queue) Front() *list.Element`：获取队首元素（不出队）
- `(*Queue) Remove(e *list.Element)`：移除指定元素（常配合 Front 使用实现出队）
- `(*Queue) Len() int`：返回队列长度

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/listx"
)

func main() {
	q := listx.NewQueue()
	q.PushBack(1)
	q.PushBack(2)

	// 出队：Front + Remove
	for e := q.Front(); e != nil; e = q.Front() {
		fmt.Println(e.Value) // 依次输出 1 2
		q.Remove(e)
	}
	fmt.Println(q.Len()) // 0
}
```

## 备注

- `Front` 仅获取队首元素而不移除，需调用 `Remove` 完成出队。
- 所有方法内部加锁，可在多 goroutine 中安全使用。

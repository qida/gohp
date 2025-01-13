package timex

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type ITimer interface {
	AddTaskByFunc(taskName string, spec string, task func(), option ...cron.Option) (cron.EntryID, error)
	AddTaskByJob(taskName string, spec string, job interface{ Run() }, option ...cron.Option) (cron.EntryID, error)
	FindCron(taskName string) (*cron.Cron, bool)
	StartTask(taskName string)
	StopTask(taskName string)
	Remove(taskName string, id int)
	Clear(taskName string)
	Close()
}

// 定时任务管理
type LcTimer struct {
	TaskList map[string]*cron.Cron
	sync.Mutex
}

// 通过函数创建方法添加任务
func (t *LcTimer) AddTaskByFunc(taskName string, spec string, task func(), option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddFunc(spec, task)
	t.TaskList[taskName].Start()
	return id, err
}

// 通过接口方式批量添加任务
func (t *LcTimer) AddTaskByJob(taskName string, spec string, job interface{ Run() }, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TaskList[taskName]; !ok {
		t.TaskList[taskName] = cron.New(option...)
	}
	id, err := t.TaskList[taskName].AddJob(spec, job)
	t.TaskList[taskName].Start()
	return id, err
}

// 根据名称获取对应的cron
func (t *LcTimer) FindCron(taskName string) (*cron.Cron, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.TaskList[taskName]
	return v, ok
}

// 开始定时任务
func (t *LcTimer) StartTask(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Start()
	}
}

// 停止定时任务
func (t *LcTimer) StopTask(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Stop()
	}
}

// 移除指定的定时任务
func (t *LcTimer) Remove(taskName string, id int) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Remove(cron.EntryID(id))
	}
}
func (t *LcTimer) Clear(taskName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.TaskList[taskName]; ok {
		v.Stop()
		delete(t.TaskList, taskName)
	}
}
func (t *LcTimer) Close() {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.TaskList {
		v.Stop()
	}
}

func NewTimerTask() ITimer {
	return &LcTimer{
		TaskList: make(map[string]*cron.Cron),
	}
}

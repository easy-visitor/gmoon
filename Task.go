package gmoon

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type TaskFunc func(params ...interface{})

var once sync.Once
var taskList chan *TaskExecutor //任务列表

var onceCron sync.Once

var taskCronList *cron.Cron

//初始化任务信息
func init() {
	chList := getTaskList()
	go func() {
		for t := range chList {
			doTask(t)
		}
	}()
}

//执行任务
func doTask(t *TaskExecutor) {
	go func() {
		defer func() {
			//如果存在任务 需要去执行改任务信息
			if t.callback != nil {
				t.callback()
			}
		}()
		t.Exec() //执行任务
	}()

}


func getCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCronList = cron.New(cron.WithSeconds())
	})
	return taskCronList
}

//获取当前任务
func getTaskList() chan *TaskExecutor {
	once.Do(func() {
		taskList = make(chan *TaskExecutor, 0)
	})
	return taskList
}

type TaskExecutor struct {
	f        TaskFunc
	callback func()
	p        []interface{}
}

func NewTaskExecutor(f TaskFunc, callback func(), p []interface{}) *TaskExecutor {
	return &TaskExecutor{f: f, callback: callback, p: p}
}

func (this *TaskExecutor) Exec() { //执行任务
	this.f(this.p...)
}

//调用任务
func Task(f TaskFunc, callback func(), params ...interface{}) {
	if f == nil {
		return
	}
	go func() {
		getTaskList() <- NewTaskExecutor(f, callback, params) //增加任务队列
	}()

}

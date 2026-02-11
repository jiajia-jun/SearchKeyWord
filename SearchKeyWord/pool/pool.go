// Package pool 实现目标：1.储存Task(Put)  2.分发 Task(Arrange)  3.等待所有Task完成(Wait)
package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	mutex       sync.Mutex
	taskChannel chan func()
	waitGroup   sync.WaitGroup
}

// NewPool 创建协程池
func NewPool() *Pool {
	return &Pool{}
}

// Put 储存目标数目的Task
func (p *Pool) Put(task func(), tasknum int) {

	p.mutex.Lock() // 加锁保证初始化唯一
	defer p.mutex.Unlock()

	if p.taskChannel == nil { //初始化Pool任务通道
		p.taskChannel = make(chan func(), tasknum+10)
	}
	p.waitGroup.Add(tasknum) //等待组加入任务数目
	for i := 1; i <= tasknum; i++ {
		t := task
		p.taskChannel <- t
	}
}

// Arrange 对Task进行分发协程
func (p *Pool) Arrange(goroutines int) {
	for i := 1; i <= goroutines; i++ {
		go func() {
			for task := range p.taskChannel {
				task() //互斥锁不该出现在这里，理论上传入的Task里面就应该有互斥锁保护，锁的业务让Task自己实现
				p.waitGroup.Done()
			}
		}()
	}
	close(p.taskChannel) //关闭任务通道
}

// Wait 等待Task全部结束
func (p *Pool) Wait() {
	p.waitGroup.Wait()
	fmt.Println("Pool Closed")
}

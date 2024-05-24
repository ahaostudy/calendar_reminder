package crontab

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	gs  = 20
	mws = 200

	taskCount      = 1000
	ready          = time.Duration(100)
	maxWaitSeconds = 5

	start  = time.Now().Add(time.Millisecond * ready)
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestCrontab(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(taskCount)
	doneCounter := int32(0)
	cancelCounter := int32(0)

	var tasks []*Task
	for i := 0; i < taskCount; i++ {
		// generate random execution time
		expectTime := start.Add(time.Second * time.Duration(random.Intn(maxWaitSeconds)))
		task := NewTask(strconv.Itoa(i), expectTime, nil, func(id string, _ any) {
			defer func() {
				atomic.AddInt32(&doneCounter, 1)
				wg.Done()
			}()
			// compare the actual execution time with the expected time at each task execution
			// limit the time difference here to no more than 10 ms
			actualTime := time.Now()
			diff := diffMilli(actualTime, expectTime)
			assert.LessOrEqual(t, diff, 10)
			//t.Logf("run task: %v, actual time: %v, expect time: %v, millisecond difference: %v\n", id, actualTime, expectTime, diff)
		})
		tasks = append(tasks, task)
	}

	// create a scheduled task scheduler
	c := NewCrontab(gs, mws)
	c.Run()
	for _, task := range tasks {
		go func(task Task) {
			// add a timed task each time you wait for a random amount of time
			time.Sleep(time.Millisecond * time.Duration(random.Intn(100)))
			c.AddTask(&task)
		}(*task)
	}

	// test delete some tasks
	for i := 0; i < taskCount/4; i++ {
		idx := random.Intn(len(tasks))
		go func(task Task) {
			time.Sleep(time.Millisecond * time.Duration(random.Intn(maxWaitSeconds*1000)))
			if c.DeleteTask(task.Key()) {
				atomic.AddInt32(&cancelCounter, 1)
				wg.Done()
				//t.Logf("delete task: %v\n", task.ID)
			}
		}(*tasks[idx])
		// delete the task
		*tasks[idx] = *tasks[len(tasks)-1]
		tasks = tasks[:len(tasks)-1]
	}

	// wait for all tasks to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	timer := time.NewTimer(time.Duration(maxWaitSeconds+1) * time.Second)
	select {
	case <-timer.C:
	case <-done:
		timer.Stop()
	}

	// determine if the number of executions is as expected
	assert.Equal(t, int32(taskCount), atomic.LoadInt32(&doneCounter)+atomic.LoadInt32(&cancelCounter))
}

func diffMilli(actual, expect time.Time) int {
	diff := int(actual.Sub(expect).Milliseconds())
	if diff < 0 {
		return -diff
	}
	return diff
}

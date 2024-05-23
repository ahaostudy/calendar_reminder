package crontab

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	gs  = 20
	mws = 200

	maxSeconds = 12
	taskCount  = 1000

	now = time.Now()
	r   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestCrontab(t *testing.T) {
	var tasks []*Task
	for i := 0; i < taskCount; i++ {
		fi := i
		dt := now.Add(time.Second * time.Duration(r.Intn(maxSeconds)))
		tasks = append(tasks, NewTask(dt, fi, func(data any) {
			actual := time.Now()
			expect := dt
			assert.LessOrEqual(t, diff(actual, expect), 1)
			//t.Logf("run %v, actual time: %v, expect time: %v\n", data, actual, expect)
		}))
	}

	c := NewCrontab(gs, mws)
	c.Run()
	for _, task := range tasks {
		go func(task Task) {
			time.Sleep(time.Millisecond * time.Duration(r.Intn(1000)))
			c.AddTask(&task)
		}(*task)
	}

	time.Sleep(time.Duration(maxSeconds) * time.Second)
}

func diff(actual, expect time.Time) int {
	if actual.After(expect) {
		return int(actual.Unix() - expect.Unix())
	}
	return int(expect.Unix() - actual.Unix())
}

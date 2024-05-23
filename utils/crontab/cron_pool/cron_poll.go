package cron_pool

// CronPoll is a coroutine pool used to execute scheduled tasks
// it is only executed in the order in which tasks are received, but supports task interruption and gives up resources
type CronPoll struct {
	Gos      int
	MaxWorks int
	works    chan *Work
	workers  chan *Worker
}

func NewCronPoll(gs, mws int) *CronPoll {
	c := &CronPoll{Gos: gs, MaxWorks: mws}
	c.workers = make(chan *Worker, gs)
	c.works = make(chan *Work, mws)
	return c
}

func (c *CronPoll) Run() {
	for i := 0; i < c.Gos; i++ {
		w := NewWorker(c)
		go w.Run()
	}

	go func() {
		for {
			work := <-c.works
			worker := <-c.workers
			worker.Push(work)
		}
	}()
}

func (c *CronPoll) AddWork(work *Work) {
	c.works <- work
}

package cron_pool

type Worker struct {
	WorkChan      chan *Work
	InterruptChan chan struct{}

	cp *CronPoll
}

func NewWorker(cp *CronPoll) *Worker {
	return &Worker{cp: cp, WorkChan: make(chan *Work), InterruptChan: make(chan struct{})}
}

func (w Worker) Run() {
	for {
		w.ready()

		work := <-w.WorkChan
		work.worker = &w

		timer := work.Wait()

		select {
		case <-timer.C:
			work.Run()
		case <-w.InterruptChan:
			timer.Stop()
		}
	}
}

func (w Worker) ready() {
	w.cp.workers <- &w
}

func (w Worker) Push(work *Work) {
	w.WorkChan <- work
}

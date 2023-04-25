package task_driver

import (
	"context"
	"github.com/pefish/go-logger"
	"sync"
	"time"
)

type Runner interface {
	Run(ctx context.Context) error
	Stop() error
	GetName() string
	GetInterval() time.Duration
}

type TaskDriver struct {
	runners   []Runner
	waitGroup sync.WaitGroup
	logger    go_logger.InterfaceLogger
}

func NewTaskDriver() *TaskDriver {
	return &TaskDriver{
		runners:   make([]Runner, 0),
		waitGroup: sync.WaitGroup{},
	}
}

func (driver *TaskDriver) SetLogger(logger go_logger.InterfaceLogger) {
	driver.logger = logger
}

func (driver *TaskDriver) Register(runner Runner) {
	driver.runners = append(driver.runners, runner)
}

func (driver *TaskDriver) RunWait(ctx context.Context) {
	for _, runner := range driver.runners {
		driver.waitGroup.Add(1)
		go func(runner Runner) {
			defer driver.waitGroup.Done()
			if runner.GetInterval() == 0 {
				err := runner.Run(ctx)
				if err != nil {
					driver.logger.ErrorF("%#v", err)
				}
				runner.Stop()
				return
			} else {
				timer := time.NewTimer(0)
				for {
					select {
					case <-timer.C:
						err := runner.Run(ctx)
						if err != nil {
							driver.logger.ErrorF("%#v", err)
						}
						timer.Reset(runner.GetInterval())
					case <-ctx.Done():
						runner.Stop()
						return
					}
				}
			}

		}(runner)
	}
	driver.waitGroup.Wait()
}

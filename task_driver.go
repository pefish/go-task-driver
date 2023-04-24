package task_driver

import (
	"context"
	"github.com/pefish/go-logger"
	"sync"
)

type Runner interface {
	Run(ctx context.Context) error
	Stop() error
	GetName() string
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
			err := runner.Run(ctx)
			if err != nil {
				driver.logger.ErrorF("%#v", err)
			}
			runner.Stop()
		}(runner)
	}
	driver.waitGroup.Wait()
}

package task_driver

import (
	"context"
	"sync"
	"time"

	i_logger "github.com/pefish/go-interface/i-logger"
)

type Runner interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Stop() error
	Name() string
	Interval() time.Duration
	Logger() i_logger.ILogger
}

type TaskDriver struct {
	runners   []Runner
	waitGroup sync.WaitGroup
}

func NewTaskDriver() *TaskDriver {
	return &TaskDriver{
		runners:   make([]Runner, 0),
		waitGroup: sync.WaitGroup{},
	}
}

func (driver *TaskDriver) Register(runner Runner) {
	driver.runners = append(driver.runners, runner)
}

func (driver *TaskDriver) RunWait(ctx context.Context) {
	for _, runner := range driver.runners {
		err := runner.Init(ctx)
		if err != nil {
			runner.Logger().ErrorF("Init failed. err: %+v", err)
			continue
		}
		driver.waitGroup.Add(1)
		go func(runner Runner) {
			defer driver.waitGroup.Done()
			if runner.Interval() == 0 {
				err := runner.Run(ctx)
				if err != nil {
					runner.Logger().ErrorF("%#v", err)
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
							runner.Logger().ErrorF("%+v", err)
						}
						timer.Reset(runner.Interval())
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

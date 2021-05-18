package task_driver

import (
	"github.com/pefish/go-logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner interface {
	// run in goroutine
	Run() error
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
		logger:    go_logger.DefaultLogger,
	}
}

func (driver *TaskDriver) SetLogger(logger go_logger.InterfaceLogger) {
	driver.logger = logger
}

func (driver *TaskDriver) Register(runner Runner) {
	driver.runners = append(driver.runners, runner)
}

func (driver *TaskDriver) RunWait(exit <- chan struct{}) {
	finished := make(chan bool)
	go func() {
		for _, runner := range driver.runners {
			driver.waitGroup.Add(1)
			go func(runner Runner) {
				defer driver.waitGroup.Done()
				err := runner.Run()
				if err != nil {
					driver.logger.ErrorF("[%s]: %#v", runner.GetName(), err)
				}
			}(runner)
		}
		driver.waitGroup.Wait()
		close(finished)
	}()

	if exit != nil {
		select {
		case <-exit:
		case <-finished:
		}
	} else {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalChan:
		case <-finished:
		}
	}

	driver.stopWait()
}

func (driver *TaskDriver) stopWait() {
	for _, runner := range driver.runners {
		driver.logger.InfoF("[%s]: stopping...", runner.GetName())
		err := runner.Stop()
		if err != nil {
			driver.logger.ErrorF("[%s]: %#v", runner.GetName(), err)
		}
		driver.logger.InfoF("[%s]: stopped", runner.GetName())
	}
}


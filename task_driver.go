package task_driver

import (
	"github.com/pefish/go-logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner interface {
	// run in goroutine, do not exit
	Run() error
	Stop() error
	GetName() string
	GetLogger() go_logger.InterfaceLogger
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

func (driver *TaskDriver) RunWait(exit <- chan struct{}) {
	finished := make(chan bool)
	go func() {
		for _, runner := range driver.runners {
			driver.waitGroup.Add(1)
			go func(runner Runner) {
				defer driver.waitGroup.Done()
				err := runner.Run()
				if err != nil {
					runner.GetLogger().ErrorF("%#v", err)
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
		runner.GetLogger().Info("stopping...")
		err := runner.Stop()
		if err != nil {
			runner.GetLogger().Error("%#v", err)
		}
		runner.GetLogger().Info("stopped")
	}
}


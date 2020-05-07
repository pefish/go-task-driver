package task_driver

import (
	"github.com/pefish/go-interface-logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner interface {
	// run in goroutine
	Run() error
	Stop() error
}

type TaskDriver struct {
	runners   map[string]Runner
	waitGroup sync.WaitGroup
	logger    go_interface_logger.InterfaceLogger
}

func NewTaskDriver() *TaskDriver {
	return &TaskDriver{
		runners:   map[string]Runner{},
		waitGroup: sync.WaitGroup{},
		logger:    go_interface_logger.DefaultLogger,
	}
}

func (driver *TaskDriver) SetLogger(logger go_interface_logger.InterfaceLogger) {
	driver.logger = logger
}

func (driver *TaskDriver) Register(name string, runner Runner) {
	driver.runners[name] = runner
}

func (driver *TaskDriver) RunWait() {
	finished := make(chan bool)
	go func() {
		for name, runner := range driver.runners {
			driver.waitGroup.Add(1)
			go func(runner Runner) {
				defer driver.waitGroup.Done()
				err := runner.Run()
				if err != nil {
					driver.logger.ErrorF("[%s]: %#v", name, err)
				}
			}(runner)
		}
		driver.waitGroup.Wait()
		close(finished)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
	case <-finished:
	}

	driver.stopWait()
}

func (driver *TaskDriver) stopWait() {
	for name, runner := range driver.runners {
		driver.logger.InfoF("[%s]: stopping...", name)
		err := runner.Stop()
		if err != nil {
			driver.logger.ErrorF("[%s]: %#v", name, err)
		}
		driver.logger.InfoF("[%s]: stopped", name)
	}
}

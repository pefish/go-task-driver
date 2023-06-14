package main

import (
	"context"
	"github.com/pefish/go-logger"
	task_driver "github.com/pefish/go-task-driver"
	"time"
)

func main() {
	driver := task_driver.NewTaskDriver()

	driver.Register(&Test{})

	ctx, cancel := context.WithCancel(context.Background())

	//exitChan := make(chan os.Signal)
	//signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	driver.RunWait(ctx)
}

type Test struct {
}

func (t *Test) GetLogger() go_logger.InterfaceLogger {
	return go_logger.Logger
}

func (t *Test) Stop() error {
	t.GetLogger().Info("stopped")
	return nil
}

func (t *Test) Run(ctx context.Context) error {
	for {
		t.GetLogger().Info(1)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (t *Test) GetName() string {
	return "hsgh"
}

func (t *Test) GetInterval() time.Duration {
	return 2 * time.Second
}

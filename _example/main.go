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

func (t *Test) Logger() go_logger.InterfaceLogger {
	return go_logger.Logger
}

func (t *Test) Stop() error {
	t.Logger().Info("stopped")
	return nil
}

func (t *Test) Init(ctx context.Context) error {
	t.Logger().Info("Inited.")
	return nil
}

func (t *Test) Run(ctx context.Context) error {
	for {
		t.Logger().Info(1)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (t *Test) Name() string {
	return "hsgh"
}

func (t *Test) Interval() time.Duration {
	return 2 * time.Second
}

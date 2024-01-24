package task_driver

import (
	"context"
	"fmt"
	go_logger "github.com/pefish/go-logger"
	"time"
)

func ExampleDriverType_Register() {
	driver := NewTaskDriver()

	driver.Register(&Test{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()
	driver.RunWait(ctx)

	// Output:
	// 1
	// 1
	// haha
	// xixi
}

type Test struct {
}

func (t *Test) Stop() error {
	fmt.Println("haha")
	time.Sleep(2 * time.Second)
	fmt.Println("xixi")
	return nil
}

func (t *Test) Init(ctx context.Context) error {
	return nil
}

func (t *Test) Run(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			fmt.Println(1)
			timer.Reset(2 * time.Second)
		}
	}
	return nil
}

func (t *Test) Name() string {
	return "hsgh"
}

func (t *Test) Interval() time.Duration {
	return 0
}

func (t *Test) Logger() go_logger.InterfaceLogger {
	return go_logger.Logger
}

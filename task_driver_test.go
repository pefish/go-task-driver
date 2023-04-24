package task_driver

import (
	"context"
	"fmt"
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

func (t *Test) GetName() string {
	return "hsgh"
}

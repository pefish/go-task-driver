package main

import (
	"context"
	"fmt"
	task_driver "github.com/pefish/go-task-driver"
	"time"
)

func main() {
	driver := task_driver.NewTaskDriver()

	driver.Register(&Test{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	driver.RunWait(ctx)
}

type Test struct {
}

func (t *Test) Stop() error {
	fmt.Println("stop")
	return nil
}

func (t *Test) Run(ctx context.Context) error {
	fmt.Println(1)
	return nil
}

func (t *Test) GetName() string {
	return "hsgh"
}

func (t *Test) GetInterval() time.Duration {
	return 2 * time.Second
}

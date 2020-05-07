package driver

import (
	"fmt"
	"time"
)

func ExampleDriverType_Register() {
	driver := NewDriver()

	driver.Register(`test`, &Test{})

	driver.RunWait()

	// Output:
	// a
	// [INFO] [test]: stopping...
	// haha
	// xixi
	// [INFO] [test]: stopped
}

type Test struct {

}


func (t *Test) Stop() error {
	fmt.Println("haha")
	time.Sleep(2 * time.Second)
	fmt.Println("xixi")
	return nil
}

func (t *Test) Run() error {
	fmt.Println(`a`)
	return nil
}


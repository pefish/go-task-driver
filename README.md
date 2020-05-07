# Go-task-driver

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/go-task-driver)

task driver for golang

## Quick start

```go
package task_driver

import (
	"fmt"
	"time"
)

func ExampleDriverType_Register() {
	driver := NewTaskDriver()

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


```


## Document

[Doc](https://godoc.org/github.com/pefish/go-task-driver)

## Security Vulnerabilities

If you discover a security vulnerability, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [Apache License](LICENSE).

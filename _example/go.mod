module _example

go 1.20

require (
	github.com/pefish/go-logger v0.5.5
	github.com/pefish/go-task-driver v0.2.2
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
)

replace github.com/pefish/go-task-driver => ../

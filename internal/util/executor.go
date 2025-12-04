package util

type Commander interface {
	Command() (args []string)
}

type Executor interface {
	Execute(Commander) ExecResult
}

type ExecResult interface {
	Code() int
	Stdout() []byte
	Stderr() []byte
	Error() error
}

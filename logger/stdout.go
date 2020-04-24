package loggers

import "fmt"

type StdOut struct{}

func (s *StdOut) SetMessage(msg string) {
	fmt.Println(msg)
}

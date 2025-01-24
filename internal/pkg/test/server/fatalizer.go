package server

import (
	"log"
	"os"
	"testing"
)

var logger = log.New(os.Stdout, "", 0)

type Fatalizer interface {
	Name() string
	Fatal(...any)
	Helper()
}

var _ Fatalizer = (testing.TB)(nil)

type fatalizer struct {
	name string
}

var _ Fatalizer = (*fatalizer)(nil)

func NewFatalizer(name string) *fatalizer {
	return &fatalizer{name: name}
}

func (f *fatalizer) Name() string {
	return f.name
}
func (*fatalizer) Fatal(a ...any) {
	logger.Fatal(a...)
}
func (*fatalizer) Helper() {}

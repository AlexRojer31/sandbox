package processes

import (
	"reflect"

	"github.com/AlexRojer31/sandbox/internal/config"
)

type IBuilder interface {
	Build(conf config.ChainConfig) IChain
}

type Builder struct{}

func (b *Builder) Build(conf config.ChainConfig) IChain {
	return NewChain(conf.Name, b.makeProcesses(conf.Name, conf.Processes))
}

func (b *Builder) makeProcesses(name string, names []string) []IProcess {
	var proc []IProcess
	pc := NewProcessCreator(name)
	for _, n := range names {
		process := b.call(pc, "Get"+n)
		if process != nil {
			switch v := process.(type) {
			case IProcess:
				proc = append(proc, v)
			}
		}
	}

	return proc
}

func (b *Builder) call(obj IProcessesCreator, methodName string, args ...any) any {
	e := reflect.ValueOf(obj)
	method := e.MethodByName(methodName)

	if !method.IsValid() {
		return nil
	}

	methodArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		methodArgs[i] = reflect.ValueOf(arg)
	}

	result := method.Call(methodArgs)
	if len(result) > 0 {
		return result[0].Interface()
	}

	return nil
}

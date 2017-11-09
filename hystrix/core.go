package hystrix

import (
	"errors"
)

type RESULT_TYPE int

const (
	SUCCESS RESULT_TYPE = 0
	FAILURE RESULT_TYPE = 1
	TIMEOUT RESULT_TYPE = 2
	REJECT  RESULT_TYPE = 3
)

type Command interface {
	Execute(args interface{}) RESULT_TYPE
	FallBack(args interface{})
}

var command map[string]Command
var circuit map[string]*Circuit

func init() {
	command = make(map[string]Command, 1000)
	circuit = make(map[string]*Circuit, 1000)
}

func RegisterCmd(name string, cmd Command) error {

	cmd, b := command[name]
	if b == true {
		return errors.New("cmd exist!" + name)
	}

	command[name] = cmd
	circuit[name] = NewCircuit(10, 50, 5)

	return nil
}

func UnRegisterCmd(name string) error {
	_, b := command[name]
	if b == false {
		return errors.New("cmd not exist!" + name)
	}

	command[name] = nil
	circuit[name] = nil

	return nil
}

func ExecuteCmd(name string, input interface{}) error {

	var result RESULT_TYPE

	cmd, b := command[name]
	if b == false {
		return errors.New("cmd not exist!" + name)
	}

	cir, b := circuit[name]
	if b == false {
		return errors.New("circuit not exist!" + name)
	}

	if cir.IsOpen() {
		cmd.FallBack(input)
	} else {
		result = cmd.Execute(input)
	}

	switch result {
	case SUCCESS:
		{
			cir.Success()
		}
	case FAILURE:
		{
			cir.Failure()
		}
	case TIMEOUT:
		{
			cir.Timeout()
		}
	case REJECT:
		{
			cir.Reject()
		}
	}

	return nil
}

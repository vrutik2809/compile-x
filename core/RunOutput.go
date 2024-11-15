package core

import "fmt"

type RunOutput struct {
	exitCode int
	output  string
}

func NewRunOutput(exitCode int, output string) RunOutput {
	return RunOutput{
		exitCode: exitCode,
		output: output,
	}
}

func (ro RunOutput) GetExitCode() int {
	return ro.exitCode
}

func (ro RunOutput) GetOutput() string {
	return ro.output
}

func (ro *RunOutput) SetExitCode(exitCode int) {
	ro.exitCode = exitCode
}

func (ro *RunOutput) SetOutput(output string) {
	ro.output = output
}

func (ro RunOutput) String() string {
	return fmt.Sprintf("{exitCode: %d, output: %s}", ro.exitCode, ro.output)
}
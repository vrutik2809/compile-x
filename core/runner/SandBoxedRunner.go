package runner

import (
	"github.com/docker/docker/client"
	"github.com/vrutik2809/compile-x/core"
	"github.com/vrutik2809/compile-x/core/executor"
)

type SandBoxedRunner struct {
	cli *client.Client
}

func NewSandBoxedRunner(cli *client.Client) SandBoxedRunner {
	return SandBoxedRunner{cli: cli}
}

func (sbr SandBoxedRunner) Run(sourceCode core.SourceCode) core.RunOutput {
	pool, err := core.GetContainerPool(sbr.cli,sourceCode.GetLanguage())
	if err != nil {
		return core.NewRunOutput(1, err.Error())
	}
	executor, err := executor.GetExecutor(sourceCode.GetLanguage())
	if err != nil {
		return core.NewRunOutput(1, err.Error())
	}
	runOutput := make(chan core.RunOutput)
	job := func(cli *client.Client,respId string) {
		result, err := executor.Execute(cli,respId, sourceCode.GetCode())
		if err != nil {
			runOutput <- core.NewRunOutput(1, err.Error())
			return
		}
		runOutput <- core.NewRunOutput(0, result)
	}
	pool.AddJob(job)
	return <-runOutput
}
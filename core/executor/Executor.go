package executor

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/vrutik2809/compile-x/core"
)

type Executor interface {
	Execute(cli *client.Client, respID string, code string) (string, error)
}

var mp map[core.Language]Executor

func GetExecutor(lang core.Language) (Executor,error) {
	executor, ok := mp[lang]
	if !ok {
		return nil, fmt.Errorf("Executor not found for language: %s", lang)
	}
	return executor, nil
}

func init() {
	mp = make(map[core.Language]Executor)
	mp[core.JAVA_22] = JavaExecutor{}
	mp[core.CPP_17_20] = CPPExecutor{}
	mp[core.PYTHON_3_12] = PythonExecutor{}
}
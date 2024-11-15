package runner

import "github.com/vrutik2809/compile-x/core"

type Runner interface {
	Run(sourceCode core.SourceCode) core.RunOutput
}
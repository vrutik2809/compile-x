package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type PythonExecutor struct {

}

func (je PythonExecutor) Execute(cli *client.Client, respID string, code string) (string, error) {
	// replace " with \" and ` with \` and replace \ with \\
	code = strings.ReplaceAll(code, `\`, `\\`)
	code = strings.ReplaceAll(code, `"`, `\"`)
	code = strings.ReplaceAll(code, "`", "\\`")
	runCmd := "timeout 1 python3 main.py"
	cmd := fmt.Sprintf("echo \"%s\" > main.py && %s",code, runCmd)
	execConfig := container.ExecOptions{
		Cmd:          []string{"bash", "-c", cmd},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}
	execResp, err := cli.ContainerExecCreate(context.Background(), respID, execConfig)
	if err != nil {
		return "", err
	}

	attachResp, err := cli.ContainerExecAttach(context.Background(), execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", err
	}
	defer attachResp.Close()

	resultCh := make(chan string)
	errCh := make(chan error)

	go func() {

		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer
		_, err = stdcopy.StdCopy(&outputBuffer, &errorBuffer, attachResp.Reader)
		if err != nil {
			errCh <- err
			return
		}
		// Wait until the process completes
		execInspect, err := cli.ContainerExecInspect(context.Background(), execResp.ID)
		if err != nil {
			errCh <- err
			return
		}
		if execInspect.ExitCode != 0 {
			if execInspect.ExitCode == 124 {
				errCh <- errors.New("execution timed out")
				return
			}
			errCh <- errors.New(outputBuffer.String())
			return
		}
		resultCh <- outputBuffer.String()
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return "", err
	}
}
package services

import (
	"context"
	"io/ioutil"
	"os/exec"
	"syscall"
)

type ExecCommandService struct {
}

func (e *ExecCommandService) Shell(ctx context.Context, request *Command) (*ShellResp, error) {
	return shell(request.Cmd, request.Args...)
}

func (e *ExecCommandService) Output(ctx context.Context, request *Command) (*OutputResp, error) {
	return output(request.Cmd, request.Args...)
}

func output(command string, args ...string) (*OutputResp, error) {
	result := &OutputResp{}

	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	result.Out = string(out)
	return result, nil
}

func shell(command string, args ...string) (*ShellResp, error) {
	result := &ShellResp{}

	cmd := exec.Command(command, args...)

	stdoutP, _ := cmd.StdoutPipe()
	stderrP, _ := cmd.StderrPipe()

	defer func() { _ = stdoutP.Close() }()
	defer func() { _ = stderrP.Close() }()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stdoutResult, err := ioutil.ReadAll(stdoutP)
	if err != nil {
		return nil, err
	}

	stderrResult, err := ioutil.ReadAll(stderrP)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			result.Code = int32(ex.Sys().(syscall.WaitStatus).ExitStatus())
		}
	}

	result.Pid = uint32(cmd.Process.Pid)
	result.Stdout = string(stdoutResult)
	result.Stderr = string(stderrResult)
	return result, nil
}

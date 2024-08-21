package tool

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// exec bash shell cmd
func ExecCmd(ctx context.Context, command string) (out string, err error) {
	var stdoutBuf bytes.Buffer
	for i := 0; i < 3; i++ {
		cmd := exec.CommandContext(ctx, "bash", "-c", command)
		cmd.Stdout = &stdoutBuf

		err = cmd.Run()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				time.Sleep(time.Millisecond * 500)
				continue
			}
			// 如果不是因为超时，直接返回err
			return "", err
		}
		// 命令执行成功
		return stdoutBuf.String(), nil
	}
	// 三次执行均超时
	return "", err
}

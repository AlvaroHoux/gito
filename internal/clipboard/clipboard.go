package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

func CopyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("wl-copy")
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("Not suported system")
	}
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
package git

import (
	"fmt"
	"os/exec"
)

func Init(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init failed: %w\n%s", err, out)
	}
	return nil
}

func IsAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

package git

import "os/exec"

func GetDiff() (string, error) {
	diffCmd := exec.Command("git", "diff", "--staged")
	diffOut, err := diffCmd.Output()

	if err != nil {
		return "", err
	}

	return string(diffOut), nil
}

func Commit(commit string) error {
	commitCmd := exec.Command("git", "commit", "-m", commit)
	if _, err := commitCmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

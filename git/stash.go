package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CreateStash(name string) (string, error) {
	var cmd *exec.Cmd

	if len(name) > 0 {
		cmd = exec.Command("git", "stash", "-u", "-m", name)
	} else {
		cmd = exec.Command("git", "stash", "-u")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	l := out.String()

	if err := cmd.Run(); err != nil {
		return l, err
	}

	return l, nil
}

func GetStashList() ([]string, error) {
	cmd := exec.Command("git", "--no-pager", "stash", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(out.String(), "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if len(line) > 0 {
			result = append(result, line)
		}
	}

	return result, nil
}

func ApplyStash(id int) (string, error) {
	cmd := exec.Command("git", "stash", "apply", fmt.Sprintf("stash@{%d}", id))
	var out bytes.Buffer
	cmd.Stdout = &out
	l := out.String()

	if err := cmd.Run(); err != nil {
		return l, err
	}

	return l, nil
}

func DropStash(id int) (string, error) {
	cmd := exec.Command("git", "stash", "drop", fmt.Sprintf("stash@{%d}", id))
	var out bytes.Buffer
	cmd.Stdout = &out
	l := out.String()

	if err := cmd.Run(); err != nil {
		return l, err
	}

	return l, nil
}

func ClearWorkspace() error {
	sh := exec.Command("git", "stash", "-u")

	if err := sh.Run(); err != nil {
		return err
	}
	dr := exec.Command("git", "stash", "drop")

	if err := dr.Run(); err != nil {
		return err
	}

	return nil
}

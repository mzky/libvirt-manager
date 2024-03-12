package utils

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func getCliPath(name string) (string, error) {
	execCmd := exec.Command("/usr/bin/which", name)

	var stdout bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stdout

	if err := execCmd.Run(); err != nil {
		logrus.Errorf("cmd run error, err=%v stderr=%v", err, stdout.String())
		return stdout.String(), err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func Exec(name string, args ...string) (string, error) {
	cli, err := getCliPath(name)
	if err != nil {
		return "", err
	}
	os.Setenv("LANG", "en_US.UTF-8") // 使用标准英文返回值，否则不同系统返回的中文不一样
	logrus.Debug(cli, args)
	execCmd := exec.Command(cli, args...)

	var stdout bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stdout

	if e := execCmd.Run(); e != nil {
		logrus.Warnf(name+" run error: %v ;%v", e, stdout.String())
		return stdout.String(), e
	}
	return stdout.String(), nil
}

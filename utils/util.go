package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func JoinVMPath(name string) string {
	return filepath.Join(VMPath, fmt.Sprintf("%s%s", name, ".qcow2"))
}

func JoinISOPath(name string) string {
	return filepath.Join(ISOPath, name)
}

func JoinImagePath(name string) string {
	return filepath.Join(ImagePath, fmt.Sprintf("%s%s", name, ".qcow2"))
}

// IsAutoStart 开机自启动虚拟机
func IsAutoStart(state bool) string {
	if state {
		return VMAutoStart
	}
	return VMNullArg
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// StringToUint64 仅用于内存修改的单位转换，默认单位M
func StringToUint64(str string) (uint64, error) {
	var multiplier float64
	switch str[len(str)-1] {
	case 'G', 'g':
		multiplier = 1 << 20 // 2^20=1024*1024
		str = str[:len(str)-1]
	case 'M', 'm':
		multiplier = 1 << 10 // 2^10=1024
		str = str[:len(str)-1]
	default:
		return 0, fmt.Errorf("unrecognized")
	}

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return uint64(num * multiplier), nil
}

type Result struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) *Result {
	return &Result{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}

func Error(data interface{}) *Result {
	return &Result{
		Status:  http.StatusExpectationFailed,
		Message: http.StatusText(http.StatusExpectationFailed),
		Data:    data,
	}
}

func CustomError(status int, data interface{}) *Result {
	return &Result{
		Status:  status,
		Message: StatusText(status),
		Data:    data,
	}
}

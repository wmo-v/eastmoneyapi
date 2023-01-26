package util

import (
	"errors"
	"os/exec"
)

// 识别验证码
func ImgOCR(path string) (string, error) {
	output, err := exec.Command("python", "./util/ocr.py", path).CombinedOutput()
	if err != nil {
		return "", errors.New(err.Error() + string(output))
	}
	return string(output), nil
}

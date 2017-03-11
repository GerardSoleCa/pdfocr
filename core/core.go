package core

import (
	"os/exec"
	"fmt"
)

func CheckDependencies() (error) {
	if _, err := exec.LookPath("pdftk"); err != nil {
		return fmt.Errorf("pdftk utility is not installed!")
	}

	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return fmt.Errorf("pdftoppm utility is not installed!")
	}
	if _, err := exec.LookPath("tesseract"); err != nil {
		return fmt.Errorf("tesseract utility is not installed!")
	}
	return nil
}

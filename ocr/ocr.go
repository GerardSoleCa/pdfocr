package ocr

import "os/exec"

func ProcessPPM(file string) {
	args := []string{"-l", "eng", file + ".ppm", file + "-new", "pdf"}
	if err := exec.Command("tesseract", args...).Run(); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"flag"
	"github.com/GerardSoleCa/pdf-ocr-go/core"
	"github.com/GerardSoleCa/pdf-ocr-go/pdf"
)

func main() {
	inPdfName := flag.String("i", "", "Choose input pdf to apply ocr")
	outPdfName := flag.String("o", "", "Output filename")
	flag.Parse()

	if err := core.CheckDependencies(); err != nil {
		panic("Dependencies not meet " + err.Error())
	}

	p := &pdf.Processor{
		InName:  *inPdfName,
		OutName: *outPdfName,
	}
	p.Process()
}

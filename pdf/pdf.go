package pdf

import (
	"os/exec"
	"regexp"
	"strconv"
	"fmt"
	"io/ioutil"
	"github.com/GerardSoleCa/pdf-ocr-go/ocr"
	"os"
	"runtime"
)

type Processor struct {
	InName   string
	OutName  string
	tmpName  string
	pages    int
	cpuCores int
}

func (p *Processor) Process() {
	defer func() {
		if r := recover(); r != nil {
			p.rmdir()
			fmt.Println("Process could not be completed... cleaning")
		}
	}()
	p.cpuCores = runtime.NumCPU()
	p.generateRandomTmpDir()
	p.mkdir()
	p.countPages()
	p.runThreaded()
	p.joinPdfs()
	p.addPdfInfo()
	p.finish()
	p.rmdir()
}

func (p *Processor) countPages() {
	args := []string{p.InName, "dump_data"}
	out, err := exec.Command("pdftk", args...).Output()
	panicIfError(err)

	rxp, err := regexp.Compile(`NumberOfPages: (\d+)`)
	group := rxp.FindStringSubmatch(string(out))

	count, err := strconv.Atoi(group[1])
	panicIfError(err)
	p.pages = count
	p.savePdfInfo(out)
	fmt.Printf("%s has %d pages\n", p.InName, p.pages)
}

func (p *Processor) savePdfInfo(info []byte) {
	err := ioutil.WriteFile(p.tmpName+"/dump_data.txt", info, 0644)
	panicIfError(err)
	fmt.Printf("%s metadata information has been saved\n", p.InName)
}

func (p *Processor) splitPage(i int) {
	//for i := 1; i < p.pages+1; i++ {
	outFileName := fmt.Sprintf(p.tmpName+"/%06d", i)
	n := strconv.Itoa(i)
	args := []string{p.InName, "cat", n, "output", outFileName + ".pdf"}
	err := exec.Command("pdftk", args...).Run()
	panicIfError(err)
	fmt.Printf("Splitting page %d\n", i)
}

func (p *Processor) generatePPM(i int) {
	outFileName := p.genFileName(i)
	args := []string{"-r", "300", outFileName + ".pdf"} //, ">", strings.Replace(file, ".pdf", ".ppm", 1)}
	out, err := exec.Command("pdftoppm", args...).Output()
	panicIfError(err)
	err = ioutil.WriteFile(outFileName+".ppm", out, 0644)
	panicIfError(err)
	fmt.Printf("Generating image from page %d\n", i)
}

func (p *Processor) ocrPage(i int) {
	outFileName := p.genFileName(i)
	ocr.ProcessPPM(outFileName)
	fmt.Printf("Applying OCR to page %d\n", i)
}

func (p *Processor) joinPdfs() {
	args := []string{"-c", "pdftk " + p.tmpName + "/*-new.pdf cat output " + p.tmpName + "/merged.pdf"}
	err := exec.Command("sh", args...).Run()
	panicIfError(err)
	fmt.Printf("Joining PDF\n")
}

func (p *Processor) addPdfInfo() {
	tn := p.tmpName
	args := []string{tn + "/merged.pdf", "update_info", tn + "/dump_data.txt", "output", tn + "/merged-data.pdf"}
	err := exec.Command("pdftk", args...).Run()
	panicIfError(err)
	fmt.Printf("Adding metadata information to joined pdf\n")

}

func (p *Processor) finish() {
	err := os.Rename(p.tmpName+"/merged-data.pdf", p.OutName)
	panicIfError(err)
	fmt.Printf("\n\nFinished! Generated PDF is %s\n", p.OutName)
}

func panicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

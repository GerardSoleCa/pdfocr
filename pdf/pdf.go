package pdf

import (
	"os/exec"
	"regexp"
	"strconv"
	"fmt"
	"io/ioutil"
	"github.com/GerardSoleCa/pdf-ocr-go/ocr"
	"os"
	"github.com/google/uuid"
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
	p.cpuCores = runtime.NumCPU()
	p.genRandomTmp()
	p.mkdir()
	p.countPages()
	for i := 1; i < p.pages+1; i++ {
		p.splitPage(i)
		p.generatePPM(i)
		p.ocrPage(i)
	}
	p.joinPdfs()
	p.applyDumpData()
	p.moveToOutname()
	p.rmdir()
}

func (p *Processor) genRandomTmp() {
	p.tmpName = uuid.New().String()
}

func (p *Processor) mkdir() {
	os.Mkdir(p.tmpName, 0755)
}

func (p *Processor) rmdir() {
	os.RemoveAll(p.tmpName)
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
}

func (p *Processor) savePdfInfo(info []byte) {
	err := ioutil.WriteFile(p.tmpName+"/dump_data.txt", info, 0644)
	panicIfError(err)
}

func (p *Processor) splitPage(i int) {
	//for i := 1; i < p.pages+1; i++ {
	outFileName := fmt.Sprintf(p.tmpName+"/%06d", i)
	n := strconv.Itoa(i)
	args := []string{p.InName, "cat", n, "output", outFileName + ".pdf"}
	err := exec.Command("pdftk", args...).Run()
	panicIfError(err)
}

func (p *Processor) generatePPM(i int) {
	outFileName := p.genFileName(i)
	args := []string{"-r", "300", outFileName + ".pdf"} //, ">", strings.Replace(file, ".pdf", ".ppm", 1)}
	out, err := exec.Command("pdftoppm", args...).Output()
	panicIfError(err)
	err = ioutil.WriteFile(outFileName+".ppm", out, 0644)
	panicIfError(err)
}

func (p *Processor) ocrPage(i int) {
	outFileName := p.genFileName(i)
	ocr.ProcessPPM(outFileName)
}

func (p *Processor) joinPdfs() {
	args := []string{"-c", "pdftk " + p.tmpName + "/*-new.pdf cat output " + p.tmpName + "/merged.pdf"}
	err := exec.Command("sh", args...).Run()
	panicIfError(err)
}

func (p *Processor) applyDumpData() {
	tn := p.tmpName
	args := []string{tn + "/merged.pdf", "update_info", tn + "/dump_data.txt", "output", tn + "/merged-data.pdf"}
	err := exec.Command("pdftk", args...).Run()
	panicIfError(err)
}

func (p *Processor) moveToOutname() {
	err := os.Rename(p.tmpName+"/merged-data.pdf", p.OutName)
	panicIfError(err)
}

func panicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (p *Processor) genFileName(i int) string {
	return fmt.Sprintf(p.tmpName+"/%06d", i)
}

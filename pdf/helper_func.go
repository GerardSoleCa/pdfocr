package pdf

import (
	"os"
	"github.com/google/uuid"
	"fmt"
)

func (p *Processor) generateRandomTmpDir() {
	p.tmpName = uuid.New().String()
}

func (p *Processor) mkdir() {
	os.Mkdir(p.tmpName, 0755)
}

func (p *Processor) rmdir() {
	os.RemoveAll(p.tmpName)
}

func (p *Processor) genFileName(i int) string {
	return fmt.Sprintf(p.tmpName+"/%06d", i)
}
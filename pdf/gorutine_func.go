package pdf

import (
	"sync"
	"fmt"
)

func (p *Processor) prepareChunkJobs() [][]int {
	var divided [][]int
	arr := make([]int, p.pages)
	for i := range arr {
		arr[i] = i + 1
	}

	chunkSize := (p.pages + p.cpuCores - 1) / p.cpuCores

	for i := 0; i < p.pages; i += chunkSize {
		end := i + chunkSize

		if end > p.pages {
			end = p.pages
		}
		divided = append(divided, arr[i:end])
	}
	return divided
}

func (p *Processor) runThreaded() {
	fmt.Printf("Starting threaded job with %d threads\n", p.cpuCores)
	chunks := p.prepareChunkJobs()
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, chunk := range chunks {
		go p.runJob(chunk, &wg)
	}
	wg.Wait()
}

func (p *Processor) runJob(chunk []int, wg *sync.WaitGroup) {
	for _, i := range chunk {
		p.splitPage(i)
		p.generatePPM(i)
		p.ocrPage(i)
	}
	wg.Done()
}

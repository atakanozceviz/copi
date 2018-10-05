package copi

import (
	"fmt"
	"os"
	"sync"
)

type Job struct {
	Src, Dst string
}

func worker(job <-chan *Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range job {
		fmt.Printf("Copy: %s\n", j.Src)
		if err := fcopy(j.Src, j.Dst); err != nil {
			fmt.Printf("Error: cannot copy %s to %s: %v\n", j.Src, j.Dst, err)
			os.Exit(1)
		}
	}
}

/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consumememory

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// CmdConsumeMemory is used by agnhost Cobra.
var CmdConsumeMemory = &cobra.Command{
	Use:   "consume-memory",
	Short: "Consumes the given amount of memory over the given amount of time",
	Long:  "Consumes the given amount of memory over the given amount of time. This command is also used by the \"resource-consumer\" subcommand.",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

var (
	workers int
	sizestr string
	vmHang  int
	timeout string
	sizeMB  = 32000000 //made match usage reported from task manager as 1 000 000 B didn't cut it
	sizeB   int64
)

func init() {
	CmdConsumeMemory.Flags().IntVarP(&workers, "workers", "m", 0, "spawn N workers spinning on malloc()/free()")
	CmdConsumeMemory.Flags().StringVar(&sizestr, "vm-bytes", "256M", "malloc B bytes per vm worker")
	CmdConsumeMemory.Flags().IntVar(&vmHang, "vm-hang", 0, "Dummy flag")
	CmdConsumeMemory.Flags().StringVarP(&timeout, "timeout", "t", "0", "timeout after N seconds")
}

func parseSize(str string) (val int64, err error) {
	if strings.HasSuffix(str, "M") {
		strval := strings.TrimSuffix(str, "M")
		val, err := strconv.ParseInt(strval, 10, 64)
		if err != nil {
			return 0, err
		}
		val = val * int64(sizeMB)
		return val, nil
	}
	return 0, nil
}

func bigBytes(bytes int64) *[]byte {
	s := make([]byte, bytes)
	return &s
}

func worker(id int, wg *sync.WaitGroup, sleep time.Duration, jobs <-chan int, results chan<- int, errors chan<- error) {
	for j := range jobs {
		s := bigBytes(sizeB)
		if s == nil {
			errors <- fmt.Errorf("pointer is nil, error on job %v", j)
		}
		time.Sleep(sleep)
		wg.Done()
	}
}

func main(cmd *cobra.Command, args []string) {
	var err error
	sizeB, err = parseSize(sizestr)
	if err != nil {
		fmt.Println("got error: ", err)
	}
	timeoutd, err := strconv.Atoi(timeout)
	if err != nil {
		fmt.Println("got error: ", err)
	}
	sleep := time.Duration(timeoutd) * time.Second
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	errors := make(chan error, 100)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		go worker(w, &wg, sleep, jobs, results, errors)
	}
	for j := 1; j <= workers; j++ {
		jobs <- j
		wg.Add(1)
	}
	close(jobs)
	wg.Wait()
	select {
	case err := <-errors:
		fmt.Println("finished with error:", err.Error())
	default:
	}
}

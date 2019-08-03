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

package consumecpu

import (
	"math"
	"time"

	"github.com/spf13/cobra"
)

// CmdConsumeCPU is used by agnhost Cobra.
var CmdConsumeCPU = &cobra.Command{
	Use:   "consume-cpu",
	Short: "TBA",
	Long:  "TBA",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

var (
	millicores  int
	durationSec int
)

func init() {
	CmdConsumeCPU.Flags().IntVar(&millicores, "millicores", 0, "millicores number")
	CmdConsumeCPU.Flags().IntVar(&durationSec, "duration-sec", 0, "duration time in seconds")
}

const sleep = time.Duration(10) * time.Millisecond

func doSomething() {
	for i := 1; i < 10000000; i++ {
		x := float64(0)
		x += math.Sqrt(0)
	}
}

func main(cmd *cobra.Command, args []string) {
	consumeCPU()
}

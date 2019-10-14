/*
Copyright 2015 The Kubernetes Authors.

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

package resconsumer

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
)

<<<<<<< Updated upstream:test/images/agnhost/resource-consumer/utils.go
var agnhostBinary, _ = filepath.Abs("./agnhost")
=======
var (
	consumeCPUBinary, _ = filepath.Abs("./consume-cpu/consume-cpu")
	consumeMemBinary, _ = filepath.Abs("./consume-memory/consume-memory")
)
>>>>>>> Stashed changes:test/images/resource-consumer/utils.go

// ConsumeCPU consumes a given number of millicores for the specified duration.
func ConsumeCPU(millicores int, durationSec int) {
	log.Printf("ConsumeCPU millicores: %v, durationSec: %v", millicores, durationSec)
	// creating new consume cpu process
<<<<<<< Updated upstream:test/images/agnhost/resource-consumer/utils.go
	arg1 := fmt.Sprintf("--millicores=%d", millicores)
	arg2 := fmt.Sprintf("--duration-sec=%d", durationSec)
	consumeCPU := exec.Command(agnhostBinary, "consume-cpu", arg1, arg2)
=======
	arg1 := fmt.Sprintf("-millicores=%d", millicores)
	arg2 := fmt.Sprintf("-duration-sec=%d", durationSec)
	consumeCPU := exec.Command(consumeCPUBinary, arg1, arg2)
>>>>>>> Stashed changes:test/images/resource-consumer/utils.go
	err := consumeCPU.Run()
	if err != nil {
		log.Printf("Error while consuming CPU: %v", err)
	}
}

// ConsumeMem consumes a given number of megabytes for the specified duration.
func ConsumeMem(megabytes int, durationSec int) {
	log.Printf("ConsumeMem megabytes: %v, durationSec: %v", megabytes, durationSec)
	megabytesString := strconv.Itoa(megabytes) + "M"
	durationSecString := strconv.Itoa(durationSec)
	// creating new consume memory process
<<<<<<< Updated upstream:test/images/agnhost/resource-consumer/utils.go
	consumeMem := exec.Command(agnhostBinary, "consume-memory", "-m", "1", "--vm-bytes", megabytesString, "--vm-hang", "0", "-t", durationSecString)
=======
	consumeMem := exec.Command(consumeMemBinary, "-m", "1", "--vm-bytes", megabytesString, "--vm-hang", "0", "-t", durationSecString)
>>>>>>> Stashed changes:test/images/resource-consumer/utils.go
	err := consumeMem.Run()
	if err != nil {
		log.Printf("Error while consuming memory: %v", err)
	}
}

// GetCurrentStatus prints out a no-op.
func GetCurrentStatus() {
	log.Printf("GetCurrentStatus")
	// not implemented
}

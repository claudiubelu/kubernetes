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
	"net/http"

	"github.com/spf13/cobra"
)

// CmdResourceConsumer is used by agnhost Cobra.
var CmdResourceConsumer = &cobra.Command{
	Use:   "resource-consumer",
	Short: "Starts an HTTP server which can accept requests for consuming CPU/memory",
	Long: `Starts an HTTP server on the given "--port" (default: 8080), which can accept requests for consuming CPU / memory in the container, or bump the value for a fake custom metric. Useful for testing autoscaling.

Endpoints:

- "/ConsumeCPU": parameters: "milicores", "durationSec". Consumes specified amount of "millicores" for "durationSec" seconds.
- "/ConsumeMemory": parameters: "megabytes", "durationSec". Consumes specified amount of "megabytes" for "durationSec" seconds.
- "/BumpMetric": parameters: "metric", "delta", "durationSec". Bumps "metric" with given name by "delta" for "durationSec" seconds.

Example:

./agnhost resource-consumer &
curl --data "millicores=300&durationSec=600" http://localhost:8080/ConsumeCPU
`,
	Args: cobra.MaximumNArgs(0),
	Run:  main,
}

var port int

func init() {
	CmdResourceConsumer.Flags().IntVar(&port, "port", 8080, "Port number.")
}

func main(cmd *cobra.Command, args []string) {
	resourceConsumerHandler := NewResourceConsumerHandler()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), resourceConsumerHandler))
}

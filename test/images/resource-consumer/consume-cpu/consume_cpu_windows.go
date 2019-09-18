// +build windows

/*
Copyright 2020 The Kubernetes Authors.

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

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	syswin "golang.org/x/sys/windows"
)

type procCPUStats struct {
	User   int64     // nanoseconds spent in user mode
	System int64     // nanoseconds spent in system mode
	Time   time.Time // when the sample was taken
	Total  int64     // total of all time fields (nanoseconds)
}

// Retrieves the amount of CPU time this process has used since it started.
func statsNow(handle syswin.Handle) (s procCPUStats) {
	var processInfo syswin.Rusage
	syswin.GetProcessTimes(handle, &processInfo.CreationTime, &processInfo.ExitTime, &processInfo.KernelTime, &processInfo.UserTime)

	s.Time = time.Now()
	s.User = processInfo.UserTime.Nanoseconds()
	s.System = processInfo.KernelTime.Nanoseconds()
	s.Total = s.User + s.System
	return s
}

// Given stats from two time points, calculates the millicores used by this
// process between the two samples.
func usageNow(first procCPUStats, current procCPUStats) float64 {
	dT := current.Time.Sub(first.Time).Nanoseconds()
	//dUser := (current.User - first.User)
	if dT == 0 {
		return 0
	}
	dUsage := (current.Total - first.Total)
	//fmt.Println("Usage: ", dUsage / 1000000, "DT: ", dT / 1000000)
	return float64(1000*dUsage) / float64(dT)
	//return 1000 * dUser / dT
}

func main() {
	flag.Parse()
	pid := os.Getpid()
	handle, _ := syswin.OpenProcess(syswin.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	defer syswin.CloseHandle(handle)

	targetMillicores := float64(*millicores)
	duration := time.Duration(*durationSec) * time.Second

	fmt.Println("pid: ", pid)

	start := time.Now()
	first := statsNow(handle)
	for time.Since(start) < duration {
		currentMillicores := usageNow(first, statsNow(handle))
		if currentMillicores < targetMillicores {
			doSomething()
		} else {
			time.Sleep(sleep)
		}
	}
}

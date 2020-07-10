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

package dns

import (
	"bytes"
	"os/exec"
	"strings"
)

const etcHostsFile = "C:/Windows/System32/drivers/etc/hosts"

func getDNSServerList() []string {
	output := runCommand("powershell", "-Command", "(Get-DnsClientServerAddress).ServerAddresses")
	if len(output) > 0 {
		return strings.Split(output, "\r\n")
	}

	panic("Could not find DNS Server list!")
}

// GetDNSSuffixList reads DNS config file and returns the list of configured DNS suffixes
func GetDNSSuffixList() []string {
	// We start with the general suffix list that apply to all network connections.
	allSuffixes := []string{}
	suffixes := getRegistryValue(netRegistry, "SearchList")
	if suffixes != "" {
		allSuffixes = strings.Split(suffixes, ",")
	}

	// Then we append the network-specific DNS suffix lists.
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, netIfacesRegistry, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		panic(err)
	}
	defer regKey.Close()
}

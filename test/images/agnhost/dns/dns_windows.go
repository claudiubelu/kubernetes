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
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	//"golang.org/x/sys/windows/registry"
)

const (
	etcHostsFile       = "C:/Windows/System32/drivers/etc/hosts"
	networkingRegistry = `System\CurrentControlSet\Services\TCPIP\Parameters`
	maxHostnameLen     = 128
	maxDomainNameLen   = 128
	maxScopeIDLen      = 256
)

// FixedInfo information: https://docs.microsoft.com/en-us/windows/win32/api/iptypes/ns-iptypes-fixed_info_w2ksp1
type FixedInfo struct {
	HostName         [maxHostnameLen + 4]byte
	DomainName       [maxDomainNameLen + 4]byte
	CurrentDNSServer *syscall.IpAddrString
	DNSServerList    syscall.IpAddrString
	NodeType         uint32
	ScopeID          [maxScopeIDLen + 4]byte
	EnableRouting    uint32
	EnableProxy      uint32
	EnableDNS        uint32
}

var (
	// GetNetworkParams can be found in iphlpapi.dll
	// see: https://docs.microsoft.com/en-us/windows/win32/api/iphlpapi/nf-iphlpapi-getnetworkparams?redirectedfrom=MSDN
	iphlpapidll          = windows.MustLoadDLL("iphlpapi.dll")
	procGetNetworkParams = iphlpapidll.MustFindProc("GetNetworkParams")
)

func getNetworkParams() *FixedInfo {
	// We don't know how big we should make the byte buffer, but the call will tell us by
	// setting the size afterwards.
	var size int
	buffer := make([]byte, 1)
	procGetNetworkParams.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	buffer = make([]byte, size)
	procGetNetworkParams.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	info := (*FixedInfo)(unsafe.Pointer(&buffer[0]))
	return info
}

func getDNSSuffixList() []string {
	/*
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, networkingRegistry, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}
	defer regKey.Close()

	searchList, _, err := regKey.GetStringValue("SearchList")
	if err != nil {
		panic(err)
	}

	return strings.Split(searchList, ",")
	*/
	output := runCommand("ipconfig", "/all")
	if len(output) > 0 {
		return []string{output}
	}

	panic("Could not find DNS search list!")
}

func getDNSServerList() []string {
	dnsServerList := []string{}
	fixedInfo := getNetworkParams()
	list := &(fixedInfo.DNSServerList)

	for list != nil {
		dnsServer := strings.TrimRight(string(list.IpAddress.String[:]), "\x00")
		dnsServerList = append(dnsServerList, dnsServer)
		list = list.Next
	}
	return dnsServerList
}

func runCommand(name string, arg ...string) string {
	var out bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(out.String())
}

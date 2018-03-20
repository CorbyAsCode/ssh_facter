package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"
)

type facterOut struct {
	Kernel                 string `json:"kernel"`
	Architecture           string `json:"architecture"`
	KernelMajorVersion     string `json:"kernelmajversion"`
	KernelRelease          string `json:"kernelrelease"`
	KernelVersion          string `json:"kernelversion"`
	OperatingSystem        string `json:"operatingsystem"`
	OSMajorRelease         string `json:"operatingsystemmajrelease"`
	OSRelease              string `json:"operatingsystemrelease"`
	OSFamily               string `json:"osfamily"`
	SwapSize               string `json:"swapsize"`
	TimeZone               string `json:"timezone"`
	UptimeDays             int    `json:"uptime_days"`
	UptimeHours            int    `json:"uptime_hours"`
	UptimeSeconds          int    `json:"uptime_seconds"`
	HardwareIsa            string `json:"hardwareisa"`
	HardwareModel          string `json:"hardwaremodel"`
	ISVirtual              bool   `json:"is_virtual"`
	MemorySize             string `json:"memorysize"`
	PhysicalProcessorCount int    `json:"physicalprocessorcount"`
	ProcessorModel         string `json:"processor0"`
	CoreCount              int    `json:"processorcount"`
	Manufacturer           string `json:"manufacturer"`
	BoardProductName       string `json:"boardproductname"`
	BoardSerial            string `json:"boardserialnumber"`
	Domain                 string `json:"domain"`
	Interfaces             string `json:"interfaces"`
	Ipaddress              string `json:"ipaddress"`
	Macaddress             string `json:"macaddress"`
	Netmask                string `json:"netmask"`
	Hostname               string `json:"hostname"`
	Model                  string `json:"productname"`
}

func sshWorker(host chan string, r chan facterOut, user string, sshKeyPath string, wg *sync.WaitGroup) {
	var facts facterOut
	sudoCmd := "sudo su - root -c "
	facterPrefix := "'facter -p -j "
	facterSlice := []string{"architecture",
		"kernel",
		"kernelmajversion",
		"kernelrelease",
		"kernelversion",
		"operatingsystem",
		"operatingsystemmajrelease",
		"operatingsystemrelease",
		"osfamily",
		"swapsize",
		"timezone",
		"uptime_days",
		"uptime_hours",
		"uptime_seconds",
		"hardwareisa",
		"hardwaremodel",
		"is_virtual",
		"memorysize",
		"physicalprocessorcount",
		"processor0",
		"processorcount",
		"domain",
		"network",
		"interfaces",
		"ipaddress",
		"macaddress",
		"netmask",
		"server_env",
		"likewisestatus",
		"likewiseversion",
		"datacenter",
		"partitions",
		"hostname",
		"manufacturer",
		"productname",
		"puppet_lastrun'"}

	facterSuffix := strings.Join(facterSlice, " ")
	cmd := sudoCmd + facterPrefix + facterSuffix

	client := &sshClient{
		IP:   <-host,
		User: user,
		Port: 22,
		Cert: sshKeyPath,
	}
	client.Connect()
	//output := client.RunCmd("cat test.json")
	output := client.RunCmd(cmd)
	client.Close()
	err := json.Unmarshal(output, &facts)
	if err != nil {
		fmt.Println(err)
	}

	r <- facts
	wg.Done()
}

func main() {

	// Parse CLI flags.
	user := flag.String("user", "", "User for ssh")
	sshKeyPath := flag.String("keypath", "", "ssh private key path")
	hostsCli := flag.String("hosts", "", "Comma-separated list of hosts")
	flag.Parse()

	// Set up variables.
	hosts := strings.Split(*hostsCli, ",")
	inputs := make(chan string, 5)
	outputs := make(chan facterOut, len(hosts))
	var wg sync.WaitGroup

	// Create ssh workers.
	for i := 1; i <= len(hosts); i++ {
		wg.Add(1)
		go sshWorker(inputs, outputs, *user, *sshKeyPath, &wg)
		fmt.Printf("Created goroutine #%d\n", i)
	}

	// Push hosts onto inputs channel.
	// Close the channel to signal that this channel is finished.
	go func() {
		for _, host := range hosts {
			inputs <- string(host)
		}
		close(inputs)
	}()

	// Wait for all sshWorkers to finish.
	wg.Wait()
	for a := 1; a <= len(hosts); a++ {
		fmt.Println(<-outputs)
		fmt.Printf("Processed #%d\n", a)
	}

}

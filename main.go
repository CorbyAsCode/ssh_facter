package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"
)

type person struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func sshWorker(host chan string, r chan person, user string, sshKeyPath string, wg *sync.WaitGroup) {
	var somebody person
	client := &sshClient{
		IP:   <-host,
		User: user,
		Port: 22,
		Cert: sshKeyPath,
	}
	client.Connect()
	output := client.RunCmd("cat test.json")
	client.Close()
	err := json.Unmarshal(output, &somebody)
	if err != nil {
		fmt.Println(err)
	}

	r <- somebody
	wg.Done()
}

func main() {

	// Parse CLI flags.
	user := flag.String("user", "", "User for ssh")
	sshKeyPath := flag.String("keypath", "", "ssh private key path")
	hostsCli := flag.String("hosts", "", "Comma-separated list of hosts")
	flag.Parse()

	// Set up variables.
	hosts := strings.Split(*hostsCli, ",") // {"10.0.0.2", "10.0.0.10"}
	inputs := make(chan string, 5)
	outputs := make(chan person, len(hosts))
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

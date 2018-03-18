package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func sshWorker(host chan string, r chan person) {
	//host := "10.0.0.2"

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
	//return c

	//fmt.Printf("Name: %s, Location: %s\n", somebody.Name, somebody.Location)

}

func concatStuff(h string, out chan<- string) {

	fmt.Println("created goroutine")
	result := h + "more"
	//fmt.Println(result)
	out <- result
}

func multiplyByTwo(in <-chan int, out chan<- int) {
	fmt.Println("Initializing goroutine...")
	num := <-in
	result := num * 2
	out <- result
}

const user = "root"
const sshKeyPath = "/Users/corbyshaner/.ssh/id_rsa"

func main() {

	/*
	  There will be an array of hosts
	  Create a function to be executed as a goroutine
	  Put each host into a channel
	  Pull each host off of the channel with a goroutine
	  Create a new client
	  Execute as normal
	  Put output onto a new channel
	  Pull outputs off of the channel and print them
	*/

	hosts := []string{"10.0.0.2", "10.0.0.10"}
	inputs := make(chan string, len(hosts))
	outputs := make(chan person, len(hosts))

	for i := 1; i < 3; i++ {
		fmt.Println("Creating goroutine")
		go func() {
			//defer wg.Done()
			var somebody person

			client := &sshClient{
				IP:   <-inputs,
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

			outputs <- somebody
		}()

	}

	for _, host := range hosts {
		fmt.Println("Adding host")
		inputs <- host
	}
	close(inputs)

	for i := 0; i < len(hosts); i++ {
		fmt.Printf("Processing #%d", i)
		fmt.Println(<-outputs)
	}

}

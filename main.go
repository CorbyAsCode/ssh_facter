package main

import (
	"encoding/json"
	"fmt"
	"sync"
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
	outputs := make(chan person, len(hosts))
	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			var somebody person

			client := &sshClient{
				IP:   host,
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
		}(host)

		/* This works
		go func(h string) {
			defer wg.Done()
			fmt.Println("created goroutine")
			result := h + "more"
			//fmt.Println(result)
			outputs <- result
			//fmt.Println("Waiting for receiver to be received.")
		}(host)
		*/
	}

	go func() {
		wg.Wait()
		fmt.Println("Closing receiver")
		close(outputs)

	}()

	for s := range outputs {
		fmt.Println(s)
	}

	//fmt.Println(<-receiver)
	//fmt.Println(<-receiver)
	//go fmt.Println(<-receiver)

	/*
		//go fmt.Println(<-receiver)
		out := make(chan int)
		in := make(chan int)

		// Create 3 `multiplyByTwo` goroutines.
		go multiplyByTwo(in, out)
		go multiplyByTwo(in, out)
		go multiplyByTwo(in, out)

		// Up till this point, none of the created goroutines actually do
		// anything, since they are all waiting for the `in` channel to
		// receive some data
		in <- 1
		in <- 2
		in <- 3

		// Now we wait for each result to come in
		fmt.Println(<-out)
		fmt.Println(<-out)
		fmt.Println(<-out)
	*/
}

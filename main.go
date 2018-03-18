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
	user := "root"
	sshKeyPath := "/Users/corbyshaner/.ssh/id_rsa"
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

func concatStuff(h <-chan string, r chan<- string) {
	fmt.Println("created goroutine")
	s := <-h
	result := s + "more"
	r <- result
}

func multiplyByTwo(in <-chan int, out chan<- int) {
	fmt.Println("Initializing goroutine...")
	num := <-in
	result := num * 2
	out <- result
}

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

	receiver := make(chan string)
	hosts := make(chan string)

	go concatStuff(hosts, receiver)
	go concatStuff(hosts, receiver)

	hosts <- "10.0.0.2"
	hosts <- "10.0.0.10"
	fmt.Println(<-receiver)
	fmt.Println(<-receiver)
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

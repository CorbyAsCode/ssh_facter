package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name     string `json:"name"`
	Location string `json:"location"`
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

	host := "10.0.0.2"
	user := "root"
	sshKeyPath := "/Users/corbyshaner/.ssh/id_rsa"
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

	//fmt.Printf("%s\n", output)

	err := json.Unmarshal(output, &somebody)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Name: %s, Location: %s\n", somebody.Name, somebody.Location)

}

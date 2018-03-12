package main

func main() {

	init()

	//hostname := &host
	client := &SSH{
		Ip:   host,
		User: user,
		Port: 22,
		Cert: sshKeyPath,
	}
	client.Connect()
	client.RunCmd("ls /")
	client.Close()

}

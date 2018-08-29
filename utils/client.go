package utils

import (
	"log"

	"k8s.io/helm/pkg/helm"
)

//Client creates a Helm client and checks the connection works
func Client(local bool) (client *helm.Client) {
	log.Println("Starting Helm client")
	if local {
		client = helm.NewClient(helm.Host("127.0.0.1:44134"))
	} else {
		client = helm.NewClient()
	}
	log.Println("Checking Helm client connection")
	_, err := client.GetVersion()
	PanicCheck(err)
	log.Println("Connection: OK")
	return
}

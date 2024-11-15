package core

import (
	"sync"

	"github.com/docker/docker/client"
)

var lock = &sync.Mutex{}

var dockerClient *client.Client

func GetDockerClient() (*client.Client,error) {
    if dockerClient != nil {
		return dockerClient, nil
	}
	lock.Lock()
	defer lock.Unlock()
	if dockerClient != nil {
		return dockerClient, nil
	}
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
    return dockerClient, nil
	//Path to the TLS certificates
    // caPath := "./ca.pem"
    // certPath := "./cert.pem"
    // keyPath := "./key.pem"

    // // Docker host address (VM External IP)
    // host := "tcp://35.200.228.129:2376"

    // // Initialize the Docker client with TLS configuration
    // cli, err := client.NewClientWithOpts(
    //     client.WithHost(host),
    //     client.WithTLSClientConfig(caPath, certPath, keyPath),
    //     client.WithAPIVersionNegotiation(),
    // )
	// if err != nil {
	// 	return nil, err
	// }
	// return cli, nil
}
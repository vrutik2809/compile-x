package core

import (
	"log"
	"sync"

	"github.com/docker/docker/client"
	"github.com/vrutik2809/compile-x/core/utils"
)

type ContainerPoolSingleton struct {
	instance *ContainerPool
	lock     *sync.Mutex
	image    string
}

var mp map[Language]*ContainerPoolSingleton

func GetContainerPool(cli *client.Client,lang Language) (*ContainerPool, error) {
	if mp[lang].instance != nil {
		return mp[lang].instance, nil
	}
	mp[lang].lock.Lock()
	defer mp[lang].lock.Unlock()
	if mp[lang].instance != nil {
		return mp[lang].instance, nil
	}
	pool, err := NewContainerPool(cli, mp[lang].image, 3)
	if err != nil {
		return nil, err
	}
	log.Println("New ContainerPool created for language: ", lang)
	mp[lang].instance = pool
	return pool, nil

}

func ContainerPoolWait() {
	var wg sync.WaitGroup
	wg.Add(len(mp))
	for lang, v := range mp {
		go func() {
			defer wg.Done()
			if v.instance != nil {
				v.instance.Wait()
				log.Println("ContainerPool terminated for language: ", lang)
			}
		}()
	}
	wg.Wait()
}

func init() {
	mp = make(map[Language]*ContainerPoolSingleton)
	mp[JAVA_22] = &ContainerPoolSingleton{
		instance: nil,
		lock:     &sync.Mutex{},
		image:    "openjdk:22-slim-bookworm",
	}
	mp[CPP_17_20] = &ContainerPoolSingleton{
		instance: nil,
		lock:     &sync.Mutex{},
		image:    "gcc:14.2.0-bookworm",
	}
	mp[PYTHON_3_12] = &ContainerPoolSingleton{
		instance: nil,
		lock:     &sync.Mutex{},
		image:    "python:3.12.7-slim-bookworm",
	}
	pullImages := func() {
		cli, err := GetDockerClient()
		if err != nil {
			log.Panicf("Failed to get docker client, error: %v", err)
		}
		var wg sync.WaitGroup
		wg.Add(len(mp))
		for _, v := range mp {
			go func() {
				defer wg.Done()
				if err := utils.PullImageIfNotExists(cli, v.image); err != nil {
					log.Panicf("Failed to pull image: %s, error: %v", v.image, err)
				}
			}()
		}
		wg.Wait()
	}
	pullImages()
	log.Println("ContainerPoolSingleton initialized")
}

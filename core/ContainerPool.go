package core

import (
	"log"
	"sync"

	"github.com/docker/docker/client"
	"github.com/vrutik2809/compile-x/core/utils"
)

type Job func(*client.Client,string)

type ContainerPool struct {
	workerQueue chan Job
	wg 		sync.WaitGroup
}

func NewContainerPool(cli *client.Client,image string, workerCount int) (*ContainerPool, error) {
	pool := &ContainerPool{
		workerQueue: make(chan Job),
	}

	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		respID, err := utils.CreateDockerContainer(cli, image)
		if err != nil {
			return nil, err
		}
		if err := utils.StartContainer(cli, respID); err != nil {
			return nil, err
		}
		go func() {
			defer pool.wg.Done()
			defer utils.RemoveContainer(cli, respID)
			defer utils.StopContainer(cli, respID)
			for job := range pool.workerQueue {
				log.Println("Job started for container: ", respID)
				job(cli,respID)
			}
		}()
	}

	return pool, nil
}

func (cp *ContainerPool) AddJob(job Job) {
	cp.workerQueue <- job
}

func (cp *ContainerPool) Wait() {
	close(cp.workerQueue)
	cp.wg.Wait()
}
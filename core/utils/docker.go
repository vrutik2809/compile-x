package utils

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	imageType "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullDockerImage(cli *client.Client, image string) error {
	out, err := cli.ImagePull(context.Background(), image, imageType.PullOptions{})
	if err != nil {
		return err
	}
    _, err = io.Copy(io.Discard, out) // using io.Discard to avoid printing to console, or os.Stdout to print
    if err != nil {
        return err
    }
    return nil
}

func IsImageExists(cli *client.Client, image string) (bool,error) {
	filters := filters.NewArgs()
	filters.Add("reference", image)
	images, err := cli.ImageList(context.Background(), imageType.ListOptions{
		Filters: filters,
	})
	if err != nil {
		return false,err
	}
	return len(images) != 0,nil
}

func PullImageIfNotExists(cli *client.Client, image string) error {
	if exists,err := IsImageExists(cli, image); err != nil {
		return err
	} else if !exists {
		log.Printf("Image: %s not found, pulling from docker registry", image)
		err = PullDockerImage(cli, image)
		if err != nil {
			return err
		}
		log.Printf("Image: %s pulled successfully", image)
	}
	return nil
}

func CreateDockerContainer(cli *client.Client, image string) (string,error) {
	if err := PullImageIfNotExists(cli, image); err != nil {
		return "",err
	}
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Tty:   true,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}, &container.HostConfig{}, nil, nil, "")
	if err != nil {
		return "",err
	}
	return resp.ID,nil
}

func StartContainer(cli *client.Client, respID string) error{
	err := cli.ContainerStart(context.Background(), respID, container.StartOptions{})
	return err
}

func StopContainer(cli *client.Client, respID string) error{
	err := cli.ContainerStop(context.Background(), respID, container.StopOptions{})
	return err
}

func RemoveContainer(cli *client.Client, respID string) error{
	err := cli.ContainerRemove(context.Background(), respID, container.RemoveOptions{})
	return err
}
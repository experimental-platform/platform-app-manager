package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strings"
)

type Dokku struct {
	Client    *docker.Client
	Container *docker.Container
}

type DokkuApp struct {
	Name          string
	ContainerType string
	ContainerId   string
	State         string
}

func NewDokku() (*Dokku, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	container, err := client.InspectContainer("dokku")
	if err != nil {
		return nil, err
	}

	return &Dokku{
		Client:    client,
		Container: container,
	}, nil
}

func (d *Dokku) raw_exec(cmd []string) (string, error) {
	exec, err := d.Client.CreateExec(docker.CreateExecOptions{
		Container:    d.Container.ID,
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	})
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	opts := docker.StartExecOptions{
		OutputStream: &stdout,
		ErrorStream:  &stderr,
	}

	err = d.Client.StartExec(exec.ID, opts)

	if err != nil {
		return "", err
	}

	if stderr.Len() != 0 {
		return "", errors.New(stderr.String())
	}
	return fmt.Sprintf("%s", &stdout), nil
}

func (d *Dokku) exec(cmd []string) (string, error) {
	full_cmd := make([]string, len(cmd)+1)
	full_cmd[0] = "dokku"
	copy(full_cmd[1:], cmd)
	return d.raw_exec(full_cmd)
}

func (d *Dokku) List() []DokkuApp {
	str, err := d.exec([]string{"ls"})
	var apps []DokkuApp
	if err == nil {
		for _, line := range strings.Split(str, "\n")[1:] {
			var appStr []string
			for _, col := range strings.Split(line, " ") {
				if col != "" {
					appStr = append(appStr, col)
				}
			}
			if length := len(appStr); length > 0 {
				var dokkuApp DokkuApp
				dokkuApp.Name = appStr[0]
				if length > 1 {
					dokkuApp.ContainerType = appStr[1]
				}
				if length > 2 {
					dokkuApp.ContainerId = appStr[2]
				}
				if length > 3 {
					dokkuApp.State = appStr[3]
				}
				apps = append(apps, dokkuApp)
			}
		}
	}
	return apps
}

func (d *Dokku) start(appName string) error {
	_, err := d.exec([]string{"ps:start", appName})
	return err
}

func (d *Dokku) stop(appName string) error {
	_, err := d.exec([]string{"ps:stop", appName})
	return err
}

func (d *Dokku) restart(appName string) error {
	_, err := d.exec([]string{"ps:restart", appName})
	return err
}

func (d *Dokku) rebuild(appName string) error {
	_, err := d.exec([]string{"ps:rebuild", appName})
	return err
}

func (d *Dokku) urls(appName string) ([]string, error) {
	str, err := d.exec([]string{"urls", appName})
	if err != nil {
		return nil, err
	}
	return strings.Split(str, "\n"), nil
}

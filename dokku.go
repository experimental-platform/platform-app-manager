package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strings"
)

type Dokku struct {
	Client *docker.Client
}

type DokkuApp struct {
	Name          string   `form:"name" json:"name" binding:"required"`
	ContainerType string   `form:"-" json:"container_type,omitempty"`
	AppType       string   `form:"-" json:"app_type",omitempty"`
	ContainerId   string   `form:"-" json:"container_id,omitempty"`
	State         string   `form:"-" json:"state,omitempty"`
	Urls          []string `form:"-" json:"urls,omitempty"`
}

func NewDokku() (*Dokku, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Dokku{
		Client: client,
	}, nil
}

func (d *Dokku) raw_exec(cmd []string) (string, error) {
	exec, err := d.Client.CreateExec(docker.CreateExecOptions{
		Container:    "dokku",
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
	str, err := d.exec([]string{"protonet:ls"})
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
					dokkuApp.AppType = appStr[2]
				}
				if length > 3 {
					dokkuApp.ContainerId = appStr[3]
				}
				if length > 4 {
					dokkuApp.State = appStr[4]
				}
				urls, err := d.urls(dokkuApp.Name)
				if err == nil {
					dokkuApp.Urls = urls
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

func (d *Dokku) destroy(appName string) error {
	_, err := d.exec([]string{"apps:destroy", appName, "--force"})
	return err
}

func (d *Dokku) logs(appName string) (string, error) {
	str, err := d.exec([]string{"logs", appName})
	return str, err
}

func (d *Dokku) urls(appName string) ([]string, error) {
	str, err := d.exec([]string{"urls", appName})
	if err != nil {
		return nil, err
	}
	result := strings.Split(str, "\n")
	for i, url := range result {
		if url == "" {
			result = append(result[:i], result[i+1:]...)
		}
	}
	return result, nil
}

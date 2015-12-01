package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

var client *Dokku
var CONTAINER_NAME = "dokku"

func main() {
	var port int
	flag.IntVar(&port, "port", 3001, "server port")
	flag.Parse()
	fmt.Println("Port: ", port)

	// disable timestamps, since we're using journald
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	var err error
	client, err = NewDokku()

	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/list", func(r render.Render) {
		apps, err := client.List()
		if err == nil {
			r.JSON(http.StatusOK, apps)
		} else {
			log.Errorf("/list: %v", err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/start", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.start(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			log.Errorf("/start '%v': %v", d.Name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/stop", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.stop(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			log.Errorf("/stop '%v': %v", d.Name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/restart", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.restart(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			log.Errorf("/restart '%v': %v", d.Name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/rebuild", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.rebuild(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			log.Errorf("/rebuild '%v': %v", d.Name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/destroy", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.destroy(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			// ignore OverlayFS error and 'no such id' error
			overlayErrorString := "Driver overlay failed to remove root filesystem"
			noSuchIDErrorString := "Error response from daemon: no such id:"
			if strings.Contains(err.Error(), overlayErrorString) || strings.Contains(err.Error(), noSuchIDErrorString) {
				r.JSON(http.StatusOK, d)
			} else {
				log.Errorf("/destroy '%v': %v", d.Name, err.Error())
				r.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	})

	m.Get("/urls/:name", func(args martini.Params, r render.Render) {
		name := args["name"]
		urls, err := client.urls(name)
		if err == nil {
			r.JSON(http.StatusOK, urls)
		} else {
			log.Errorf("/urls/%v: %v", name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Get("/logs/:name", func(args martini.Params, r render.Render) {
		name := args["name"]
		str, err := client.logs(name)
		if err == nil {
			r.JSON(http.StatusOK, struct {
				Message string `json:"msg"`
			}{
				str,
			})
		} else {
			log.Errorf("/logs/%v: %v", name, err.Error())
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	http.Handle("/", m)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

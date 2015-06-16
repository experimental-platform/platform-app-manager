package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

var client *Dokku
var CONTAINER_NAME = "dokku"

func main() {
	var port int
	flag.IntVar(&port, "port", 3001, "server port")
	flag.Parse()
	fmt.Println("Port: ", port)

	var err error
	client, err = NewDokku()

	if err != nil {
		panic(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/list", func(r render.Render) {
		r.JSON(http.StatusOK, client.List())
	})

	m.Post("/start", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.start(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/stop", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.stop(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/restart", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.restart(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/rebuild", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.rebuild(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Post("/destroy", binding.Bind(DokkuApp{}), func(d DokkuApp, r render.Render) {
		err := client.destroy(d.Name)
		if err == nil {
			r.JSON(http.StatusOK, d)
		} else {
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	m.Get("/urls/:name", func(args martini.Params, r render.Render) {
		name := args["name"]
		urls, err := client.urls(name)
		if err == nil {
			r.JSON(http.StatusOK, urls)
		} else {
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
			r.JSON(http.StatusInternalServerError, err.Error())
		}
	})

	http.Handle("/", m)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

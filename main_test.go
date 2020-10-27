package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jimlawless/whereami"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var a App

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	setConfigs()
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {}

func checkResponseCode(expected, actual int) {
	if expected != actual {
		log.Println(fmt.Printf("Expected response code %d. Got %d\n", expected, actual))
	}
}

func executeRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	req.Header.Add("api-key", "DNU7vhMsXWEymmxt")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	return resp
}

func TestGetPost(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:9000/api/v1/post/1", nil)
	response := executeRequest(req)

	checkResponseCode(http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	if body := string(body); body == "" {
		t.Errorf("Expected an non empty string. Got %s", body)
	}
}

func TestGetPosts(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:9000/api/v1/posts", nil)
	response := executeRequest(req)

	checkResponseCode(http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	if body := string(body); body == "" {
		t.Errorf("Expected an non empty string. Got %s", body)
	}
}

func TestGetPostByUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:9000/api/v1/post/1/user", nil)
	response := executeRequest(req)

	checkResponseCode(http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	if body := string(body); body == "" {
		t.Errorf("Expected an non empty string. Got %s", body)
	}
}

func TestCreatePost(t *testing.T) {
	test := &Post{
		UserId: 1,
		Id:     2,
		Title:  "Spring Boot is cooler",
		Body:   "Spring Boot makes it easy to create stand-alone, production-grade Spring based Applications that you can \"just run\"",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(test)
	req, _ := http.NewRequest("POST", "http://localhost:9000/api/v1/post", buf)
	response := executeRequest(req)
	defer response.Body.Close()

	checkResponseCode(201, response.StatusCode)
}

func TestUpdatePost(t *testing.T) {
	test := &Post{
		UserId: 1,
		Id:     1,
		Title:  "Go is is awesome",
		Body:   "Go was developed by Google",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(test)
	req, _ := http.NewRequest("PUT", "http://localhost:9000/api/v1/post/1", buf)
	response := executeRequest(req)

	checkResponseCode(200, response.StatusCode)
}

func TestDeletePost(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://localhost:9000/api/v1/post/1", nil)
	response := executeRequest(req)

	checkResponseCode(200, response.StatusCode)
}


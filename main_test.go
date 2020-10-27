package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jimlawless/whereami"
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

func TestCreatePosts(t *testing.T) {
	tests := []Post{}

	test1 := Post{
		UserId: 1,
		Id:     1,
		Title:  "Node is awesome",
		Body:   "Node.js is a JavaScript runtime built on Chrome's V8 JavaScript engine.",
	}

	test2 := Post{
		UserId: 1,
		Id:     2,
		Title:  "Spring Boot is cooler",
		Body:   "Spring Boot makes it easy to create stand-alone, production-grade Spring based Applications that you can \"just run\".",
	}

	test3 := Post{
		UserId: 2,
		Id:     3,
		Title:  "Go is faster",
		Body:   "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
	}

	test4 := Post{
		UserId: 3,
		Id:     4,
		Title:  "'What about me?' -Rails",
		Body:   "Ruby on Rails makes it much easier and more fun. It includes everything you need to build fantastic applications, and you can learn it with the support of our large, friendly community.",
	}

	tests = append(tests, test1)
	tests = append(tests, test2)
	tests = append(tests, test3)
	tests = append(tests, test4)

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(tests)
	req, _ := http.NewRequest("POST", "http://localhost:9000/api/v1/posts", buf)
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

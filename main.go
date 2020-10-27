package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jimlawless/whereami"
	"github.com/rs/cors"
)

func main() {
	defer recoverPanic()
	setConfigs()
	routes()
}

var config Configs
var posts Posts

type Configs struct {
	domain     string
	port       string
	protocol   string
	apiPath    string
	apiVersion string
	apiKey     string
}

type Post struct {
	UserId int    `json:"userid" validate:"min=1,required"`
	Id     int    `json:"id" validate:"min=1,required"`
	Title  string `json:"title" validate:"regexp=^[a-zA-Z]*$,required"`
	Body   string `json:"body" validate:"regexp=^[a-zA-Z0-9]*$,required"`
}

type Posts struct {
	Collection []Post `json:"posts"`
}

func setConfigs() {
	config = Configs{}
	config.domain = "localhost"
	config.port = ":9000"
	config.apiPath = "/api/"
	config.apiVersion = "v1"
	config.protocol = "http://"
	config.apiKey = "DNU7vhMsXWEymmxt"
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-key") != config.apiKey {
			http.Error(w, "Access Denied", http.StatusForbidden)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func routes() {
	target := config.protocol + config.domain
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{target},
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})
	path := config.apiPath + config.apiVersion
	router := mux.NewRouter()
	router.Use(Authenticate)
	router.HandleFunc(path+"/post/{id}", getPost).Methods("GET")
	router.HandleFunc(path+"/post/{id}/user", getPostByUser).Methods("GET")
	router.HandleFunc(path+"/post", createPost).Methods("POST")
	router.HandleFunc(path+"/post/{id}", updatePost).Methods("PUT")
	router.HandleFunc(path+"/post/{id}", deletePost).Methods("DELETE")
	router.HandleFunc(path+"/posts", getPosts).Methods("GET")

	err := http.ListenAndServe(config.port, corsOpts.Handler(router))
	if err != nil {
		log.Fatal(err)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	var filteredUsers Posts
	parts := strings.Split(r.URL.Path, "/")
	postId, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	for _, v := range posts.Collection {
		if v.Id == postId {
			filteredUsers.Collection = append(filteredUsers.Collection, v)
		}
	}
	output(w, filteredUsers)
}

func getPostByUser(w http.ResponseWriter, r *http.Request) {
	var filteredUsers Posts
	parts := strings.Split(r.URL.Path, "/")
	userId, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	for _, v := range posts.Collection {
		if v.UserId == userId {
			filteredUsers.Collection = append(filteredUsers.Collection, v)
		}
	}
	output(w, filteredUsers)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	initialLength := len(posts.Collection)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}
	defer r.Body.Close()

	var post Post
	err = json.Unmarshal([]byte(body), &post)
	if err != nil {
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Character)"
			log.Println(whereami.WhereAmI(), fmt.Sprintf("Invalid character at offset %v\n %s", t.Offset, jsn))
		case *json.UnmarshalTypeError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Type)"
			log.Println(whereami.WhereAmI(), fmt.Sprintf("Invalid value at offset %v\n %s", t.Offset, jsn))
		default:
			log.Println(err.Error())
		}
	}

	newLength := 0

	if validPost(post) {
		posts.Collection = append(posts.Collection, post)
		newLength = len(posts.Collection)
	}

	if initialLength < newLength {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(400)
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	postId, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}
	defer r.Body.Close()

	var post Post
	err = json.Unmarshal([]byte(body), &post)
	if err != nil {
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Character)"
			log.Println(whereami.WhereAmI(), fmt.Sprintf("Invalid character at offset %v\n %s", t.Offset, jsn))
		case *json.UnmarshalTypeError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Type)"
			log.Println(whereami.WhereAmI(), fmt.Sprintf("Invalid value at offset %v\n %s", t.Offset, jsn))
		default:
			log.Println(err.Error())
		}
	}
	var filteredUsers Posts

	var valid bool
	var statusCode int
	if validPost(post) {
		valid = true
		statusCode = 200
	} else {
		valid = false
		statusCode = 400
	}

	for _, v := range posts.Collection {
		if v.Id == postId && valid {
			filteredUsers.Collection = append(filteredUsers.Collection, post)
		} else {
			filteredUsers.Collection = append(filteredUsers.Collection, v)
		}
	}

	posts = filteredUsers
	w.WriteHeader(statusCode)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	postId, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	var filteredUsers Posts
	for _, v := range posts.Collection {
		if v.Id != postId {
			filteredUsers.Collection = append(filteredUsers.Collection, v)
		}
	}

	posts = filteredUsers
	w.WriteHeader(200)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	output(w, posts)
}

func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error)
		log.Println(whereami.WhereAmI(), err.Error())

		var l *net.TCPListener
		file, err := l.File()
		if err != nil {
			log.Println(whereami.WhereAmI(), err.Error())
		}

		path := os.Args
		args := []string{"-graceful"}

		cmd := exec.Command(path[0], args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = []*os.File{file}

		err2 := cmd.Start()
		if err2 != nil {
			log.Println(whereami.WhereAmI(), err2.Error())
		} else {
			log.Println(whereami.WhereAmI(), "Restarted...")
		}
	}
}

func output(w http.ResponseWriter, u Posts) {
	s := serialize(u)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, s)
}

func serialize(u Posts) string {
	out, err := json.Marshal(u)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}

	return string(out)
}

func validPost(post Post) bool {
	if post.Id < 1 {
		return false
	}

	if post.UserId < 1 {
		return false
	}

	if post.Body == "" {
		return false
	}

	if post.Title == "" {
		return false
	}

	return true
}

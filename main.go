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

var errorMsg ErrorMsg
var config Configs
var posts Posts

type ErrorMsg struct {
	Msg string
}

type Configs struct {
	domain     string
	port       string
	protocol   string
	apiPath    string
	apiVersion string
	apiKey     string
}

type Post struct {
	UserID int    `json:"userid" validate:"min=1,required"`
	ID     int    `json:"id" validate:"min=1,required"`
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

func authenticate(h http.Handler) http.Handler {
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
	router.Use(authenticate)
	router.HandleFunc(path+"/post/{id}", getPost).Methods("GET")
	router.HandleFunc(path+"/post/{id}/user", getPostByUser).Methods("GET")
	router.HandleFunc(path+"/post", createPost).Methods("POST")
	router.HandleFunc(path+"/posts", createPosts).Methods("POST")
	router.HandleFunc(path+"/post/{id}", updatePost).Methods("PUT")
	router.HandleFunc(path+"/post/{id}", deletePost).Methods("DELETE")
	router.HandleFunc(path+"/posts", getPosts).Methods("GET")

	err := http.ListenAndServe(config.port, corsOpts.Handler(router))
	if err != nil {
		log.Fatal(err)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	var filteredPosts Posts
	parts := strings.Split(r.URL.Path, "/")
	postID, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	for _, v := range posts.Collection {
		if v.ID == postID {
			filteredPosts.Collection = append(filteredPosts.Collection, v)
		}
	}

	if len(filteredPosts.Collection) == 1 {
		output(w, filteredPosts)
	} else {
		e := ErrorMsg{Msg: "Post not found"}
		printError(w, e)
	}
}

func getPostByUser(w http.ResponseWriter, r *http.Request) {
	var filteredPosts Posts
	parts := strings.Split(r.URL.Path, "/")
	userID, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	for _, v := range posts.Collection {
		if v.UserID == userID {
			filteredPosts.Collection = append(filteredPosts.Collection, v)
		}
	}

	if 0 < len(filteredPosts.Collection) {
		output(w, filteredPosts)
	} else {
		e := ErrorMsg{Msg: "Posts not found"}
		printError(w, e)
	}
}

func createPosts(w http.ResponseWriter, r *http.Request) {
	initialLength := len(posts.Collection)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}
	defer r.Body.Close()

	var postsReq []Post
	err = json.Unmarshal([]byte(body), &postsReq)

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

		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	newLength := 0

	for _, post := range postsReq {
		valid := uniquePost(post)
		if validPost(post) && valid {
			posts.Collection = append(posts.Collection, post)
			newLength = len(posts.Collection)
		}
	}

	if initialLength < newLength {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(400)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	initialLength := len(posts.Collection)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
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

		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	newLength := 0
	unique := uniquePost(post)
	if !unique {
		e := ErrorMsg{Msg: "Id integrity violation"}
		printError(w, e)
		return
	}

	valid := validPost(post)
	if !valid {
		e := ErrorMsg{Msg: "Invalid Post"}
		printError(w, e)
		return
	}

	if valid && unique {
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
	postID, err := strconv.Atoi(parts[4])
	if err != nil {
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
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

		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	errorMsg := ""
	var filteredPosts Posts
	var valid bool
	var statusCode int

	if validPost(post) {
		valid = true
		statusCode = 200
	} else {
		valid = false
		statusCode = 400
		errorMsg = "Invalid Post"
	}

	for _, v := range posts.Collection {
		if v.ID == postID && valid {
			filteredPosts.Collection = append(filteredPosts.Collection, post)
		} else {
			filteredPosts.Collection = append(filteredPosts.Collection, v)
		}
	}
	log.Println(errorMsg)
	posts = filteredPosts

	if errorMsg != "" {
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
	} else {
		w.WriteHeader(statusCode)
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	initialLength := len(posts.Collection)
	parts := strings.Split(r.URL.Path, "/")
	postID, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
		e := ErrorMsg{Msg: err.Error()}
		printError(w, e)
		return
	}

	var filteredPosts Posts
	for _, v := range posts.Collection {
		if v.ID != postID {
			filteredPosts.Collection = append(filteredPosts.Collection, v)
		}
	}

	newLength := len(filteredPosts.Collection)

	posts = filteredPosts
	if newLength < initialLength {
		w.WriteHeader(200)
	} else {
		e := ErrorMsg{Msg: "Post not deleted"}
		printError(w, e)
	}
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	if 0 < len(posts.Collection) {
		output(w, posts)
	} else {
		e := ErrorMsg{Msg: "Posts not created"}
		printError(w, e)
	}
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

func printError(w http.ResponseWriter, e ErrorMsg) {
	out, err := json.Marshal(e)
	if err != nil {
		log.Println(whereami.WhereAmI(), err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	fmt.Fprintf(w, string(out))
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
	if post.ID < 1 {
		return false
	}

	if post.UserID < 1 {
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

func uniquePost(p Post) bool {
	for _, post := range posts.Collection {
		if p.ID == post.ID {
			return false
		}
	}
	return true
}

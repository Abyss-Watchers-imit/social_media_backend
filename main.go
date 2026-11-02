package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var posts []Post

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/posts", MakePost).Methods("POST")
	r.HandleFunc("/posts/{id}/like", LikePost).Methods("POST")
	r.HandleFunc("/feed", GetFeed).Methods("GET")

	http.ListenAndServe(":8080", r)
}

type Post struct {
	Id      int
	Author  string
	Content string
}

func MakePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responcePayload := struct {
		PostId  int    "json:'post_id'"
		Author  string "json:'author'"
		Content string "json:'content'"
	}{
		PostId:  post.Id,
		Author:  post.Author,
		Content: post.Content,
	}
	err := json.NewEncoder(w).Encode(responcePayload)

	posts = append(posts, post)

	if err != nil {
		http.Error(w, "Failed encoding json", http.StatusInternalServerError)
		return
	}
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf("{'post_liked_id':%s}", postID)
	w.Write([]byte(response))
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Failed encoding posts", http.StatusInternalServerError)
	}
}

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Repo struct {
	Full_name string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, , r.URL.Path[1:])
}

func getRepos() *http.Response {
	var data = interface{}
	resp, _ := http.Get("https://api.github.com/repositories")
	return json.Unmarshal(ioutil.ReadAll(resp.Body), &data)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Repo struct {
	Full_name string
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := make(chan Repo)
	go getAllRepos(c)
	for repo := range c {
		fmt.Fprintf(w, "%s\n", repo.Full_name)
	}
}

func getLanguageForRep(repo string) map[string]int {
	languages := make(map[string]int)
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/languages", repo))
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &languages)
	return languages
}

func getAllRepos(c chan Repo) {
	next := "https://api.github.com/repositories"
	for {
		fmt.Println(next)
		repos, header := getRepos(next)
		if header == "" {
			break
		}

		for _, repo := range repos {
			c <- repo
		}

		next = header
	}
}

func parseHeader(linkHeader string) string {
	re := regexp.MustCompile("<(.*?)>; rel=\"next\"")
	return re.FindStringSubmatch(linkHeader)[1]
}

func getRepos(url string) ([]Repo, string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var repos []Repo
	json.Unmarshal(body, &repos)
	for _, repo := range repos {
		fmt.Printf("Got %s\n", repo.Full_name)
	}

	return repos, parseHeader(resp.Header.Get("Link"))
}

func main() {
	fmt.Println("Started")
	langs := getLanguageForRep("austinbv/dino")
	fmt.Println(len(langs))
	for language, bytes := range langs {
		fmt.Printf("%s: %v\n", language, bytes)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

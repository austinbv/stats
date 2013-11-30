package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"sort"
	"github_status/github"
)

type Data struct {
	Language map[string]int
	Limit int
}

type Repo struct {
	Full_name string
}

func getLanguageForRep(repo string) map[string]int {
	languages := make(map[string]int)
	resp, err := http.Get(addUserInfo(fmt.Sprintf("https://api.github.com/repos/%s/languages", repo)))
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &languages)
	return languages
}

func getAllRepos(c chan Data) {
	next := "https://api.github.com/repositories"
	for {
		if addUserInfo(next) == "" {
			break
		}

		repos, header := getRepos(addUserInfo(next))
		if header.RateLimitRemaining == 0 {
			fmt.Printf("waiting unil %v\n", header.RateLimitReset)
			time.Sleep(header.RateLimitReset.Sub(time.Now()))
		}

		for _, repo := range repos {
			for lang, bytes := range getLanguageForRep(repo.Full_name) {
				c <- Data{Language: map[string]int{lang: bytes}, Limit: header.RateLimitRemaining}
			}
		}

		next = header.Next
	}
}

func addUserInfo(url_string string) string {
	parsed_url, _ := url.Parse(url_string)
	parsed_url.User = url.UserPassword("austinbv", "b528d870c0e40a0ff1603c8b3c93fb6b5e7d1cb9")
	return parsed_url.String()
}

func getRepos(url string) ([]Repo, GitHubHeader) {
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

	return repos, github.ParseHeader(resp.Header)
}

func clear() {
	for i := 0; i < 500; i++ {
		fmt.Print("\033[1A")
		fmt.Print("\033[80X")
	}
}

func main() {
	languages := make(map[string]int)
	c := make(chan Data, 500)
	limit := 0

	go getAllRepos(c)
	go func(languages map[string]int, limit *int) {
		for {
			keys := make([]string, len(languages))
			i := 0
			for key := range languages {
				keys[i] = key
				i++
			}

			sum := 0.0
			for _, v := range languages {
				sum += float64(v)
			}

			clear()
			fmt.Printf("Limit:\t%v\n", *limit)
			fmt.Println("____________")
			sort.Strings(keys)
			for _, v := range keys {
				fmt.Printf("%s:\t%v%%\n", v, int(float64(languages[v])/sum*100))
			}

			time.Sleep(1*time.Second)
		}
	}(languages, &limit)

	for data := range c {
		limit = data.Limit
		for a, z := range data.Language {
			languages[a] += z
		}
	}
}

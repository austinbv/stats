package github

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Repo struct {
	Full_name string
}

func GetRepos(url string) ([]Repo, GitHubHeader) {
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

	return repos, ParseHeader(resp.Header)
}

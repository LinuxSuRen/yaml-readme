package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

func PrintStargazers(owner, repo string) (output string) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/stargazers", owner, repo)

	var (
		users     []map[string]interface{}
		companies map[string]int = make(map[string]int)
		err       error
	)

	if users, err = ghRequestAsSlice(api); err == nil {
		for _, user := range users {
			api := fmt.Sprintf("https://api.github.com/users/%s", user["login"])

			var (
				data map[string]interface{}
			)
			if data, err = ghRequestAsMap(api); err == nil {
				name := data["company"]
				if name == nil {
					continue
				}
				if count, ok := companies[name.(string)]; ok {
					companies[name.(string)] = count + 1
				} else {
					companies[name.(string)] = 1
				}
			}
		}
	}

	companies = getTopN(companies, 5)
	for name, count := range companies {
		output = output + fmt.Sprintf("<tr><td>%s</td><td>%d</td><tr>", name, count)
	}
	return fmt.Sprintf("<table>%s</table>", output)
}

func getTopN(values map[string]int, count int) (result map[string]int) {
	tops := []int{}
	for _, v := range values {
		tops = append(tops, v)
	}

	sort.Slice(tops, func(i, j int) bool {
		return tops[i] > tops[j]
	})

	if len(tops) < count {
		count = len(tops)
	} else {
		tops = tops[:count]
	}

	result = make(map[string]int, count)
	for _, v := range tops {
		for k, vul := range values {
			if vul == v {
				result[k] = vul
			}
		}
	}
	return result
}

// PrintContributors from a GitHub repository
func PrintContributors(owner, repo string) (output string) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)

	var (
		contributors []map[string]interface{}
		err          error
	)

	if contributors, err = ghRequestAsSlice(api); err == nil {
		var text string
		group := 6
		for i := 0; i < len(contributors); {
			next := i + group
			if next > len(contributors) {
				next = len(contributors)
			}
			text = text + "<tr>" + generateContributor(contributors[i:next]) + "</tr>"
			i = next
		}

		output = fmt.Sprintf(`<table>%s</table>
`, text)
	}
	return
}

// PrintPages prints the repositories which enabled pages
func PrintPages(owner string) (output string) {
	api := fmt.Sprintf("https://api.github.com/users/%s/repos?type=owner&per_page=100&sort=updated&username=%s", owner, owner)

	var (
		repos []map[string]interface{}
		err   error
	)

	if repos, err = ghRequestAsSlice(api); err == nil {
		var text string
		for i := 0; i < len(repos); i++ {
			repo := strings.TrimSpace(generateRepo(repos[i]))
			if repo != "" {
				text = text + repo + "\n"
			}
		}

		output = fmt.Sprintf(`||||
|---|---|---|
%s`, strings.TrimSpace(text))
	}
	return
}

func ghRequest(api string) (data []byte, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)

	if req, err = http.NewRequest(http.MethodGet, api, nil); err == nil {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			token = os.Getenv("GH_TOKEN")
		}
		if token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
		}

		if resp, err = http.DefaultClient.Do(req); err == nil && resp.StatusCode == http.StatusOK {
			data, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

func ghRequestAsSlice(api string) (data []map[string]interface{}, err error) {
	var byteData []byte
	if byteData, err = ghRequest(api); err == nil {
		err = json.Unmarshal(byteData, &data)
	}
	return
}

func ghRequestAsMap(api string) (data map[string]interface{}, err error) {
	var byteData []byte
	if byteData, err = ghRequest(api); err == nil {
		err = json.Unmarshal(byteData, &data)
	}
	return
}

var pageRepoTemplate = `
{{if eq .has_pages true}}  
|{{.name}}|![GitHub Repo stars](https://img.shields.io/github/stars/{{.owner.login}}/{{.name}}?style=social)|[view](https://{{.owner.login}}.github.io/{{.name}}/)| 
{{end}}
`

func generateRepo(repo interface{}) (output string) {
	var tpl *template.Template
	var err error
	if tpl, err = template.New("repo").Parse(pageRepoTemplate); err == nil {
		buf := bytes.NewBuffer([]byte{})
		if err = tpl.Execute(buf, repo); err == nil {
			output = buf.String()
		}
	}
	return
}

func generateContributor(contributors []map[string]interface{}) (output string) {
	var tpl *template.Template
	var err error
	if tpl, err = template.New("contributors").Parse(contributorsTpl); err == nil {
		buf := bytes.NewBuffer([]byte{})
		if err = tpl.Execute(buf, contributors); err == nil {
			output = buf.String()
		}
	}
	return
}

var contributorsTpl = `{{- range $i, $val := .}}
	<td align="center">
		<a href="{{$val.html_url}}">
			<img src="{{$val.avatar_url}}" width="100;" alt="{{$val.login}}"/>
			<br />
			<sub><b>{{$val.login}}</b></sub>
		</a>
	</td>
{{- end}}
`

// GitHubUsersLink parses a text and try to make the potential GitHub IDs be links
func GitHubUsersLink(ids, sep string) (links string) {
	if sep == "" {
		sep = " "
	}

	splits := strings.Split(ids, sep)
	var items []string
	for _, item := range splits {
		items = append(items, GithubUserLink(strings.TrimSpace(item), false))
	}

	// having additional whitespace it's an ASCII character
	if sep == "," {
		sep = sep + " "
	}
	links = strings.Join(items, sep)
	return
}

// GithubUserLink makes a GitHub user link
func GithubUserLink(id string, bio bool) (link string) {
	link = id
	if strings.Contains(id, " ") { // only handle the valid GitHub ID
		return
	}

	// return the original text if there are Markdown style link exist
	if hasLink(id) {
		if bio {
			return GithubUserLink(GetIDFromGHLink(id), bio)
		}
		return
	}

	api := fmt.Sprintf("https://api.github.com/users/%s", id)

	var (
		err  error
		data map[string]interface{}
	)
	if data, err = ghRequestAsMap(api); err == nil {
		link = fmt.Sprintf("[%s](%s)", data["name"], data["html_url"])
		if bioText, ok := data["bio"]; ok && bio && bioText != nil {
			link = fmt.Sprintf("%s (%s)", link, bioText)
		}
	}
	return
}

// GitHubEmojiLink returns a Markdown style link or empty
func GitHubEmojiLink(user string) (output string) {
	if user != "" {
		output = Link(":octocat:", fmt.Sprintf("https://github.com/%s", user))
	}
	return
}

// GetIDFromGHLink return the GitHub ID from a link
func GetIDFromGHLink(link string) string {
	reg, _ := regexp.Compile("\\[.*\\]\\(.*/|\\)")
	return reg.ReplaceAllString(link, "")
}

// PrintUserAsTable generates a table for a GitHub user
func PrintUserAsTable(id string) (result string) {
	api := fmt.Sprintf("https://api.github.com/users/%s", id)

	result = `|||
|---|---|
`

	var (
		err  error
		data map[string]interface{}
	)
	if data, err = ghRequestAsMap(api); err == nil {
		result = result + addWithEmpty("Name", "name", data) +
			addWithEmpty("Location", "location", data) +
			addWithEmpty("Bio", "bio", data) +
			addWithEmpty("Blog", "blog", data) +
			addWithEmpty("Twitter", "twitter_username", data) +
			addWithEmpty("Organization", "company", data)
	}
	return
}

func addWithEmpty(title, key string, data map[string]interface{}) (result string) {
	if val, ok := data[key]; ok && val != "" {
		desc := val
		switch key {
		case "twitter_username":
			desc = fmt.Sprintf("[%s](https://twitter.com/%s)", val, val)
		}
		result = fmt.Sprintf(`| %s | %s |
`, title, desc)
	}
	return
}

// hasLink determines if there are Markdown style links
func hasLink(text string) (ok bool) {
	reg, _ := regexp.Compile(".*\\[.*\\]\\(.*\\)")
	ok = reg.MatchString(text)
	return
}

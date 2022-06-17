package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

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

func ghRequest(api string) (data []byte, err error) {
	var (
		resp *http.Response
	)

	if resp, err = http.Get(api); err == nil && resp.StatusCode == http.StatusOK {
		data, err = ioutil.ReadAll(resp.Body)
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

// GithubUserLink makes a GitHub user link
func GithubUserLink(id string, bio bool) (link string) {
	api := fmt.Sprintf("https://api.github.com/users/%s", id)

	var (
		err  error
		data map[string]interface{}
	)
	if data, err = ghRequestAsMap(api); err == nil {
		link = fmt.Sprintf("[%s](%s)", data["name"], data["html_url"])
		if bio {
			link = fmt.Sprintf("%s (%s)", link, data["bio"])
		}
	}
	return
}

package main

import (
	"bytes"
	"github.com/h2non/gock"
	"io/ioutil"
	"net/http"
	"testing"
)
import "github.com/stretchr/testify/assert"

func Test_sortBy(t *testing.T) {
	type args struct {
		items  []map[string]interface{}
		sortBy string
	}
	tests := []struct {
		name   string
		args   args
		verify func([]map[string]interface{}, *testing.T)
	}{{
		name: "normal",
		args: args{
			items: []map[string]interface{}{{
				"name": "b",
			}, {
				"name": "c",
			}, {
				"name": "a",
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": "a",
			}, data[0])
		},
	}, {
		name: "number values",
		args: args{
			items: []map[string]interface{}{{
				"name": "12",
			}, {
				"name": "13",
			}, {
				"name": "11",
			}, {
				"name": "1",
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": "1",
			}, data[0])
		},
	}, {
		name: "slice values",
		args: args{
			items: []map[string]interface{}{{
				"name": []string{},
			}, {
				"name": []string{},
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": []string{},
			}, data[0])
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortBy(tt.args.items, tt.args.sortBy, true)
			tt.verify(tt.args.items, t)
		})
	}
}

func Test_generateTOC(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name    string
		args    args
		wantToc string
	}{{
		name: "simple text",
		args: args{
			txt: `## Good`,
		},
		wantToc: `- [Good](#good)
`,
	}, {
		name: "multiple levels of the titles",
		args: args{
			txt: `## Good
content
### Better`,
		},
		wantToc: `- [Good](#good)
 - [Better](#better)
`,
	}, {
		name: "has whitespace between title",
		args: args{
			txt: `## Good
## This is good`,
		},
		wantToc: `- [Good](#good)
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantToc, generateTOC(tt.args.txt), "generateTOC(%v)", tt.args.txt)
		})
	}
}

func Test_printContributor(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name       string
		args       args
		prepare    func()
		wantOutput func() string
	}{{
		name: "normal case",
		args: args{
			owner: "linuxsuren",
			repo:  "yaml-readme",
		},
		prepare: func() {
			gock.New("https://api.github.com").
				Get("/repos/linuxsuren/yaml-readme/contributors").
				Reply(http.StatusOK).
				File("data/yaml-readme.json")
		},
		wantOutput: func() string {
			data, _ := ioutil.ReadFile("data/yaml-readme-contributors.txt")
			return string(data)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.prepare()
			assert.Equalf(t, tt.wantOutput(), printContributors(tt.args.owner, tt.args.repo), "printContributors(%v, %v)", tt.args.owner, tt.args.repo)
		})
	}
}

func Test_printStarHistory(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "simple",
		args: args{
			owner: "linuxsuren",
			repo:  "yaml-readme",
		},
		want: `[![Star History Chart](https://api.star-history.com/svg?repos=linuxsuren/yaml-readme&type=Date)](https://star-history.com/#linuxsuren/yaml-readme&Date)`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, printStarHistory(tt.args.owner, tt.args.repo), "printStarHistory(%v, %v)", tt.args.owner, tt.args.repo)
		})
	}
}

func Test_getFuncMap(t *testing.T) {
	funcMap := getFuncMap("")
	assert.NotNil(t, funcMap["printToc"])
	assert.NotNil(t, funcMap["printHelp"])
	assert.NotNil(t, funcMap["printContributors"])
	assert.NotNil(t, funcMap["printStarHistory"])
	assert.NotNil(t, funcMap["printVisitorCount"])

	buf := bytes.NewBuffer([]byte{})
	printFunctions(buf)
	for k := range funcMap {
		assert.Contains(t, buf.String(), k)
	}
}

func Test_dataRender(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "bool type with true value",
		args: args{
			data: true,
		},
		want: ":white_check_mark:",
	}, {
		name: "bool type with false value",
		args: args{
			data: false,
		},
		want: ":x:",
	}, {
		name: "normal string value fake",
		args: args {
			data:"fake",
		},
		want:"fake",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, dataRender(tt.args.data), "dataRender(%v)", tt.args.data)
		})
	}
}

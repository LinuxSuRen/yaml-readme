package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	for k, val := range funcMap {
		assert.Contains(t, buf.String(), k)
		assert.NotNil(t, val)

		valType := reflect.TypeOf(val)
		numOut := valType.NumOut()
		assert.True(t, numOut > 0 && numOut < 3)

		params := make([]reflect.Value, valType.NumIn())
		for i := 0; i < valType.NumIn(); i++ {
			switch valType.In(i).Kind() {
			case reflect.Int:
				params[i] = reflect.ValueOf(1)
			case reflect.Bool:
				params[i] = reflect.ValueOf(true)
			case reflect.String:
				fallthrough
			default:
				params[i] = reflect.ValueOf("")
			}
		}

		reflect.ValueOf(val).Call(params)

		assert.Equal(t, reflect.String, valType.Out(0).Kind())
		if numOut == 2 {
			assert.Equal(t, reflect.Interface, valType.Out(1).Kind())
		}
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
		args: args{
			data: "fake",
		},
		want: "fake",
	}, {
		name: "struct parameter",
		args: args{
			data: struct {
			}{},
		},
		want: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, dataRender(tt.args.data), "dataRender(%v)", tt.args.data)
		})
	}
}

func Test_renderTemplateToString(t *testing.T) {
	type args struct {
		tplContent string
		object     interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    assert.ErrorAssertionFunc
	}{{
		name: "built-in function",
		args: args{
			tplContent: `{{render true}}`,
		},
		wantOutput: ":white_check_mark:",
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			return false
		},
	}, {
		name: "Sprig function",
		args: args{
			tplContent: `{{ "hello!" | upper | repeat 5 }}`,
		},
		wantOutput: "HELLO!HELLO!HELLO!HELLO!HELLO!",
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			return false
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := renderTemplateToString(tt.args.tplContent, tt.args.object)
			if !tt.wantErr(t, err, fmt.Sprintf("renderTemplateToString(%v, %v)", tt.args.tplContent, tt.args.object)) {
				return
			}
			assert.Equalf(t, tt.wantOutput, gotOutput, "renderTemplateToString(%v, %v)", tt.args.tplContent, tt.args.object)
		})
	}
}

func Test_newRootCommand(t *testing.T) {
	cmd := newRootCommand()
	flags := []string{"pattern", "template", "include-header", "sort-by", "group-by", "print-functions", "print-variables"}
	for _, flag := range flags {
		assert.NotNil(t, cmd.Flag(flag))
	}
}

func TestCommand(t *testing.T) {
	tests := []struct {
		name         string
		flags        []string
		hasError     bool
		expectOutput string
	}{{
		name:     "print variables",
		flags:    []string{"--print-variables"},
		hasError: false,
		expectOutput: `filename
parentname
fullpath`,
	}, {
		name:     "print functions",
		flags:    []string{"--print-functions"},
		hasError: false,
		expectOutput: `gh
ghEmoji
ghID
ghs
gstatic
link
linkOrEmpty
printContributors
printGHTable
printHelp
printPages
printStarHistory
printStargazers
printToc
printVisitorCount
render
twitterLink
youTubeLink`,
	}, {
		name:     "normal case",
		flags:    []string{"--template", "function/data/README.tpl", "--pattern", "function/data/*.yaml", "--sort-by", "zh"},
		hasError: false,
		expectOutput: `> This file was generated by [README.tpl](README.tpl) via [yaml-readme](https://github.com/LinuxSuRen/yaml-readme), please don't edit it directly!

|中文名称|英文名称|JD|
|---|---|---|
|zh|en|jd|
|zh|en|jd|`,
	}, {
		name:     "invalid group feature template with group-by flag",
		flags:    []string{"--template", "function/data/README.tpl", "--pattern", "function/data/*.yaml", "--group-by", "year"},
		hasError: true,
	}, {
		name:     "valid group feature template",
		flags:    []string{"--template", "function/data/README-group.tpl", "--pattern", "function/data/*.yaml", "--group-by", "year", "--include-header=false"},
		hasError: false,
		expectOutput: `
Year: 2021
| Zh | En |
|---|---|
| zh | en |

Year: 2022
| Zh | En |
|---|---|
| zh | en |
`,
	}, {
		name:     "group by a string",
		flags:    []string{"--template", "function/data/README-group.tpl", "--pattern", "function/data/*.yaml", "--group-by", "zh", "--include-header=false"},
		hasError: false,
		expectOutput: `
Year: zh
| Zh | En |
|---|---|
| zh | en |
| zh | en |
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			buf := bytes.NewBuffer([]byte{})
			cmd.SetOut(buf)
			cmd.SetArgs(tt.flags)

			err := cmd.Execute()
			if tt.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)

				assert.Equal(t, tt.expectOutput, buf.String())
			}
		})
	}
}

func Test_sortMetadata(t *testing.T) {
	type args struct {
		items       []map[string]interface{}
		sortByField string
	}
	tests := []struct {
		name   string
		args   args
		verify func(t *testing.T, items []map[string]interface{})
	}{{
		name: "normal case",
		args: args{
			items: []map[string]interface{}{{
				"name": "c",
			}, {
				"name": "b",
			}, {
				"name": "a",
			}},
			sortByField: "name",
		},
		verify: func(t *testing.T, items []map[string]interface{}) {
			assert.Equal(t, map[string]interface{}{"name": "a"}, items[0])
			assert.Equal(t, map[string]interface{}{"name": "b"}, items[1])
			assert.Equal(t, map[string]interface{}{"name": "c"}, items[2])
		},
	}, {
		name: "sort with descending",
		args: args{
			items: []map[string]interface{}{{
				"name": "c",
			}, {
				"name": "b",
			}, {
				"name": "a",
			}},
			sortByField: "!name",
		},
		verify: func(t *testing.T, items []map[string]interface{}) {
			assert.Equal(t, map[string]interface{}{"name": "c"}, items[0])
			assert.Equal(t, map[string]interface{}{"name": "b"}, items[1])
			assert.Equal(t, map[string]interface{}{"name": "a"}, items[2])
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortMetadata(tt.args.items, tt.args.sortByField)
		})
	}
}

func Test_loadTemplate(t *testing.T) {
	type args struct {
		templateFile  string
		includeHeader bool
	}
	tests := []struct {
		name          string
		args          args
		wantReadmeTpl func() string
		wantErr       assert.ErrorAssertionFunc
	}{{
		name: "normal case",
		args: args{
			templateFile:  "function/data/README.tpl",
			includeHeader: false,
		},
		wantReadmeTpl: func() string {
			data, _ := ioutil.ReadFile("function/data/README.tpl")
			return string(data)
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}, {
		name: "fake file",
		args: args{
			templateFile:  "fake",
			includeHeader: false,
		},
		wantReadmeTpl: func() string {
			data, _ := ioutil.ReadFile("function/data/README.tpl")
			return string(data)
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}, {
		name: "include header",
		args: args{
			templateFile:  "function/data/README.tpl",
			includeHeader: true,
		},
		wantReadmeTpl: func() string {
			data, _ := ioutil.ReadFile("function/data/README.tpl")
			return `> This file was generated by [README.tpl](README.tpl) via [yaml-readme](https://github.com/LinuxSuRen/yaml-readme), please don't edit it directly!

` + string(data)
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}, {
		name: "has metadata",
		args: args{
			templateFile:  "function/data/README-with-metadata.tpl",
			includeHeader: false,
		},
		wantReadmeTpl: func() string {
			return `a fake template`
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReadmeTpl, err := loadTemplate(tt.args.templateFile, tt.args.includeHeader)
			if !tt.wantErr(t, err, fmt.Sprintf("loadTemplate(%v, %v)", tt.args.templateFile, tt.args.includeHeader)) {
				return
			}
			assert.Equalf(t, tt.wantReadmeTpl(), gotReadmeTpl, "loadTemplate(%v, %v)", tt.args.templateFile, tt.args.includeHeader)
		})
	}
}

func Test_loadMetadata(t *testing.T) {
	type args struct {
		pattern string
		groupBy string
	}
	tests := []struct {
		name          string
		args          args
		wantItems     []map[string]interface{}
		wantGroupData map[string][]map[string]interface{}
		wantErr       assert.ErrorAssertionFunc
	}{{
		name: "normal case",
		args: args{
			pattern: "function/data/*.yaml",
			groupBy: "year",
		},
		wantItems: []map[string]interface{}{{
			"en": "en", "filename": "item-2022", "fullpath": "function/data/item-2022.yaml", "jd": "jd", "parentname": "data", "zh": "zh", "year": 2022,
		}, {
			"en": "en", "filename": "item", "fullpath": "function/data/item.yaml", "jd": "jd", "parentname": "data", "zh": "zh", "year": 2021,
		}},
		wantGroupData: map[string][]map[string]interface{}{
			"2021": {{
				"en": "en", "filename": "item", "fullpath": "function/data/item.yaml", "jd": "jd", "parentname": "data", "zh": "zh", "year": 2021,
			}},
			"2022": {{
				"en": "en", "filename": "item-2022", "fullpath": "function/data/item-2022.yaml", "jd": "jd", "parentname": "data", "zh": "zh", "year": 2022,
			}},
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, gotGroupData, err := loadMetadata(tt.args.pattern, tt.args.groupBy)
			if !tt.wantErr(t, err, fmt.Sprintf("loadMetadata(%v, %v)", tt.args.pattern, tt.args.groupBy)) {
				return
			}
			assert.Equalf(t, tt.wantItems, gotItems, "loadMetadata(%v, %v)", tt.args.pattern, tt.args.groupBy)
			assert.Equalf(t, tt.wantGroupData, gotGroupData, "loadMetadata(%v, %v)", tt.args.pattern, tt.args.groupBy)
		})
	}
}

func Test_printVariables(t *testing.T) {
	tests := []struct {
		name       string
		wantStdout string
	}{{
		name: "normal case",
		wantStdout: `filename
parentname
fullpath`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			printVariables(stdout)
			assert.Equalf(t, tt.wantStdout, stdout.String(), "printVariables(%v)", stdout)
		})
	}
}

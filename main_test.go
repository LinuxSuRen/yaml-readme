package main

import (
	"bytes"
	"fmt"
	"reflect"
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

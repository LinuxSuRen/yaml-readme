package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type option struct {
	pattern       string
	templateFile  string
	includeHeader bool
	sortBy        string
	groupBy       string

	printFunctions bool
	printVariables bool
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	if o.printFunctions {
		printFunctions(cmd.OutOrStdout())
		return
	}

	if o.printVariables {
		printVariables(cmd.OutOrStdout())
		return
	}

	var items []map[string]interface{}
	groupData := make(map[string][]map[string]interface{})

	// find YAML files
	var files []string
	var data []byte
	if files, err = filepath.Glob(o.pattern); err == nil {
		for _, metaFile := range files {
			if data, err = ioutil.ReadFile(metaFile); err != nil {
				cmd.PrintErrf("failed to read file [%s], error: %v\n", metaFile, err)
				continue
			}

			metaMap := make(map[string]interface{})
			if err = yaml.Unmarshal(data, metaMap); err != nil {
				cmd.PrintErrf("failed to parse file [%s] as a YAML, error: %v\n", metaFile, err)
				continue
			}

			// skip this item if there is a 'ignore' key is true
			if val, ok := metaMap["ignore"]; ok {
				if ignore, ok := val.(bool); ok && ignore {
					continue
				}
			}

			filename := strings.TrimSuffix(filepath.Base(metaFile), filepath.Ext(metaFile))
			parentname := filepath.Base(filepath.Dir(metaFile))

			metaMap["filename"] = filename
			metaMap["parentname"] = parentname
			metaMap["fullpath"] = metaFile

			if val, ok := metaMap[o.groupBy]; ok && val != "" {
				var strVal string
				switch val.(type) {
				case string:
					strVal = val.(string)
				case int:
					strVal = strconv.Itoa(val.(int))
				}

				if _, ok := groupData[strVal]; ok {
					groupData[strVal] = append(groupData[strVal], metaMap)
				} else {
					groupData[strVal] = []map[string]interface{}{
						metaMap,
					}
				}
			}

			items = append(items, metaMap)
		}
	}

	if o.sortBy != "" {
		descending := true
		if strings.HasPrefix(o.sortBy, "!") {
			o.sortBy = strings.TrimPrefix(o.sortBy, "!")
			descending = false
		}
		sortBy(items, o.sortBy, descending)
	}

	// load readme template
	var readmeTpl string
	if data, err = ioutil.ReadFile(o.templateFile); err != nil {
		fmt.Printf("failed to load README template, error: %v\n", err)
		readmeTpl = `
|中文名称|英文名称|JD|
|---|---|---|
{{- range $val := .}}
|{{$val.zh}}|{{$val.en}}|{{$val.jd}}|
{{end}}
`
	}
	if o.includeHeader {
		readmeTpl = fmt.Sprintf("> This file was generated by [%s](%s) via [yaml-readme](https://github.com/LinuxSuRen/yaml-readme), please don't edit it directly!\n\n",
			filepath.Base(o.templateFile), filepath.Base(o.templateFile))
	}
	readmeTpl = readmeTpl + string(data)

	// generate readme file
	var tpl *template.Template
	if tpl, err = template.New("readme").Funcs(getFuncMap(readmeTpl)).Parse(readmeTpl); err != nil {
		return
	}

	// render it with grouped data
	if o.groupBy != "" {
		err = tpl.Execute(os.Stdout, groupData)
	} else {
		err = tpl.Execute(os.Stdout, items)
	}
	return
}

func printVariables(stdout io.Writer) {
	_, _ = stdout.Write([]byte(`filename
parentname
fullpath`))
}

func printFunctions(stdout io.Writer) {
	funcMap := getFuncMap("")
	for k := range funcMap {
		_, _ = stdout.Write([]byte(fmt.Sprintf("%s\n", k)))
	}
}

func getFuncMap(readmeTpl string) template.FuncMap {
	return template.FuncMap{
		"printHelp": func(cmd string) (output string) {
			var err error
			var data []byte
			if data, err = exec.Command(cmd, "--help").Output(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "failed to run command", cmd)
			} else {
				output = fmt.Sprintf(`%s
%s
%s`, "```shell", string(data), "```")
			}
			return
		},
		"printToc": func() string {
			return generateTOC(readmeTpl)
		},
		"printContributors": func(owner, repo string) template.HTML {
			return template.HTML(printContributors(owner, repo))
		},
		"printStarHistory": func(owner, repo string) string {
			return printStarHistory(owner, repo)
		},
		"printVisitorCount": func(id string) string {
			return fmt.Sprintf(`![Visitor Count](https://profile-counter.glitch.me/%s/count.svg)`, id)
		},
		"render": dataRender,
	}
}

func sortBy(items []map[string]interface{}, sortBy string, descending bool) {
	sort.SliceStable(items, func(i, j int) (compare bool) {
		left, ok := items[i][sortBy].(string)
		if !ok {
			return false
		}
		right, ok := items[j][sortBy].(string)
		if !ok {
			return false
		}

		compare = strings.Compare(left, right) < 0
		if !descending {
			compare = !compare
		}
		return
	})
}

func generateTOC(txt string) (toc string) {
	items := strings.Split(txt, "\n")
	for i := range items {
		item := items[i]

		var prefix string
		var tag string
		if strings.HasPrefix(item, "## ") {
			tag = strings.TrimPrefix(item, "## ")
			prefix = "- "
		} else if strings.HasPrefix(item, "### ") {
			tag = strings.TrimPrefix(item, "### ")
			prefix = " - "
		} else {
			continue
		}

		// not support those titles which have whitespaces
		tag = strings.TrimSpace(tag)
		if len(strings.Split(tag, " ")) > 1 {
			continue
		}

		toc = toc + fmt.Sprintf("%s[%s](#%s)\n", prefix, tag, strings.ToLower(tag))
	}
	return
}

func printContributors(owner, repo string) (output string) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)

	var (
		resp *http.Response
		err  error
	)

	if resp, err = http.Get(api); err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	var contributors []map[string]interface{}
	if err = json.Unmarshal(data, &contributors); err != nil {
		return
	}

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
	return
}

func generateContributor(contributors []map[string]interface{}) (output string) {
	var tpl *template.Template
	var err error
	if tpl, err = template.New("contributors").Parse(contributorsTpl); err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte{})
	if err = tpl.Execute(buf, contributors); err != nil {
		return
	}
	output = buf.String()
	return
}

func printStarHistory(owner, repo string) string {
	return fmt.Sprintf(`[![Star History Chart](https://api.star-history.com/svg?repos=%[1]s/%[2]s&type=Date)](https://star-history.com/#%[1]s/%[2]s&Date)`,
		owner, repo)
}

func dataRender(data interface{}) string {
	switch val := data.(type) {
	case bool:
		if val {
			return ":white_check_mark:"
		} else {
			return ":x:"
		}
	case string:
		return val
	}
	return ""
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

func main() {
	opt := &option{}
	cmd := cobra.Command{
		Use:   "yaml-readme",
		Short: "A helper to generate a README file from Golang-based template",
		RunE:  opt.runE,
	}
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "items/*.yaml",
		"The glob pattern with Golang spec to find files")
	flags.StringVarP(&opt.templateFile, "template", "t", "README.tpl",
		"The template file which should follow Golang template spec")
	flags.BoolVarP(&opt.includeHeader, "include-header", "", true,
		"Indicate if include a notice header on the top of the README file")
	flags.StringVarP(&opt.sortBy, "sort-by", "", "",
		"Sort the array data descending by which field, or sort it ascending with the prefix '!'. For example: --sort-by !year")
	flags.StringVarP(&opt.groupBy, "group-by", "", "",
		"Group the array data by which field")
	flags.BoolVarP(&opt.printFunctions, "print-functions", "", false,
		"Print all the functions and exit")
	flags.BoolVarP(&opt.printVariables, "print-variables", "", false,
		"Print all the variables and exit")

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type option struct {
	pattern      string
	templateFile string
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var items []map[string]interface{}

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

			items = append(items, metaMap)
		}
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
	readmeTpl = string(data)

	// generate readme file
	var tpl *template.Template
	if tpl, err = template.New("readme").Parse(readmeTpl); err != nil {
		return
	}
	err = tpl.Execute(os.Stdout, items)
	return
}

func main() {
	opt := &option{}
	cmd := cobra.Command{
		Use:  "yaml-readme",
		RunE: opt.runE,
	}
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "items/*.yaml",
		"The glob pattern with Golang spec to find files")
	flags.StringVarP(&opt.templateFile, "template", "t", "README.tpl",
		"The template file which should follow Golang template spec")
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

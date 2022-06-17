[![codecov](https://codecov.io/gh/LinuxSuRen/yaml-readme/branch/master/graph/badge.svg?token=mnFyeD2IQ7)](https://codecov.io/gh/LinuxSuRen/yaml-readme)

A helper to generate the READE file automatically.

## Get started

Install it via [hd](https://github.com/LinuxSuRen/http-downloader/):

```shell
hd i yaml-readme
```

## Usage

```shell
# yaml-readme -h
Usage:
  yaml-readme [flags]

Flags:
  -h, --help              help for yaml-readme
  -p, --pattern string    The glob pattern with Golang spec to find files (default "items/*.yaml")
  -t, --template string   The template file which should follow Golang template spec (default "README.tpl")
```

### Available variables:

| Name         | Usage                                                                                           |
|--------------|-------------------------------------------------------------------------------------------------|
| `filename`   | The filename of a particular item file. For example, `items/good.yaml`, the filename is `good`. |
| `parentname` | The parent directory name. For example, `items/good.yaml`, the parent name is `items`.          |
| `fullpath`   | The related file path of each items.                                                            |

### Available functions

| Name                | Usage                                              | Description                                                             |
|---------------------|----------------------------------------------------|-------------------------------------------------------------------------|
| `printHelp`         | `{{printHelp 'hd'}}`                               | Print the help text of a command                                        |
| `printToc`          | `{{printToc}}`                                     | Print the [TOC](https://en.wikipedia.org/wiki/TOC) of the template file |
| `printContributors` | `{{printContributors "linuxsuren" "yaml-readme"}}` | Print all the contributors of an repository                             |
| `printStarHistory`  | `{{printStarHistory "linuxsuren" "yaml-readme"}}`  | Print the star history of an repository                                 |
| `printVisitorCount` | `{{printVisitorCount "repo-id"}}`                  | Print the visitor count chart of an repository                          |
| `render`            | `{{render true}}`                                  | Make the value be readable, turn `true` to `:white_check_mark:`         |
| `gh`                | `{{gh 'linuxsuren'}}`                              | Render a GitHub user to be a link                                       | 

### Ignore particular items

In case you want to ignore some particular items, you can put a key `ignore` with value `true`. Let's see the following sample:

```yaml
name: rick
ignore: true
```

## Use in GitHub actions

You could copy the following sample YAML, and change some variables according to your needs.
```yaml
name: generator

on:
  push:
    branches: [ master ]

  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    if: "!contains(github.event.head_commit.message, 'ci skip')"

    steps:
      - uses: actions/checkout@v3
      - name: Update readme
        uses: linuxsuren/yaml-readme@v0.0.6
        env:
          GH_TOKEN: ${{ secrets.GH_SECRETS }}
        with:
          pattern: 'config/*/*.yml'
          username: linuxsuren
          org: linuxsuren
          repo: hd-home
```

### Samples

Below is a simple template sample:
```gotemplate
The total number of tools is: {{len .}}
| Name | Latest | Download |
|---|---|---|
{{- range $val := .}}
| {{$val.name}} | {{$val.latest}} | {{$val.download}} |
{{- end}}
```

Below is a grouped data sample:
```gotemplate
{{- range $key, $val := .}}
Year: {{$key}}
| Name | Age |
|---|---|
{{- range $item := $val}}
| {{$item.name}} | {{$item.age}} |
{{- end}}
{{end}}
```

You could use the following command to render it:
```shell
yaml-readme --group-by year
```

Assume there is a complex YAML like this:
```yaml
metadata:
  annotations:
    group/key: 'a value'
```

then you can use the following template:
```gotemplate
{{index $item.metadata.annotations "group/key"}}
```

This plugin could help you to maintain a complex README file.

## Prerequisite

Before get started, please install `yaml-readme` to you system. You can use [hd](https://github.com/LinuxSuRen/http-downloader/) or others.

```shell
hd i yaml-readme
```

## Get started

Put the following code in the first line of the [Go template](https://pkg.go.dev/text/template) file:

```
#!yaml-readme -p 'data/financing/*.yaml' --output financing.md
```

See also [this example file](https://github.com/LinuxSuRen/open-source-best-practice/blob/master/data/financing/financing.tpl).

then press `Ctrl+Shift+P` and type `yaml-readme` command to generate the Markdown file specific with `--output`.

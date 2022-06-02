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

| Name | Usage |
|---|---|
| `filename` | The filename of a particular item file. For example, `items/good.yaml`, the filename is `good`. |
| `parentname` | The parent directory name. For example, `items/good.yaml`, the parent name is `items`. |

### Ignore particular items

In case you want to ignore some particular items, you can put a key `ignore` with value `true`. Let's see the following sample:

```yaml
name: rick
ignore: true
```

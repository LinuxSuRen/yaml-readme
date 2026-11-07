#!/bin/sh

while [ $# -gt 0 ]; do
  case "$1" in
    --pattern=*)
      pattern="${1#*=}"
      ;;
    --username=*)
      username="${1#*=}"
      ;;
    --org=*)
      org="${1#*=}"
      ;;
    --repo=*)
      repo="${1#*=}"
      ;;
    --sortby=*)
      sortby="${1#*=}"
      ;;
    --groupby=*)
      groupby="${1#*=}"
      ;;
    --output=*)
      output="${1#*=}"
      ;;
    --template=*)
      template="${1#*=}"
      ;;
    --push=*)
      push="${1#*=}"
      ;;
    --tool=*)
      tool="${1#*=}"
      ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid argument.*\n"
      printf "***************************\n"
      exit 1
  esac
  shift
done

if [ "$tool" != "" ]
then
  echo "start to install tool $tool"
  hd i "$tool"
fi

# Render to a temp file first: on a render failure the CLI may print
# usage/help text to stdout, which must never end up as the README.
tmp_output=$(mktemp)
if ! yaml-readme -p "$pattern" --sort-by "$sortby" --group-by "$groupby" --template "$template" > "$tmp_output"; then
  echo "******************************"
  echo "* Error: render failed, the README is left untouched."
  echo "******************************"
  rm -f "$tmp_output"
  exit 1
fi
mv "$tmp_output" "$output"

if [ "$push" = "true" ]
then
  git config --local user.email "${username}@users.noreply.github.com"
  git config --local user.name "${username}"
  git add .

  git commit -m "Auto commit by bot, ci skip"
  git push https://${username}:${GH_TOKEN}@github.com/${org}/${repo}.git
fi

#!/bin/sh

yaml-readme -p "$1" --sort-by "$5" --group-by "$5" > "$6"

git config --local user.email "LinuxSuRen@users.noreply.github.com"
git config --local user.name "rick"
git add .

git commit -m "Auto commit by rick's bot, ci skip"
git push https://${2}:${GH_TOKEN}@github.com/${3}/${4}.git HEAD:master

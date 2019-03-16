#!/usr/bin/env bash

set -o errexit
set -o nounset

go build ./cmd/identicon/

readme=README.md
cmd="./identicon"
themes="rainbow cool"

printf "Experiments with [identicon](https://en.wikipedia.org/wiki/Identicon) generation.\n\n" > $readme

for theme in $themes; do
    for i in {1..9}; do
        name=examples/$theme-$i.png
        $cmd -theme=$theme -input=$name > $name
        printf "<img src=\"$name\" width=\"80\"> " >> $readme
    done
    printf "\n\n" >> $readme
done

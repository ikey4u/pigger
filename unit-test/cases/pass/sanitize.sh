#!/usr/bin/env bash

for i in *
do
    if [[ -d $i ]]
    then
        cd $i && rm -rf css js tpl index.html.txt images
        cd ..
    fi
done

#!/usr/bin/env bash

dir=$1

function sanitize() {
    for i in "$1"/*
    do
        if [[ -d $i ]]
        then
            echo "[+] Sanitize $i ..."
            rm -rf $i/css $i/js $i/tpl $i/index.html.txt $i/images
            if [[ -e $i/index.html ]]
            then
                sed -i '/headtitle/d' $i/index.html
                sed -i '/headinfo/d' $i/index.html
                sed -i '/lastupdate/d' $i/index.html
            fi
        fi
    done
}

case "$dir" in
    cases)
        sanitize cases/pass
        ;;
    edge)
        sanitize edge
        ;;
    *)
        echo "Usage: $0 {cases|edge}"
esac

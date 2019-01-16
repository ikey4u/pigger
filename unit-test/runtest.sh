#! /bin/bash

action=$1

function test_cases() {
    declare -a names
    cd cases && mkdir -p pass

    echo "[+] Generating rendering files ..."
    for i in *.txt
    do
        name=$(echo $i | cut -f1  -d".")
        names=("${names[@]}" $name)
        if [[ -e $name ]]; then rm -rf $name; fi
        pigger $i
    done

    echo "[+] Run unit testing ..."
    for f in ${names[@]}
    do
        if [[ -e "pass/$f" ]]
        then
            sed -i '/headtitle/d' $f/index.html
            sed -i '/headinfo/d' $f/index.html
            sed -i '/lastupdate/d' $f/index.html
            cursum=$(gmd5sum $f/index.html | cut -f1 -d' ')
            oksum=$(gmd5sum pass/$f/index.html | cut -f1 -d' ')
            if [[ $cursum == $oksum ]]
            then
                echo "$f.txt ⇒ [✓]"
            else
                echo "$f.txt ⇒ [✗]"
            fi
        else
            echo "[!] $f is missed!"
        fi
        rm -rf $f
    done
    cd ..
}

function test_edge() {
    cd edge
    for f in *.txt
    do
        echo "[+] $f"
        pigger $f
    done
    cd ..
}

case "$action" in
    cases)
        test_cases
        ;;
    edge)
        test_edge
        ;;
    *)
        echo "Usage: $0 {cases|edge}"
esac

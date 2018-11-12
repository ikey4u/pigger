package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "log"
)

var (
    newsite = flag.String("newsite", "", "Create a new site")
	help       = flag.Bool("h", false, "show this help")
)

func setupSite() {
}

func usage() {
    fmt.Println("Usage: pigger [flags]", "")
    fmt.Println("Flags: ", "")
	flag.PrintDefaults()
}

func main() {
    flag.Usage = usage
    flag.Parse()

    if *help {
        usage()
        os.Exit(2)
    }

    if *newsite != "" {
        root, err := filepath.Abs(*newsite)
        if err != nil {
            log.Fatal("Cannot resovle " + *newsite)
        }

        /*
         * Well, a not so elegant way to create structure
         */
        os.MkdirAll(filepath.Join(root, "sys", "etc", "css"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "sys", "etc", "themes"), os.ModePerm)

        os.MkdirAll(filepath.Join(root, "sys", "www", "images"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "sys", "www", "videos"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "sys", "tmp"), os.ModePerm)

        os.MkdirAll(filepath.Join(root, "etc", "css"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "etc", "themes"), os.ModePerm)

        os.MkdirAll(filepath.Join(root, "home", "assets", "images"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "home", "assets", "videos"), os.ModePerm)
        os.Mkdir(filepath.Join(root, "home", "draft"), os.ModePerm)

        fmt.Println("The new site is located in " + root)
    }
}

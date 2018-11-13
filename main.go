package main

import (
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "log"
    "strings"
    "bufio"
)

var (
    newsite = flag.String("newsite", "", "Create a new site")
	help       = flag.Bool("h", false, "show this help")
    build = flag.Bool("b", false, "build all notes")
)

func usage() {
    fmt.Println("Usage: pigger [flags]", "")
    fmt.Println("Flags: ", "")
	flag.PrintDefaults()
}

func analyzeLine(line string) (string){
    richline := ""
    runeline := []rune(line)
    for i := 0; i < len(runeline); i++ {
        switch runeline[i] {
            case ' ' :
                if i + 1 < len(runeline) && runeline[i + 1] == ' ' {
                    remain := string(runeline[i + 2:])
                    idx := strings.Index(remain, "  ")
                    if idx != -1 {
                        blk := string(remain[0: idx])
                        richline += `&nbsp;<code class="language-clike">` + blk + "</code>&nbsp;"
                        i += len("  ") * 2 + len(blk) - 1
                    } else {
                        richline += string(runeline[i])
                    }
                }
            default:
                richline += string(runeline[i])
        }
    }
    return richline
}

func buildNotes() {
    curdir, err := os.Getwd()
    if err != nil {
        log.Fatal("Cannot get current directory!")
    }
    docdir := filepath.Join(curdir, "home")
    files, err := ioutil.ReadDir(docdir)
    if err != nil {
        log.Fatal("Cannot read dir ")
    }
    for _, f := range files {
        ext := filepath.Ext(f.Name())
        if ext == ".md" || ext == ".txt" {
            doc, _ := os.Open(filepath.Join(docdir, f.Name()))
            defer doc.Close()
            htmlname := strings.TrimSuffix(f.Name(), ext) + ".html"
            out, _ := os.Create(filepath.Join(curdir, "sys", "www", htmlname))
            title := "Test Title"
            out.WriteString(` <!DOCTYPE html> <html width="97%"lang="en"> <head> <meta charset="UTF-8">` + "\n")
            out.WriteString("<title>" + title + "</title>" + "\n")
            out.WriteString(`<link href="css/prism.css" rel="stylesheet" />` + "\n")
            out.WriteString(`<link href="css/normalize.css" rel="stylesheet" />` + "\n")
            out.WriteString("</head>" + "\n")
            out.WriteString(`<body style="margin:5%">` + "\n")

            defer out.Close()

            scanner := bufio.NewScanner(doc)

            hungry := true
            food := ""
            preline := ""
            blkend := false
            gap := false
            endmark := ""
            pretext := false

            for scanner.Scan() {
                line := scanner.Text()
                bareline := strings.TrimSpace(scanner.Text())

                if len(bareline) == 0 {
                    gap = true
                }

                if gap {
                    if preline != "" {
                        blkend = true
                    }
                }

                if blkend {
                    gtmark := strings.Index(food, ">")
                    if gtmark != -1 {
                        out.WriteString(food + endmark + "\n")
                        // endtag := "</" + string([]rune(food)[1:gtmark + 1])
                        // taglen := (len(endtag) - 1 ) * 2
                        // if len(strings.TrimSpace(strings.Replace(food, "\n", "", -1))) > taglen {
                            // out.WriteString(food + endtag + "\n")
                        // }
                    }

                    hungry = true
                    food = ""
                    blkend = false
                    gap = false
                    endmark = ""
                    if pretext {
                        pretext = false
                    }
                }

                if hungry {
                    // <pre></pre>
                    if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
                        idx := strings.Index(line, "//:")
                        highlights := "language-clike"
                        if idx != -1 {
                            highlights = "language-" + strings.TrimSpace(line[idx + 3:])
                        }
                        food = `<pre><code class="` + highlights + `">`
                        endmark = "</code></pre>"
                        if idx == -1 {
                            food += bareline + "\n"
                        }
                        hungry = false
                        pretext = true
                    // h4
                    } else if strings.HasPrefix(bareline, "####") {
                        h4 := strings.TrimLeft(bareline, "####")
                        food = "<h4>" + analyzeLine(h4)
                        endmark = "</h4>"
                        blkend = true
                    // h3
                    } else if strings.HasPrefix(bareline, "###") {
                        h3 := strings.TrimLeft(bareline, "###")
                        food = "<h3>" + analyzeLine(h3)
                        endmark = "</h3>"
                        blkend = true
                    // h2
                    } else if strings.HasPrefix(bareline, "##") {
                        h2 := strings.TrimLeft(bareline, "##")
                        food = "<h2>" + analyzeLine(h2)
                        endmark = "</h2>"
                        blkend = true
                    // h1
                    } else if strings.HasPrefix(bareline, "#") {
                        h1 := strings.TrimLeft(bareline, "#")
                        food = "<h1>" + analyzeLine(h1)
                        endmark = "</h1>"
                        blkend = true
                    // <p></p>
                    } else if preline == "" {
                        food = "<p>" + analyzeLine(line) + "\n"
                        endmark = "</p>"
                        hungry = false
                    }
                } else {
                    if pretext {
                        food += bareline + "\n"
                    } else {
                        food += analyzeLine(line) + "\n"
                    }
                }
                preline = bareline
            }
            out.WriteString(`<script src="css/prism.js"></script>` + "\n")
            out.WriteString("</body>" + "\n")
            out.WriteString("</html>" + "\n")
            out.Sync()
        }
    }
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

    if *build {
        buildNotes()
    }
}

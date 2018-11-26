package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "bytes"
    "strings"
    "html"
    "github.com/gobuffalo/packr"
    "flag"
    "os/user"
    "path"
    "path/filepath"
)

type pigconf struct {
    style string
    // private variables
    imgout_ string // Where images located in when output
}
var pc pigconf

type postmeta struct {
    Title string
    Date string
    Author string
}
type post struct {
    postmetas []postmeta
}

func getHeadline(block []byte) string {
    headline := make(map[string]string)
    lines := bytes.Split(block, []byte{0xa})
    if string(lines[0]) != "---" || string(lines[len(lines) - 1]) != "---" {
        log.Fatal("Wrong meta format!\n")
    }
    // Remove the first and last "---" from headline,
    // you know that the slice in go is really silly, it should not support negative index!
    lines = lines[1:len(lines) - 1]
    for _, line := range lines {
        s := string(line)
        info := strings.Split(s, ":")
        if len(info) < 2 {
            log.Fatal("The format of <", s , "> is not correct!\n")
        }
        key := strings.ToLower(strings.TrimSpace(info[0]))
        val := strings.TrimSpace(info[1])
        headline[key] = val
    }

    // may use template later
    headhtml := ""
    headhtml += "<!DOCTYPE html>\n"
    headhtml += `<html width="97%" lang="en">` + "\n"
    headhtml += `<head>` + "\n"
    headhtml += `<meta charset="UTF-8">` + "\n"
    headhtml += "<title>" + headline["title"] + "</title>" + "\n"
    headhtml += `<link href="` + pc.style + `/css/prism.css" rel="stylesheet" />` + "\n"
    headhtml += `<link href="` + pc.style + `/css/normalize.css" rel="stylesheet" />` + "\n"
    headhtml += `<link href="` + pc.style + `/css/pigger.css" rel="stylesheet" />` + "\n"
    headhtml += "</head>" + "\n"
    headhtml += `<body style="margin: 1% 5% 1% 5%;">` + "\n"
    headhtml += `<section style="padding-top: 20px; padding-bottom: 5px; color: #fff; text-align: center; background-image: linear-gradient(120deg, #224a73, #0d4027);">` + "\n"
    headhtml += `<h1 style="font-size: 2.25rem;">` + "\n"
    headhtml += headline["title"]
    headhtml += `</h1>` + "\n"
    headhtml += `<h3 style="font-weight: normal; opacity: 0.7; font-size: 1.15rem;">` + "\n"
    headhtml += headline["date"]
    headhtml += ` by ` + headline["author"] + "\n"
    headhtml += `</h3>` + "\n"
    headhtml += `</section>`

    return headhtml
}

func renderLine(block []byte) (string){
    htmlline := ""
    line := []rune(string(block))
    // for rune, len returns the number of character
    // i is the index of unicode character
    for i := 0; i < len(line); i++ {
        switch line[i] {
            case '`' :
                remain := string(line[i + 1:])
                // well, the fuck! idx is the byte index but not unicode!
                idx := strings.IndexRune(remain, '`')
                // fmt.Printf("idx = %d\n", idx)
                if idx != -1 {
                    blk := string(remain[0: idx])
                    htmlline += `&nbsp;<code class="language-clike">` + html.EscapeString(blk) + "</code>&nbsp;"
                    // notice that we should calculate the number of unicode and accumulate
                    i += 1 + len([]rune(blk))
                    // fmt.Printf("i = %d\n", i)
                } else {
                    htmlline += html.EscapeString(string(line[i]))
                }
            case '@':
                if i + 1 < len(line) && line[i + 1] == '[' {
                    remain := string(line[i+1:])
                    idx := strings.IndexRune(remain, ']')
                    // find ']' and there is at least one character in '[]'
                    if idx != -1 && idx > 1{
                        blk := string(remain[1:idx])
                        // link
                        if strings.HasPrefix(blk, "http") || strings.HasPrefix(blk, "ftp") {
                            link := blk
                            if len(blk) > 32 {
                                link = blk[0:32] + "..."
                            }
                            htmlline += fmt.Sprintf("<a href=\"%s\">%s</a>", blk, link)
                        // image
                        } else {
                            // copy image to destination dir
                            inimg := expandPath(blk)
                            outimg := filepath.Join(pc.imgout_, path.Base(blk))
                            // fmt.Printf("inimg: %s outimg: %s\n", inimg, outimg)
                            // avoid copy same image to itself
                            if inimg != outimg {
                                imgdata, _ := ioutil.ReadFile(blk)
                                ioutil.WriteFile(outimg, imgdata, os.ModePerm)
                            }
                            htmlline += fmt.Sprintf("<img src=\"images/%s\"/>", path.Base(blk))
                        }
                        i += 2 + len(blk)
                    } else {
                        htmlline += html.EscapeString(string(line[i]))
                    }
                } else {
                    htmlline += html.EscapeString(string(line[i]))
                }
            default:
                htmlline += html.EscapeString(string(line[i]))
        }
    }
    return htmlline
}

func renderPara(block []byte) (string) {
    lines := bytes.Split(block, []byte{0xa})
    para := "<p>"
    for _, line := range lines {
        para += renderLine(line)
    }
    return para + "</p>"
}

type Stack struct {
    data []string;
    l int;
}

func NewStack() *Stack {
    stack := new(Stack)
    stack.data = make([]string, 0)
    stack.l = 0
    return stack
}

func (s *Stack) Push(item string) {
    s.data = append(s.data, item)
    s.l += 1
}

func (s *Stack) Pop() (string) {
    if s.l > 0 {
        item := s.data[s.l - 1]
        s.data = s.data[0:s.l - 1]
        s.l -= 1
        return item
    } else {
        return ""
    }
}

func (s *Stack) Size() (int) {
    return s.l;
}

func (s *Stack) Print() {
    for idx, item := range s.data {
        fmt.Printf("[%d] = %s\n", idx, item)
    }
}

func renderList(btlines [][]byte) (string) {
    stack := NewStack()
    listhtml := "<ul>"
    stack.Push("</ul>")
    indent := 0
    firstitem := true
    for i := 0; i < len(btlines); i++ {
        line := strings.TrimRight(string(btlines[i]), " ")
        // If there should a item, then it must be the first
        wantidx := len(line) - len(strings.TrimLeft(line, " "))
        idx := strings.Index(line, "- ")
        // if the idx is not found or the item indicator is not the first
        if idx == -1 || line[wantidx:wantidx + 2] != "- " {
            space := len(line) - len(strings.TrimLeft(line, " "))
            if space >= 8 {
                codeblk := make([]byte, 0, 64)
                for j := i; j < len(btlines); j++ {
                    space = len(btlines[j]) - len(bytes.TrimLeft(btlines[j], " "))
                    if space >= 8 {
                        codeblk = append(codeblk, btlines[j]...)
                        codeblk = append(codeblk, 0xa)
                    }
                    if space < 8 || j == len(btlines) - 1 {
                        tmp := renderCode(codeblk, 8)
                        listhtml += tmp
                        if j == len(btlines) - 1 {
                            i = j
                        } else {
                            i = j - 1
                        }
                        break
                    }
                }
            } else {
                listhtml += renderLine(([]byte(line)))
            }
        } else {
            brk := ""
            if line[len(line) - 1] == ':' {
                line = line[0:len(line) - 1]
                brk = "<br/>"
            }
            if idx / 4 == indent {
                if !firstitem {
                    listhtml += stack.Pop()
                    firstitem = false
                }
                listhtml += "<li>"
                listhtml += renderLine([]byte(line[idx + 2:])) + brk
                stack.Push("</li>")
            } else if idx / 4 > indent {
                listhtml += "<ul>"
                stack.Push("</ul>")
                listhtml += "<li>"
                listhtml += renderLine([]byte(line[idx + 2:])) + brk
                stack.Push("</li>")
                indent = idx / 4
            } else {
                for j := idx / 4; j < indent; j++ {
                    listhtml += stack.Pop()
                    listhtml += stack.Pop()
                }
                listhtml += stack.Pop()
                listhtml += "<li>"
                listhtml += renderLine([]byte(line[idx + 2:])) + brk
                stack.Push("</li>")
                indent = idx / 4
            }
        }
    }
    for stack.Size() > 0 {
        listhtml += stack.Pop()
    }
    return listhtml
}

func renderTitle(block []byte) string {
    line := string(block)
    level := 0
    for idx, ch := range line {
        if ch != '#' {
            level = idx + 1
            break
        }
    }
    return fmt.Sprintf("<h%d>%s</h%d>", level, line[level:], level)
}

func renderCode(block []byte, outindent int) string {
    btlines := bytes.Split(block, []byte{0xa})
    idx := strings.Index(string(btlines[0]), "//:")
    highlights := "language-clike"
    if idx != -1 {
        highlights = "language-" + strings.TrimSpace(string(btlines[0])[idx + 3:])
    }
    code := fmt.Sprintf("<pre><code class=\"%s\">", highlights)
    for no, btline := range btlines {
        // skip highlight line
        if idx != -1 && no == 0 {
            continue
        }
        // if the last line is empty, we skip it
        if no == len(btlines) - 1 && len(bytes.TrimSpace(btline)) == 0 {
            continue
        }
        if outindent > len(btline) {
            outindent = 0
        }
        line := html.EscapeString(string(btline[outindent:]))
        code += line + "\n"
    }
    code += "</code></pre>"
    return code
}

func renderBlock(block []byte) string {
    block = bytes.TrimPrefix(block, []byte{0xa})
    lines := bytes.Split(block, []byte{0xa})
    flag := string(lines[0])
    if len(flag) >= 1 && flag[0] == '#' {
        return renderTitle(block)
    }

    if len(flag) >= 3 && flag == "---" {
        headline := getHeadline(block)
        return headline
    }

    if len(flag) >= 2 && flag[0:2] == "- " {
        items := renderList(bytes.Split(block, []byte{0xa}))
        return items
    }

    if len(flag) >= 8 && flag[0:8] == "        " {
        code := renderCode(block, 8)
        return code
    }

    if len(flag) >= 4 && flag[0:4] == "    " {
        code := renderCode(block, 4)
        return code
    }

    return renderPara(block)
}

func renderFile(infile string, outfile string) {
    input, err := ioutil.ReadFile(infile)
    if err != nil {
        log.Fatal("Cannot read file!")
    }
    blocks := bytes.Split(input[0:], []byte{0xa, 0xa})

    dochtml := ""
    for _, block := range blocks {
        dochtml += renderBlock(block) + "\n"
    }
    dochtml += `<script src="` + pc.style + `/js/prism.js"></script>` + "\n"
    dochtml += "</body>" + "\n"
    dochtml += "</html>" + "\n"

    out, _ := os.Create(outfile)
    defer out.Close()
    out.WriteString(dochtml)
}

func expandPath(p string) (out string){
    if strings.HasPrefix(p, "~") {
        usr, _ := user.Current()
        if len(p) > 1 {
            out = filepath.Join(usr.HomeDir, p[1:])
        } else {
            out = usr.HomeDir
        }
    } else {
        out, _ = filepath.Abs(p)
    }
    return out
}

func unpackResource(box packr.Box, unpack2dir string) {
    if _, err := os.Stat(unpack2dir); os.IsNotExist(err) {
        os.MkdirAll(unpack2dir, os.ModePerm)
    }
    resource := [...]string{"normalize.css", "pigger.css", "prism.css", "prism.js"}
    cssdir := filepath.Join(unpack2dir, "css"); os.Mkdir(cssdir, os.ModePerm)
    jsdir := filepath.Join(unpack2dir, "js"); os.Mkdir(jsdir, os.ModePerm)
    for _, f := range resource {
        if strings.HasSuffix(f, ".css") {
            fout, _ := os.Create(filepath.Join(cssdir, f))
            txt, _ := box.FindString("css/" + f)
            fout.WriteString(txt)
        } else if strings.HasSuffix(f, ".js") {
            fout, _ := os.Create(filepath.Join(jsdir, f))
            txt, _ := box.FindString("js/" + f)
            fout.WriteString(txt)
        }
    }
}

func main() {
    // pack static resources
    box := packr.NewBox("./etc")
    // set cmd argument options
    var outbase string
    flag.StringVar(&outbase, "o", "", "(optional) The output directory.")
    var cutoff bool
    flag.BoolVar(&cutoff, "x", false, "(optional) Cut off css and js files.")
    var style string
    flag.StringVar(&style, "style", "", "(optional) Specify a remote style root directory.")
    help := flag.Bool("h", false, "(optional) Show this help.")
    flag.Usage = func() {
        fmt.Printf("Usage: %s [[OPTIONS] <infile>]|[ACTIONS PARAMS]\nOPTIONS:\n", os.Args[0])
        flag.PrintDefaults()
        fmt.Printf("ACTIONS:\n")
        fmt.Printf("  build: Build all files\n")
        fmt.Printf("  new <sitename>: Create a new site\n")
    }
    flag.Parse()
    // check cmd args
    if *help || flag.NArg() == 0 {
        flag.Usage()
        os.Exit(0)
    }
    switch flag.Arg(0) {
    case "build":
        fmt.Printf("Build all files ...\n")
    case "new":
        if flag.NArg() != 2 {
            flag.Usage()
            log.Fatal("You forget input the name for the site, see the help above.\n")
        }
        sitedir := expandPath(flag.Arg(1))
        fmt.Printf("Create new site %s ...\n", sitedir)
        // unpackResource will create the dir if it does not exist
        unpackResource(box, filepath.Join(sitedir, "posts", "pigger"))
        os.MkdirAll(filepath.Join(sitedir, "images"), os.ModePerm)
        os.MkdirAll(filepath.Join(sitedir, "home"), os.ModePerm)
        // create pigger configuration pigger.json
        piggerconf, err := os.Create(filepath.Join(sitedir, "posts", "pigger.json"))
        defer piggerconf.Close()
        if err != nil {
            log.Fatal("Cannot create pigger config file!\n")
        }
        fmt.Printf("Good! The new site is created successfully and could be found at %s!\n", sitedir)
    default:
        infile := expandPath(flag.Arg(0))
        // test if input file is exist
        if _, err := os.Stat(infile); os.IsNotExist(err) {
            log.Fatal("Input file is not exist!\n")
        }
        // prepare input and output
        _, fname := path.Split(infile)
        barename := strings.TrimRight(fname, path.Ext(fname))
        if outbase == "" {
            outbase, _ = filepath.Abs(".")
        } else {
            outbase = expandPath(outbase)
        }
        if _, err := os.Stat(outbase); os.IsNotExist(err) {
            os.MkdirAll(outbase, os.ModePerm)
        }
        outdir := filepath.Join(outbase, barename);os.Mkdir(outdir, os.ModePerm)
        outfile := filepath.Join(outdir, "index.html")
        // unpack static resources
        if cutoff {
            if style == "" {
                pc.style = "../pigger"
                unpackResource(box, expandPath(filepath.Join(outdir, pc.style)))
            } else {
                pc.style = style
            }
        } else {
            unpackResource(box, outdir)
            pc.style = "."
        }
        // render file
        pc.imgout_ = filepath.Join(outdir, "images"); os.Mkdir(pc.imgout_, os.ModePerm)
        renderFile(infile, outfile)
        fmt.Printf("Save file into %s\n", outfile)
    }
}

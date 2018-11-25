package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "bytes"
    "strings"
    "html"
)

func sentry(err error, msg string) {
    if err != nil {
        log.Fatal(msg)
    }
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
    headhtml += `<link href="css/prism.css" rel="stylesheet" />` + "\n"
    headhtml += `<link href="css/normalize.css" rel="stylesheet" />` + "\n"
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
    // fmt.Printf("len(line) = %d\n", len(line))
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

func main() {
    input, err := ioutil.ReadFile(os.Args[1])
    sentry(err, "Cannot read file!")
    blocks := bytes.Split(input[0:], []byte{0xa, 0xa})

    dochtml := ""
    for _, block := range blocks {
        dochtml += renderBlock(block) + "\n"
    }
    dochtml += `<script src="css/prism.js"></script>` + "\n"
    dochtml += "</body>" + "\n"
    dochtml += "</html>" + "\n"

    out, _ := os.Create("test.html")
    defer out.Close()
    out.WriteString(dochtml)
    fmt.Printf("Save file into %s!\n", "test.html")
}

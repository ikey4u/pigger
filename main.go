package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "bytes"
    "strings"
)

func sentry(err error, msg string) {
    if err != nil {
        log.Fatal(msg)
    }
}

func getHeadline(block []byte) (map[string]string) {
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
    return headline
}

func renderPara(block []byte) (string) {
    para := ""
    return para
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
    html := "<ul>"
    stack.Push("</ul>")
    html += "<li>" + string(btlines[0])[2:]
    stack.Push("</li>")
    indent := 0
    for i := 1; i < len(btlines); i++ {
        line := string(btlines[i])
        // TODO:"- "  may be not the first one
        idx := strings.Index(line, "- ")
        if idx == -1 {
            html += line
        } else if idx / 4 == indent {
            html += stack.Pop()
            html += "<li>"
            html += line[idx + 2:]
            stack.Push("</li>")
        } else if idx / 4 > indent {
            html += "<ul>"
            stack.Push("</ul>")
            html += "<li>"
            html += line[idx + 2:]
            stack.Push("</li>")
            indent = idx / 4
        } else {
            for j := idx / 4; j < indent; j++ {
                html += stack.Pop()
                html += stack.Pop()
            }
            html += stack.Pop()
            html += "<li>"
            html += line[idx + 2:]
            stack.Push("</li>")
            indent = idx / 4
        }
    }
    for stack.Size() > 0 {
        html += stack.Pop()
    }
    return html
}

func renderBlock(block []byte) {
    lines := bytes.Split(block, []byte{0xa})
    flag := string(lines[0])

    if len(flag) >= 3 && flag == "---" {
        headline := getHeadline(block)
        fmt.Println(headline)
        return
    }

    if len(flag) >= 2 && flag[0:2] == "- " {
        items := renderList(bytes.Split(block, []byte{0xa}))
        fmt.Println(items)
    }
}

func main() {
    input, err := ioutil.ReadFile(os.Args[1])
    sentry(err, "Cannot read file!")
    blocks := bytes.Split(input[0:], []byte{0xa, 0xa})

    for _, block := range blocks {
        renderBlock(block)
    }
}

package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "bytes"
    "strings"
    "runtime"
    "path/filepath"
    "log"
    "time"

    "../pig"
)

const LATEST_VERSION_URL = "https://raw.githubusercontent.com/ikey4u/pigger/master/LATEST"

const UNIX_USAGE = `
1. Use echo $SHELL to check what shell you are using

2. Run the command given below

- zsh
        echo export PATH=$HOME/.local/pigger:$PATH >> $HOME/.zshrc

- bash

    - [mac]

        echo export PATH=$HOME/.local/pigger:$PATH >> $HOME/.bash_profile

    - [linux]

        echo export PATH=$HOME/.local/pigger:$PATH >> $HOME/.bashrc

3. Open a new command window, run 'pigger -v'
`

const VERSION = "1.0.0"

const USAGE = `
pigmgr [-v|-V|-h]

    -v Print pigmgr information
    -V Print version number
    -h Print this help and exit
`

func main() {
    // a dirty command line parser
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "-v":
            fmt.Printf("pigmgr version %s\n", VERSION)
        case "-V":
            fmt.Printf("%s", VERSION)
        case "-h":
            fmt.Printf(USAGE)
        default:
            fmt.Printf(USAGE)
            os.Exit(1)
        }
      os.Exit(0)
    }

    _, err := exec.LookPath("pigger")
    var curversion []byte
    if err != nil {
        curversion = []byte("0.0.0")
        fmt.Printf("[+] Start to install pigger ...\n")
    } else {
        cmd := exec.Command("pigger", "-V")
        curversion, err = cmd.Output()
        if err != nil {
            // v1.0.2, v1.0.1, v1.0.0 and before versions have no version number
            fmt.Printf("[!] pigger is too old! Try to install the latest pigger ...\n")
            curversion = []byte("0.0.0")
        } else {
            fmt.Printf("[+] The current pigger version is %s.\n", string(curversion))
        }
    }

    fd, err := ioutil.TempFile("", "pigger_version")
    defer os.Remove(fd.Name())
    pig.Download(LATEST_VERSION_URL, fd.Name())
    data, _ := ioutil.ReadFile(fd.Name())
    version := string(bytes.Split(data, []byte{0xa})[0])
    var major, minor, patch, curmajor, curminor, curpatch int
    fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
    fmt.Sscanf(string(curversion), "%d.%d.%d", &curmajor, &curminor, &curpatch)

    if major >= curmajor && minor >= curminor && patch > curpatch {
        fmt.Printf("[+] The latest version is %s, do you want to install? (y/n) ", version)

        var ans string
        for true {
            fmt.Scan(&ans)
            if ans != "y" && ans != "n" {
                fmt.Printf("Wrong choice: %s!\n", ans)
            } else {
                break
            }
        }

        if ans == "n" {
            fmt.Printf("Goodbye!\n")
            os.Exit(0)
        } else {
            dlurl := ""
            dlurl += "https://github.com/ikey4u/pigger/releases/download"
            dlurl += "/v" + strings.TrimSpace(version)
            dlurl += fmt.Sprintf("/pigger_%s_%s", runtime.GOOS, runtime.GOARCH)
            if runtime.GOOS == "windows" {
                dlurl += ".exe"
            }

            fmt.Printf("[+] Start to download ...\n")
            pigbin, _ := ioutil.TempFile("", "pigger")
            defer os.Remove(pigbin.Name())
            pig.Download(dlurl, pigbin.Name())

            fmt.Printf("[+] Installing ....\n")
            homedir := pig.SysHomedir()
            if homedir == "" {
                log.Println("[x] Cannot get user home directory!\n")
                os.Exit(1)
            }

            pighome := filepath.Join(homedir, ".local", "pigger")
            err = os.MkdirAll(pighome, os.ModePerm)
            if err != nil {
                log.Printf("[x] Cannot make pigger home: %s\n", pighome)
                os.Exit(1)
            }

            pathval := os.Getenv("PATH")
            pigbindst := filepath.Join(pighome, "pigger")
            pigmgrbindst := filepath.Join(pighome, "pigmgr")
            if runtime.GOOS == "windows" {
                if !strings.Contains(pathval, pighome) {
                    addpathcmd := `@wmic environment where 'UserName="<system>" and name = "Path"' set VariableValue="%PATH%;` + pighome + `"`
                    pigscript := "pigger.bat"
                    ioutil.WriteFile(pigscript, []byte(addpathcmd), os.ModePerm)
                    defer os.Remove(pigscript)
                    err := exec.Command(pigscript).Run()
                    if err != nil {
                        log.Printf("[x] Cannot add pigger to environment!\n")
                        log.Println(err)
                        log.Printf("[...] Pigmgr will exit in 20 seconds ...\n")
                        time.Sleep(20 * time.Second)
                        os.Exit(1)
                    }
                }

                pigbindst += ".exe"
                pigmgrbindst += ".exe"

                // refresh windows environment
                winexeurl := "https://raw.githubusercontent.com/ikey4u/pigger/master/cutils/refreshwin.exe"
                winexe := "refreshwin.exe"
                pig.Download(winexeurl, winexe)
                defer os.Remove(winexe)
                err = exec.Command(winexe).Run()
                if err != nil {
                    log.Printf("[x] Cannot refresh windows environment variables, please reboot!\n")
                    log.Println(err)
                }
            } else {
                if !strings.Contains(pathval, pighome) {
                    fmt.Printf(UNIX_USAGE)
                }
            }
            // copy pigger
            data, err := ioutil.ReadFile(pigbin.Name())
            err = ioutil.WriteFile(pigbindst, data, os.ModePerm)
            if err != nil {
                log.Printf("[x] Cannot install pigger!\n")
                log.Println(err)
            } else {
                fmt.Printf("[OK] Congratulations! Fireup a new command windows and run pigger!\n")
            }
            exec.Command("chmod +x " + pigbindst).Run()

            // copy myself(pigmgr)
            pigmgrbin, _ := filepath.Abs(os.Args[0])
            if pigmgrbin != pigmgrbindst {
                data, _ = ioutil.ReadFile(pigmgrbin)
                ioutil.WriteFile(pigmgrbindst, data, os.ModePerm)
                exec.Command("chmod +x " + pigmgrbindst).Run()
            }
        }
    } else {
        fmt.Printf("You have the latest pigger version!\n")
    }
}

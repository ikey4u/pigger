package pig

import (
    "io"
    "net/http"
    "os"
    "runtime"
)

func Download(url string, topath string) error {
    fd, err := os.Create(topath)
    if err != nil {
        return err
    }
    defer fd.Close()

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    _, err = io.Copy(fd, resp.Body)
    if err != nil {
        return err
    }
    return err
}

func SysHomedir() string {
    envhome := ""
    switch rtos := runtime.GOOS; rtos { // rtos: runtime operating system
    case "windows":
        envhome = "USERPROFILE"
    case "linux", "darwin":
        envhome = "HOME"
    }
    return os.Getenv(envhome)
}

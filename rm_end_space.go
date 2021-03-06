package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
)

const (
    VERSION = "0.1.0"
)

const (
    retOK = iota
    retFail
)

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func isProcessOK(err error) {
    if err != nil {
        fmt.Println("     [FAIL]")
    } else {
        fmt.Println("     [OK]")
    }
}

func execCmd(cmd string, shell bool) (out []byte, err error) {
    fmt.Printf("run command: %s", cmd)
    if shell {
        out, err = exec.Command("bash", "-c", cmd).Output()
        isProcessOK(err)
    } else {
        out, err = exec.Command(cmd).Output()
        isProcessOK(err)
    }
    return out, err
}

func dealFileWithWhiteList(filename string, cmd string, suffixs []string) {
    if cap(suffixs) > 0 {
        ext := filepath.Ext(filename)

        if !stringInSlice(ext, suffixs) {
            fmt.Printf("skip deal with file of: %s\n", filename)
            return
        }
    }

    cmd += filename
    execCmd(cmd, true)
}

func dealDirWithWhiteList(path string, cmd string, suffixs []string) {
    fmt.Printf("deal with dir of: %s\n", path)
    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }
        if f.IsDir() {
            if strings.HasPrefix(f.Name(), ".") {
                return filepath.SkipDir
            } else {
                return nil
            }
        } else {
            if !strings.HasPrefix(f.Name(), ".") {
                dealFileWithWhiteList(path, cmd, suffixs)
            }
        }
        return nil
    })

    if err != nil {
        fmt.Printf("filepath.Walk() returned %v\n", err)
    }
}

func isExists(file string) (ret bool, err error) {
    // equivalent to Python's `if not os.path.exists(filename)`
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return false, err
    } else {
        return true, nil
    }
}

func main() {

    var cmd string
    var suffixArray []string

    op_path := flag.String("d", "", "dir to recursive remove spaces at the end of the line.")
    op_file := flag.String("f", "", "file name for remove spaces at the end of line.")
    suffixs := flag.String("s", "", "white list of file suffixs for deal.")
    version := flag.Bool("v", false, "show version")

    flag.Parse()

    if *version {
        fmt.Printf("%s: %s\n", os.Args[0], VERSION)
        os.Exit(retOK)
    }

    switch runtime.GOOS {
    case "windows":
        fmt.Printf("[Error] not supported under windows.\n")
        os.Exit(retFail)
    case "darwin", "freebsd":
        cmd = "/usr/bin/sed -i \"\" \"s/[ ]*$//g\" "
    default:
        cmd = "sed -i \"s/[ \t]*$//g\" "
    }

    *op_path = strings.TrimSpace(*op_path)
    *op_file = strings.TrimSpace(*op_file)
    *suffixs = strings.TrimSpace(*suffixs)

    if *suffixs != "" {
        suffixArray = strings.Split(*suffixs, ",")
    }

    if *op_path == "" && *op_file == "" {
        fmt.Printf("[Error] path or file must provide one.\n\n")
        flag.Usage()
        os.Exit(retFail)
    } else if *op_file != "" {
        if _, err := isExists(*op_file); err == nil {
            if *op_file, err = filepath.Abs(*op_file); err != nil {
                panic(err)
            }
            //dealWithFile(*op_file, cmd)
            dealFileWithWhiteList(*op_file, cmd, suffixArray)
        }
    } else if *op_path != "" {
        if _, err := isExists(*op_path); err == nil {
            if *op_path, err = filepath.Abs(*op_path); err != nil {
                panic(err)
            }
            dealDirWithWhiteList(*op_path, cmd, suffixArray)
        }
    }
}
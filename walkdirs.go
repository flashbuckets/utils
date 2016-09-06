package main

import (
    "flag"
    "fmt"
    "os"
    "bytes"
    "path/filepath"
)

var dirList []string

func walker(path string, f os.FileInfo, err error) error {

    if f.IsDir() && path[0] != '.' {
        // skip "doc" directory
        if len(path) >= 3 && path[0:3] == "doc" {
            return nil
        }

        dirList = append(dirList, path)
    }

    return nil
} 

func generateCMakeList(dirs []string) {

    resultFileName := "CMakeLists.txt"
    cmakeVersion := "cmake_minimum_required(VERSION 3.5)\n"

    var tempFileBuffer bytes.Buffer
    
    tempFileBuffer.WriteString(cmakeVersion)

    root := fmt.Sprintf("aux_source_directory(. SRC_LIST)\n")
    tempFileBuffer.WriteString(root)

    for _, dir := range dirs {
        sourceDir := fmt.Sprintf("aux_source_directory(./%s SRC_LIST)\n", dir)
        includeDir := fmt.Sprintf("INCLUDE_DIRECTORIES(./%s)\n", dir)
        tempFileBuffer.WriteString(sourceDir)
        tempFileBuffer.WriteString(includeDir)
    }

    execName := fmt.Sprintf("add_executable(${PROJECT_NAME} ${SRC_LIST})")
    tempFileBuffer.WriteString(execName)

    file, err := os.Create(resultFileName)
    defer file.Close()
    if err != nil {
        fmt.Println("create \"%s\" failed, %s", resultFileName, err.Error())
        return
    }

    file.Write(tempFileBuffer.Bytes())
}

func GetDirList(path string) error {
    err := filepath.Walk(path, walker)
    return err
}

func DumpDirList(dirs []string) {
    for _, dir := range dirs {
        fmt.Printf("%s\n", dir)
    }
}

func main() {
    flag.Parse()
    base := flag.Arg(0)
    GetDirList(base)
    //DumpDirList(dirList)
    generateCMakeList(dirList)
}


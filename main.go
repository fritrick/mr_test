package main

import (
    "fmt"
    "bufio"
    "os"
    "flag"
    "net/http"
    "io/ioutil"
    "strings"
)

func main() {
    typePtr := flag.String("type", "", "Input type for parse ( must be url|file )")
    flag.Parse();
    if *typePtr != "url" && *typePtr != "file" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    var listItems []string
    messages := make(chan int, 3)

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        var url = scanner.Text()
        listItems = append(listItems, url)
        go func() { messages <- RoutineCounterHttp(url, "Go") }() //RoutineCounterHttp(url, "Go")

    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
    var msg int
    for i := 0; i < len(listItems); i++ {
        msg <- messages
        fmt.Printls(msg)
    }
}




func RoutineCounterHttp(path string, substring string) int {
    result := MakeHttpRequest(path)
    return CountSubstrings(result, substring)
}

func RoutineCounterFile(path string, substring string) int {
    result := ReadFileAsString(path)
    return CountSubstrings(result, substring)
}

func CountSubstrings(text string, substring string) int {
    return strings.Count(text, substring)
}

func MakeHttpRequest(url string) string {
    resp, _ := http.Get(url)
    bytes, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    return string(bytes)
}

func ReadFileAsString(filename string) string {
    b, err := ioutil.ReadFile(filename) // just pass the file name
    if err != nil {
        fmt.Print(err)
        return ""
    }
    return string(b)
}

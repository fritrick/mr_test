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

const (
    K=5
)

type fn func(path string, substring string) int

func main() {
    typePtr := flag.String("type", "", "Input type for parse ( must be url|file )")
    flag.Parse();
    if *typePtr != "url" && *typePtr != "file" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    var listItems []string
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        var url = scanner.Text()
        listItems = append(listItems, url)
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }

    var callback fn
    switch *typePtr {
        case "url":
            callback = CounterHttp
        case "file":
            callback = CounterFile
    }

    ProceedQueue(listItems, callback)
}


func ProceedQueue(listItems []string, f fn) {
    counts := make(chan int)
    listLength := GetWorkersCount(len(listItems))
    // var wg sync.WaitGroup

    for j := 0; j < listLength; j++ {
        // wg.Add(1)
        go func(j int) {
            currentUrl := listItems[j]
            currentCount := f(currentUrl, "Go")
            fmt.Printf("Count for %s: %d \n", currentUrl, currentCount)
            counts <- currentCount
            // wg.Done()
        }(j)
    }




    totalCounter := 0
    for i := 0; i < len(listItems); i++ {
        totalCounter = totalCounter + (<- counts)
    }
    fmt.Println("Total: ", totalCounter)
    // wg.Wait()
}


func GetWorkersCount(length int) int {
    var workersCount int
    if length < K {
        workersCount = length
    } else {
        workersCount = K
    }
    return workersCount
}

func CounterHttp(path string, substring string) int {
    result := MakeHttpRequest(path)
    return CountSubstrings(result, substring)
}

func CounterFile(path string, substring string) int {
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

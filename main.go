package main

import (
    "fmt"
    "bufio"
    "os"
    "flag"
    "net/http"
    "io/ioutil"
    "strings"
    "sync"
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
    items := make(chan string)
    totalCounter := 0

    workersCount := GetWorkersCount(len(listItems))
    var wg sync.WaitGroup

    for j := 0; j < workersCount; j++ {
        wg.Add(1)
        go func() {
            for v := range items{
                currentCount := f(v, "Go")
                totalCounter = totalCounter + currentCount
                fmt.Printf("Count for %s: %d \n", v, currentCount)
            }
            wg.Done()
        }()
    }

    for i := 0; i < len(listItems); i++ {
		items <- listItems[i]
	}
	close(items)
    wg.Wait()
    fmt.Println("Total: ", totalCounter)
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

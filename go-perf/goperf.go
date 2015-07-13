/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package main

import (
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "math/rand"
    "net/http"
    "os"
    "os/signal"
    "io"
    "time"
)

var template = `{"line": "%s","source": "generated","tag": null,"logtypeid": 5,"timings": []}`
var templateSize = len(template) - 2
var url = flag.String("url", "", "Endpoint url")
var conns = flag.Int("c", 0, "Num connections")

func parseAndValidateArgs() {
    flag.Parse()
    if *url == "" {
        fmt.Println("--url flag is required")
        os.Exit(1)
    }

    if *conns == 0 {
        fmt.Println("--c flag is required")
        os.Exit(1)
    }
}

func main() {
    parseAndValidateArgs()
    ctrlc := make(chan os.Signal, 1)
    signal.Notify(ctrlc, os.Interrupt)

    endpoint := &EndpointInfo{"", *url, 1700}

    counts := make(chan int, *conns)
    go func() {
        count := 0
        start := time.Now()
        for _ = range counts {
            count++
            elapsed := time.Since(start)
            if elapsed.Seconds() >= 1 {
                fmt.Println(fmt.Sprintf("Per Second %d", count))
                count = 0
                start = time.Now()
            }
        }
    }()

    message := fmt.Sprintf(template, randomString(endpoint.Size-templateSize))
    for i := 0; i < *conns; i++ {
        go produce(endpoint.Framework, endpoint.Endpoint, message, counts)
    }

    <-ctrlc
    fmt.Println("Shutdown triggered")
}

func produce(framework string, endpoint string, message string, counts chan int) {
    tr := &http.Transport{}
    client := &http.Client{Transport: tr}

    messageBytes := []byte(message)
    for {
        request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(messageBytes))
        if err != nil {
            fmt.Println("Create request error:", err)
            continue
        }
        request.Header.Set("Content-Type", "application/json")
        res, err := client.Do(request)
        if err != nil {
            fmt.Println("POST error:", err)
            continue
        }
        io.Copy(ioutil.Discard, res.Body)
        res.Body.Close()

        counts <- 1
    }
}

func randomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)[:n]
}

type EndpointInfo struct {
    Framework string
    Endpoint  string
    Size      int
}

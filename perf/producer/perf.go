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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
)

var template = `{"line": "%s","source": "generated","tag": null,"logtypeid": 5,"timings": []}`
var templateSize = len(template) - 2
var config = flag.String("config", "", "Config json file")

func parseAndValidateArgs() {
	flag.Parse()
	if *config == "" {
		fmt.Println("--config flag is required")
		os.Exit(1)
	}
}

func main() {
	parseAndValidateArgs()
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	contents, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}
	endpoints := make([]*EndpointInfo, 0)
	json.Unmarshal(contents, &endpoints)

	for _, endpoint := range endpoints {
		message := fmt.Sprintf(template, randomString(endpoint.Size-templateSize))
		go produce(endpoint.Framework, endpoint.Endpoint, message)
	}

	<-ctrlc
	fmt.Println("Shutdown triggered")
}

func produce(framework string, endpoint string, message string) {
	messageBytes := []byte(message)
	for {
		_, err := http.Post(endpoint, "application/json", bytes.NewBuffer(messageBytes))
		if err != nil {
			panic(err)
		}
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

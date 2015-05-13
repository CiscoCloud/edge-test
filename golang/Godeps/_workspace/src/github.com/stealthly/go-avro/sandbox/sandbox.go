package main

import (
    "fmt"
    "encoding/json"
)

var rawJson = `[[{"text":"hey"}],"1231212412423235","channelName"]`

type PubNubMessage struct {
    Tags map[string]string
    Id string
    Channel string
}

func (pn *PubNubMessage) UnmarshalJSON(data []byte) error {
    js := make([]interface{}, 0)
    err := json.Unmarshal(data, &js)
    if err != nil {
        return err
    }

    pn.Tags = make(map[string]string)
    for key, value := range js[0].([]interface{})[0].(map[string]interface{}) {
        pn.Tags[key] = value.(string)
    }
    pn.Id = js[1].(string)
    pn.Channel = js[2].(string)

    fmt.Printf("%#v\n", js)

    return nil
}

func main() {
    pubNub := &PubNubMessage{}
    json.Unmarshal([]byte(rawJson), pubNub)

    fmt.Println(pubNub)
}

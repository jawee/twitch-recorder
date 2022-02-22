package configuration

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Configuration struct {
    ClientId string `json:"client-id"`
    ClientSecret string `json:"client-secret"`
    StreamersString string `json:"streamers"`
    Streamers []string
}


func New() *Configuration {
    var configuration Configuration
    pwd, _ := os.Getwd()

    log.Println("Loading configuration from " + pwd + "/config.json")
    file, err := os.Open(pwd + "/config.json") 
    if err != nil { 
        return nil 
    }  
    decoder := json.NewDecoder(file) 
    err = decoder.Decode(&configuration) 
    if err != nil {  
       return nil 
    }

    if configuration.ClientId == "" || configuration.ClientSecret == "" {
        log.Println("ClientId or ClientSecret is empty")
        return nil
    }

    streamers := strings.Replace(configuration.StreamersString, " ", "", -1)
    configuration.Streamers = strings.Split(streamers, ",")
    return &configuration
}

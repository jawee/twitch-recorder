package configuration

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
)

type Configuration struct {
    ClientId string `json:"client-id"`
    ClientSecret string `json:"client-secret"`
    StreamersString string `json:"streamers"`
    Streamers []string
}

type ConfigurationProvider interface {
    GetConfigurationJson() ([]byte, error)
}

type FileConfigurationProvider struct {
}

func (f *FileConfigurationProvider)GetConfigurationJson() ([]byte, error) {
    pwd := "/config"
    log.Println("FileConfigurationProvider.GetConfigurationJson. Loading configuration from " + pwd + "/config.json")
    path := path.Join(pwd, "config.json")
    file, err := os.Open(path)

    if err != nil {
        log.Printf("FileConfigurationProvider.GetConfigurationJson. Error opening configuration file: %s", err)
        return nil, err
    }
    bytes := make([]byte, 1024)

    readTotal, err := file.Read(bytes)

    if err != nil {
        log.Printf("FileConfigurationProvider.GetConfigurationJson. Error reading configuration file: %s", err)
        return nil, err
    }
    
    return bytes[:readTotal], nil
}


func New(configProvider ConfigurationProvider) *Configuration {
    log.Println("New. Loading configuration")    
    // pwd, _ := os.Getwd()
    bytes, err := configProvider.GetConfigurationJson()

    if err != nil {
        log.Printf("Error getting configuration from provider: %s", err)
        return nil
    }

    log.Println("Printing bytes")
    log.Println(string(bytes))

    var configuration *Configuration
    err = json.Unmarshal(bytes, &configuration) 
    if err != nil {  
        log.Printf("Error unmarshalling configuration: %s", err)
        return nil 
    }

    if configuration.ClientId == "" || configuration.ClientSecret == "" {
        log.Println("ClientId or ClientSecret is empty")
        return nil
    }

    streamers := strings.Replace(configuration.StreamersString, " ", "", -1)
    configuration.Streamers = strings.Split(streamers, ",")
    return configuration
}

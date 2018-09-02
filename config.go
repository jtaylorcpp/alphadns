package alphadns

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DomainConfig struct {
	Name      string   `yaml:"name"`
	Addresses []string `yaml:"addresses"`
}

type DNSServerConfig struct {
	ServerAddress string         `yaml:"serverAddress"`
	ServerPort    string         `yaml:"serverPort"`
	TTL           uint32         `yaml:"ttl"`
	Domains       []DomainConfig `yaml:"domains"`
}

// Parse a already read yaml file into the ConfigFile
func parseConfigFile(file string) *DNSServerConfig {
	newConfig := &DNSServerConfig{}
	err := yaml.Unmarshal([]byte(file), newConfig)
	if err != nil {
		log.Fatal("Unable to parse provided config file")
	}
	return newConfig
}

func DNSServerFromConfig(config *DNSServerConfig) *DNSServer {
	newServer := &DNSServer{
		ServerAddress: config.ServerAddress,
		ServerPort:    config.ServerPort,
		DNSRecords: &DNSGraph{
			Roots: make(map[Domain]*DNSNode),
		},
		TTL: config.TTL,
	}

	for _, domainConfig := range config.Domains {
		log.Println("domain server adding record: ", domainConfig)
		newServer.AddRecord(domainConfig.Name, domainConfig.Addresses)
	}

	return newServer
}

// Load a yaml file and parse out DNS servers
func LoadFromFile(filePath string) *DNSServer {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to parse yaml file")
	}

	newConfig := &DNSServerConfig{}
	err = yaml.Unmarshal(yamlFile, newConfig)

	if err != nil {
		log.Fatal(err.Error())
	}

	newServer := DNSServerFromConfig(newConfig)

	return newServer
}

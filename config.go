package alphadns

import (
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
		DNSRecords:    make(map[Domain]DNSNode),
	}

	for _, domainConfig := range config.Domains {
		log.Println("domain server adding record: ", domainConfig)
		newServer.AddRecord(domainConfig.Name, domainConfig.Addresses)
	}

	return newServer
}

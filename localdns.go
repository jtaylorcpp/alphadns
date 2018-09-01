package alphadns

import (
	"log"
	"net"
)

type Domain string

type DNSNode interface {
	GetIPAddresses([]Domain) []*net.IP
}

type DNSLeaf struct {
	Name      Domain
	Addresses []*net.IP
}

func (dl *DNSLeaf) GetIPAddresses(domains []Domain) []*net.IP {
	if len(domains) == 1 && domains[0] == dl.Name {
		return dl.Addresses
	} else if len(domains) == 1 && dl.Name == "*" {
		return dl.Addresses
	} else {
		log.Printf("Leaf Node %v found while searching for %v\n", dl.Name, domains)
		return []*net.IP{}
	}
}

type DNSBranch struct {
	Name       Domain
	SubDomains map[Domain]DNSNode
}

func (db *DNSBranch) GetIPAddresses(domains []Domain) []*net.IP {
	log.Print("looking up: ", domains)
	if domains[0] == db.Name && len(domains) > 1 {
		log.Print("searching for: ", domains[1:])
		if moreDomains, ok := db.SubDomains[domains[1]]; ok {
			if len(domains) > 1 {
				return moreDomains.GetIPAddresses(domains[1:])
			} else {
				return moreDomains.GetIPAddresses([]Domain{})
			}
		} else if moreDomains, ok := db.SubDomains["*"]; ok {
			if len(domains) > 1 {
				return moreDomains.GetIPAddresses(domains[1:])
			} else {
				return moreDomains.GetIPAddresses([]Domain{})
			}
		} else {
			return []*net.IP{}
		}
	} else {
		return []*net.IP{}
	}
}

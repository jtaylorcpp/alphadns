package alphadns

import (
	"log"
	"net"
	"strings"
)

type DNSServer struct {
	ServerAddress string
	ServerPort    string
	DNSRecords    map[Domain]DNSNode
}

func (ds *DNSServer) AddRecord(domain string, ips []string) {
	domains := reverse(strings.Split(domain, "."))
	log.Println("adding domain record: ", domains, ips)

	// traverse tree and add new leaf
	var current DNSBranch
	for idx, domain := range domains {
		// get current node
		switch idx {
		case 0:
			//root node
			if node, ok := ds.DNSRecords[Domain(domain)]; ok {
				// root node exists
				current = node
			} else {
				newRoot := &DNSBranch{
					Name:       Domain(domain),
					SubDomains: make(map[Domain]DNSNode),
				}
				current = newRoot
				ds.DNSRecords[Domain(domain)] = newRoot
			}
		case len(domains) - 1:
			//leaf node
			if leaf, ok := current.SubDomains[Domain(domain)]; ok {
				// leaf already exists
			} else {
				newLeaf := &DNSLeaf{
					Name:      Domain(domain),
					Addresses: make([]*net.IP, len(ips)),
				}
				for ipIdx, ip := range ips {
					newIP := net.ParseIP(ip).To4()
					newLeaf.Addresses[ipIdx] = &newIP
				}
				current.SubDomains[Domain(domain)] = newLeaf
			}
		default:
			//branches between root and leaf
			if node, ok := current.SubDomains[Domain(domain)]; ok {
				current = node
			} else {
				newBranch := &DNSBranch{
					Name:       Domain(domain),
					SubDomains: make(map[Domain]DNSNode),
				}
				current.SubDomains[Domain(domain)] = newBranch
				current = newBranch

			}
		}
	}
}

func reverse(s []string) []string {
	newSlice := make([]string, len(s))
	for idx, val := range s {
		newSlice[len(s)-idx-1] = val
	}
	return newSlice
}

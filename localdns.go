package alphadns

import (
	"log"
	"net"
	"strings"
)

type Domain string

// A graph structure used to manage multiple DNS entries/records in both a
// memory and computaitonally convienient way
type DNSGraph struct {
	Roots map[Domain]*DNSNode
}

// Node for the DNSGraph. Although one stuct, there are two different uses of
// this struct. 1) as a node which will have childrens and therefore
// potentially no addresses. 2) as a child node with 1:N addresses assigned.
type DNSNode struct {
	Name       Domain
	SubDomains map[Domain]*DNSNode
	Addresses  []*net.IP
}

// Takes a domain string ("test.example.com") and a set of IP addresses
// (["10.0.0.1","10.0.0.2"]) and will build out the necessary paths in
// DNSGraph to make them available.
func (dg *DNSGraph) AddRecord(domainString string, ips []string) {
	domains := reverse(strings.Split(domainString, "."))
	log.Println("adding record: ", domainString, ips)
	var currentNode *DNSNode
	if root, ok := dg.Roots[Domain(domains[0])]; ok {
		log.Println("using root: ",root)
		currentNode = root
	} else {
		newRoot := &DNSNode{
			Name:       Domain(domains[0]),
			SubDomains: make(map[Domain]*DNSNode),
			Addresses:  []*net.IP{},
		}
		dg.Roots[Domain(domains[0])] = newRoot
		currentNode = newRoot
		log.Println("adding new root: ",newRoot)
	}

	for idx, domain := range domains[1:] {
		log.Println("evaluating domain: ",domain)
		switch idx {
		case len(domains) - 2:
			newNode := &DNSNode{
				Name:       Domain(domain),
				SubDomains: map[Domain]*DNSNode{},
				Addresses:  make([]*net.IP, len(ips)),
			}
			for ipIdx, ip := range ips {
				nextIP := net.ParseIP(ip).To4()
				newNode.Addresses[ipIdx] = &nextIP
			}
			log.Println("adding leaf to node: ",currentNode)
			log.Println("new leaf: ",newNode)
			currentNode.SubDomains[Domain(domain)] = newNode

		default:
			if node, ok := currentNode.SubDomains[Domain(domain)]; ok {
				currentNode = node
			} else {
				newNode := &DNSNode{
					Name:       Domain(domain),
					SubDomains: make(map[Domain]*DNSNode),
					Addresses:  []*net.IP{},
				}
				currentNode.SubDomains[Domain(domain)] = newNode
				currentNode = newNode
			}
		}
	}
}

// Take a string domain name ("test.example.com") and returns a set of
// assigned IP addresses. This allows for a singular IP or a round robin DNS
// to be set up.
func (dg *DNSGraph) GetIPAddresses(domainString string) ([]*net.IP, bool) {
	domains := removeWhitespace(
		reverse(strings.Split(domainString, ".")))
	log.Println("looking up domains: ", domains)
	if root, ok := dg.Roots[Domain(domains[0])]; ok {
		log.Println("found root: ",root)
		var current *DNSNode = root
		for _, domain := range domains[1:] {
			if node, ok := current.SubDomains[Domain(domain)]; ok {
				// if subdomain exists
				log.Println("next subdomain found: ",node)
				current = node
			} else {
				// subdomain doesnt esist
				if node, ok := current.SubDomains[Domain("*")]; ok {
					// wildcard
					log.Println(
						"wildcard subdomain found: ",
						node)
					current = node
				} else {
					// no matching domain
					log.Println("no further domains")
					return []*net.IP{}, false
				}
			}
		}
		return current.Addresses, true
	} else {
		log.Println("domain doesnt exist: ", domains)
		return []*net.IP{}, false
	}
}

func reverse(s []string) []string {
	newSlice := make([]string, len(s))
	for idx, val := range s {
		newSlice[len(s)-idx-1] = val
	}
	return newSlice
}

func removeWhitespace(s []string) []string {
	count := 0
	for _, item := range s {
		if len(strings.TrimSpace(item)) == 0 {
			count += 1
		} else {
			break
		}
	}

	newSlice := make([]string, len(s) - count)
	for idx, item := range s[count:] {
		newSlice[idx] = item
	}

	return newSlice
}

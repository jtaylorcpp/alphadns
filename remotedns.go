package alphadns

import (
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type RemoteResolver struct {
	Timeout time.Duration
}

type RemoteResponse struct {
	msg        *dns.Msg
	nameserver string
	ttl        time.Duration
}

func (rr *RemoteResolver) Resolve(nameservers []string, net string, r *dns.Msg) (*dns.Msg, bool) {

	if len(nameservers) == 0 {
		return &dns.Msg{}, false
	}

	if net == "udp" {
		r.SetEdns0(65535, true)
	}

	if len(r.Question) > 1 {
		log.Println("remote resolver has more than 1 question: ",
			r.Question)
	}

	resolved := make(chan *RemoteResponse, len(nameservers))
	var wg sync.WaitGroup

	for _, ns := range nameservers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &dns.Client{
				Net:          net,
				ReadTimeout:  rr.Timeout,
				WriteTimeout: rr.Timeout,
			}

			msg, rtt, err := client.Exchange(r, ns)

			if err != nil {
				log.Println("error: ", ns, err.Error())
			}

			resolved <- &RemoteResponse{
				nameserver: ns,
				msg:        msg,
				ttl:        rtt,
			}

			return

		}()
	}

	wg.Wait()
	close(resolved)

	for resp := range resolved {
		if resp.msg != nil && resp.msg.Rcode != dns.RcodeSuccess {

		}
	}

	return &dns.Msg{}, true
}

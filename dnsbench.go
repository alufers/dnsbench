package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/miekg/dns"
)

var servers = []string{
	"8.8.8.8",
	"8.8.4.4",
	"1.1.1.1",
}

var sites = []string{
	"microsoft.com",
	"google.com",
	"grzywok.eu",
	"github.com",
	"qsdafcnjdhuiohdfsiopfsdvbiyasedfxd8dfvwqeoawdddssdsddddddd.github.io",
}

var repeats = 10

var stats = map[string][]time.Duration{}

func runBenchmark(target, server string) (time.Duration, error) {

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)

	r, t, err := c.Exchange(&m, server+":53")
	if err != nil {
		return 0, err
	}
	// log.Printf("Took %v", t)
	if len(r.Answer) == 0 {
		return 0, errors.New("empty anwser")
	}
	return t, nil
}

func main() {
	for i := 0; i < repeats; i++ {
		for _, serv := range servers {
			fmt.Printf("Testing server %s:\n", serv)
			for _, site := range sites {
				dur, err := runBenchmark(site, serv)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("  %s - %v\n", site, dur.String())
				stats[serv+": "+site] = append(stats[serv+": "+site], dur)
			}
			fmt.Printf("---\n\n")
		}
	}

	fmt.Printf("Average results for %d repeats, %d servers and %d domain names\n", repeats, len(servers), len(sites))
	lines := []string{}
	for statName, durations := range stats {
		var sum time.Duration

		for _, dur := range durations {
			sum += dur
		}

		avg := sum / time.Duration(len(durations))

		lines = append(lines, fmt.Sprintf("%s - %s\n", statName, avg.String()))
	}

	sort.Strings(lines)

	for _, l := range lines {
		fmt.Print(l)
	}

}

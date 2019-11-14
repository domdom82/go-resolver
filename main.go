package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var ns = flag.String("ns", "", "name server to use (default: use /etc/resolv.conf)")
	var duration = flag.Duration("duration", 1*time.Minute, "how long do you want to query the name server? (default: 1m)")
	var rate = flag.Int("rate", 1, "how many DNS queries per second? (default: 1")

	flag.Parse()

	if flag.NArg() == 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <host>\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Options are:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var host = flag.Arg(0)

	if *ns != "" {
		//TODO
	}

	tickCount := 0
	tick := time.Duration(float64(1) / float64(*rate) * float64(time.Second))
	tickChan := time.Tick(tick)
	endChan := time.After(*duration)

	for {
		select {
		case <-endChan:
			fmt.Printf("\nExiting.\n")
			os.Exit(0)
		case <-tickChan:
			tickCount++
			tStart := time.Now()
			addrs, err := net.LookupHost(host)
			tEnd := time.Now()
			tDuration := tEnd.Sub(tStart)

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v  %v (count=%d, duration=%v)\n", tEnd.Format(time.RFC3339), err, tickCount, tDuration)
			} else {
				fmt.Printf("\r%v  %v  (count=%d, duration=%v)", time.Now().Format(time.RFC3339), addrs, tickCount, tDuration)
			}
		}
	}

}

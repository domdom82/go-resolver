package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var godns = flag.Bool("godns", true, "use go dns resolver (true) or native getaddrinfo (false)")
	var duration = flag.Duration("duration", 10*time.Second, "how long do you want to query the name server?")
	var rate = flag.Int("rate", 1, "how many DNS queries per second?")

	flag.Parse()

	if flag.NArg() == 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <host>\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Options are:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var host = flag.Arg(0)

	net.DefaultResolver.PreferGo = *godns

	tickCount := int64(0)
	tickTotal := int64(*duration) / int64(time.Second) * int64(*rate)
	tick := time.Duration(float64(1) / float64(*rate) * float64(time.Second))
	tickChan := time.Tick(tick)
	endChan := time.After(*duration)
	var tStart, tEnd time.Time
	var tDuration time.Duration
	var tDurationTotal time.Duration

	for {
		select {
		case <-endChan:
			fmt.Printf("\nExiting.\n")
			fmt.Printf("Resolved %d/%d queries in %s. Actual rate was %.2f/s with an average query duration of %s\n",
				tickCount, tickTotal, *duration, float64(tickCount)/float64(*duration/time.Second), tDurationTotal/time.Duration(tickCount))
			os.Exit(0)
		case <-tickChan:
			tickCount++
			tStart = time.Now()
			addrs, err := net.LookupHost(host)
			tEnd = time.Now()
			tDuration = tEnd.Sub(tStart)
			tDurationTotal += tDuration

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v  %v (count=%d, duration=%v)\n", tEnd.Format(time.RFC3339), err, tickCount, tDuration)
			} else {
				fmt.Printf("\r%v  %v  (count=%d, duration=%v)", time.Now().Format(time.RFC3339), addrs, tickCount, tDuration)
			}
		}
	}

}

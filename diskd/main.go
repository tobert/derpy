package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var listenFlag string

func init() {
	flag.StringVar(&listenFlag, "listen", "127.0.0.1:8000", "")
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", listenFlag)
	if err != nil {
		log.Fatalf("Could not listen on address '%s': %s\n", listenFlag, err)
	}
	log.Printf("Listening on address: %s\n", listenFlag)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}
		go printStats(conn)
	}
}

func printStats(conn io.Writer) {
	for {
		select {
		case <-time.After(time.Second):
			stats := ReadDiskstats()
			for _, st := range stats {
				//                t  0  1  2  3  4  5  6  7  8  9 10 11 12 13
				fmt.Fprintf(conn, "%d,%d,%d,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
					st.Time.UnixNano(), st.Major, st.Minor, st.Name, // t,0,1,2
					st.ReadComplete, st.ReadMerged, st.ReadSectors, st.ReadMs, // 3,4,5,6
					st.WriteComplete, st.WriteMerged, st.WriteSectors, st.WriteMs, // 7,8,9,10
					st.IOPending, st.IOMs, st.IOQueueMs) // 11,12,13
			}
		}
	}
}

package main

/*
 * Copyright 2014 Albert P. Tobey <atobey@datastax.com> @AlTobey
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var cass *gocql.Session
var portFlag int
var ksFlag, cqlFlag string

type Metric struct {
	Time          time.Time     // timestamp
	Major         uint          //  0: major dev no
	Minor         uint          //  1: minor dev no
	Name          string        //  2: device name
	ReadComplete  uint64        //  3: reads completed
	ReadMerged    uint64        //  4: writes merged
	ReadSectors   uint64        //  5: sectors read
	ReadMs        uint          //  6: ms spent reading
	WriteComplete uint64        //  7: reads completed
	WriteMerged   uint64        //  8: writes merged
	WriteSectors  uint64        //  9: sectors read
	WriteMs       uint          // 10: ms spent writing
	IOPending     uint          // 11: number of IOs currently in progress
	IOMs          uint          // 12: jiffies_to_msecs(part_stat_read(hd, io_ticks))
	IOQueueMs     uint          // 13: jiffies_to_msecs(part_stat_read(hd, time_in_queue))
	Duration      time.Duration // only used in deltas
	RawTime       int64         // raw timestamp from C*
}

func init() {
	flag.StringVar(&ksFlag, "ks", "derpy", "keyspace containing stats")
	flag.StringVar(&cqlFlag, "cql", "127.0.0.1", "CQL address")
	flag.IntVar(&portFlag, "port", 8880, "HTTP port")
}

func main() {
	flag.Parse()

	cluster := gocql.NewCluster(cqlFlag)
	cluster.Keyspace = ksFlag
	cluster.Consistency = gocql.Quorum

	var err error
	cass, err = cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("Error creating Cassandra session: %v", err))
	}
	defer cass.Close()

	r := mux.NewRouter()
	r.HandleFunc("/1.0/metric/", MetricsHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", portFlag), nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := []Metric{}
	/*
		metrics := []Metric{
			Metric{time.Now(), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second * 2), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second * 3), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second * 4), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second * 5), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
			Metric{time.Now().Add(time.Second * 6), 2, 3, "sda", 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16},
		}
	*/

	query := `SELECT time,major,minor,name,rcp,rmg,rsc,rms,wcp,wmg,wsc,wms,iop,ioms,ioq FROM metrics WHERE name=?`
	iq := cass.Query(query, "sda").Iter()
	for {
		m := Metric{}
		ok := iq.Scan(
			&m.RawTime,
			&m.Major,
			&m.Minor,
			&m.Name,
			&m.ReadComplete,
			&m.ReadMerged,
			&m.ReadSectors,
			&m.ReadMs,
			&m.WriteComplete,
			&m.WriteMerged,
			&m.WriteSectors,
			&m.WriteMs,
			&m.IOPending,
			&m.IOMs,
			&m.IOQueueMs)

		if ok {
			m.Time = time.Unix(m.RawTime/1000000000, m.RawTime%1000000000)
			metrics = append(metrics, m)
		} else {
			break
		}
	}
	if err := iq.Close(); err != nil {
		log.Printf("Query failed: %s\n", err)
		return
	}

	jsonOut(w, r, metrics)
}

func jsonOut(w http.ResponseWriter, r *http.Request, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

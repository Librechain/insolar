//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	logDir         = ".artifacts/launchnet/logs/"
	statLogMessage = "stat_log_message"

	typeSent         = "sent"
	typeReply        = "reply"
	typeCallStarted  = "cr_call_started"
	typeCallReturned = "cr_call_returned"

	tracePerBucketSample = 3
)

var pattern = regexp.MustCompile(".*output.log$")

var ignoreTracePrefixes = []string{
	"main_",
	"pulse_",
}

type StatLog struct {
	StatType    string    `json:"stat_type"`
	TraceID     string    `json:"traceid"`
	Message     string    `json:"message"`
	MessageType string    `json:"message_type"`
	ReplyTimeMS float32   `json:"reply_time_ms"`
	Time        time.Time `json:"time"`
}

type Stats struct {
	lock sync.RWMutex
	// Trace id to stat.
	stats map[string]*TraceStats
}

func NewStats() *Stats {
	return &Stats{stats: map[string]*TraceStats{}}
}

func (s *Stats) GetOrCreate(trace string) *TraceStats {
	s.lock.Lock()
	defer s.lock.Unlock()

	if stat, ok := s.stats[trace]; ok {
		return stat
	}

	stat := NewTraceStats(trace)
	s.stats[trace] = stat
	return stat
}

type TraceStats struct {
	sync.RWMutex
	TraceID                   string
	First, CallReturned, Last time.Time
	// Message type to reply times.
	ReplyTimings map[string][]float32
	// Message type to sent count.
	SentCounts map[string]uint64
}

func NewTraceStats(trace string) *TraceStats {
	return &TraceStats{
		TraceID:      trace,
		ReplyTimings: map[string][]float32{},
		SentCounts:   map[string]uint64{},
	}
}

func (s *TraceStats) CallDuration() time.Duration {
	if s.CallReturned.IsZero() {
		return 0
	}
	return s.CallReturned.Sub(s.First)
}

func (s *TraceStats) TotalDuration() time.Duration {
	return s.Last.Sub(s.First)
}

func (s *TraceStats) Parse(log StatLog) {
	s.Lock()
	defer s.Unlock()

	if s.First.IsZero() || s.First.After(log.Time) {
		s.First = log.Time
	}
	if s.Last.IsZero() || s.Last.Before(log.Time) {
		s.Last = log.Time
	}

	switch log.StatType {
	case typeSent:
		s.SentCounts[log.MessageType]++
	case typeReply:
		s.ReplyTimings[log.MessageType] = append(s.ReplyTimings[log.MessageType], log.ReplyTimeMS)
	case typeCallReturned:
		s.CallReturned = log.Time
	}
}

func main() {
	var parsedCount uint64

	shouldParse := func(log StatLog) bool {
		if log.Message != statLogMessage {
			return false
		}
		if log.TraceID == "" {
			return false
		}
		for _, i := range ignoreTracePrefixes {
			if strings.HasPrefix(log.TraceID, i) {
				return false
			}
		}

		return true
	}

	parseFile := func(stats *Stats, filename string) {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			if !bytes.Contains(line, []byte(statLogMessage)) {
				continue
			}

			log := StatLog{}
			err = json.Unmarshal(line, &log)
			if err != nil {
				continue
			}

			if !shouldParse(log) {
				continue
			}

			atomic.AddUint64(&parsedCount, 1)
			stats.GetOrCreate(log.TraceID).Parse(log)
		}
	}

	var files []string
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !pattern.MatchString(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	checkError(err)

	stats := NewStats()
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, path := range files {
		path := path
		go func() {
			parseFile(stats, path)
			wg.Done()
		}()
	}
	wg.Wait()

	aggregate := MetaAggregate{
		TotalTraces: len(stats.stats),
		TotalLogs:   parsedCount,
		Aggregates: []Aggregator{
			NewAggSent(),
			NewAggReply(),
			&AggTraceTimes{},
		},
	}
	aggregate.Aggregate(stats)
	out := bufio.NewWriter(os.Stdout)
	_, err = out.Write([]byte(aggregate.String()))
	checkError(err)
	err = out.Flush()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Aggregator interface {
	fmt.Stringer
	Aggregate(*Stats)
}

type MetaAggregate struct {
	TotalTraces int
	TotalLogs   uint64
	Aggregates  []Aggregator
}

func (a *MetaAggregate) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Analysed %d logs in %d traces.\n", a.TotalLogs, a.TotalTraces))
	for _, a := range a.Aggregates {
		b.WriteString(a.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (a *MetaAggregate) Aggregate(s *Stats) {
	for _, a := range a.Aggregates {
		a.Aggregate(s)
	}
}

type AggSent struct {
	// Message type to count.
	AVGCount map[string]float64 `json:"avg_count"`
}

func (a *AggSent) String() string {
	b := strings.Builder{}

	type cnt struct {
		msgType string
		count   float32
	}
	counts := make([]cnt, 0, len(a.AVGCount))
	for msgType, count := range a.AVGCount {
		counts = append(counts, cnt{msgType: msgType, count: float32(count)})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	b.WriteString("[Average count per trace]\n")
	for _, c := range counts {
		b.WriteString("    " + c.msgType + ": " + fmt.Sprintf("%.1f", c.count) + "\n")
	}

	return b.String()
}

func (a *AggSent) Aggregate(logStats *Stats) {
	var aggs []*AggSent
	for _, stats := range logStats.stats {
		agg := NewAggSent()
		for msgType, count := range stats.SentCounts {
			agg.AVGCount[msgType] = float64(count)
		}
		aggs = append(aggs, agg)
	}

	for _, agg := range aggs {
		for msgType, count := range agg.AVGCount {
			a.AVGCount[msgType] += count
		}
	}

	avgDivider := len(aggs)
	for msgType := range a.AVGCount {
		a.AVGCount[msgType] /= float64(avgDivider)
	}
}

func NewAggSent() *AggSent {
	return &AggSent{AVGCount: map[string]float64{}}
}

type AggReply struct {
	// Message type to reply time.
	AVGReplyTime map[string]float64 `json:"avg_reply_time"`
}

func NewAggReply() *AggReply {
	return &AggReply{AVGReplyTime: map[string]float64{}}
}

func (a *AggReply) String() string {
	type rep struct {
		msgType   string
		replyTime int
	}
	reps := make([]rep, 0, len(a.AVGReplyTime))
	for msgType, replyTime := range a.AVGReplyTime {
		reps = append(reps, rep{msgType: msgType, replyTime: int(replyTime)})
	}
	sort.Slice(reps, func(i, j int) bool {
		return reps[i].replyTime > reps[j].replyTime
	})

	b := strings.Builder{}
	b.WriteString("[Average reply times per trace, ms]\n")
	for _, r := range reps {
		b.WriteString("    " + r.msgType + ": " + fmt.Sprintf("%d", r.replyTime) + "\n")
	}

	return b.String()
}

func (a *AggReply) Aggregate(logStats *Stats) {
	var aggs []*AggReply
	for _, stats := range logStats.stats {
		agg := NewAggReply()
		for msgType, timings := range stats.ReplyTimings {
			var summ float64
			for _, t := range timings {
				summ += float64(t)
			}
			agg.AVGReplyTime[msgType] = summ
		}
		aggs = append(aggs, agg)
	}

	for _, agg := range aggs {
		for msgType, replyTime := range agg.AVGReplyTime {
			a.AVGReplyTime[msgType] += replyTime
		}
	}

	avgDivider := len(aggs)
	for msgType := range a.AVGReplyTime {
		a.AVGReplyTime[msgType] /= float64(avgDivider)
	}
}

const (
	traceTimeBucketCount = 5
)

var traceTimeBuckets = [traceTimeBucketCount]time.Duration{
	time.Second,
	10 * time.Second,
	20 * time.Second,
	40 * time.Second,
	time.Minute,
}

type AggTraceTimes struct {
	callTimes  [traceTimeBucketCount + 1]uint64
	totalTimes [traceTimeBucketCount + 1]uint64

	callSamples  [traceTimeBucketCount + 1][]*TraceStats
	totalSamples [traceTimeBucketCount + 1][]*TraceStats
}

func (a *AggTraceTimes) String() string {
	writeCallSample := func(b *strings.Builder, bucket int) {
		if len(a.callSamples[bucket]) == 0 {
			return
		}
		for _, t := range a.callSamples[bucket] {
			b.WriteString(fmt.Sprintf("   %s (%s)\n", t.TraceID, t.CallDuration()))
		}
		b.WriteString("   . . .\n")
	}

	writeTotalSample := func(b *strings.Builder, bucket int) {
		if len(a.totalSamples[bucket]) == 0 {
			return
		}
		for _, t := range a.totalSamples[bucket] {
			b.WriteString(fmt.Sprintf("   %s (%s)\n", t.TraceID, t.TotalDuration()))
		}
		b.WriteString("   . . .\n")
	}

	b := strings.Builder{}
	b.WriteString("[Call return percentiles]\n")
	for i := 0; i < len(a.callTimes)-1; i++ {
		b.WriteString(fmt.Sprintf("< %s (%d) \n", traceTimeBuckets[i], a.callTimes[i]))
		writeCallSample(&b, i)
	}
	b.WriteString(fmt.Sprintf(
		"> %s (%d) \n",
		traceTimeBuckets[len(traceTimeBuckets)-1],
		a.callTimes[len(a.callTimes)-1]),
	)
	writeCallSample(&b, len(a.callTimes)-1)
	b.WriteString("\n")

	b.WriteString("[Total time percentiles]\n")
	for i := 0; i < len(a.totalTimes)-1; i++ {
		b.WriteString(fmt.Sprintf("< %s (%d) \n", traceTimeBuckets[i], a.totalTimes[i]))
		writeTotalSample(&b, i)
	}
	b.WriteString(fmt.Sprintf(
		"> %s (%d) \n",
		traceTimeBuckets[len(traceTimeBuckets)-1],
		a.callTimes[len(a.totalTimes)-1]),
	)
	writeTotalSample(&b, len(a.totalTimes)-1)

	return b.String()
}

func (a *AggTraceTimes) add(stat *TraceStats) {
	addCall := func(dur time.Duration) int {
		for i, t := range traceTimeBuckets {
			if dur < t {
				a.callTimes[i]++
				return i
			}
		}
		return len(traceTimeBuckets)
	}

	addTotal := func(dur time.Duration) int {
		for i, t := range traceTimeBuckets {
			if dur < t {
				a.totalTimes[i]++
				return i
			}
		}
		return len(traceTimeBuckets)
	}

	if stat.CallDuration() != 0 {
		bucket := addCall(stat.CallDuration())
		if len(a.callSamples[bucket]) < tracePerBucketSample {
			a.callSamples[bucket] = append(a.callSamples[bucket], stat)
		}
	}
	bucket := addTotal(stat.TotalDuration())
	if len(a.totalSamples[bucket]) < tracePerBucketSample {
		a.totalSamples[bucket] = append(a.totalSamples[bucket], stat)
	}
}

func (a *AggTraceTimes) Aggregate(logStats *Stats) {
	for _, stats := range logStats.stats {
		a.add(stats)
	}
}

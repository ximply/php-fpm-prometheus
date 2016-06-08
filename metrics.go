package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Metrics contains the status data collected from PHP-FPM.
type Metrics struct {
	StartSince         int `help:"Seconds since FPM start" type:"counter" name:"start_since"`
	AcceptedConn       int `help:"Total of accepted connections" type:"counter" name:"accepted_conn"`
	ListenQueue        int `help:"Number of connections that have been initiated but not yet accepted" type:"gauge" name:"listen_queue"`
	MaxListenQueue     int `help:"Max. connections the listen queue has reached since FPM start" type:"counter" name:"max_listen_queue"`
	ListenQueueLength  int `help:"Maximum number of connections that can be queued" type:"gauge" name:"listen_queue_length"`
	IdleProcesses      int `help:"Idle process count" type:"gauge" name:"idle_processes"`
	ActiveProcesses    int `help:"Active process count" type:"gauge" name:"active_processes"`
	TotalProcesses     int `help:"Total process count" type:"gauge" name:"total_processes"`
	MaxActiveProcesses int `help:"Maximum active process count" type:"counter" name:"max_active_processes"`
	MaxChildrenReached int `help:"Number of times the process limit has been reached" type:"counter" name:"max_children_reached"`
	SlowRequests       int `help:"Number of requests that exceed request_slowlog_timeout" type:"counter" name:"slow_requests"`
}

// NewMetricsFromMatches creates a new Metrics instance and populates it with given data.
func NewMetricsFromMatches(matches [][]string) *Metrics {
	metrics := &Metrics{}
	metrics.populateFromMatches(matches)
	return metrics
}

func (m *Metrics) populateFromMatches(matches [][]string) {
	for _, match := range matches {
		key := match[1]
		value := match[2]
		switch key {
		case "start since":
			m.StartSince, _ = strconv.Atoi(value)
		case "accepted conn":
			m.AcceptedConn, _ = strconv.Atoi(value)
		case "listen queue":
			m.ListenQueue, _ = strconv.Atoi(value)
		case "max listen queue":
			m.MaxListenQueue, _ = strconv.Atoi(value)
		case "listen queue len":
			m.ListenQueueLength, _ = strconv.Atoi(value)
		case "idle processes":
			m.IdleProcesses, _ = strconv.Atoi(value)
		case "active processes":
			m.ActiveProcesses, _ = strconv.Atoi(value)
		case "total processes":
			m.TotalProcesses, _ = strconv.Atoi(value)
		case "max active processes":
			m.MaxActiveProcesses, _ = strconv.Atoi(value)
		case "max children reached":
			m.MaxChildrenReached, _ = strconv.Atoi(value)
		case "slow requests":
			m.SlowRequests, _ = strconv.Atoi(value)
		}
	}
}

func (m *Metrics) WriteTo(w io.Writer) {
	typ := reflect.TypeOf(*m)
	val := reflect.ValueOf(*m)
	buf := &bytes.Buffer{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("name")
		buf.WriteString(fmt.Sprintf("# HELP %s %s\n", name, field.Tag.Get("help")))
		buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", name, field.Tag.Get("type")))
		buf.WriteString(fmt.Sprintf("%s %d\n", name, val.Field(i).Int()))
	}

	io.Copy(w, buf)
}

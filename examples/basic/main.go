package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	userName := ""
	password := ""
	// Create a new client using an InfluxDB server base URL and an authentication token
	// For authentication token supply a string in the form: "username:password" as a token. Set empty value for an unauthenticated server
	client := influxdb2.NewClient("http://localhost:8086", fmt.Sprintf("%s:%s", userName, password))
	// Get the blocking write client
	// Supply a string in the form database/retention-policy as a bucket. Skip retention policy for the default one, use just a database name (without the slash character)
	// Org name is not used
	writeAPI := client.WriteAPIBlocking("", "test/autogen")
	// create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	// Write data
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Printf("Write error: %s\n", err.Error())
	}

	client.Close()
}

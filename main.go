package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		glog.Fatalf("Error parsing env file: %v", err)
	}
	glog.Infof("Successfully loaded the env values")
}

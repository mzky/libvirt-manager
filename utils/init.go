package utils

import (
	"flag"
	"fmt"
	"os"
)

func Flag() {
	Port = *flag.String("p", "9100", "service port") //nolint
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println("SourceVersion: " + SourceVersion)
		fmt.Println("BuildTime: " + BuildTime)
		fmt.Println("Version: " + Version)
		os.Exit(0)
	}
}

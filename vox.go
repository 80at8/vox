// polly.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func displayCommandLine(serverAddress string) {

	fmt.Printf("Polly [http://%v]\n\n", serverAddress)
	fmt.Println("=== Endpoint Information ===\n")
	fmt.Printf("http://%v/convert\n", serverAddress)
	fmt.Printf("http://%v/showhelp\n", serverAddress)

}

// Global Variables.

var downloadServerURL string
var outputPath string

func main() {

	var port = flag.Int("port", 8080, "Port to listen on")
	var ip = flag.String("ip", "localhost", "IP to bind to")
	var path = flag.String("outdir", "./mp3/", "Path to save converted files to")
	var serverAddress bytes.Buffer

	flag.Parse()

	_, err := serverAddress.WriteString(*ip + ":" + strconv.Itoa(*port))

	if err != nil {
		panic(err)
	}

	downloadServerURL = serverAddress.String()
	outputPath = *path


	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	fs := http.FileServer(http.Dir(outputPath))

	http.Handle("/mp3/", http.StripPrefix("/mp3", fs))
	http.HandleFunc("/convert", convert)
	http.HandleFunc("/showhelp", showHTMLHelp)

	displayCommandLine(downloadServerURL)

	if err := http.ListenAndServe(serverAddress.String(), nil); err != nil {
		log.Fatal(err)
	}
}

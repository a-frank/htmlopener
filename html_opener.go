package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"io/ioutil"
	"os/exec"
	"fmt"
	"runtime"
)

var filePath string

func main() {
	filePath = determineFilePath()
	if filePath == "" {
		fmt.Errorf("can't get the file path for the operating system")
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/openHtml", OpenHtml).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func OpenHtml(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(writer, "Can't read the body of the request", http.StatusBadRequest)
		return
	}

	fileErr := ioutil.WriteFile(filePath, body, 0644)
	if fileErr != nil {
		log.Printf("Error writing html: %v", fileErr)
		http.Error(writer, "Can't write HTML from body of the request", http.StatusBadRequest)
		return
	}

	runErr := openBrowser(filePath)
	if runErr != nil {
		log.Printf("Error writing opening page: %v", runErr)
		http.Error(writer, "Can't open HTML", http.StatusInternalServerError)
		return
	}
}

func openBrowser(urlOrFilePath string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", urlOrFilePath).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", urlOrFilePath).Run()
	case "darwin":
		return exec.Command("open", urlOrFilePath).Run()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func determineFilePath() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return "/tmp/website.html"
	case "windows":
		return "C:\\tmp\\website.html"
	default:
		return ""
	}
}

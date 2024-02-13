package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/LucasNT/transmission-automation/internals/config"
	bitTorrentImplementation "github.com/LucasNT/transmission-automation/internals/externals/bit_torrent_implementations"
	TorrentCompletedHandler "github.com/LucasNT/transmission-automation/internals/externals/torrent_completed_handler"
	torrentdownloadedinformation "github.com/LucasNT/transmission-automation/internals/externals/torrent_downloaded_information"
	TorrentEntryReader "github.com/LucasNT/transmission-automation/internals/externals/torrent_entry_reader"
	"github.com/LucasNT/transmission-automation/internals/interfaces"
	useCases "github.com/LucasNT/transmission-automation/internals/use_cases"

	_ "embed"

	"github.com/sirupsen/logrus"
)

//go:embed root.tmpl
var root string
var transmission interfaces.BitTorrentclient

const http_port = "8080"
const CONFIG_PATH string = "./config.yaml"

func executeUseCase(byteReader io.Reader) {
	var reader TorrentEntryReader.CsvTorrentEntryReader
	var torrentDownloadedInformationVoid torrentdownloadedinformation.TorrentDownloadedInformationVoid
	var copy TorrentCompletedHandler.TorrentCompletedHandlerCopy
	reader = TorrentEntryReader.NewCsvTorrentEntryReader(byteReader)

	copy, err := TorrentCompletedHandler.NewTorrentCompletedHandlerCopy(config.Config.CopyHandler.TorrentPath, config.Config.CopyHandler.DestinyPath)
	if err != nil {
		logrus.Errorf("Faild to create a copy completed handler: %v", err)
		return
	}

	torrentDownloadedInformationVoid = torrentdownloadedinformation.TorrentDownloadedInformationVoid{}

	if err := useCases.ExecProgramn(transmission, copy, reader, torrentDownloadedInformationVoid, 1*time.Minute); err != nil {
		logrus.Errorf("Failed to execute use case: %v", err)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("root")
	tmpl.Parse(root)
	err := tmpl.Execute(w, "")
	if err != nil {
		http.NotFound(w, r)
	}
}

func serverError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "503 Internal Server Error")
}

func uploadCsvHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		logrus.Info("Upload request received")
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			logrus.Errorf("Faild to parse Multi part Form: %v", err)
			serverError(w, r)
			return
		}
		file, header, err := r.FormFile("uploadCsv")
		if err != nil {
			logrus.Errorf("Failed to read file: %v", err)
			serverError(w, r)
			return
		}
		defer file.Close()
		if header.Header.Get("Content-Type") == "text/csv" {
			fileContent, err := io.ReadAll(file)
			if err != nil {
				serverError(w, r)
				return
			}
			logrus.Info("Starting new thread for the use case")
			go executeUseCase(bytes.NewReader(fileContent))
			logrus.Info("Started new thread for the use case")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			serverError(w, r)
			return
		}
	default:
		http.NotFound(w, r)
		logrus.Error("Wrong method")
		return
	}
}

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error when oppened the log file")
	}
	defer logFile.Close()
	multi := io.MultiWriter(logFile, os.Stderr)
	logrus.SetOutput(multi)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Programn Started")

	if err = config.LoaderConfigs(CONFIG_PATH); err != nil {
		logrus.Fatal("Failed to load config file ", err)
	}
	logrus.Info("Loaded Settings")

	logrus.Info("Connecting to the bitTorrentClient")
	endpoint, err := url.Parse(config.Config.Url)
	if err != nil {
		logrus.Fatal("Failed to parse url ", err)
	}
	endpoint.User = url.UserPassword(config.Config.Username, config.Config.Password)

	err = fmt.Errorf("dummy")
	for err != nil {
		transmission, err = bitTorrentImplementation.NewTransmision(endpoint, nil)
		if err != nil {
			logrus.Fatal("Failed to Connect to bitTorrent Client ", err)
		}
	}
	defer transmission.Close()

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/", rootHandler)
	serverMux.HandleFunc("/upload_csv", uploadCsvHandler)

	s := &http.Server{
		Addr:         "0.0.0.0:" + http_port,
		Handler:      serverMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logrus.Infof("Web server started at port %s", http_port)
	err = s.ListenAndServe()
	if err != nil {
		logrus.Fatalf("Failed to start web server, %v", err)
	}

}

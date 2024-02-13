package main

import (
	"fmt"
	"io"
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
	log "github.com/sirupsen/logrus"
)

const CONFIG_PATH string = "./config.yaml"

func main() {
	var err error
	var bitTorrent interfaces.BitTorrentclient
	var torrentHandler interfaces.TorrentCompletedHandler
	var reader interfaces.TorrentEntryReader

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error when oppened the log file")
	}
	defer logFile.Close()
	multi := io.MultiWriter(logFile, os.Stderr)
	log.SetOutput(multi)
	log.SetLevel(log.DebugLevel)
	log.Info("Programn Started initialization")

	if len(os.Args) < 2 {
		log.Fatal("Need the path of the csv file")
	}

	if err = config.LoaderConfigs(CONFIG_PATH); err != nil {
		log.Fatal("Failed to load config file ", err)
	}
	log.Info("Loaded Settings")

	log.Info("Connecting to the bitTorrentClient")
	endpoint, err := url.Parse(config.Config.Url)
	if err != nil {
		log.Fatal("Failed to parse url ", err)
	}
	endpoint.User = url.UserPassword(config.Config.Username, config.Config.Password)

	bitTorrent, err = bitTorrentImplementation.NewTransmision(endpoint, nil)
	if err != nil {
		log.Fatal("Failed to Connect to bitTorrent Client ", err)
	}
	defer bitTorrent.Close()

	torrentHandler, err = TorrentCompletedHandler.NewTorrentCompletedHandlerCopy(config.Config.CopyHandler.TorrentPath, config.Config.CopyHandler.DestinyPath)
	if err != nil {
		log.Fatal("Failed to Create the TorrentCompletedHandler ", err)
	}
	log.Info("Connected to the bitTorrentClient")

	var csvFilePath string = os.Args[1]

	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal("Failed to Open the csv file ", err)
	}

	defer file.Close()

	reader = TorrentEntryReader.NewCsvTorrentEntryReader(file)

	log.Info("Programn finished initialization")
	err = useCases.ExecProgramn(bitTorrent, torrentHandler, reader, torrentdownloadedinformation.TorrentDownloadedInformationVoid{}, 1*time.Minute)
	if err != nil {
		log.Fatal("Programn failed ", err)
	}
}

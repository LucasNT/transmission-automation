package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/LucasNT/transmission-automation/config"
	bitTorrentImplementation "github.com/LucasNT/transmission-automation/externals/bit_torrent_implementations"
	TorrentCompletedHandler "github.com/LucasNT/transmission-automation/externals/torrent_completed_handler"
	CsvTorrentEntryReader "github.com/LucasNT/transmission-automation/externals/torrent_entry_reader"
	"github.com/LucasNT/transmission-automation/interfaces"
	useCases "github.com/LucasNT/transmission-automation/use_cases"
)

const CONFIG_PATH string = "./config.yaml"

func main() {
	var err error
	var bitTorrent interfaces.BitTorrentclient
	var torrentHandler interfaces.TorrentCompletedHandler
	var reader interfaces.TorrentEntryReader

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Need at least one argument")
		os.Exit(1)
	}

	if err = config.LoaderConfigs(CONFIG_PATH); err != nil {
		panic(err)
	}

	endpoint, err := url.Parse(config.Config.Url)
	if err != nil {
		panic(err)
	}
	endpoint.User = url.UserPassword(config.Config.Username, config.Config.Password)

	bitTorrent, err = bitTorrentImplementation.NewTransmision(endpoint, nil)
	torrentHandler, err = TorrentCompletedHandler.NewTorrentCompletedHandlerCopy(config.Config.CopyHandler.TorrentPath, config.Config.CopyHandler.DestinyPath)

	if err != nil {
		panic(err)
	}
	defer bitTorrent.Close()

	var csvFilePath string = os.Args[1]

	file, err := os.Open(csvFilePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader = CsvTorrentEntryReader.NewCsvTorrentEntryReader(file)

	fmt.Println("ola mundo")
	err = useCases.ExecProgramn(bitTorrent, torrentHandler, reader, 1*time.Minute)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"net/url"
	"os"

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
	var torrentHandler interfaces.TorrentCompletedHandler
	var bitTorrent interfaces.BitTorrentclient
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

	tempDir, err := os.MkdirTemp("", "mockTransmission")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)
	bitTorrent, err = bitTorrentImplementation.NewBitTorrentMock(tempDir)
	torrentHandler, err = TorrentCompletedHandler.NewTorrentCompletedHandlerCopy(tempDir, tempDir)

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
	err = useCases.ExecProgramn(bitTorrent, torrentHandler, reader, 0)
	if err != nil {
		panic(err)
	}

}

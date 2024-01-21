package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/LucasNT/transmission-automation/config"
	bitTorrentImplementation "github.com/LucasNT/transmission-automation/externals/bit_torrent_implementations"
	TorrentCompletedHandler "github.com/LucasNT/transmission-automation/externals/torrent_completed_handler"
	CsvTorrentEntryReader "github.com/LucasNT/transmission-automation/externals/torrent_entry_reader"
	"github.com/LucasNT/transmission-automation/interfaces"
)

const CONFIG_PATH string = "./config.yaml"

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Need at least one argument")
		os.Exit(1)
	}

	if err = config.LoaderConfigs(CONFIG_PATH); err != nil {
		panic(err)
	}

	fmt.Println(config.Config)

	endpoint, err := url.Parse(config.Config.Url)
	if err != nil {
		panic(err)
	}
	endpoint.User = url.UserPassword(config.Config.Username, config.Config.Password)

	var copy interfaces.TorrentCompletedHandler
	var bitTorrent interfaces.BitTorrentclient
	var reader interfaces.TorrentEntryReader
	tempDir, err := os.MkdirTemp("", "mockTransmission")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)
	bitTorrent, err = bitTorrentImplementation.NewBitTorrentMock(tempDir)
	copy, err = TorrentCompletedHandler.NewTorrentCompletedHandlerCopy(tempDir, tempDir)

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

	magnetLink, handlerString, errReadTorrent := reader.ReadTorrentEntry()

	for errReadTorrent == nil {
		tr_id, err := bitTorrent.TorrentAdd(magnetLink)
		if err != nil {
			panic(err)
		}
		fmt.Println(tr_id)
		percent := float64(0)
		for percent != 1 {
			time.Sleep(1 * time.Second)
			percent, err = bitTorrent.GetTorrentPercentComplete(tr_id)
			if err != nil {
				panic(err)
			}
			fileName, err := bitTorrent.GetTorrentName(tr_id)
			if err != nil {
				panic(err)
			}
			fmt.Println(fileName, percent)
		}
		listFileName, err := bitTorrent.GetTorrentFiles(tr_id)
		if err != nil {
			panic(err)
		}
		cmd, err := copy.CreateExec(handlerString)
		if err != nil {
			panic(err)
		}

		ok, err := cmd(listFileName)

		if !ok {
			fmt.Println(err)
		}

		if err != nil {
			var auxErr *exec.Error
			if errors.As(err, &auxErr) {
				fmt.Println("é um erro")
			}
			fmt.Println(err)
		}

		magnetLink, handlerString, errReadTorrent = reader.ReadTorrentEntry()
	}
	if errors.Is(errReadTorrent, interfaces.ErrNoTorrentEntry) {
		fmt.Println("fim da execulção")
	} else if errReadTorrent != nil {
		panic(errReadTorrent)
	}

}

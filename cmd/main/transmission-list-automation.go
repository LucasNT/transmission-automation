package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/LucasNT/transmission-automation/config"
	bitTorrentImplementation "github.com/LucasNT/transmission-automation/externals/bit_torrent_implementations"
	TorrentCompletedHandler "github.com/LucasNT/transmission-automation/externals/torrent_completed_handler"
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
	var bitTorrent interfaces.BitTorrentclient
	var torrentHandler interfaces.TorrentCompletedHandler
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

	reader := csv.NewReader(file)

	l, err := reader.Read()
	fmt.Println(l, err)

	for l != nil && err == nil {
		tr_id, err := bitTorrent.TorrentAdd(l[0])
		if err != nil {
			panic(err)
		}
		fmt.Println(tr_id)
		percent := float64(0)
		for percent != 1 {
			time.Sleep(1 * time.Minute)
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
		cmd, err := torrentHandler.CreateExec(l[1])

		if err != nil {
			fmt.Printf("Erro ao criar o comando de copiar %s", err.Error())
		}

		_, err = cmd(listFileName)

		if err != nil {
			fmt.Printf("Erro ao copiar o arquivo: %s", err.Error())
		}

		l, err = reader.Read()
	}
	if err != nil {
		panic(err)
	}

}

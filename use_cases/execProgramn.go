package useCases

import (
	"errors"
	"fmt"
	"time"

	"github.com/LucasNT/transmission-automation/interfaces"
)

func ExecProgramn(bitTorrentClient interfaces.BitTorrentclient, torrentCompletedHandler interfaces.TorrentCompletedHandler, torrentEntryReader interfaces.TorrentEntryReader, sleepTime time.Duration) error {

	magnetLink, handlerString, errReadTorrent := torrentEntryReader.ReadTorrentEntry()

	for errReadTorrent == nil {
		tr_id, err := bitTorrentClient.TorrentAdd(magnetLink)
		if err != nil {
			return err
		}
		fmt.Println(tr_id)
		percent := float64(0)
		for percent != 1 {
			time.Sleep(sleepTime)
			percent, err = bitTorrentClient.GetTorrentPercentComplete(tr_id)
			if err != nil {
				return err
			}
			fileName, err := bitTorrentClient.GetTorrentName(tr_id)
			if err != nil {
				return err
			}
			fmt.Println(fileName, percent)
		}
		listFileName, err := bitTorrentClient.GetTorrentFiles(tr_id)
		if err != nil {
			return err
		}
		cmd, err := torrentCompletedHandler.CreateExec(handlerString)

		if err != nil {
			fmt.Printf("Erro ao criar o comando de copiar %s", err.Error())
		}

		_, err = cmd(listFileName)

		if err != nil {
			fmt.Printf("Erro ao copiar o arquivo: %s", err.Error())
		}

		magnetLink, handlerString, errReadTorrent = torrentEntryReader.ReadTorrentEntry()
	}
	if errors.Is(errReadTorrent, interfaces.ErrNoTorrentEntry) {
		fmt.Println("fim da execulção")
	} else if errReadTorrent != nil {
		return errReadTorrent
	}
	return nil

}

package useCases

import (
	"errors"
	"time"

	"github.com/LucasNT/transmission-automation/internals/interfaces"
	log "github.com/sirupsen/logrus"
)

func ExecProgramn(bitTorrentClient interfaces.BitTorrentclient, torrentCompletedHandler interfaces.TorrentCompletedHandler, torrentEntryReader interfaces.TorrentEntryReader, sleepTime time.Duration) error {

	magnetLink, handlerString, errReadTorrent := torrentEntryReader.ReadTorrentEntry()

	for errReadTorrent == nil {
		tr_id, err := bitTorrentClient.TorrentAdd(magnetLink)
		log.Info("Torrent added successfully")
		if err != nil {
			return err
		}
		log.Debug(tr_id)
		percent := float64(0)
		for percent != 1 {
			time.Sleep(sleepTime)
			percent, err = bitTorrentClient.GetTorrentPercentComplete(tr_id)
			if err != nil {
				return err
			}
			torrentName, err := bitTorrentClient.GetTorrentName(tr_id)
			if err != nil {
				return err
			}
			log.Debugf("Torrent: '%s' is %f Downloaded", torrentName, percent)
		}
		listFileName, err := bitTorrentClient.GetTorrentFiles(tr_id)
		if err != nil {
			return err
		}
		cmd, err := torrentCompletedHandler.CreateExec(handlerString)

		if err != nil {
			log.Errorf("Erro ao criar o comando de copiar %s", err.Error())
		}

		_, err = cmd(listFileName)

		if err != nil {
			log.Errorf("Erro ao copiar o arquivo: %s", err.Error())
		}
		log.Info("Torrent Handled completed successfully")

		magnetLink, handlerString, errReadTorrent = torrentEntryReader.ReadTorrentEntry()
	}
	if errors.Is(errReadTorrent, interfaces.ErrNoTorrentEntry) {
		log.Info("fim da execulção")
	} else if errReadTorrent != nil {
		return errReadTorrent
	}
	return nil

}

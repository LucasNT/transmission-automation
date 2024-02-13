package torrentdownloadedinformation

import "math"

type TorrentDownloadedInformationVoid struct {
}

func (t TorrentDownloadedInformationVoid) GetTorrentInformation() (string, float64, error) {
	return "", math.NaN(), nil
}

func (t TorrentDownloadedInformationVoid) SetTorrentInformation(filename string, percentage float64) error {
	return nil
}

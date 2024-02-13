package interfaces

type TorrentDownloadedInformation interface {
	GetTorrentInformation() (string, float64, error)
	SetTorrentInformation(string, float64) error
}

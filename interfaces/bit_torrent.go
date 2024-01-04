package interfaces

type BitTorrentclient interface {
	TorrentAdd(magnet_link string) (int64, error)
	GetTorrentPercentComplete(id int64) (float64, error)
	GetTorrentName(id int64) (string, error)
	GetTorrentFiles(id int64) ([]string, error)
	Close() error
}

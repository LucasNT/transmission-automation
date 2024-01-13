package interfaces

type TorrentCompletedHandler interface {
	Exec(fileNames []string) error
}

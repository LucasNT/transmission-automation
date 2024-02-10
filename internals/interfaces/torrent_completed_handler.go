package interfaces

type TorrentCompletedHandler interface {
	CreateExec(config string) (func(fileName []string) (bool, error), error)
}

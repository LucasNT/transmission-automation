package interfaces

import "fmt"

var ErrNoTorrentEntry error = fmt.Errorf("No more entries")

type TorrentEntryReader interface {
	ReadTorrentEntry() (string, string, error)
}

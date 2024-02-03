package TorrentEntryReader

type ChannelTorrentEntryReader struct {
	ch chan []string
}

func (c ChannelTorrentEntryReader) ReadTorrentEntry() (string, string, error) {
	slice := <-c.ch
	return slice[0], slice[1], nil
}

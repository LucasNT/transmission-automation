package bitTorrentImplementation

import "fmt"

type torrent struct {
	name    string
	percent float64
}

type Mock struct {
	torrentListPercent map[int64]torrent
	id                 *int64
}

func NewBitTorrentMock() (Mock, error) {
	var id int64 = 0
	return Mock{
		torrentListPercent: make(map[int64]torrent),
		id:                 &id,
	}, nil
}

func (m Mock) TorrentAdd(magnet_link string) (int64, error) {
	var name string
	(*m.id)++
	name = fmt.Sprintf("Torrent %d", len(m.torrentListPercent))
	t := torrent{
		name:    name,
		percent: 0,
	}
	m.torrentListPercent[*m.id] = t
	return *m.id, nil
}

func (m Mock) GetTorrentName(id int64) (string, error) {
	item, ok := m.torrentListPercent[id]
	if ok {
		return item.name, nil
	} else {
		return "", fmt.Errorf("Invalid id")
	}
}

func (m Mock) GetTorrentPercentComplete(id int64) (float64, error) {
	item, ok := m.torrentListPercent[id]
	if ok {
		item.percent = item.percent + 0.2
		m.torrentListPercent[id] = item
		return item.percent, nil
	} else {
		return -1, fmt.Errorf("Invalid id")
	}
}

func (m Mock) Close() error {
	for k := range m.torrentListPercent {
		delete(m.torrentListPercent, k)
	}
	return nil
}

func (m Mock) GetTorrentFiles(id int64) ([]string, error) {
	v, ok := m.torrentListPercent[id]
	if ok {
		return []string{v.name}, nil
	} else {
		return []string{}, nil
	}
}

package bitTorrentImplementation

import (
	"fmt"
	"os"
	"path"
)

type torrent struct {
	name    string
	percent float64
}

type Mock struct {
	torrentListPercent map[int64]torrent
	id                 *int64
	path               string
}

func NewBitTorrentMock(path string) (Mock, error) {
	var id int64 = 0
	return Mock{
		torrentListPercent: make(map[int64]torrent),
		id:                 &id,
		path:               path,
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
	file, err := os.Create(path.Join(m.path, t.name))
	if err != nil {
		return 0, fmt.Errorf("Mock failed to create file: %w", err)
	}
	file.Close()
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
	for k, v := range m.torrentListPercent {
		os.Remove(path.Join(m.path, v.name))
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

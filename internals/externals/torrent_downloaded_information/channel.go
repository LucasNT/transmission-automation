package torrentdownloadedinformation

type TorrentDownloadedInformationChannel struct {
	ch chan informationStruct
}

type informationStruct struct {
	filname     string
	percertange float64
}

func NewTorrentDownloadedInformation(bufferSize int32) (TorrentDownloadedInformationChannel, error) {
	tr := TorrentDownloadedInformationChannel{}
	tr.ch = make(chan informationStruct, bufferSize)
	return tr, nil
}

func (t TorrentDownloadedInformationChannel) GetTorrentInformation() (string, float64, error) {
	ret := <-t.ch
	return ret.filname, ret.percertange, nil
}

func (t TorrentDownloadedInformationChannel) SetTorrentInformation(desc string, percentage float64) error {
	aux := informationStruct{filname: desc, percertange: percentage}
	t.ch <- aux
	return nil
}

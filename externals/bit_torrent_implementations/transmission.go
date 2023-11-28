package bitTorrentImplementation

import (
    "github.com/hekmon/transmissionrpc/v3"
    "net/url"
    "context"
)

type Transmission struct {
    tbt *transmissionrpc.Client;
}

func NewTransmision( url *url.URL, extra *transmissionrpc.Config ) (Transmission, error) {
    tbt, err := transmissionrpc.New(url, extra);
    if err != nil {
        return Transmission{}, err
    }
    return Transmission{
        tbt: tbt,
    }, nil
}

func (t Transmission) hasSuportedVersion() (ok bool, err error) {
    ok, _, _, err = t.tbt.RPCVersion(context.TODO())
    return
}

func (t Transmission) TorrentAdd (magnet_link string) (int64, error) {
    tr, err := t.tbt.TorrentAdd(context.TODO(), transmissionrpc.TorrentAddPayload{
        Filename: &magnet_link,
    })
    if  err != nil {
        return -1, err;
    }
    var id int64 = *tr.ID;
    return id,nil;
}

func (t Transmission) GetTorrentPercentComplete(id int64) (float64, error) {
    listTr, err := t.tbt.TorrentGet(context.TODO(), []string{"percentComplete"}, []int64{id});
    if err != nil {
        return -1, err;
    }
    return *(listTr[0].PercentComplete), nil;
}

func (t Transmission) GetTorrentName(id int64) (string , error) {
    listTr, err := t.tbt.TorrentGet(context.TODO(), []string{"name"}, []int64{id});
    if err != nil {
        return "", err;
    }
    return *(listTr[0].Name), nil;
}

func (t Transmission) GetTorrentFiles(id int64) ([]string , error) {
    listTr, err := t.tbt.TorrentGet(context.TODO(), []string{"files"}, []int64{id});
    if err != nil {
        return nil, err;
    }
    var filenameList []string;
    for _, val := range listTr[0].Files {
        filenameList = append(filenameList, val.Name);
    }
    return filenameList, nil;
}

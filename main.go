package main 

import (
    "net/url"
    "fmt"
    "encoding/csv"
    "os"
    "os/exec"
    "time"
    "github.com/LucasNT/transmission-list-automation/config"
    "github.com/LucasNT/transmission-list-automation/bit_torrent_implementations"
    "github.com/LucasNT/transmission-list-automation/interfaces"
)

const CONFIG_PATH string = "./config.yaml"

func main () {
    var err error;

    if len(os.Args) < 1 {
        fmt.Fprintf(os.Stderr, "Need at least one argument" )
        os.Exit(1);
    }
    
    if err = config.LoaderConfigs(CONFIG_PATH); err != nil {
        panic( err )
    }

    endpoint, err := url.Parse(config.Config.Url);
    if err != nil {
        panic(err);
    }
    endpoint.User = url.UserPassword(config.Config.Username, config.Config.Password);
    var bitTorrent interfaces.BitTorrentclient;
    bitTorrent,err = bitTorrentImplementation.NewTransmision(endpoint, nil );

    

    var csvFilePath string = os.Args[1]

    file,err := os.Open(csvFilePath);
    if err != nil {
        panic(err)
    }

    reader := csv.NewReader(file)

    l,err := reader.Read();
    fmt.Println(l, err)

    for l != nil && err == nil {
        tr_id, err := bitTorrent.TorrentAdd(l[0]);
        if err != nil {
            panic(err);
        }
        fmt.Println(tr_id)
        percent := float64(0)
        for percent != 1 {
            time.Sleep(1 * time.Minute);
            percent,err = bitTorrent.GetTorrentPercentComplete(tr_id)
            if err != nil {
                panic(err);
            }
            fileName, err := bitTorrent.GetTorrentName(tr_id);
            if err != nil {
                panic(err);
            }
            fmt.Println(fileName, percent)
        }
        listFileName, err := bitTorrent.GetTorrentFiles(tr_id);
        if err != nil {
            panic(err);
        }
        cmd := exec.Command("cp", "--reflink=auto", "/data/torrents/" + listFileName[0], l[1])

        err = cmd.Run()

        if err != nil {
            panic(err)
        }

        l,err = reader.Read();
    }
    if err != nil {
        panic(err)
    }


}

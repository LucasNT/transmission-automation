package main 

import (
    "net/url"
    "context"
    "fmt"
    "encoding/csv"
    "os"
    "os/exec"
    "time"
    "github.com/hekmon/transmissionrpc/v3"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Password string;
    Url string;
    Username string
}

func main () {
    var config Config
    file,err := os.Open("./config.yaml");
    if err != nil {
        panic(err)
    }
    decode := yaml.NewDecoder(file);
    decode.Decode(&config);
    fmt.Print(config)
    endpoint, err := url.Parse(config.Url)
    if err != nil {
        panic(err)
    }
    endpoint.User = url.UserPassword(config.Username, config.Password);
    tbt, err := transmissionrpc.New(endpoint, nil)
    if err != nil {
        panic(err)
    }
    fmt.Println("a")
    ok, serverVersion, serverMinimumVersion, err := tbt.RPCVersion(context.TODO())
    if err != nil {
        panic(err)
    }
    if !ok {
        panic(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
        serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
    }
    fmt.Printf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
    serverVersion, transmissionrpc.RPCVersion)

    var csvFilePath string = os.Args[1]

    file,err = os.Open(csvFilePath);
    if err != nil {
        panic(err)
    }

    reader := csv.NewReader(file)

    l,err := reader.Read();
    fmt.Println(l, err)

    for l != nil && err == nil {
        tr, err := tbt.TorrentAdd(context.TODO(), transmissionrpc.TorrentAddPayload{
            Filename: &l[0],
        })
        if err != nil {
            panic(err);
        }
        time.Sleep(1 * time.Minute)
        id  := *tr.ID
        fmt.Println(id)
        listTr,err := tbt.TorrentGet(context.TODO(), []string{"files","percentComplete"}, []int64{id})
        for *(listTr[0].PercentComplete) != 1 {
            fmt.Println(*(listTr[0].PercentComplete))
            time.Sleep(1 * time.Minute)
            listTr,err = tbt.TorrentGet(context.TODO(), []string{"files","percentComplete"}, []int64{id})
            fmt.Println(listTr[0].Files[0].Name)
        }
        cmd := exec.Command("cp", "--reflink=auto", "/data/torrents/" + listTr[0].Files[0].Name, l[1])

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

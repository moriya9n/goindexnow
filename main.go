package main

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"

    "github.com/oxffaa/gopher-parse-sitemap"
)

type IndexNowRequest struct {
    Host string `json:"host"`
    Key string `json:"key"`
    KeyLocation string `json:"keyLocation"`
    UrlList []string `json:"urlList"`
}

func locationFromSitemap(url string) ([]string, error) {
    var locations []string
    err := sitemap.ParseFromSite(url, func (e sitemap.Entry) error {
        fmt.Println(e.GetLocation())
        locations = append(locations, e.GetLocation())
        return nil
    })
    if err != nil {
        return nil, err
    }
    return locations, nil
}

func buildIndexNowRequest(key string, keyLocation string, host string, locations []string) ([]byte, error) {
    req := IndexNowRequest{host, key, keyLocation, locations}
    body, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    return body, nil
}

func requestIndexNow(searchEngineUrl string, request []byte) error {
    resp, err := http.Post(searchEngineUrl, "application/json; charset=utf-8", bytes.NewReader(request))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode > 200 {
        return errors.New(fmt.Sprintf("error: %d", resp.StatusCode))
    }

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(string(respBody))
    return nil
}

func main() {
    args := os.Args[1:]
    if len(args) != 5 {
        fmt.Printf("usage: %s search_engine_url host key keyLocation sitemap_url\n", os.Args[0])
        os.Exit(1)
    }
    searchEngineUrl := args[0]
    host := args[1]
    indexNowKey := args[2]
    indexNowKeyLocation := args[3]
    sitemapUrl := args[4]

    locations, err := locationFromSitemap(sitemapUrl)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    if len(locations) == 0 {
        fmt.Println("no locations in the sitemap")
        os.Exit(2)
    }

    req, err := buildIndexNowRequest(indexNowKey, indexNowKeyLocation, host, locations)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    err = requestIndexNow(searchEngineUrl, req)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("IndexNow Success")
}

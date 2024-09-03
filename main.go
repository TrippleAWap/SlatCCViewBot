package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	proxyListPath := "./proxyList.txt"

	data, err := os.ReadFile(proxyListPath)
	if err != nil {
		log.Fatalln(err)
	}
	proxyList := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	userIds := []int{59717}
	batching := 100
	var wg sync.WaitGroup
	for i := 0; i < int(math.Ceil(float64(len(proxyList)/batching))); i++ {
		wg.Add(1)
		go func(proxyList []string) {
			defer wg.Done()
			for _, userId := range userIds {
				for _, proxy := range proxyList {
					go func(userId int, proxy string) {
						_ = ViewProfile(userId, proxy)
					}(userId, proxy)
				}
			}
		}(proxyList[i*batching:])
	}
	wg.Wait()
	fmt.Println("DONE")
}

func ViewProfile(userId int, proxyUrl string) error {
	proxyURL, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	resp, err := client.Post("https://slat.cc/api/users/"+strconv.Itoa(userId)+"/views", "application/json", nil)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	println(string(body))
	return nil
}

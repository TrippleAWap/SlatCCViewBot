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
)

func main() {
	proxyListPath := "./proxyList.txt"

	data, err := os.ReadFile(proxyListPath)
	if err != nil {
		log.Fatalln(err)
	}
	proxyList := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	var userId int
	print("Enter target user id: ")
	_, err = fmt.Scan(&userId)
	if err != nil {
		log.Fatalln(err)
	}
	userIds := []int{userId}
	batching := 100

	fmt.Printf("Loaded %v proxies\nview botting %v account(s)\n", len(proxyList), len(userIds))
	for i := 0; i < int(math.Ceil(float64(len(proxyList)/batching))); i++ {
		go func(proxyList []string) {
			for _, userIdV := range userIds {
				for _, proxy := range proxyList {
					go func(proxy string) {
						_ = ViewProfile(userIdV, proxy)
					}(proxy)
				}
			}
		}(proxyList[i*batching:])
	}
	for {

	}
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

	if strings.Contains(string(body), "Successfully added view") {
		println(string(body))
	}
	return nil
}

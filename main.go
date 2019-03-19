package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	RequestTimeout int = 1
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

func main() {
	var (
		err      error
		template = `{ "jsonrpc": "2.0", "method": "sendrawtransaction", "params": ["%s"], "id": %d }`
	)

	connNum := flag.Int("c", 1, "number of simultaneous connections")
	txNum := flag.Int("t", 100, "total amount of transactions")
	filename := flag.String("f", "./raw.txs", "file with raw transactions")
	endpoint := flag.String("url", "http://127.0.0.1:30334/", "connection url")

	flag.Parse()

	rawTxs, err := readTxFromFile(*filename)
	if err != nil {
		log.Fatal(err)
	}
	if len(rawTxs) < *txNum {
		log.Fatal("there is not enough raw transactions in file")
	}

	var txPerConn int
	txPerConn = *txNum / *connNum

	// prepare all requests beforehand
	req := make([]*http.Request, 0, len(rawTxs))
	for i, body := range rawTxs {
		r, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer([]byte(fmt.Sprintf(template, body, i+1))))
		r.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatal(err)
		}
		req = append(req, r)
	}

	wg := sync.WaitGroup{}
	wg.Add(*connNum)

	log.Printf("Sending %d transactions (%d tx in %d connections) to %s\n", *txNum, txPerConn, *connNum, *endpoint)
	totalStart := time.Now()

	for i := 0; i < *connNum; i++ {
		go func(ind int) {
			defer wg.Done()

			var response *http.Response
			httpClient := createHTTPClient()
			from := ind * txPerConn
			to := from + txPerConn

			start := time.Now()
			for j := from; j < to; j++ {
				response, _ = httpClient.Do(req[j])
				if err != nil && response == nil {
					log.Fatalf("Error sending request to API endpoint. %+v", err)
				}
				ioutil.ReadAll(response.Body)
			}
			log.Println("connection ", i, ") time: ", time.Since(start))

		}(i)

	}

	wg.Wait()
	totalTime := time.Since(totalStart)
	log.Println("---")
	sec := float64(totalTime) / float64(time.Second)
	log.Println("Total time: ", totalTime, "Approximate TPS: ", float64(*txNum)/sec)
	log.Println("For accurate TPS values use external connection sniffers, e.g. Wireshark")
}

func readTxFromFile(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

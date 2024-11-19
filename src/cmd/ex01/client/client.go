package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"go_day04/pkg/types"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	CertPath            = "../../../cert/client-cert/cert.pem"
	KeyPath             = "../../../cert/client-cert/key.pem"
	RootCertificatePath = "../../../cert/minica.pem"
)

func main() {
	candyType := flag.String("k", "", "accepts two-letter abbreviation for the candy type")
	count := flag.String("c", "", "count of candy to buy")
	money := flag.String("m", "", "money of candy to buy")
	flag.Parse()
	if *candyType == "" || *count == "" || *money == "" {
		flag.Usage()
		return
	}
	c, err := strconv.Atoi(*count)
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(*money)
	if err != nil {
		log.Fatal(err)
	}

	rootCA, err := os.ReadFile(RootCertificatePath)
	if err != nil {
		log.Fatalf("reading cert failed: %v", err)
	}
	rootCAPool := x509.NewCertPool()
	rootCAPool.AppendCertsFromPEM(rootCA)
	cert, err := tls.LoadX509KeyPair(CertPath, KeyPath)
	if err != nil {
		log.Fatalf("loading client pair key failed: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      rootCAPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	data := &types.Order{
		CandyType:  *candyType,
		Money:      m,
		CandyCount: c,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	req, err := client.Post("https://candy.tld:3333/buy_candy", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	res, _ := io.ReadAll(req.Body)
	if req.StatusCode != 201 {
		fmt.Print(string(res))
	} else {
		resp := struct {
			Change int    `json:"change"`
			Thanks string `json:"thanks"`
		}{}
		err = json.Unmarshal(res, &resp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s Your change is %d\n", resp.Thanks, resp.Change)
	}
	//resp, err := io.ReadAll(req.Body)
	//fmt.Println(resp)
	//if
	//query := fmt.Sprintf(`'{"money": %d, "candyType": %s, "candyCount": %d}'`, m, *candyType, c)
	//request, err := http.NewRequest(http.MethodPost, "https://candy.tld:3333/buy_candy", strings.NewReader(query))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//request.Header.Set("Content-Type", "application/json")
	////request.Header.Set("candyType", *candyType)
	////request.Header.Set("money", *money)
	////request.Header.Set("candyCount", *count)
	////request.
	////
	//response, err := client.Do(request)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer response.Body.Close()
	//if response.StatusCode != http.StatusOK {
	//	fmt.Errorf("server returned non-200 status: %s", response.Status)
	//}

}

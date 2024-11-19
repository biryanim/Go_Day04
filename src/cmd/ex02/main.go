package main

/*
#include "cow.c"
*/
import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"go_day04/pkg/types"
	"log"
	"net/http"
	"os"
	"unsafe"
)

const (
	CertPath            = "../../cert/candy.tld/cert.pem"
	KeyPath             = "../../cert/candy.tld/key.pem"
	RootCertificatePath = "../../cert/minica.pem"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buy_candy", BuyCandyHandler)
	clientCa, err := os.ReadFile(RootCertificatePath)
	if err != nil {
		log.Fatalf("Reading cert failed: %v", err)
	}
	cert, err := tls.LoadX509KeyPair(CertPath, KeyPath)
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCa)
	server := &http.Server{
		Handler: mux,
		Addr:    ":3333",
		TLSConfig: &tls.Config{
			ClientCAs:    clientCAPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
		},
	}
	log.Println("Starting server on port 3333")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func BuyCandyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order types.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	candyCost, ok := types.Candies[order.CandyType]
	if !ok || order.CandyCount < 0 || order.Money < 0 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if candyCost*order.CandyCount > order.Money {
		http.Error(w, fmt.Sprintf("You need %d more money!", candyCost*order.CandyCount-order.Money), http.StatusPaymentRequired)
		return
	}

	cstr := C.CString("Thank you!")
	defer C.free(unsafe.Pointer(cstr))

	res := C.ask_cow(cstr)
	gostr := C.GoString(res)
	response := map[string]interface{}{
		"change": order.Money - candyCost*order.CandyCount,
		"thanks": gostr,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

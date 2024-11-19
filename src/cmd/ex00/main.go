package main

import (
	"encoding/json"
	"fmt"
	"go_day04/pkg/types"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/buy_candy", BuyCandyHandler)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Listening on port 3333")
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
	response := map[string]interface{}{
		"change": order.Money - candyCost*order.CandyCount,
		"thanks": "Thank you!",
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

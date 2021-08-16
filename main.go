package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"
)

type Payload struct {
	No       string `json:"no"`
	To       string `json:"to"`
	ShrPulsa int    `json:"share_pulsa"`
}

func main() {
	http.HandleFunc("/test", RouterSendPulsa)
	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func RouterSendPulsa(w http.ResponseWriter, r *http.Request) {
	var (
		Pulsa = 10000
	)

	pay, code, err := BodyRequest(r)
	if err != nil || code >= http.StatusBadRequest {
		log.Info("Error while body request")
		Response(w, code, "Gagal Mengirim Pulsa")
		return
	}

	code, msg, err := SharePulsa(pay, Pulsa)
	if err != nil || code >= http.StatusBadRequest {
		log.Infof("%v", msg)
		Response(w, code, msg)
		return
	}

	Response(w, code, msg)
	return

}

func BodyRequest(r *http.Request) (*Payload, int, error) {
	var (
		pay Payload
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Info("Error while readAll")
		return nil, http.StatusBadRequest, err
	}

	err = json.Unmarshal(body, &pay)
	if err != nil {
		log.Info("Error while unmarshal")
		return nil, http.StatusBadRequest, err
	}

	return &pay, http.StatusOK, nil
}

func SharePulsa(pay *Payload, Pulsa int) (int, string, error) {
	if pay.No == pay.To {
		log.Info("Tidak Bisa dikirim ke nomer sendiri")
		return http.StatusBadRequest, "Gagal Mengirim Pulsa", errors.New("Failed")
	}

	if pay.ShrPulsa >= Pulsa {
		log.Info("Gagal Kirim Pulsa")
		return http.StatusBadRequest, "Pulsa Anda Tidak Mencukupi", errors.New("Failed")
	}

	return http.StatusOK, "Berhasil Mengirim Pulsa", nil

}

func Response(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

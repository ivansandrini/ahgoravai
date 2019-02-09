package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ahgoraURL = "https://www.ahgora.com.br/batidaonline/verifyIdentification"
)

type Ponto struct {
	Identity     string  `json:"identity"`
	Account      string  `json:"account"`
	Password     string  `json:"password"`
	Logon        bool    `json:"logon"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Accuracy     int     `json:"accuracy"`
	TimestampLoc int64   `json:"timestamp_loc"`
	Provider     string  `json:"provider"`
	Offline      bool    `json:"offline"`
	Timestamp    int64   `json:"timestamp"`
	Origin       string  `json:"origin"`
}

type Handler struct {
	url   string
	mPath string
}

func NewHandler(url string, mPath string) *Handler {
	return &Handler{
		url,
		mPath,
	}
}

func (h Handler) Hponto(w http.ResponseWriter, r *http.Request) ([]string, error) {
	ms := lMat(h.mPath)

	var res []string
	for _, m := range ms {
		p, _ := bPonto(h.url, m)
		res = append(res, p)
	}

	return res, nil
}

func lMat(mPath string) []string {
	file, err := os.Open(mPath)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)

	var m []string
	for s.Scan() {
		m = append(m, s.Text())
	}

	if err := s.Err(); err != nil {
		log.Print(err)
	}

	return m
}

func bPonto(url string, m string) (string, error) {
	d := time.Now().Unix()
	p := Ponto{
		Identity:     os.Getenv("AHGORA_IDENTITY"),
		Account:      m,
		Password:     m,
		Logon:        false,
		Longitude:    -48.879195599999996,
		Latitude:     -26.241964400000004,
		Accuracy:     28,
		TimestampLoc: d,
		Provider:     "network/wifi",
		Offline:      false,
		Timestamp:    d,
		Origin:       "chr",
	}
	pJSON, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERRO ao bater ponto. matrícula %v", m), err
	}

	pResp, err := http.Post(url+"/batidaonline/verifyIdentification", "application/json", bytes.NewBuffer(pJSON))
	if err != nil {
		return fmt.Sprintf("ERRO ao bater ponto. matrícula %v", m), err
	}

	defer pResp.Body.Close()
	if pResp.StatusCode == http.StatusOK {
		return fmt.Sprintf("Ponto batido com SUCESSO. matrícula: %v", m), nil
	}

	return fmt.Sprintf("ERRO %v ao bater ponto. matrícula %v", pResp.StatusCode, m), nil
}

func main() {
	lambda.Start(NewHandler(ahgoraURL, "matriculas"))
}

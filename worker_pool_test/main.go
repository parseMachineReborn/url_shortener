package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	workerpool "github.com/parseMachineReborn/worker_pool/worker_pool"
)

const URL = "http://localhost:8080"
const CONTENT_TYPE = "application/json"
const LETTERS = "abcdefghijklmnopqrstuvwxyz"

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	input := struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}{
		Email: "pookich@mail.com",
		Pass:  "123456",
	}
	jsonedInput, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
		return
	}

	if err := SignUp(client, bytes.NewBuffer(jsonedInput)); err != nil {
		log.Printf("SIGN UP ERROR: %s", err)
		return
	}
	jwtToken, err := Login(client, bytes.NewBuffer(jsonedInput))
	if err != nil {
		log.Printf("LOGIN ERROR: %s", err)
		return
	}

	wp := workerpool.NewWorkerPool[string](5)
	resCh, err := wp.Start(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for res := range resCh {
			log.Printf("Результат сокращенной ссылки: %s", res)
		}
	}()

	for i := 0; i < 100000; i++ {
		url := struct {
			Address string `json:"address"`
		}{
			Address: randomURL(),
		}
		marshaledUrl, err := json.Marshal(url)
		if err != nil {
			log.Println(err)
			return
		}
		wp.AddTask(func() string { return Shorten(client, bytes.NewBuffer(marshaledUrl), jwtToken) })
	}

	err = wp.Stop()
	if err != nil {
		fmt.Printf("ПРОБЛЕМА ПРИ ЗАКРЫТИИ ПУЛА: %s", err)
	}

	<-ctx.Done()
}

func SignUp(c *http.Client, buf io.Reader) error {
	resp, err := c.Post(URL+"/signup", CONTENT_TYPE, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func Login(c *http.Client, buf io.Reader) (string, error) {
	resp, err := c.Post(URL+"/login", CONTENT_TYPE, buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respRes map[string]string
	fmt.Println(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&respRes)
	if err != nil {
		return "", err
	}

	return respRes["token"], nil
}

func Shorten(c *http.Client, buf io.Reader, jwtToken string) string {
	req, err := http.NewRequest("POST", URL+"/shorten", buf)
	if err != nil {
		log.Println(err)
		return ""
	}
	req.Header.Set("Content-Type", CONTENT_TYPE)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var shortenUrl string
	err = json.NewDecoder(resp.Body).Decode(&shortenUrl)
	if err != nil {
		return ""
	}

	return shortenUrl
}

func randomURL() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = LETTERS[rand.IntN(len(LETTERS))]
	}
	return "https://test.com/" + string(b)
}

package main

import (
	"bufio"
	"log"
	"os"
	"sync"

	"github.com/valyala/fasthttp"
)

func checkToken(token string) {
	tokenToCheck := token

	req := fasthttp.AcquireRequest()
	req.Header.Set("Authorization", tokenToCheck)
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://canary.discordapp.com/api/v8/users/@me")
	res := fasthttp.AcquireResponse()

	fasthttp.Do(req, res)
	fasthttp.ReleaseRequest(req)
	statusCode := int(res.StatusCode())
	fasthttp.ReleaseResponse(res)

	if statusCode == 200 {
		file, err := os.OpenFile("./output.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		_, err = file.WriteString(tokenToCheck + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		wg.Add(1)
		go checkToken(scanner.Text())
	}
	wg.Wait()
}

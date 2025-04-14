package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	sketchyURL := "http://localhost:3000/galleries/2/images/../images/Resting.jpg"
	resp, err := http.Get(sketchyURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

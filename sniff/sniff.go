package main

import (
	"fmt"
	"github.com/marketvibe/mimesniff"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://probablyprogramming.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	res := mimesniff.NewHttpResponse(resp)
	meta := mimesniff.NewResourceMetadata(res)

	fmt.Println(mimesniff.DetermineEncoding(meta.ParsedType, res.Header()))
	body, err := ioutil.ReadAll(meta.DecodeBody())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

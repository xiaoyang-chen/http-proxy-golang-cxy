package forward_proxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	local8081 = "127.0.0.1:8081"
)

// ForwardPxy 是正向代理，通过发送给代理服务器，
// 由代理服务器代为请求，并将响应回传
func ForwardPxy() {
	http.HandleFunc("/", handleHttp2)
	http.ListenAndServe(local8081, nil)
}

func handleHttp(w http.ResponseWriter, req *http.Request) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	for k, v := range res.Header {
		for _, v2 := range v {
			w.Header().Add(k, v2)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

var client = &http.Client{}

func handleHttp2(w http.ResponseWriter, req *http.Request) {
	errDeal := func(err error, statusCode int) {
		log.Println(err)
		w.WriteHeader(statusCode)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errDeal(err, http.StatusBadRequest)
		return
	}
	newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
	if err != nil {
		errDeal(err, http.StatusBadRequest)
		return
	}
	for k, v := range req.Header {
		for _, v2 := range v {
			newReq.Header.Add(k, v2)
		}
	}

	res, err := client.Do(newReq)
	if err != nil {
		errDeal(err, http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	for k, v := range res.Header {
		for _, v2 := range v {
			w.Header().Add(k, v2)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

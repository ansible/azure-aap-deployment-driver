package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HttpRequester struct {
	client *http.Client
}

type HttpRequest struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    interface{}
}

type HttpResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

func NewHttpRequesterWithCertificate(certPEMString, privkeyPEMString string) (*HttpRequester, error) {
	cert, err := tls.X509KeyPair([]byte(certPEMString), []byte(privkeyPEMString))
	if err != nil {
		log.Printf("Couldn't parse certificate and/or key. %v\n", err)
		return nil, err
	}
	// setup tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	// setup transport
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	// client
	return &HttpRequester{
		client: &http.Client{
			Transport: httpTransport,
			Timeout:   30 * time.Second,
		},
	}, nil
}

func (requester *HttpRequester) MakeJSONRequest(ctx context.Context, request HttpRequest) (*HttpResponse, error) {
	// create body
	var bodyBuffer bytes.Buffer
	if err := json.NewEncoder(&bodyBuffer).Encode(request.Body); err != nil {
		log.Printf("Couldn't encode body. %v\n", err)
		return nil, err
	}
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, request.Url, &bodyBuffer)
	if err != nil {
		log.Printf("Couldn't prepare HTTP request. %v\n", err)
		return nil, err
	}
	// add or update content type header
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["Content-Type"] = "application/json"
	// add all header to request
	for h, v := range request.Headers {
		httpRequest.Header.Add(h, v)
	}
	// this following block can be wrapped in go routine if this was not supposed to be blocking
	httpResponse, err := requester.client.Do(httpRequest)
	if err != nil {
		log.Printf("Couldn't send HTTP request. %v\n", err)
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	return &HttpResponse{
		StatusCode: httpResponse.StatusCode,
		Headers:    map[string][]string(httpResponse.Header),
		Body:       bodyBytes,
	}, nil
}

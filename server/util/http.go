package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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
	Body    *bytes.Buffer
}

type HttpResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

func NewHttpRequester() *HttpRequester {
	// client
	return newRequester(nil)
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
	return newRequester(httpTransport), nil
}

func newRequester(transport *http.Transport) *HttpRequester {
	return &HttpRequester{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func EncodeAsJSON(body interface{}) (*bytes.Buffer, error) {
	var bodyBuffer bytes.Buffer
	if err := json.NewEncoder(&bodyBuffer).Encode(body); err != nil {
		log.Printf("Couldn't encode body. %v\n", err)
		return nil, err
	}
	return &bodyBuffer, nil
}

func EncodeAsWWWFormURLEncoding(body map[string]string) (*bytes.Buffer, error) {
	values := url.Values{}
	for n, v := range body {
		values.Add(n, v)
	}
	return bytes.NewBufferString(values.Encode()), nil
}

func (requester *HttpRequester) MakeRequestWithJSONBody(ctx context.Context, method string, url string, headers map[string]string, body interface{}) (*HttpResponse, error) {
	bodyEncoded, err := EncodeAsJSON(body)
	if err != nil {
		log.Fatalf("Could not encode body as JSON. %v", err)
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return requester.MakeRequest(ctx, HttpRequest{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    bodyEncoded,
	})
}

func (requester *HttpRequester) MakeRequestWithWWWFormUrlEncodedBody(ctx context.Context, method string, url string, headers map[string]string, body map[string]string) (*HttpResponse, error) {
	bodyEncoded, err := EncodeAsWWWFormURLEncoding(body)
	if err != nil {
		log.Fatalf("Could not encode body as WWW Form Url Encoded. %v", err)
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return requester.MakeRequest(ctx, HttpRequest{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    bodyEncoded,
	})
}

func (requester *HttpRequester) MakeRequest(ctx context.Context, request HttpRequest) (*HttpResponse, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, request.Url, request.Body)
	if err != nil {
		log.Printf("Couldn't prepare HTTP request. %v\n", err)
		return nil, err
	}
	// if there are headers add them all to request
	if request.Headers != nil {
		for h, v := range request.Headers {
			httpRequest.Header.Add(h, v)
		}
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

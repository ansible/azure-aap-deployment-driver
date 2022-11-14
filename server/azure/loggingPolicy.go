package azure

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	log "github.com/sirupsen/logrus"
)

type LoggingPolicy struct {
	LogPrefix string
}

func (m *LoggingPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Mutate/process request.
	start := time.Now()
	// Forward the request to the next policy in the pipeline.
	res, err := req.Next()
	if res == nil {
		// End of policies
		return res, err
	}
	// Mutate/process response.
	// Return the response & error back to the previous policy in the pipeline.
	bodyString, _ := io.ReadAll(res.Body)
	record := struct {
		URL        string
		Duration   time.Duration
		Response   string
		StatusCode int
		Headers    map[string][]string
	}{
		URL:        req.Raw().URL.RequestURI(),
		Duration:   time.Duration(time.Since(start).Milliseconds()),
		Response:   string(bodyString),
		StatusCode: res.StatusCode,
		Headers:    res.Header,
	}
	b, _ := json.Marshal(record)
	log.Printf("%s %s\n", m.LogPrefix, b)
	return res, err
}

// This provides a policy added to a ClientOptions instance
// that can be used with any Azure client to make it log
// the http requests it sends out plus responses with
// status codes and headers.
func GetClientOptionsWithLogging() *arm.ClientOptions {
	opts := arm.ClientOptions{}
	lp := &LoggingPolicy{}
	lp.LogPrefix = "URL:"
	opts.PerCallPolicies = append(opts.PerCallPolicies, lp)
	return &opts
}

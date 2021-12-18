package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Api struct {
	Name        string
	Limits      Limits
	Credentials map[string]Credential
	Endpoints   map[string]Endpoint
}

type Limits struct {
	Calls    int32
	Frame    int32
	NextCall time.Time
}

type Endpoint struct {
	Path       string
	Method     string
	Block      bool
	Parameters map[string]struct {
		Type  string
		Value string
	}
}

type Credential struct {
	Secret string
	Public string
}

type Response struct {
	Method     		  string
	Endpoint   		  string
	Status     		  int
	Parameters 		  map[string]string
	Value      		  interface{}
	RawValue   		  []byte
	DNSStart   		  time.Time
	DNSDone    		  time.Duration
	TLSHandshakeStart time.Time
	TLSHandshakeDone  time.Duration
	ConnectStart	  time.Time
	ConnectDone		  time.Duration
	GotFirstResByte	  time.Time
	ExactStartTime	  time.Time
	TotalTime		  time.Duration
}

var configs = make(map[string]Api)
var loadOnce sync.Once

func loadAll() {
	loadOnce.Do(func() {
		files, err := ioutil.ReadDir("./example_config")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".yml" {
				conf := Api{}
				content, _ := ioutil.ReadFile(fmt.Sprintf("./example_config/%s", file.Name()))
				if err := yaml.Unmarshal(content, &conf); err != nil {
					panic(err)
				}
				apiName := strings.Split(file.Name(), ".")[0]
				configs[apiName] = conf
			}
		}
	})
}

func GetApiConfig(name string) (*Api, error) {
	loadAll()
	config, ok := configs[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Requested api (%s) does not exist", name))
	}
	return &config, nil
}

// CallWGet Calls the api endpoint with the GET http method,
// the function accepts a credentials name, and an arguments map
func (i *Api) CallWGet(endpoint string, credentials string, args map[string]string) (Response, error) {
	creds := i.Credentials[credentials]
	endp := i.Endpoints[endpoint]
	result := Response {
		Method:     "GET",
		Endpoint:   endp.Path,
		Parameters: args,
	}
	if endp.Method == "GET" {
		endp.prepGetCredentials(creds)
		endp.addGetParams(args)
		result.Endpoint = endp.Path
		request := endp.createRequestWithContext(&result)
		endp.makeRoundTrip(request, &result)
		return result, nil
	}
	return Response{}, errors.New("could not resolve API endpoint method")
}

func (i *Api) BulkCallWGet(endpoint string, credentials string, argsArr []map[string]string) <-chan Response {
	results := make(chan Response, len(argsArr))
	// concurrently go through the bulk
	go func() {
		// for each argument map
		for _, arg := range argsArr {
			// check if the NextCall is not zero'd out and that now is before the next call
			if !i.Limits.NextCall.IsZero() && time.Now().Before(i.Limits.NextCall) {
				// generate a new timer
				nanoDur := i.Limits.NextCall.Sub(time.Now())
				timer := time.NewTimer(nanoDur)
				// Wait for the timer channel to send a message to the goroutine
				<-timer.C
			}
			// Register the current call so the limiter knows when to make the next one
			i.Limits.registerCall()
			go func(query map[string]string) {
				res, _ := i.CallWGet(endpoint, credentials, query)
				results <- res
			}(arg)
		}
	}()
	return results
}

func (e *Endpoint) createRequestWithContext(resp *Response) *http.Request {
	req, _ := http.NewRequest(e.Method, e.Path, nil)
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			resp.DNSStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			resp.DNSDone = time.Since(resp.DNSStart)
		},
		TLSHandshakeStart: func() {
			resp.TLSHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			resp.TLSHandshakeDone = time.Since(resp.TLSHandshakeStart)
		},
		ConnectStart: func(network, addr string) {
			resp.ConnectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			resp.ConnectDone = time.Since(resp.ConnectStart)
		},
		GotFirstResponseByte: func() {
			resp.GotFirstResByte = time.Now()
		},
	}
	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}

// makeRoundTrip Makes a call to the endpoint given a request and response
// objects, it uses a prepared request object to make the call and a response
// object to store the information about the request
func (e *Endpoint) makeRoundTrip(request *http.Request, response *Response) {
	response.ExactStartTime = time.Now()
	if netResp, err := http.DefaultTransport.RoundTrip(request); err != nil {
		log.Fatal(err)
	} else {
		defer netResp.Body.Close()
		body, _ := ioutil.ReadAll(netResp.Body)
		var jsonResp interface{}
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			log.Fatal(err)
		}
		response.Status = netResp.StatusCode
		response.Value = jsonResp
		response.RawValue = body
	}
	response.TotalTime = time.Since(response.ExactStartTime)
}

func (e *Endpoint) prepGetCredentials(cred Credential) {
	args := make(map[string]string)
	for name, param := range e.Parameters {
		switch param.Type {
			case "credentials":
				if param.Value == "secret" {
					args[name] = cred.Secret
				} else if param.Value == "public" {
					args[name] = cred.Public
				}
		}
	}
	// Add key to request
	e.addGetParams(args)
}

func (e *Endpoint) addGetParams(args map[string]string) {
	if len(args) > 0 {
		currentQuery, _ := url.Parse(e.Path)
		if len(currentQuery.Query()) > 0 {
			e.Path += "&"
		} else {
			e.Path += "?"
		}
		query := url.Values{}
		for key, value := range args {
			query.Add(key, value)
		}
		e.Path += query.Encode()
	}
}

func (l *Limits) registerCall() {
	var startingTime time.Time
	if l.NextCall.After(time.Now()) {
		startingTime = l.NextCall
	} else {
		startingTime = time.Now()
	}
	nextMilliseconds := int64(math.Round(l.MillisecondsBetweenCalls()))
	l.NextCall = startingTime.Add(time.Millisecond * time.Duration(nextMilliseconds)) // add a duration to starting time
}

func (l *Limits) MillisecondsBetweenCalls() float64 {
	var numerator, denominator int32
	numerator = l.Frame
	denominator = l.Calls
	milliseconds := (float64(numerator) / float64(denominator)) * 1000 // Turn rate to milliseconds
	return milliseconds
}

func GetArgumentsFromSlice(params []string) []map[string]string {
	var arguments = make([]map[string]string, len(params))
	if params != nil {
		for i, val := range params {
			// Key, Value
			// kv := strings.SplitN(val, ":", 2)
			var paramMap map[string]string
			if err := json.Unmarshal([]byte(val), &paramMap); err != nil {
				panic(err)
			}
			arguments[i] = paramMap
			// arguments[kv[0]] = kv[1]
		}
	}
	return arguments
}

// Print prints outs the response struct with the status code and the endpoint
// you can also pass a boolean value to print the value returned from the api or not
// make sure to mark verbose as false in production environments
func (r *Response) Print(verbose bool) string {
	var outputColor color.Attribute
	switch {
	case r.Status >= 200 && r.Status < 300:
		outputColor = color.FgGreen
		break
	case r.Status >= 300 && r.Status < 400:
		outputColor = color.FgYellow
		break
	case r.Status >= 400:
		outputColor = color.FgRed
		break
	}
	statusColorSprintf := color.New(outputColor).SprintfFunc()
	durationColorSprintf := color.New(color.BgBlack, color.FgWhite).SprintfFunc()
	var printArr []interface{}
	printFormat := "%s %s %s"
	printArr = append(
		printArr,
		statusColorSprintf("[%s:%d]", r.Method, r.Status),
		durationColorSprintf("%s", r.TotalTime),
		r.Endpoint,
	)
	if verbose {
		printFormat += "\n%s"
		printArr = append(printArr, r.RawValue)
	}
	return fmt.Sprintf(printFormat, printArr...)
}

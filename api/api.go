package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"net/http"
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
	Method     string
	Endpoint   string
	Status     int
	Parameters map[string]string
	Value      interface{}
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

// Call the api endpoint with the GET http method,
// the function accepts a credentials name, and an arguments map
func (i *Api) CallWGet(endpoint string, credentials string, args map[string]string) (Response, error) {
	creds := i.Credentials[credentials]
	endp := i.Endpoints[endpoint]
	if endp.Method == "GET" {
		endp.prepGetCredentials(creds)
		endp.addGetParams(args)
		resp, _ := http.Get(endp.Path)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var jsonResp interface{}
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			return Response{}, err
		}
		return Response{
			Method:     "GET",
			Endpoint:   endp.Path,
			Status:     resp.StatusCode,
			Parameters: args,
			Value:      jsonResp,
		}, nil
	}
	return Response{}, errors.New("could not resolve API endpoint method")
}

func (i *Api) BulkCallWGet(endpoint string, credentials string, argsArr []map[string]string) <-chan Response {
	results := make(chan Response, len(argsArr))
	// concurrently go though the bulk
	go func() {
		// for each argument map
		for _, arg := range argsArr {
			// check if the NextCall is not zero'd out and that now is before the next call
			if !i.Limits.NextCall.IsZero() && time.Now().Before(i.Limits.NextCall) {
				// generate a new timer
				nanoDur := i.Limits.NextCall.Sub(time.Now())
				timer := time.NewTimer(nanoDur)
				<-timer.C
			}
			// register the call
			i.Limits.registerCall()
			go func() {
				res, _ := i.CallWGet(endpoint, credentials, arg)
				results <- res
			}()
		}
	}()
	return results
}

func (e *Endpoint) prepGetCredentials(cred Credential) {
	args := make(map[string]string)
	for name, param := range e.Parameters {
		if param.Type != "input" {
			switch param.Type {
			case "credentials":
				if param.Value == "secret" {
					args[name] = cred.Secret
				} else if param.Value == "public" {
					args[name] = cred.Public
				}
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

func GetArgumentsFromSlice(params []string) map[string]string {
	var arguments = make(map[string]string)
	if params != nil {
		for _, val := range params {
			// Key, Value
			kv := strings.SplitN(val, ":", 2)
			arguments[kv[0]] = kv[1]
		}
	}
	return arguments
}

func (r *Response) Print() string {
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
	colorSprintf := color.New(outputColor).SprintfFunc()
	return fmt.Sprintf("%s %s", colorSprintf("[%s:%d]", r.Method, r.Status), r.Endpoint)
}

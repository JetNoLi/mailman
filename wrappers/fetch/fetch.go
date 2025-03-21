package fetch

import (
	"bytes"
	"encoding/json"
	"ewails/wrappers/common"
	"fmt"
	"net/http"
	"net/url"
)

type Api struct {
	*ApiOptions
	Url string
}

func New(url string, optsFns ...ApiOptsFn) *Api {
	baseOptions := &ApiOptions{
		Headers: map[string]string{
			// "Content-Type": "application/json",
		},
		Client: &http.Client{},
	}

	options := common.CreateOptions(baseOptions, optsFns...)

	return &Api{
		Url:        url,
		ApiOptions: options,
	}
}

func CreateQuery[T comparable](Url string, method string, options *QueryOptions[T]) (*http.Request, error) {
	params := url.Values{}

	for key, val := range options.Query {
		params.Add(key, val)
	}

	u, err := url.Parse(Url)

	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	u.RawQuery = params.Encode()

	body := []byte("")

	if options.Body != nil {
		body, err = json.Marshal(options.Body)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for key, val := range options.Headers {
		req.Header.Add(key, val)
	}

	return req, err
}

func (a *Api) ExecuteRequest(method string, Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	baseOptions := &QueryOptions[any]{
		ApiOptions: *a.ApiOptions,
	}

	options := common.CreateOptions(baseOptions, optFns...)
	path := common.ComposePath(a.Url, Url)

	req, err := CreateQuery(path[:len(path)-1], method, options)

	fmt.Println(req)

	if err != nil {
		return nil, err
	}

	return options.Client.Do(req)
}

func (a *Api) Get(Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	return a.ExecuteRequest("GET", Url, optFns...)
}

func (a *Api) Post(Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	return a.ExecuteRequest("POST", Url, optFns...)
}

func (a *Api) Put(Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	return a.ExecuteRequest("PUT", Url, optFns...)
}

func (a *Api) Patch(Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	return a.ExecuteRequest("PATCH", Url, optFns...)
}

func (a *Api) Delete(Url string, optFns ...QueryOptsFn) (*http.Response, error) {
	return a.ExecuteRequest("DELETE", Url, optFns...)
}

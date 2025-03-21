package fetch

import (
	"bytes"
	"encoding/json"
	"ewails/wrappers/common"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"maps"
)

type ApiOptions struct {
	Headers map[string]string
	Client  *http.Client
}

type QueryOptions[T comparable] struct {
	ApiOptions
	Query map[string]string
	Body  *T
}

type ApiOptsFn = common.OptFns[ApiOptions]
type QueryOptsFn = common.OptFns[QueryOptions[any]]

// Api Options
func WithBaseHeaders(headers map[string]string) ApiOptsFn {
	return func(options *ApiOptions) {

		if options.Headers == nil {
			options.Headers = headers
			return
		}

		maps.Copy(options.Headers, headers)
	}
}

func WithClient(client *http.Client) ApiOptsFn {
	return func(options *ApiOptions) {
		options.Client = client
	}
}

// Query Options
func WithHeaders(headers map[string]string) QueryOptsFn {
	return func(options *QueryOptions[any]) {

		if options.Headers == nil {
			options.Headers = headers
			return
		}

		maps.Copy(options.Headers, headers)
	}
}

func WithQuery(query map[string]string) QueryOptsFn {
	return func(options *QueryOptions[any]) {

		if options.Query == nil {
			options.Query = query
			return
		}

		maps.Copy(options.Query, query)
	}
}

func WithBody(body any) QueryOptsFn {
	return func(options *QueryOptions[any]) {
		options.Body = &body
	}
}

// API Calls
func AsJson[T any](res *http.Response, err error) (data *T, e error) {
	data = new(T)

	if err != nil {
		return nil, err
	}

	// defer res.Body.Close()

	// bodyCopy := &res.Body

	// json.NewDecoder(res.Body).Decode(&data)

	// Copy the body to a buffer first
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return nil, err
	}

	// Log the raw response body (as a string)
	fmt.Println("Response Body:", string(body))

	// Recreate the body reader using the copied body so we can decode
	res.Body = io.NopCloser(bytes.NewReader(body))

	// Decode the JSON response body directly into the data structure
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		log.Fatal("Error decoding JSON:", err)
		return nil, err
	}

	// Optionally log the decoded data to check
	fmt.Printf("Decoded Data: %+v\n", data)

	return data, err
}

func AsUrlEncoded[T UrlEncodedResponse](res *http.Response, err error) (T, error) {
	var data T

	if err != nil {
		return data, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error reading response body:", err)
		return data, err
	}

	// Log the raw response body
	fmt.Println("Response Body:", string(body))

	// Parse the URL-encoded form data
	params, err := url.ParseQuery(string(body))

	if err != nil {
		log.Fatal("Error parsing URL-encoded body:", err)
		return data, err
	}

	fmt.Println("creating params")

	// Create an instance of T correctly
	if reflect.TypeOf(data).Kind() == reflect.Pointer {
		// If T is already a pointer, use it directly
		data = reflect.New(reflect.TypeOf(data).Elem()).Interface().(T)
	} else {
		// If T is a struct, create a pointer to it
		ptr := reflect.New(reflect.TypeOf(data))
		data = ptr.Interface().(T)
	}

	err = data.FromParams(params)

	if err != nil {
		fmt.Println("error creating params", err.Error())
		return data, err
	}

	return data, nil
}

type UrlEncodedResponse interface {
	FromParams(params url.Values) error
}

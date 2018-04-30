package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var apiURL = "https://kv-api.dyson-sphere.com/"

var prettyServerError = map[int]string{
	500: "internal server error",
	501: "not implemented",
	502: "bad gateway",
	503: "service unavailable",
	504: "gateway timeout",
	505: "HTTP version not support",
	506: "variant also negotiates",
	507: "insufficient storage",
	508: "loop detected",
	510: "not extended",
	511: "network authentication required",
}

func CreateStore(storeName string) error {
	return simpleRequestWrapper(makeRequest("PUT", storeURL(storeName), nil))
}

func SetKey(storeName string, keyName string, keyValue string) error {
	return simpleRequestWrapper(makeRequest("PUT", keyURL(storeName, keyName), strings.NewReader(keyValue)))
}

func DeleteStore(storeName string) error {
	return simpleRequestWrapper(makeRequest("DELETE", storeURL(storeName), nil))
}

func DeleteKey(storeName string, keyName string) error {
	return simpleRequestWrapper(makeRequest("DELETE", keyURL(storeName, keyName), nil))
}

func GetKey(storeName string, keyName string) ([]byte, error) {
	resp, err := makeRequest("GET", keyURL(storeName, keyName), nil)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, errors.New(fmt.Sprintf("reading response body (%s)", err))
		}

		return body, nil
	}

	return []byte{}, errorFromResponse(resp)
}

func IsNotFoundError(err error) bool {
	return err.Error() == "not found"
}

type apiError struct {
	Error string `json:"error"`
}

func errorFromResponse(resp *http.Response) error {
	if resp.StatusCode >= 500 {
		return errors.New(fmt.Sprintf("server error (%s)", prettyServerError[resp.StatusCode]))
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		if resp.StatusCode == 404 {
			return errors.New("not found")
		}

		if resp.StatusCode == 413 {
			return errors.New("payload too large")
		}

		var apiErr apiError
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("unknown error (response code %d)", resp.StatusCode))
		}
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			return errors.New(fmt.Sprintf("unknown error (response code %d)", resp.StatusCode))
		}

		return errors.New(apiErr.Error)
	}

	return errors.New(fmt.Sprintf("unknown error (%v)", resp))
}

func storeURL(storeName string) string {
	return apiURL + "/stores/" + storeName
}

func keyURL(storeName string, keyName string) string {
	return storeURL(storeName) + "/keys/" + keyName
}

func simpleRequestWrapper(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errorFromResponse(resp)
}

func makeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &http.Response{}, err
	}

	client := &http.Client{}
	return client.Do(req)
}

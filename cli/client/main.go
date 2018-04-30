package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var apiURL = "https://kv-api.dyson-sphere.com/"

func CreateStore(name string) error {
	req, err := http.NewRequest("PUT", apiURL+"/stores/"+name, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errorFromResponse(resp)
}

type apiError struct {
	Error string `json:"error"`
}

func errorFromResponse(resp *http.Response) error {
	if resp.StatusCode >= 500 {
		return errors.New(fmt.Sprintf("server error (%s)", prettyServerError(resp.StatusCode)))
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var apiErr apiError
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("unknown error (response code %s)", resp.StatusCode))
		}
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			return errors.New(fmt.Sprintf("unknown error (response code %s)", resp.StatusCode))
		}

		return errors.New(apiErr.Error)
	}

	return errors.New(fmt.Sprintf("unknown error (%v)", resp))
}

func prettyServerError(code int) string {
	codes := map[int]string{
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

	return codes[code]
}

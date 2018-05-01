package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyURL(t *testing.T) {
	expected := apiURL() + "/stores/foo/keys/bar"
	assert.Equal(t, expected, keyURL("foo", "bar"))
}

func TestStoreURL(t *testing.T) {
	expected := apiURL() + "/stores/foo"
	assert.Equal(t, expected, storeURL("foo"))
}

func TestIsNotFoundError(t *testing.T) {
	assert.Equal(t, true, IsNotFoundError(errors.New("not found")))
	assert.Equal(t, false, IsNotFoundError(errors.New("blah")))
}

func TestErrorFromResponse(t *testing.T) {
	var tests = []struct {
		response    *http.Response
		expectedErr error
	}{
		{
			response:    &http.Response{StatusCode: 500},
			expectedErr: errors.New("server error (internal server error)"),
		},
		{
			response:    &http.Response{StatusCode: 501},
			expectedErr: errors.New("server error (not implemented)"),
		},
		{
			response: &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error": "this is an api error"}`)),
			},
			expectedErr: errors.New("this is an api error"),
		},
		{
			response: &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(strings.NewReader("blah blah")),
			},
			expectedErr: errors.New("unknown error (response code 400)"),
		},
		{
			response:    &http.Response{StatusCode: 404},
			expectedErr: errors.New("not found"),
		},
		{
			response:    &http.Response{StatusCode: 413},
			expectedErr: errors.New("payload too large"),
		},
		{
			response:    &http.Response{StatusCode: 100},
			expectedErr: errors.New(fmt.Sprintf("unknown error (%v)", &http.Response{StatusCode: 100})),
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedErr, errorFromResponse(tc.response))
	}

}

func TestSimpleRequestWrapper(t *testing.T) {
	var tests = []struct {
		response      *http.Response
		err           error
		errorExpected bool
	}{
		{
			response:      &http.Response{StatusCode: 200},
			err:           nil,
			errorExpected: false,
		},
		{
			response:      &http.Response{StatusCode: 200},
			err:           errors.New(""),
			errorExpected: true,
		},
		{
			response:      &http.Response{StatusCode: 201},
			err:           nil,
			errorExpected: false,
		},
		{
			response:      &http.Response{StatusCode: 201},
			err:           errors.New(""),
			errorExpected: true,
		},
		{
			response: &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error": "this is an api error"}`)),
			},
			err:           nil,
			errorExpected: true,
		},
		{
			response: &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error": "this is an api error"}`)),
			},
			err:           errors.New(""),
			errorExpected: true,
		},
		{
			response:      &http.Response{StatusCode: 500},
			errorExpected: true,
		},
		{
			response:      &http.Response{StatusCode: 500},
			err:           errors.New(""),
			errorExpected: true,
		},
	}

	for _, tc := range tests {
		if tc.errorExpected {
			assert.Error(t, simpleRequestWrapper(tc.response, tc.err))
		} else {
			assert.NoError(t, simpleRequestWrapper(tc.response, tc.err))
		}
	}
}

/********************************************************************************
 *  Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

// SampleConsulResponse a slimmed down example of the contents returned by Consul with varying statuses.
const SampleConsulResponse = `
{
    "service:kong": {
        "Node": "edgex-core-consul",
        "CheckID": "service:kong",
        "Name": "Health Check: kong",
        "Status": "passing",
        "Notes": "",
        "Output": "Get http://kong:8001/status: dial tcp: lookup kong on 127.0.0.11:53: no such host",
        "ServiceID": "kong",
        "ServiceName": "kong",
        "ServiceTags": [],
        "Definition": {},
        "CreateIndex": 0,
        "ModifyIndex": 0
    },
    "service:kong-db": {
        "Node": "edgex-core-consul",
        "CheckID": "service:kong-db",
        "Name": "Health Check: kong-db (postgres)",
        "Status": "critical",
        "Notes": "",
        "Output": "fork/exec /consul/scripts/kong-db-checker.sh: no such file or directory",
        "ServiceID": "kong-db",
        "ServiceName": "kong-db",
        "ServiceTags": [],
        "Definition": {},
        "CreateIndex": 0,
        "ModifyIndex": 0
    }
}`

func mockGetter(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(SampleConsulResponse)),
	}, nil
}

func mockErrorGetter(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewBufferString("Some failure message")),
	}, errors.New("test error")
}

func TestCheckServiceStatus(t *testing.T) {
	tests := []struct {
		name           string
		getter         HTTPGetFunc
		serviceName    string
		status         string
		expectedResult bool
		expectError    bool
	}{
		{
			"Passing status which matches",
			mockGetter,
			"kong",
			"passing",
			true,
			false,
		},
		{
			"Passing status which does not match",
			mockGetter,
			"kong",
			"Good",
			false,
			false,
		},
		{
			"Critical status",
			mockGetter,
			"kong-db",
			"Passing",
			false,
			false,
		},
		{
			"HTTP error",
			mockErrorGetter,
			"", // Due to the HTTP error this is not necessary as it is never referenced
			"", // Due to the HTTP error this is not necessary as it is never referenced
			false,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CheckServiceStatus(
				test.getter,
				"URL is not needed since we are mocking the HTTP client",
				test.serviceName,
				test.status)

			if test.expectError && err == nil {
				t.Error("Expected error, but non was encountered")
			}

			if !test.expectError && err != nil {
				t.Errorf("Encountered unexpected error: %s", err.Error())
			}

			if test.expectedResult != test.expectedResult {
				t.Errorf("Expected: %t , but encountered: %t", test.expectedResult, result)
			}
		})
	}
}

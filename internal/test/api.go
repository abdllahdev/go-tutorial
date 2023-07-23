package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name             string
	Method, URL      string
	Body             interface{}
	ExpectedStatus   int
	ExpectedResponse map[string]interface{}
}

// Endpoint tests an HTTP endpoint using the given APITestCase spec.
func Endpoint(t *testing.T, router *mux.Router, testCase APITestCase) {
	t.Run(
		testCase.Name,
		func(t *testing.T) {
			// Create a new request
			jsonData, err := json.Marshal(testCase.Body)
			if err != nil {
				fmt.Println("Error:", err)
			}

			req, err := http.NewRequest(testCase.Method, testCase.URL, bytes.NewBuffer(jsonData))
			if err != nil {
				t.Errorf("Error creating new request: %s", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create new recorder
			res := httptest.NewRecorder()

			// Send request
			router.ServeHTTP(res, req)

			// Check status code
			assert.Equal(t, testCase.ExpectedStatus, res.Code, "status mismatch")

			// Compare expected and actual responses
			if testCase.ExpectedResponse != nil && res.Body != nil {
				var actualResponse map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&actualResponse)
				if err != nil {
					t.Errorf("Error decoding response body, %v", err)
				}

				assertMapEq(t, testCase.ExpectedResponse, actualResponse, "response mismatch")
			}
		})
}

func assertMapEq(t *testing.T, expected, actual interface{}, msg string) {
	expectedMap, ok1 := expected.(map[string]interface{})
	actualMap, ok2 := actual.(map[string]interface{})

	if ok1 && ok2 {
		for field, expectedValue := range expectedMap {
			actualValue := actualMap[field]

			if reflect.TypeOf(expectedValue).Kind() == reflect.Map && reflect.TypeOf(actualValue).Kind() == reflect.Map {
				expectedMap, ok1 = expectedValue.(map[string]interface{})
				actualMap, ok2 = actualValue.(map[string]interface{})

				if ok1 && ok2 {
					assertMapEq(t, expectedMap, actualMap, msg)
				}
			} else if reflect.TypeOf(expectedValue).Kind() == reflect.Slice && reflect.TypeOf(actualValue).Kind() == reflect.Slice {
				expectSlice, ok1 := expectedValue.([]interface{})
				actualSlice, ok2 := actualValue.([]interface{})

				if ok1 && ok2 {
					for i := range expectSlice {
						log.Printf("expectedSlice: %v, actualSlice: %v", expectSlice[i], actualSlice[i])
						assertMapEq(t, expectSlice[i], actualSlice[i], msg)
					}
				}
			} else {
				assert.EqualValues(t, expectedValue, actualValue, msg)
			}
		}
	} else {
		assert.EqualValues(t, expected, actual, msg)
	}
}

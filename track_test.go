package dhl

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(t *testing.T, statusCode int, fixture string) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			body, err := os.Open(fixture)
			if err != nil {
				t.Fatalf("Couldn't open textfixture %v: %v", fixture, err)
			}

			return &http.Response{
				StatusCode: statusCode,
				// Send response to be tested
				Body: ioutil.NopCloser(body),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}),
	}
}

func TestShipments(t *testing.T) {
	client := NewTestClient(t, 200, "testdata/track/ok.json")
	service := NewTrackingService(client)
	ans, err := service.Shipments("11111111111111111XXXXXXX")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	gotStatus := ans.Shipments[0].Status.Status
	wantStatus := "PROCESSED AT LOCAL DISTRIBUTION CENTER"
	if gotStatus != wantStatus {
		t.Errorf("Got: %v; Want: %v", gotStatus, wantStatus)
	}

	gotTime := ans.Shipments[0].Status.Timestamp
	wantTime := time.Date(2020, 4, 20, 12, 12, 00, 0, time.FixedZone("UTC+6", +6*60*60))
	if !gotTime.Equal(wantTime) {
		t.Errorf("Got: %v; Want: %v", gotTime, wantTime)
	}
}

func TestShipments_ErrorHandling(t *testing.T) {
	client := NewTestClient(t, 404, "testdata/track/404.json")
	service := NewTrackingService(client)
	_, err := service.Shipments("11111111111111111XXXXXXX")
	want := "No shipment with given tracking number found."
	if err.Error() != want {
		t.Fatalf("Got: %v; Want: %v", err, want)
	}
}

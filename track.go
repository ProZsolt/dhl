package dhl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// ProofOfDelivery is proof of delivery
type ProofOfDelivery struct {
	Timestamp   time.Time
	DocumentURL string
	signed      Person
}

// Shipment information including full shipment details and complete list of shipment events.
type Shipment struct {
	ID                         string
	Service                    string // enum: freight, express, parcel-de, parcel-nl, parcel-pl, dsc, dgf, ecommerce
	Origin                     Place
	Destination                Place
	Status                     ShipmentEvent
	EstimatedTimeOfDelivery    interface{}
	EstimatedDeliveryTimeFrame struct {
		EstimatedFrom    interface{}
		EstimatedThrough interface{}
	}
	EstimatedTimeOfDeliveryRemark interface{}
	Details                       ShipmentDetails
	Events                        []ShipmentEvent
}

// ShipmentDetails represents a detailed information about one shipmen
type ShipmentDetails struct {
	Carrier             Organization
	Product             Product
	Receiver            Person
	Sender              Person
	ProofOfDelivery     ProofOfDelivery
	TotalNumberOfPieces int
	PieceIds            []string
	Weight              interface{} // http://schema.org/weight
	Volume              interface{} // http://schema.org/cargoVolume
	LoadingMeters       float32
	Dimensions          struct {
		Width  interface{} // http://schema.org/QuantitativeValue
		Height interface{} // http://schema.org/QuantitativeValue
		Length interface{} // http://schema.org/QuantitativeValue
	}
	References []struct {
		Number string
		Type   string // enum: customer-reference, customer-confirmation-number, local-tracking-number, ecommerce-number, housebill, masterbill, container-number, domestic-consignment-id
	}
	DGFRoutes []DGFRoute `json:"dgf:routes"`
}

// ShipmentEvent is an event in shipment delivery; also known as Milestone, Checkpoint, Status History Entry or http://schema.org/DeliveryEvent
type ShipmentEvent struct {
	timestamp   time.Time
	location    Place
	StatusCode  string // enum: pre-transit, transit, delivered, failure, unknown
	Status      string
	Description string
	Remark      string
	NextSteps   string
}

//Shipments is a list of shipments matching the input query
type Shipments struct {
	Shipments                      []Shipment
	PossibleAdditionalShipmentsURL []string
}

// ProblemDetail is a definition of RFC7807 Problem Detail for HTTP APIs
type ProblemDetail struct {
	Type     string
	Title    string
	Status   int
	Detail   string
	Instance string
}

// DGFAirport is an airport model description
type DGFAirport struct {
	LocationName string `json:"dgf:locationName"`
	LocationCode string `json:"dgf:locationCode"`
	CountryCode  string
}

// DGFLocation is a definition of a place - location
type DGFLocation struct {
	LocationName string `json:"dgf:locationName"`
}

// DGFRoute is a definition of a route
type DGFRoute struct {
	VesselName             string      `json:"dgf:vesselName"`
	VoyageFlightNumber     string      `json:"dgf:voyageFlightNumber"`
	AirportOfDeparture     DGFAirport  `json:"dgf:airportOfDeparture"`
	AirportOfDestination   DGFAirport  `json:"dgf:airportOfDestination"`
	EstimatedDepartureDate time.Time   `json:"dgf:estimatedDepartureDate"`
	EstimatedArrivalDate   time.Time   `json:"dgf:estimatedArrivalDate"`
	PlaceOfAcceptance      DGFLocation `json:"dgf:placeOfAcceptance"`
	PortOfLoading          DGFLocation `json:"dgf:portOfLoading"`
	PortOfUnloading        DGFLocation `json:"dgf:portOfUnloading"`
	PlaceOfDelivery        DGFLocation `json:"dgf:placeOfDelivery"`
}

// DataType is the basic data types such as Integers, Strings, etc.
type DataType interface{}

// Person (see: http://schema.org/Person)
type Person struct {
	FamilyName string
	GivenName  string
	Name       string
}

// Place model description. https://gs1.org/voc/Place
type Place struct {
	Address struct {
		CountryCode     string
		PostalCode      string
		AddressLocality string
		StreetAddress   string
	}
}

// Organization model description. https://gs1.org/voc/Organization
type Organization struct {
	OrganizationName string
}

// Product used for the shipment. https://gs1.org/voc/Product
type Product struct {
	ProductName string
}

// TrackingService is a sercice for tracking calls
type TrackingService struct {
	client  *http.Client
	baseURL string
}

// NewTrackingService creates a new TrackingService.
// It uses the provided http.Client for requests.
func NewTrackingService(client *http.Client) TrackingService {
	return TrackingService{
		client:  client,
		baseURL: "https://api-eu.dhl.com/track",
	}
}

// Shipments retrieves the tracking information for shipments(s).
// The shipments are identified using the required trackingNumber query parameter.
func (ts TrackingService) Shipments(trackingNumber string) (Shipments, error) {
	var shipments Shipments

	req, err := http.NewRequest(http.MethodGet, ts.baseURL+"/shipments", nil)
	if err != nil {
		return shipments, err
	}

	q := req.URL.Query()
	q.Add("trackingNumber", trackingNumber)
	req.URL.RawQuery = q.Encode()

	resp, err := ts.client.Do(req)
	if err != nil {
		return shipments, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return shipments, err
	}

	if resp.StatusCode != 200 {
		var pd ProblemDetail
		err = json.Unmarshal(body, &pd)
		if err != nil {
			return shipments, err
		}
		return shipments, errors.New(pd.Detail)
	}

	err = json.Unmarshal(body, &shipments)
	return shipments, err
}

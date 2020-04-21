# DHL

DHL API support for GoLang.

This is WIP, currently only supports Shipment Tracking API

Download it using

``go get github.com/prozsolt/dhl``

Import it into your project

```go
package main

import "github.com/prozsolt/dhl"
```

### Authentication
Authentication onto DHL API is using the Consumer Key 

* From the [My Apps](https://developer.dhl.com/user/apps) screen, click on the name of your app.
The Details screen appears. 
* If you have access to more than one API, click the name of the relevant API. 
Note: The APIs are listed under the “Credentials” section. 
* Click the Show link below the asterisks that is hiding the Consumer Key. 
The Consumer Key appears.    

```go
client := NewClient("YOUR_CONSUMER_KEY")
```

### Shipment Tracking API

For more information go [here](https://developer.dhl.com/api-reference/shipment-tracking).

```go
client := NewClient("YOUR_CONSUMER_KEY")
service := NewTrackingService(client)
shipments, err := service.Shipments("YOUR_TRACKING_NUMBER")
if err != nil {
  fmt.Println(err)
}
fmt.Println(Shipments[0].Status.Status)
```
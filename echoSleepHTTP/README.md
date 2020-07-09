#EchoSleepHTTP Service

The EchoSleep service will echo all HTTP body payloads, query parameters, and path parameters.

The response wrapper will be used:
``````
{
uri:           "",
pathParams:    [""],
queryParams:   [{key: "value"}],
headerParams:  [{key:"value"}].
requestBody:   {""}
}
``````

The EchoSleep service respect the following optional request parameters:

* X-Sleep  (examples: 1s 1m 1hr)  use time.ParseDuration(sleepTime) convention 
* X-

The EchoSleep service will return the following response headers:

* X-serverCorrelation: uid 



## To Echo

curl -X GET http://localhost/echosleepHttp 




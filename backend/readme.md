`export TRAFIKVERKET_AUTH_KEY=...`
`export FUNCTION_TARGET=HelloHTTP`
`go run local/main.go`
http://localhost:8080/

## Set cors for bucket
`gsutil cors set bucket-cors.json gs://live-weather-eefc5.appspot.com`
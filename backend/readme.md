`export TRAFIKVERKET_AUTH_KEY=...`
`export FUNCTION_TARGET=HelloHTTP`
`go run local/main.go`
http://localhost:8080/

## Build deploy image locally
`pack build imageName --builder gcr.io/buildpacks/builder:v1`
Run image locally
`docker run -p 8080:8080 imageName`

## Pull pipeline image from gcr.io
`gcloud auth configure-docker europe-north1-docker.pkg.dev`
`docker pull imageName`

## Set cors for bucket
`gsutil cors set bucket-cors.json gs://live-weather-eefc5.appspot.com`
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"

go build -buildvcs=false -o ./dist/reader
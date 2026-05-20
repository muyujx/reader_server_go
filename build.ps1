$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"

go build -buildvcs=false -o ./dist/reader


docker build . -t reader_server
if (!$?)
{
    Write-Error 'docker build fail!'
    exit 1
}

New-Item -ItemType Directory -Path "./docker" -ErrorAction SilentlyContinue

docker image save reader_server -o docker/reader_server.image
if (!$?)
{
    Write-Error 'docker image save fail'
    exit 1
}

param (
    [string]$command
)

switch ($command) {
    "run" { go run main.go }
    "build" { go build -o ./bin/api ./cmd/api/main.go }
    "test" { go test -v ./... }
    default { Write-Host "Usage: ./build.ps1 [run|build|test]" }
}

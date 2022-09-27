set -e
mkdir -p bin

ROOT=$PWD
cd "$ROOT/src" && go test ./...
cd "$ROOT/src" && go build -o ../bin/forthly
cd "$ROOT/src" && GOOS=js GOARCH=wasm go build -o ../bin/forthly.wasm

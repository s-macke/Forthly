set -e
mkdir -p bin

cd src && go test ./... && cd ..
cd src && go build -o ../bin/forthly && cd ..
cd src && GOOS=js GOARCH=wasm go build -o ../bin/forthly.wasm && cd ..


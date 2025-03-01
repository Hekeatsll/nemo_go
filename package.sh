#!/usr/bin/env bash

echo "package nemo for macos,linux and windows..."
echo "MAKS SURE \"RunMode = Release\" in pkg/conf/config.go"
rm -rf release/*
echo > log/runtime.log
echo > log/access.log

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o server_darwin_amd64 cmd/server/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o worker_darwin_amd64 cmd/worker/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o daemon_worker_darwin_amd64 cmd/daemon_worker/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o server_linux_amd64 cmd/server/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o worker_linux_amd64 cmd/worker/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o daemon_worker_linux_amd64 cmd/daemon_worker/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o server_windows_amd64.exe cmd/server/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o worker_windows_amd64.exe cmd/worker/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o daemon_worker_windows_amd64.exe cmd/daemon_worker/main.go

tar -cvzf release/nemo_darwin_amd64.tar \
  --exclude=thirdparty/xray/xray_linux_amd64 \
  --exclude=thirdparty/xray/xray_windows_amd64.exe \
  --exclude=thirdparty/nuclei/nuclei_linux_amd64 \
  --exclude=thirdparty/nuclei/nuclei_windows_amd64.exe \
  --exclude=thirdparty/fingerprinthub/observer_ward_linux_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_windows_amd64.exe \
  server_darwin_amd64 worker_darwin_amd64 daemon_worker_darwin_amd64 version.txt \
  conf log thirdparty web

tar -cvzf release/nemo_linux_amd64.tar \
  --exclude=thirdparty/xray/xray_darwin_amd64 \
  --exclude=thirdparty/xray/xray_windows_amd64.exe \
  --exclude=thirdparty/nuclei/nuclei_darwin_amd64 \
  --exclude=thirdparty/nuclei/nuclei_windows_amd64.exe \
  --exclude=thirdparty/fingerprinthub/observer_ward_darwin_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_windows_amd64.exe \
  server_linux_amd64 worker_linux_amd64 daemon_worker_linux_amd64 package_worker.sh version.txt \
  conf log thirdparty web docker* Dockerfile*

tar -cvzf release/nemo_windows_amd64.tar \
  --exclude=thirdparty/xray/xray_darwin_amd64 \
  --exclude=thirdparty/xray/xray_linux_amd64 \
  --exclude=thirdparty/nuclei/nuclei_darwin_amd64 \
  --exclude=thirdparty/nuclei/nuclei_linux_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_darwin_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_linux_amd64 \
  server_windows_amd64.exe worker_windows_amd64.exe version.txt \
  conf log thirdparty web

tar -cvzf release/worker_linux_amd64.tar \
  --exclude=thirdparty/xray/xray_darwin_amd64 \
  --exclude=thirdparty/xray/xray_windows_amd64.exe \
  --exclude=thirdparty/nuclei/nuclei_darwin_amd64 \
  --exclude=thirdparty/nuclei/nuclei_windows_amd64.exe \
  --exclude=thirdparty/fingerprinthub/observer_ward_darwin_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_windows_amd64.exe \
  --exclude=conf/app.conf --exclude=server.yml \
  worker_linux_amd64 daemon_worker_linux_amd64 conf log thirdparty version.txt

tar -cvzf release/worker_darwin_amd64.tar \
  --exclude=thirdparty/xray/xray_linux_amd64 \
  --exclude=thirdparty/xray/xray_windows_amd64.exe \
  --exclude=thirdparty/nuclei/nuclei_linux_amd64 \
  --exclude=thirdparty/nuclei/nuclei_windows_amd64.exe \
  --exclude=thirdparty/fingerprinthub/observer_ward_linux_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_windows_amd64.exe \
  --exclude=conf/app.conf --exclude=server.yml \
  worker_darwin_amd64 daemon_worker_darwin_amd64 conf log thirdparty version.txt

tar -cvzf release/worker_windows_amd64.tar \
  --exclude=thirdparty/xray/xray_darwin_amd64 \
  --exclude=thirdparty/xray/xray_linux_amd64 \
  --exclude=thirdparty/nuclei/nuclei_darwin_amd64 \
  --exclude=thirdparty/nuclei/nuclei_linux_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_darwin_amd64 \
  --exclude=thirdparty/fingerprinthub/observer_ward_linux_amd64 \
  --exclude=conf/app.conf --exclude=server.yml \
  worker_windows_amd64.exe daemon_worker_windows_amd64.exe conf log thirdparty version.txt

rm -f server_darwin_amd64 worker_darwin_amd64 daemon_worker_darwin_amd64 \
  server_linux_amd64 worker_linux_amd64 daemon_worker_linux_amd64 \
  server_windows_amd64.exe worker_windows_amd64.exe daemon_worker_windows_amd64.exe \

echo "package done..."

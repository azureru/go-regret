#!/bin/bash

gox -osarch="linux/amd64" -osarch="linux/arm" -osarch="darwin/amd64" -osarch="windows/amd64"
rm *.zip
zip go-regret_darwin_amd64.zip go-regret_darwin_amd64
zip go-regret_linux_amd64.zip go-regret_linux_amd64
zip go-regret_linux_arm.zip go-regret_linux_arm
zip go-regret_windows_amd64.zip go-regret_windows_amd64.exe

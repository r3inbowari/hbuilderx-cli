#!/bin/bash

# author r3inbowari
# date 2021/10/30
# 编译前请确保已安装 git

packageName=cmd
appName=meiwobuxing
buildVersion=v1.7.7
major=1
minor=7
patch=7
mode=REL

goVersion=$(go version)
gitHash=$(git show -s --format=%H)
buildTime=$(git show -s --format=%cd)

echo ===================================================
echo "                 Go build running"
echo ===================================================
echo $goVersion
echo build hash $gitHash
echo build time $buildTime
echo build tag $buildVersion
echo ===================================================

cd cmd
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

GOOS=windows
GOARCH=amd64
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=${major}' -X 'main.minor=${minor}'-X 'main.patch=${patch}' -X 'main.releaseVersion=${buildVersion}' -X 'main.mode=${mode}' -X 'main.goVersion=${goVersion}' -X 'main.gitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}
md5 ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=darwin
GOARCH=amd64
go env -w GOOS=darwin
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=${major}' -X 'main.minor=${minor}'-X 'main.patch=${patch}' -X 'main.releaseVersion=${buildVersion}' -X 'main.mode=${mode}' -X 'main.goVersion=${goVersion}' -X 'main.gitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}
md5 ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
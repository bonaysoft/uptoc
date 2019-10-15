# uptoc

[![Build Status](https://travis-ci.org/saltbo/uptoc.svg?branch=master)](https://travis-ci.org/saltbo/uptoc)&nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/saltbo/uptoc)](https://goreportcard.com/report/github.com/saltbo/uptoc)&nbsp;
[![Coverage Status](https://coveralls.io/repos/github/saltbo/uptoc/badge.svg?branch=master)](https://coveralls.io/github/saltbo/uptoc?branch=master)

`uptoc` is tool to upload the dist file for the cloud engine.

## Run environment
- Mac
- Linux
- Windows

## Dependent libraries 
- cli (github.com/urfave/cli) 
- oss (github.com/aliyun/aliyun-oss-go-sdk/oss)

## Install

Download the appropriate binary for your platform from the [Releases](https://github.com/saltbo/uptoc/releases) page, or:

```sh
curl -sSf http://uptoc.saltbo.cn/install.sh | sh
```

## Usage
```sh
uptoc --endpoint oss-cn-beijing.aliyuncs.com --bucket demo-bucket --access_key LTAI4FxxxxxxxBXmS3 --access_secret Vt1FZgxxxxxxxxxxxxKp380AI /opt/blog/public
```


## Contact us
- [Author Blog](https://saltbo.cn).

## Author
- [Saltbo](https://github.com/saltbo)

## License
- [MIT](https://github.com/saltbo/uptoc/blob/master/LICENSE)

# uptoc

[![Build Status](https://travis-ci.org/saltbo/uptoc.svg?branch=master)](https://travis-ci.org/saltbo/uptoc)&nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/saltbo/uptoc)](https://goreportcard.com/report/github.com/saltbo/uptoc)&nbsp;
[![Coverage Status](https://coveralls.io/repos/github/saltbo/uptoc/badge.svg?branch=master)](https://coveralls.io/github/saltbo/uptoc?branch=master)

`uptoc` is a cli tool for deploying files to the cloud storage.

## Run environment
- Mac
- Linux
- Windows

## Dependent libraries 
- cli (github.com/urfave/cli) 
- oss (github.com/aliyun/aliyun-oss-go-sdk/oss)

## Install

Download the appropriate binary for your platform from the [Releases](https://github.com/saltbo/uptoc/releases) page, or:

```bash
curl -sSf http://uptoc.saltbo.cn/install.sh | sh
```

## Usage
```bash
uptoc --endpoint oss-cn-beijing.aliyuncs.com --access_key LTAI4FxxxxxxxBXmS3 --access_secret Vt1FZgxxxxxxxxxxxxKp380AI --bucket demo-bucket /opt/blog/public
```

And the access-key and access-secret support settings by the system environment
```bash
export UPTOC_UPLOADER_KEYID=LTAI4FxxxxxxxBXmS3
export UPTOC_UPLOADER_KEYSECRET=Vt1FZgxxxxxxxxxxxxKp380AI

uptoc --endpoint oss-cn-beijing.aliyuncs.com --bucket blog-bucket /opt/blog/public
```

So you can use it like this for the travis
```yaml
after_success:
  - curl -sSf http://uptoc.saltbo.cn/install.sh | sh
  - uptoc --endpoint uploader-cn-zhangjiakou.aliyuncs.com --bucket blog-bucket public
```

## Contact us
- [Author Blog](https://saltbo.cn).

## Author
- [Saltbo](https://github.com/saltbo)

## License
- [MIT](https://github.com/saltbo/uptoc/blob/master/LICENSE)

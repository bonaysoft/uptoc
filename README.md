# Uptoc

[![](https://github.com/saltbo/uptoc/workflows/build/badge.svg)](https://github.com/saltbo/uptoc/actions?query=workflow%3Abuild)
[![](https://codecov.io/gh/saltbo/uptoc/branch/master/graph/badge.svg)](https://codecov.io/gh/saltbo/uptoc)
[![](https://api.codacy.com/project/badge/Grade/88817db9b3b04c0293c9d001d574a5ef)](https://app.codacy.com/manual/saltbo/uptoc?utm_source=github.com&utm_medium=referral&utm_content=saltbo/uptoc&utm_campaign=Badge_Grade_Dashboard)
[![](https://img.shields.io/github/v/release/saltbo/uptoc.svg)](https://github.com/saltbo/uptoc/releases)
[![](https://img.shields.io/github/license/saltbo/uptoc.svg)](https://github.com/saltbo/uptoc/blob/master/LICENSE)

English | [🇨🇳中文](https://saltbo.cn/uptoc)

## Run environment
- Mac
- Linux
- Windows

## Support Driver 
- Aliyun OSS
- Tencent COS
- Qiniu Kodo
- Google Storage
- AWS S3

## Install

Download the appropriate binary for your platform from the [Releases](https://github.com/saltbo/uptoc/releases) page, or:

```bash
curl -sSf http://uptoc.saltbo.cn/install.sh | sh
```

## Usage

### Basic
```bash
uptoc --driver oss --region cn-beijing --access_key LTAI4FxxxxxxxBXmS3 --access_secret Vt1FZgxxxxxxxxxxxxKp380AI --bucket demo-bucket /opt/blog/public
```

And the access-key and access-secret support settings by the system environment
```bash
export UPTOC_UPLOADER_AK=LTAI4FxxxxxxxBXmS3
export UPTOC_UPLOADER_SK=Vt1FZgxxxxxxxxxxxxKp380AI

uptoc --driver oss --region cn-beijing --bucket blog-bucket /opt/blog/public
```

### Github Actions
```yml
steps:
  - name: Deploy
    uses: saltbo/uptoc@master
    with:
      driver: oss
      region: cn-zhangjiakou
      bucket: saltbo-blog
      dist: public
    env:
      UPTOC_UPLOADER_AK: ${{ secrets.UPTOC_UPLOADER_KEYID }}
      UPTOC_UPLOADER_SK: ${{ secrets.UPTOC_UPLOADER_KEYSECRET }}
```
### Similar Travis 
```yaml
after_success:
  - curl -sSf http://uptoc.saltbo.cn/install.sh | sh
  - uptoc --region cn-zhangjiakou --bucket blog-bucket public
```

## Args Examples
| driver | bucket | region | region enum |
| -----  | --------- | ------ | ---- |
| oss    | ut-uptoc  | cn-hangzhou | [Regions](https://help.aliyun.com/document_detail/31837.html?spm=a2c4g.11186623.2.12.5fdb25b7xyEcuF#concept-zt4-cvy-5db)  |
| cos    | ut-uptoc-1255970412 | ap-shanghai  |  [Regions](https://cloud.tencent.com/document/product/436/6224)  |
| qiniu  | ut-uptoc  | cn-east-1 |  [Regions](https://developer.qiniu.com/kodo/manual/4088/s3-access-domainname)  |
| google | ut-uptoc  | auto  | - |
| s3     | ut-uptoc  | ap-northeast-1  |  [Regions](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints)  |


## Contact us
- [Author Blog](https://saltbo.cn).

## Author
- [Saltbo](https://github.com/saltbo)

## License
- [MIT](https://github.com/saltbo/uptoc/blob/master/LICENSE)

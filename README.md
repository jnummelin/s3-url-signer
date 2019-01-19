[![Build Status](https://cloud.drone.io/api/badges/jnummelin/s3-url-signer/status.svg)](https://cloud.drone.io/jnummelin/s3-url-signer)

# s3-url-signer

Simple tool to pre-sign any S3 request

## Why?

The AWS official CLI client does not support creating other than `GET` URL as pre-signed.

## Usage

```
s3-url-signer -region eu-central-1 -verb PUT -bucket my-bucket -key foo
```

## Installing

Check [releases](releases/) for latest pre-build binaries.


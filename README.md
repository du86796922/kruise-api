# kruise-api

Schema of the API types that are served by Kruise.

## Purpose

This library is the canonical location of the Kruise API definition.

We recommend using the go types in this repo. You may serialize them directly to JSON.

## Where does it come from?

`kruise-api` is synced from [https://github.com/openkruise/kruise/tree/master/apis](https://github.com/openkruise/kruise/tree/master/apis).
Code changes are made in that location, merged into `openkruise/kruise` and later synced here.

## Things you should NOT do

[https://github.com/openkruise/kruise/tree/master/apis](https://github.com/openkruise/kruise/tree/master/apis) is synced to here.
All changes must be made in the former. The latter is read-only.

## generate new client
If you upgrade the client go version, please modify the corresponding go.mod file
generate client
```shell
make generate
```
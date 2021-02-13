# UserAPI

## How to

### requirements

```
go 1.15
direnv
gcp service account key
```

### Run

if using inmemorydb instead of firestore, you can skip this step

```sh
cp .envrc.sample .envrc
# create firestore and sa key from gcp console
# fill out .envrc
```

```
go run main.go
```

### Deploy

- enable cloud run and container registry
- dispatch deploy workflow from github actions

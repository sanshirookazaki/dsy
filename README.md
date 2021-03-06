# dsy

![test](https://github.com/sanshirookazaki/dsy/workflows/test/badge.svg)

dsy is a library for fixtures inspired by [timakin/dsmock](https://github.com/timakin/dsmock).

It is able to insert data in YAML formatted files to google cloud datastore.

## Install

```
go get -u github.com/sanshirookazaki/dsy
```

## Getting started

### YAML configuration

A fixture YAML file must contains the information about your datastore entity with the keys. For example:

```
kind: User
key: ID

entities:
  - ID: 1
    Name: Ryu
    Likes: Martial arts
  - ID: 2
    Name: Ken
    Likes: Pasta
```

### Upsert data

UpsertFile upserts a fixture file.

UpsertDir upserts fixture files found recorsively search for directory (.yaml or .yml extension).


```
ctx := context.Background()
client, _ := datastore.NewClient(ctx, "projectID")

err = dsy.UpsertFile(ctx, client, "path/to/fixture/file.yaml")
err = dsy.UpsertDir(ctx, client, "path/to/fixture")
```

For information on creating a client, see [cloud.google.com/go/datastore](https://godoc.org/cloud.google.com/go/datastore).

## YAML configuration details

If there is a property of key type, you must specify ```keys``` in YAML array.

In addition, set ```kind``` and ```id``` or ```name``` in entities. (if it doesn't, set kind and load value into id or name automaticaly.) From the top, it interpret as the parent.

If options include ```noIndex``` then the field will not be indexed. For example:

```
kind: Article
key: ID
keys:
  - Writer
  - Role
noIndex:
  - Title

entities:
  - ID: 1
    Writer:
      - kind: Company # parent
        name: ABC
      - kind: Writer
        id: 10
    Role: Manager # kind: Article, name: Manager
    Title: Business
```

This library supports a variety of data types for property values, as shown below.

```timeFormat``` is predefined layouts for use in time.Parse (default RFC3339), see [golang.org/pkg/time](https://golang.org/pkg/time/#pkg-constants)

```timeLocation``` is taken to be a location name corresponding to a file in the IANA Time Zone database (default UTC).

```
kind: Types
key: IntType
keys:
  - KeyType
timeFormat: RFC3339
timeLocation: UTC

entities:
  - IntType: 123
    FloatType: 1.23
    StringType: hello world
    BoolType: false
    NullType: null
    ArrayType: [1, 2, 3]
    EmbedType: {"Name": "Ken", "ID": 2}
    TimeType: 2006-01-02T15:04:05+07:00
    GeoPointType: {"Lat": 35.6809591, "Lng": 139.7673068}
    KeyType:
      - kind: Types
        name: key
```

## License
The MIT License (MIT)

Copyright 2020 Sanshiro Okazaki

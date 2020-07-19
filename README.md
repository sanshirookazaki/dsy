# dsy

dsy is a library for fixtures inspired by [timakin/dsmock](https://github.com/timakin/dsmock).

This is able to insert data in YAML formatted files to google cloud datastore.

## Install

```
go get -u github.com/sanshirookazaki/dsy
```

## Getting started

### YAML configuration

A fixture YAML file must contains the information about your datastore entity with the keys: schema and entities. For example:

```
scheme:
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

## Other configuring the fixture

If there is a property of key type, you must specify ```keys``` in scheme in YAML array.

In addition, set ```kind``` and ```id``` or ```name``` in entities. (if it doesn't, load value automaticaly.) From the top, it interpret as the parent.

If options include ```noIndex``` then the field will not be indexed. For example:

```
scheme:
  kind: Article
  key: ID
  keys:
    - Writer
    - Role
  noIndex:
    - Title

entities:
  ID: 1
  Writer:
    - kind: Company # parent
      name: ABC
    - kind: Writer
      id: 10
  Role: Manager # kind: Role, Name: Manager
  Title: Business
```

This library supports a variety of data types for property values, as shown below.

Datetime type is RFC 3339 formatted.

```
scheme:
  kind: Types
  key: ID
  keys:
    - Key

entities:
  ID: 1
  Int: 123
  Float: 1.23
  String: hello world
  Bool: false
  Null: null
  Array: [1, 2, 3]
  Key:
    - kind: Types
      name: key
  Embed: {"Name": "Ken", "ID": 2}
  Time: 2006-01-02T15:04:05+07:00
  GeoPoint: {"Lat": 35.6809591, "Lng": 139.7673068}
```

## License
The MIT License (MIT)

Copyright 2020 Sanshiro Okazaki

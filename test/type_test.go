package test

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
	"github.com/sanshirookazaki/dsy"
)

type Types struct {
	ID           *datastore.Key     `datastore:"__key__"`
	IntType      int                `datastore:"IntType"`
	FloatType    float64            `datastore:"FloatType"`
	StringType   string             `datastore:"StringType"`
	BoolType     bool               `datastore:"BoolType"`
	NullType     *int               `datastore:"NullType"`
	ArrayType    []int              `datastore:"ArrayType"`
	EmbedType    interface{}        `datastore:"EmbedType"`
	TimeType     time.Time          `datastore:"TimeType"`
	GeoPointType datastore.GeoPoint `datastore:"GeoPointType"`
	KeyType      *datastore.Key     `datastore:"KeyType"`
}

func TestType(t *testing.T) {
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	err = dsy.UpsertDir(ctx, client, "../testdata")
	if err != nil {
		t.Fatal(err)
	}

	k := datastore.IDKey("Type", 1, nil)
	ts := new(Types)
	if err := client.Get(ctx, k, ts); err != nil {
		t.Fatal(err)
	}

	pk := datastore.IDKey("Parent", 123456, nil)
	ck := datastore.NameKey("Child", "key", pk)
	tt, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	if err != nil {
		t.Fatal(err)
	}
	em := &datastore.Entity{
		Key: nil,
		Properties: []datastore.Property{
			datastore.Property{Name: "Integer", Value: int64(123456)},
			datastore.Property{Name: "String", Value: "value"},
		},
	}

	exts := &Types{
		ID:           k,
		IntType:      123456,
		FloatType:    0.123456,
		StringType:   "hello world",
		BoolType:     false,
		NullType:     nil,
		ArrayType:    []int{1, 2, 3},
		EmbedType:    em,
		TimeType:     tt,
		GeoPointType: datastore.GeoPoint{Lat: 35.6809591, Lng: 139.7673068},
		KeyType:      ck,
	}

	if !cmp.Equal(exts.EmbedType, ts.EmbedType) {
		t.Fatalf("testing: unexpected value of Types %v %v", exts, ts)
	}
}

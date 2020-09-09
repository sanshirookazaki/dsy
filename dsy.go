package dsy

import (
	"context"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"cloud.google.com/go/datastore"
)

const (
	BatchSize = 500
)

func UpsertFile(ctx context.Context, client *datastore.Client, filename string) (err error) {
	parser := NewParser()
	err = parser.ReadFile(filename)
	if err != nil {
		return err
	}

	entities, err := parser.Parse()
	if err != nil {
		return err
	}

	allPage := int(math.Ceil(float64(len(entities)) / float64(batchSize)))
	for page := 0; page < allPage; page++ {

		from := page * BatchSize
		to := (page + 1) * BatchSize
		if to > len(entities) {
			to = len(entities)
		}

		keys, values := getKeysValues(entities, from, to)
		if _, err = client.PutMulti(ctx, keys, values); err != nil {
			return err
		}
	}

	return nil
}

func UpsertDir(ctx context.Context, client *datastore.Client, dirname string) (err error) {
	datadir, err := filepath.Abs(dirname)
	if err != nil {
		return err
	}
	paths := dirwalk(datadir)

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		for _, file := range paths {
			parser := NewParser()
			err = parser.ReadFile(file)
			if err != nil {
				return err
			}

			entities, err := parser.Parse()
			if err != nil {
				return err
			}

			keys, values := getKeysValues(entities)

			if _, err = tx.PutMulti(keys, values); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}

		// .yaml or .yml
		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}

func getKeysValues(entities []*datastore.Entity, from, to int) (keys []*datastore.Key, values []interface{}) {
	for _, e := range entities[from:to] {
		keys = append(keys, e.Key)
		props := datastore.PropertyList(e.Properties)
		values = append(values, &props)
	}

	return keys, values
}

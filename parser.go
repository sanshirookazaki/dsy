package dsy

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/mitchellh/mapstructure"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Data struct {
	Kind     string   `yaml:"kind,omitempty"`
	Key      string   `yaml:"key,omitempty"`
	Keys     []string `yaml:"keys,omitempty"`
	NoIndex  []string `yaml:"noIndex,omitempty"`
	Entities []Entity `yaml:"entities,omitempty"`
}

type Entity map[string]interface{}

type Parser struct {
	Data *Data
}

func NewParser() *Parser {
	return &Parser{
		Data: &Data{},
	}
}

func (p *Parser) ReadFile(filename string) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	d := &Data{}
	if err = yaml.Unmarshal([]byte(source), d); err != nil {
		return err
	}
	p.Data = d
	return nil
}

func (p *Parser) Parse() (entities []*datastore.Entity, err error) {
	d := *p.Data

	for _, e := range d.Entities {
		entity, err := p.ParseEntity(e)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (p *Parser) ParseEntity(e Entity) (entity *datastore.Entity, err error) {
	var key *datastore.Key
	var properties []datastore.Property

	d := *p.Data

	for name, value := range e {
		if name == d.Key {
			// parse key
			k, err := p.ParseKeyList(value)
			if err != nil {
				return nil, err
			}
			key = k
		} else if funk.Contains(d.Keys, name) {
			// parse keytype property
			keyProp, err := p.ParseKeyList(value)
			if err != nil {
				return nil, err
			}

			var noIndex bool
			if funk.Contains(d.NoIndex, name) {
				noIndex = true
			}

			property := &datastore.Property{
				Name:    name,
				Value:   keyProp,
				NoIndex: noIndex,
			}
			properties = append(properties, *property)
		} else {
			// parse property
			property, err := p.ParseProperty(name, value)
			if err != nil {
				return nil, err
			}
			properties = append(properties, *property)
		}
	}

	entity = &datastore.Entity{
		Key:        key,
		Properties: properties,
	}
	return entity, nil
}

func (p *Parser) ParseProperty(name string, value interface{}) (property *datastore.Property, err error) {
	noIndexes := p.Data.NoIndex
	var noIndex bool
	if funk.Contains(noIndexes, name) {
		noIndex = true
	}

	// parse time.Time
	s, ok := value.(string)
	if ok {
		time, err := time.Parse(time.RFC3339, s)
		if err == nil {
			value = time
		}
	}

	if val, ok := value.(map[interface{}]interface{}); ok {
		// map[string]interface{}
		r := funk.Map(val, func(k interface{}, v interface{}) (string, interface{}) {
			return fmt.Sprintf("%s", k), v
		})

		// parse GeoPoint
		var geoType datastore.GeoPoint
		if satisfyFields(r.(map[string]interface{}), geoType) {
			err = mapstructure.Decode(r, &geoType)
			if err != nil {
				return nil, err
			}
			value = geoType
		} else {
			// parse embed entity
			props := make([]datastore.Property, 0)
			for name, v := range val {
				n, ok := name.(string)
				if !ok {
					return nil, fmt.Errorf("can not parse '%v' as embed property name.", name)
				}
				props = append(props, datastore.Property{
					Name:  n,
					Value: v,
				})
			}
			value = &datastore.Entity{
				Properties: props,
			}
		}
	}

	property = &datastore.Property{
		Name:    name,
		Value:   value,
		NoIndex: noIndex,
	}
	return property, nil
}

func (p *Parser) ParseKeyList(value interface{}) (key *datastore.Key, err error) {
	d := *p.Data

	keys, ok := value.([]interface{})
	if !ok {
		return p.ParseKey(d.Kind, value, nil)
	}

	var parent *datastore.Key
	for _, k := range keys {
		v, _ := k.(map[interface{}]interface{})
		if v["id"] != nil && v["name"] != nil {
			return nil, fmt.Errorf("either id or name: id=%v name=%v", v["id"], v["name"])
		} else if v["id"] == nil && v["name"] == nil {
			return nil, fmt.Errorf("id or name is required: %v", v["kind"])
		}
		var val interface{}

		if v["id"] != nil {
			val = v["id"]
		} else if v["name"] != nil {
			val = v["name"]
		}

		key, err = p.ParseKey(v["kind"], val, parent)
		if err != nil {
			return nil, err
		}
		parent = key
	}

	return key, nil
}

func (p *Parser) ParseKey(kind interface{}, value interface{}, parent *datastore.Key) (key *datastore.Key, err error) {
	k, ok := kind.(string)
	if !ok && kind != nil {
		return nil, fmt.Errorf("kind should be string: %v %v", reflect.TypeOf(kind), kind)
	}

	switch v := value.(type) {
	case string:
		key = datastore.NameKey(k, v, parent)
	case int64:
		key = datastore.IDKey(k, v, parent)
	case int:
		key = datastore.IDKey(k, int64(v), parent)
	case int32:
		key = datastore.IDKey(k, int64(v), parent)
	default:
		return nil, fmt.Errorf("key should be string or integer: %v %v", reflect.TypeOf(v), v)
	}

	return key, nil
}

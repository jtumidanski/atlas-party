package response

import (
	"atlas-party/json"
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

func NewDataContainer(singleResource bool) *DataContainer {
	return &DataContainer{SingleResource: singleResource, Data: make([]*DataBody, 0)}
}

type DataContainer struct {
	SingleResource bool
	Data           []*DataBody `json:"data"`
	Included       []*DataBody `json:"included"`
}

func (d *DataContainer) AddData(id uint32, theType string, attributes interface{}, relationships map[string]*Relationship) {
	d.Data = append(d.Data, &DataBody{
		Id:            strconv.Itoa(int(id)),
		Type:          theType,
		Attributes:    attributes,
		Relationships: relationships,
	})
}

func (d *DataContainer) AddIncluded(id uint32, theType string, attributes interface{}) {
	d.Included = append(d.Included, &DataBody{
		Id:            strconv.Itoa(int(id)),
		Type:          theType,
		Attributes:    attributes,
		Relationships: nil,
	})
}

func (d *DataContainer) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString("{")

	if d.SingleResource && len(d.Data) == 0 {
		err := writeAttribute(b, *d, "Data", nil)
		if err != nil {
			return nil, err
		}
	} else if d.SingleResource && len(d.Data) > 0 {
		err := writeAttribute(b, *d, "Data", d.Data[0])
		if err != nil {
			return nil, err
		}
	} else {
		err := writeAttribute(b, *d, "Data", d.Data)
		if err != nil {
			return nil, err
		}
	}
	if len(d.Included) > 0 {
		err := writeAttribute(b, *d, "Included", d.Included)
		if err != nil {
			return nil, err
		}
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

type DataBody struct {
	Id            string                   `json:"id"`
	Type          string                   `json:"type"`
	Attributes    interface{}              `json:"attributes"`
	Relationships map[string]*Relationship `json:"relationships"`
}

func writeAttribute(b *bytes.Buffer, root interface{}, name string, value interface{}) error {
	t := reflect.TypeOf(root)
	f, found := t.FieldByName(name)
	if found {
		if len(b.Bytes()) > 1 {
			b.WriteString(",")
		}
		b.WriteString("\"")
		b.WriteString(f.Tag.Get("json"))
		b.WriteString("\":")
		err := json.ToJSON(value, b)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("attribute not found")
}

func (d *DataBody) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString("{")
	err := writeAttribute(b, *d, "Type", d.Type)
	if err != nil {
		return nil, err
	}
	err = writeAttribute(b, *d, "Id", d.Id)
	if err != nil {
		return nil, err
	}
	err = writeAttribute(b, *d, "Attributes", d.Attributes)
	if err != nil {
		return nil, err
	}

	if len(d.Relationships) != 0 {
		err = writeAttribute(b, *d, "Relationships", d.Relationships)
		if err != nil {
			return nil, err
		}
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

type Relationship struct {
	ToOneType bool
	Links     RelationshipLinks  `json:"links"`
	Data      []RelationshipData `json:"data"`
}

func (r *Relationship) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString("{")
	err := writeAttribute(b, *r, "Links", r.Links)
	if err != nil {
		return nil, err
	}
	if r.ToOneType && len(r.Data) == 0 {
		err = writeAttribute(b, *r, "Data", nil)
		if err != nil {
			return nil, err
		}
	} else if r.ToOneType {
		err = writeAttribute(b, *r, "Data", r.Data[0])
		if err != nil {
			return nil, err
		}
	} else {
		err = writeAttribute(b, *r, "Data", r.Data)
		if err != nil {
			return nil, err
		}
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

type RelationshipLinks struct {
	Self    string `json:"self"`
	Related string `json:"related"`
}

type RelationshipData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}
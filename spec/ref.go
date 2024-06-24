package spec

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type (
	Reference struct {
		Ref         string `json:"$ref" yaml:"$ref"`
		Summary     string `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}

	ItemOrRef[T any] struct {
		Item T
		Ref  Reference
	}
)

func Ref(section string, key string) Reference {
	return Reference{
		Ref: fmt.Sprintf("#/components/%s/%s", section, key),
	}
}

func (r ItemOrRef[T]) MarshalJSON() ([]byte, error) {
	if r.Ref.Ref != "" {
		return json.Marshal(r.Ref)
	}

	return json.Marshal(r.Item)
}

func (r ItemOrRef[T]) MarshalYAML() (any, error) {
	var (
		res yaml.Node
		err error
	)

	if r.Ref.Ref != "" {
		err = res.Encode(r.Ref)
	} else {
		err = res.Encode(r.Item)
	}

	return res, err
}

func (r *ItemOrRef[T]) UnmarshalJSON(raw []byte) error {
	var err error

	if err = json.Unmarshal(raw, &r.Ref); err != nil || r.Ref.Ref == "" {
		err = json.Unmarshal(raw, &r.Item)
	}

	return err
}

func (r *ItemOrRef[T]) UnmarshalYAML(node *yaml.Node) error {
	var err error

	if err = node.Decode(&r.Ref); err != nil || r.Ref.Ref == "" {
		err = node.Decode(&r.Item)
	}

	return err
}

var (
	_ json.Marshaler   = ItemOrRef[Schema]{}
	_ yaml.Marshaler   = ItemOrRef[Schema]{}
	_ json.Unmarshaler = (*ItemOrRef[Schema])(nil)
	_ yaml.Unmarshaler = (*ItemOrRef[Schema])(nil)
)

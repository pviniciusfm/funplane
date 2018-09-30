package v1alpha1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	yaml2 "gopkg.in/yaml.v2"
	"reflect"
	"time"
)

// Make creates a new instance of the proto message
func Make(messageName string) (proto.Message, error) {
	pbt := proto.MessageType(messageName)
	if pbt == nil {
		return nil, fmt.Errorf("unknown type %q", messageName)
	}
	return reflect.New(pbt.Elem()).Interface().(proto.Message), nil
}

// FromJSONMap converts from a generic map to a proto message using canonical JSON encoding
// JSON encoding is specified here: https://developers.google.com/protocol-buffers/docs/proto3#json
func FromJSONMap(data interface{}, messageName string) (proto.Message, error) {
	// Marshal to YAML bytes
	str, err := yaml2.Marshal(data)
	if err != nil {
		return nil, err
	}

	out, err := FromYAML(string(str), messageName)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FromYAML recovers the original protobuff message from yaml file
func FromYAML(yml string, messageName string) (proto.Message, error) {
	pb, err := Make(messageName)
	if err != nil {
		return nil, err
	}
	if err = ApplyYAML(yml, pb); err != nil {
		return nil, err
	}
	return pb, nil
}

// ApplyYAML unmarshals a YAML string into a proto message
func ApplyYAML(yml string, pb proto.Message) error {
	js, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}

	return jsonpb.UnmarshalString(string(js), pb)
}

// ToJSON marshals a proto to canonical JSON
func ToJSON(msg proto.Message) (string, error) {
	return ToJSONWithIndent(msg, "")
}

// ToJSONWithIndent marshals a proto to canonical JSON with pretty printed string
func ToJSONWithIndent(msg proto.Message, indent string) (string, error) {
	if msg == nil {
		return "", errors.New("unexpected nil message")
	}

	// Marshal from proto to json bytes
	m := jsonpb.Marshaler{Indent: indent}
	return m.MarshalToString(msg)
}

func ToYAML(msg proto.Message) (string, error) {
	js, err := ToJSON(msg)
	if err != nil {
		return "", err
	}
	yml, err := yaml.JSONToYAML([]byte(js))
	return string(yml), err
}

// Duration is an abstract representation from time.Duration
// that is meant to support the serialization of time based strings
type Duration struct {
	time.Duration
}

// MarshalJSON represents Duration struct into JSON format
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON converts Duration serialized format into time.Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

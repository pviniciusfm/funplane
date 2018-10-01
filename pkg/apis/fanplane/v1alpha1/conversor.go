package v1alpha1

import (
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

//Transforms highlevel interface into concrete typeInstance
func ToStruct(spec interface{}, typeInstance interface{} ) (error) {
	rawSpec, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(rawSpec, typeInstance)
}

// ReadableDuration is an abstract representation from time.Duration
// that is meant to support the serialization of time based strings
// ReadableDuration is a duration type that is serialized to JSON in human readable format.
type ReadableDuration time.Duration

func NewReadableDuration(dur time.Duration) *ReadableDuration {
	d := ReadableDuration(dur)
	return &d
}

func (d *ReadableDuration) String() string {
	return d.Duration().String()
}

func (d *ReadableDuration) Duration() time.Duration {
	if d == nil {
		return time.Duration(0)
	}
	return time.Duration(*d)
}

func (d *ReadableDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Duration().String())), nil
}

func (d *ReadableDuration) UnmarshalJSON(raw []byte) error {
	if d == nil {
		return fmt.Errorf("cannot unmarshal to nil pointer")
	}

	str := string(raw)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return fmt.Errorf("must be enclosed with quotes: %s", str)
	}
	dur, err := time.ParseDuration(str[1 : len(str)-1])
	if err != nil {
		return err
	}
	*d = ReadableDuration(dur)
	return nil
}

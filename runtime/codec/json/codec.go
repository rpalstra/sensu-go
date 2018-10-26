package json

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/sensu/sensu-go/internal/apis/meta"
)

type JsonCodec struct{}

type surrogate struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata"`
	Spec            *json.RawMessage `json:"spec"`
}

func (JsonCodec) Decode(data []byte, objPtr interface{}) error {
	return jsoniter.Unmarshal(data, objPtr)
}

func (JsonCodec) Encode(objPtr interface{}) ([]byte, error) {
	return jsoniter.Marshal(objPtr)
}

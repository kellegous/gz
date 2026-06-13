package gz

import (
	"io"

	"github.com/kellegous/poop"
	"google.golang.org/protobuf/encoding/protojson"
)

var marshaler = protojson.MarshalOptions{
	UseProtoNames: true,
	Indent:        "  ",
}

func (b *Branch) WriteJSONTo(w io.Writer) error {
	json, err := marshaler.Marshal(b)
	if err != nil {
		return poop.Chain(err)
	}
	_, err = w.Write(json)
	return poop.Chain(err)
}

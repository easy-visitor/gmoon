package gmoon

import (
	"encoding/json"
	"fmt"
)

type Model interface {
	String() string
}
type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return Models(b)
}

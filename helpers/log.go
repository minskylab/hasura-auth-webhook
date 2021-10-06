package helpers

import (
	"encoding/json"
	"fmt"
)

func Log(i interface{}) {
	fmt.Printf("%+v\n", i)
	res, _ := json.Marshal(i)
	fmt.Println(string(res))
}

package utils

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type orderPreference struct {
	Order []string `json:"order"`
}

const prefFilename = "orderPref.yaml"

//OrderPref returns a slice containing the preferred ordering of Release names.
// If it doesn't find any defined, it returns an empty slice.
func OrderPref() (releaseOrderPref []string) {
	if _, err := os.Stat(prefFilename); !os.IsNotExist(err) {
		dat, err := ioutil.ReadFile(prefFilename)
		PanicCheck(err)
		var o orderPreference
		PanicCheck(yaml.Unmarshal(dat, &o))
		releaseOrderPref = o.Order
	}
	return
}

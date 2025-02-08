package appx

import (
	"github.com/qida/gohp/idx"
)

var SNOW_FLAK *idx.SnowFlakeJS

func init() {
	SNOW_FLAK = idx.NewSnowFlakeJS()
}

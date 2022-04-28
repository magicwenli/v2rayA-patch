package controller

import (
	"fmt"

	"github.com/magicwenli/v2rayA-patch/pkg/util/log"
)

func logError(err interface{}) error {
	e := fmt.Errorf("%v", err)
	log.Error("%v", e)
	return e
}

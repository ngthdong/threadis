package core

import "github.com/ngthdong/threadis/internal/datastructures"

var dictStore *datastructures.Dict

func init() {
	dictStore = datastructures.NewDict()
}

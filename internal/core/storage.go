package core

import "github.com/ngthdong/threadis/internal/datastructures"

var dictStore *datastructures.Dict
var zsetStore map[string]*datastructures.ZSet
var setStore map[string]*datastructures.Set

func init() {
	dictStore = datastructures.NewDict()
	zsetStore = make(map[string]*datastructures.ZSet)
	setStore = make(map[string]*datastructures.Set)
}

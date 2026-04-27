package datastructures_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ngthdong/threadis/internal/datastructures"
)

func TestSetGet(t *testing.T) {
	d := datastructures.NewDict()

	obj := d.NewObj("key", "value", 0)
	d.Set("key", obj)

	res := d.Get("key")
	assert.NotNil(t, res)
	assert.Equal(t, "value", res.Value)
}

func TestDelete(t *testing.T) {
	d := datastructures.NewDict()

	obj := d.NewObj("key", "value", 0)
	d.Set("key", obj)

	ok := d.Del("key")
	assert.True(t, ok)

	res := d.Get("key")
	assert.Nil(t, res)
}

func TestExpire(t *testing.T) {
	d := datastructures.NewDict()

	obj := d.NewObj("key", "value", 50) 
	d.Set("key", obj)

	time.Sleep(60 * time.Millisecond)

	res := d.Get("key")
	assert.Nil(t, res) 
}

func TestNoExpire(t *testing.T) {
	d := datastructures.NewDict()

	obj := d.NewObj("key", "value", 0)
	d.Set("key", obj)

	time.Sleep(50 * time.Millisecond)

	res := d.Get("key")
	assert.NotNil(t, res)
}

func TestOverwrite(t *testing.T) {
	d := datastructures.NewDict()

	obj1 := d.NewObj("key", "v1", 0)
	d.Set("key", obj1)

	obj2 := d.NewObj("key", "v2", 0)
	d.Set("key", obj2)

	res := d.Get("key")
	assert.Equal(t, "v2", res.Value)
}

func TestExpireReset(t *testing.T) {
	d := datastructures.NewDict()

	obj := d.NewObj("key", "value", 50)
	d.Set("key", obj)

	obj2 := d.NewObj("key", "value2", 0)
	d.Set("key", obj2)

	time.Sleep(60 * time.Millisecond)

	res := d.Get("key")

	assert.Nil(t, res)
}
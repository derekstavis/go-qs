package qs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	hash, err := Unmarshal("")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{})
	}

	hash, err = Unmarshal("foo")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": nil})
	}

	hash, err = Unmarshal("foo=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": ""})
	}

	hash, err = Unmarshal("foo=bar")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar"})
	}

	hash, err = Unmarshal("foo=\"bar\"")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "\"bar\""})
	}

	hash, err = Unmarshal("foo=bar&foo=quux")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "quux"})
	}

	hash, err = Unmarshal("foo&foo=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": ""})
	}

	hash, err = Unmarshal("foo=1&bar=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "1", "bar": "2"})
	}

	hash, err = Unmarshal("&foo=1&&bar=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "1", "bar": "2"})
	}

	hash, err = Unmarshal("foo&bar=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": nil, "bar": ""})
	}

	hash, err = Unmarshal("foo=bar&baz=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar", "baz": ""})
	}

	hash, err = Unmarshal("my+weird+field=q1%212%22%27w%245%267%2Fz8%29%3F")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"my weird field": `q1!2"'w$5&7/z8)?`})
	}

	hash, err = Unmarshal("a=b&pid%3D1234=1023")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"pid=1234": "1023", "a": "b"})
	}

	hash, err = Unmarshal("foo[]")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{nil}})
	}

	hash, err = Unmarshal("foo[]=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{""}})
	}

	hash, err = Unmarshal("foo[]=bar")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"bar"}})
	}

	hash, err = Unmarshal("foo[]=1&foo[]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"1", "2"}})
	}

	hash, err = Unmarshal("foo=bar&baz[]=1&baz[]=2&baz[]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar", "baz": []interface{}{"1", "2", "3"}})
	}

	hash, err = Unmarshal("foo[]=bar&baz[]=1&baz[]=2&baz[]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"bar"}, "baz": []interface{}{"1", "2", "3"}})
	}

	hash, err = Unmarshal("x[y][z]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": "1"}}})
	}

	hash, err = Unmarshal("x[y][z][]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1"}}}})
	}

	hash, err = Unmarshal("x[y][z]=1&x[y][z]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": "2"}}})
	}

	hash, err = Unmarshal("x[y][z][]=1&x[y][z][]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1", "2"}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}}}})
	}

	hash, err = Unmarshal("x[y][][z][]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": []interface{}{"1"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][w]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "2"}}}})
	}

	hash, err = Unmarshal("x[y][][v][w]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"v": map[string]interface{}{"w": "1"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][v][w]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "v": map[string]interface{}{"w": "2"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][z]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}, map[string]interface{}{"z": "2"}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][w]=a&x[y][][z]=2&x[y][][w]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "a"}, map[string]interface{}{"z": "2", "w": "3"}}}})
	}

	hash, err = Unmarshal("foo=bar&baz=qwy?djc=dk&baz1=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar", "baz": "qwy?djc=dk", "baz1": "3"})
	}

	hash, err = Unmarshal("x[y]=1&x[y]z=2")
	assert.Error(t, err)

	hash, err = Unmarshal("x[y]=1&x[]=1")
	assert.Error(t, err)

	hash, err = Unmarshal("x[y]=1&x[y][][w]=2")
	assert.Error(t, err)

}

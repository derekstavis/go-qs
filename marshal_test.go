package qs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	testables := []map[string]interface{}{
		map[string]interface{}{"foo": "bar"},
		map[string]interface{}{"foo": "1", "bar": "2"},
		map[string]interface{}{"my weird field": "q1!2\"'w$5&7/z8)?"},
		map[string]interface{}{"foo": []interface{}{nil}},
		map[string]interface{}{"foo": []interface{}{""}},
		map[string]interface{}{"foo": []interface{}{"bar"}},
		map[string]interface{}{"foo": nil, "bar": ""},
		map[string]interface{}{"foo": "bar", "baz": ""},
		map[string]interface{}{"foo": []interface{}{"1", "2"}},
		map[string]interface{}{"foo": "bar", "baz": []interface{}{"1", "2", "3"}},
		map[string]interface{}{"foo": []interface{}{"bar"}, "baz": []interface{}{"1", "2", "3"}},
		map[string]interface{}{"foo": []interface{}{"1", "2"}},
		map[string]interface{}{"foo": "bar", "baz": []interface{}{"1", "2", "3"}},
		map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": "1"}}},
		map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1"}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1", "2"}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": []interface{}{"1"}}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "2"}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"v": map[string]interface{}{"w": "1"}}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "v": map[string]interface{}{"w": "2"}}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}, map[string]interface{}{"z": "2"}}}},
		map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "a"}, map[string]interface{}{"z": "2", "w": "3"}}}},
	}

	for _, v := range testables {
		querystring, err := Marshal(v)

		if assert.NoError(t, err) {
			hash, err := Unmarshal(querystring)

			if assert.NoError(t, err) {
				assert.Equal(t, hash, v)
			}
		}
	}
}

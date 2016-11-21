package qs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

var nameRegex = pcre.MustCompile(`\A[\[\]]*([^\[\]]+)\]*`, 0)
var objectRegex1 = pcre.MustCompile(`^\[\]\[([^\[\]]+)\]$`, 0)
var objectRegex2 = pcre.MustCompile(`^\[\](.+)$`, 0)

func Unmarshal(qs string) (map[string]interface{}, error) {
	components := strings.Split(qs, "&")
	params := map[string]interface{}{}

	for _, c := range components {
		tuple := strings.Split(c, "=")

		for i, item := range tuple {
			if unesc, err := url.QueryUnescape(item); err == nil {
				tuple[i] = unesc
			}
		}

		key := ""

		if len(tuple) > 0 {
			key = tuple[0]
		}

		value := interface{}(nil)

		if len(tuple) > 1 {
			value = tuple[1]
		}

		if err := normalizeParams(params, key, value); err != nil {
			return nil, err
		}
	}

	return params, nil
}

func normalizeParams(params map[string]interface{}, key string, value interface{}) error {
	nameMatcher := nameRegex.MatcherString(key, 0)
	k := nameMatcher.GroupString(1)
	after := ""

	if pos := nameRegex.FindIndex([]byte(key), 0); len(pos) == 2 {
		after = key[pos[1]:]
	}

	objectMatcher1 := objectRegex1.MatcherString(after, 0)
	objectMatcher2 := objectRegex2.MatcherString(after, 0)

	if k == "" {
		return nil

	} else if after == "" {
		params[k] = value
		return nil

	} else if after == "[]" {
		ival, ok := params[k]

		if !ok {
			params[k] = []interface{}{value}
			return nil
		}

		array, ok := ival.([]interface{})

		if !ok {
			return fmt.Errorf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival)
		}

		params[k] = append(array, value)
		return nil

	} else if objectMatcher1.Matches() || objectMatcher2.Matches() {

		childKey := ""

		if objectMatcher1.Matches() {
			childKey = objectMatcher1.GroupString(1)
		} else if objectMatcher2.Matches() {
			childKey = objectMatcher2.GroupString(1)
		}

		if childKey != "" {
			ival, ok := params[k]

			if !ok {
				params[k] = []interface{}{}
				ival = params[k]
			}

			array, ok := ival.([]interface{})

			if !ok {
				return fmt.Errorf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival)
			}

			if length := len(array); length > 0 {
				if hash, ok := array[length-1].(map[string]interface{}); ok {
					if _, ok := hash[childKey]; !ok {
						normalizeParams(hash, childKey, value)
						return nil
					}
				}
			}

			newHash := map[string]interface{}{}
			normalizeParams(newHash, childKey, value)
			params[k] = append(array, newHash)

			return nil
		}
	}

	ival, ok := params[k]

	if !ok {
		params[k] = map[string]interface{}{}
		ival = params[k]
	}

	hash, ok := ival.(map[string]interface{})

	if !ok {
		return fmt.Errorf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival)
	}

	if err := normalizeParams(hash, after, value); err != nil {
		return err
	}

	return nil
}

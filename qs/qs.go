package qs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

var nameRegex pcre.Regexp = pcre.MustCompile(`\A[\[\]]*([^\[\]]+)\]*`, 0)
var objectRegex1 pcre.Regexp = pcre.MustCompile(`^\[\]\[([^\[\]]+)\]$`, 0)
var objectRegex2 pcre.Regexp = pcre.MustCompile(`^\[\](.+)$`, 0)

func Parse(qs string) map[string]interface{} {
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

		value := ""

		if len(tuple) > 1 {
			value = tuple[1]
		}

		normalizeParams(params, key, value)
	}

	return params
}

func normalizeParams(params map[string]interface{}, key string, value interface{}) {
	nameMatcher := nameRegex.MatcherString(key, 0)
	k := nameMatcher.GroupString(1)
	after := ""

	if pos := nameRegex.FindIndex([]byte(key), 0); len(pos) == 2 {
		after = key[pos[1]:]
	}

	objectMatcher1 := objectRegex1.MatcherString(after, 0)
	objectMatcher2 := objectRegex2.MatcherString(after, 0)

	fmt.Printf("key: %s, after: %s, value: %s\n", k, after, value)

	if k == "" {
		return
	} else if after == "" {
		params[k] = value
		return
	} else if after == "[]" {
		ival, ok := params[k]

		if !ok {
			params[k] = []interface{}{value}
			return
		}

		array, ok := ival.([]interface{})

		if !ok {
			panic(fmt.Sprintf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival))
		}

		params[k] = append(array, value)
		return
	} else if objectMatcher1.Matches() || objectMatcher2.Matches() {

		childKey := ""

		if objectMatcher1.Matches() {
			childKey = objectMatcher1.GroupString(1)
		}

		if objectMatcher2.Matches() {
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
				panic(fmt.Sprintf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival))
			}

			if length := len(array); length > 0 {
				if hash, ok := array[length-1].(map[string]interface{}); ok {
					if ival, ok := hash[childKey]; ok {
						if childHash, ok := ival.(map[string]interface{}); ok {
							normalizeParams(childHash, childKey, value)
							return
						}
					}
				}
			}

			fmt.Println(childKey)

			newHash := map[string]interface{}{}
			normalizeParams(newHash, childKey, value)
			params[k] = append(array, newHash)

			return
		}
	}

	ival, ok := params[k]

	if !ok {
		params[k] = map[string]interface{}{}
	}

	hash, ok := ival.(map[string]interface{})

	if !ok {
		panic(fmt.Sprintf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival))
	}

	fmt.Println("kldjhfkjshdfkjsdhf")
	normalizeParams(hash, after, value)

}

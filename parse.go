package qs

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var enableConvertArrays = false

func ConvertArrays(enable bool) {
	enableConvertArrays = enable
}

func getPos(component string) int {
	pos := strings.Index(component, "]=")
	if pos == -1 {
		pos = strings.Index(component, "=")
	} else {
		pos++
	}
	return pos
}

func splitKeyValue(component string) (key string, value string, err error) {
	pos := getPos(component)
	if pos == -1 {
		key = component
		value = ""
	} else {
		key = component[0:pos]
		value = component[pos+1:]
	}

	key, err = url.QueryUnescape(key)
	if err != nil {
		return "", "", err
	}

	value, err = url.QueryUnescape(value)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func Parse(qs string) (interface{}, error) {
	components := strings.Split(qs, "&")

	params := map[string]interface{}{}

	for _, c := range components {

		key, value, err := splitKeyValue(c)
		if err != nil {
			return nil, err
		}

		if err := normalizeParams2(params, key, value); err != nil {
			return nil, err
		}
	}

	if !enableConvertArrays {
		return params, nil
	}

	return convertArrays(params), nil
}

func normalizeParams2(params map[string]interface{}, key string, value interface{}) error {

	nameMatcher := nameRegex.MatcherString(key, 0)
	k := nameMatcher.GroupString(1)
	after := ""

	pos := nameRegex.FindIndex([]byte(key), 0)
	if len(pos) == 2 {
		after = key[pos[1]:]
	}

	objectMatcher1 := objectRegex1.MatcherString(after, 0)
	objectMatcher2 := objectRegex2.MatcherString(after, 0)

	/*
		fmt.Printf("\n>>>>>>>----------------------------\n")
		pretty.Println("params:", params)
		pretty.Println("key:", key)
		pretty.Println("k:", k)
		pretty.Println("value:", value)
		pretty.Println("pos:", pos)
		pretty.Println("after:", after)
		pretty.Println("obj1", objectMatcher1.GroupString(1))
		pretty.Println("obj2", objectMatcher2.GroupString(1))
		fmt.Printf("<<<<<<<----------------------------\n\n")
	*/
	if k == "" {
		params = map[string]interface{}{}
		return nil

	} else if after == "" {

		ival, ok := params[k]
		if !ok {
			// key not exist yet, create it
			params[k] = value
			return nil
		}

		// key exists
		switch i := ival.(type) {
		case []interface{}:
			params[k] = append(i, value)
		case string:
			params[k] = []interface{}{i, value}
		case map[string]interface{}:
			array, ok := toArray(i)
			if ok {
				array = append(array, value)
				params[k] = array
			} else {
				str, ok := value.(string)
				if ok {
					i[str] = ""
					params[k] = i
				} else {
					return fmt.Errorf("Unexpected value type: %T", value)
				}
			}
		default:
			return fmt.Errorf("Unexpected type: %T", ival)
		}

		return nil

	} else if after == "[]" {

		ival, ok := params[k]
		if !ok {
			// key not exist yet, create it
			params[k] = []interface{}{value}
			return nil
		}

		// key exists
		switch i := ival.(type) {
		case []interface{}:
			params[k] = append(i, value)
		case string:
			params[k] = []interface{}{i, value}
		case map[string]interface{}:
			minorIndex := getNextIndex(i)
			i[strconv.Itoa(minorIndex)] = value
			params[k] = i
		default:
			return fmt.Errorf("Unexpected type: %T", ival)
		}

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

	switch i := ival.(type) {
	case map[string]interface{}:
		return normalizeParams(i, after, value)
	case string:
		// string -> array
		params[k] = []interface{}{i, value}
		return nil
	case []interface{}:
		// array -> object
		newKey := after[1 : len(after)-1]
		hash := toHash(i)
		hash[newKey] = value
		params[k] = hash
		return nil
	default:
		return fmt.Errorf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival)
	}
}

func indexArray(in map[string]interface{}) ([]int, bool) {

	if len(in) == 0 {
		return nil, false
	}

	arr := []int{}
	for key := range in {
		i, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return nil, false
		}
		arr = append(arr, int(i))
	}
	sort.Ints(arr)
	return arr, true
}

func toHash(arr []interface{}) map[string]interface{} {

	hash := map[string]interface{}{}
	for i, value := range arr {
		hash[strconv.Itoa(i)] = value
	}
	return hash
}

func getNextIndex(hash map[string]interface{}) int {

	arr := []int{}
	for key := range hash {

		// filter only integers
		n, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			continue
		}

		arr = append(arr, int(n))
	}

	if len(arr) == 0 {
		return 0
	}

	sort.Ints(arr)
	return arr[len(arr)-1] + 1
}

func toArray(in map[string]interface{}) ([]interface{}, bool) {

	indexArr, ok := indexArray(in)
	if !ok {
		return nil, false
	}

	arr := []interface{}{}
	for _, index := range indexArr {
		key := strconv.Itoa(index)
		arr = append(arr, in[key])
	}

	return arr, true
}

func convertArrays(in map[string]interface{}) interface{} {

	arr, ok := toArray(in)
	if ok { // array
		for i, value := range arr {
			switch v := value.(type) {
			case map[string]interface{}:
				arr[i] = convertArrays(v)
			}
		}
		return arr
	}

	// map
	for key, value := range in {
		switch v := value.(type) {
		case map[string]interface{}:
			in[key] = convertArrays(v)
		}
	}
	return in
}

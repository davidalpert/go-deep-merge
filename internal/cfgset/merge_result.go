package cfgset

import (
	"fmt"
	"sort"
	"strings"
)

type MergeResult struct {
	AppDir      string
	MergeBySlug map[string]map[string]interface{}
}

var (
	DefaultPathSeparator  = "/"
	DefaultValueSeparator = ":"
)

func (r *MergeResult) FlattenToMap() map[string]string {
	return r.FlattenToMapWithSep(DefaultPathSeparator)
}

func (r *MergeResult) FlattenToMapWithSep(sep string) map[string]string {
	result := make(map[string]string)
	for slug, merge := range r.MergeBySlug {
		for k, v := range RecursiveFlattenToMapWithSep(r.AppDir+sep+slug, merge, sep) {
			result[k] = v
		}
	}
	return result
}

func FlattenedToString(flattened map[string]string) string {
	strResult := make([]string, 0)
	for k, v := range flattened {
		strResult = append(strResult, k+DefaultValueSeparator+" "+v)
	}
	sort.Strings(strResult)
	return strings.Join(strResult, "\n") + "\n"
}

func RecursiveFlattenToMapWithSep(prefix string, v interface{}, sep string) map[string]string {
	result := make(map[string]string)
	switch vv := v.(type) {
	case []interface{}:
		for i, vvv := range vv {
			for ik, iv := range RecursiveFlattenToMapWithSep(prefix+sep+fmt.Sprintf("%d", i), vvv, sep) {
				result[ik] = iv
			}
		}
	case map[string]interface{}:
		for k, v := range vv {
			for ik, iv := range RecursiveFlattenToMapWithSep(prefix+sep+k, v, sep) {
				result[ik] = iv
			}
		}
	default:
		result[prefix] = fmt.Sprintf("%v", vv)
	}
	return result
}

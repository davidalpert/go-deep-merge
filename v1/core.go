package v1

import (
	"fmt"
	"sort"
	"strings"
)

// Merge deep merges the src map into dest map with the default options and returns a new map of merged values
func Merge(src, dest map[string]interface{}) (map[string]interface{}, error) {
	return MergeWithOptions(src, dest, NewConfig())
}

// MergeWithOptions deep merges the src map into dest map with the given options and returns a new map of merged values
func MergeWithOptions(src, dest map[string]interface{}, options *Config) (map[string]interface{}, error) {
	result, err := deepMerge(src, dest, options)
	if err != nil {
		return nil, err
	}

	switch v := result.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		return nil, fmt.Errorf("unexpected result %T from deepMerge", v)
	}
}

// deepMerge is a recursive function ported from the ruby deep_merge library
func deepMerge(src, dest interface{}, o *Config) (interface{}, error) {
	overwriteUnmergeable := !o.PreserveUnmergeables
	arraySplitChar := ""
	if o.UnpackArrays != nil {
		arraySplitChar = *o.UnpackArrays
	}

	if o.KnockoutPrefix != nil && *o.KnockoutPrefix == "" {
		return nil, fmt.Errorf("o.KnockoutPrefix cannot be an empty string")
	}

	if o.KnockoutPrefix != nil && !overwriteUnmergeable {
		return nil, fmt.Errorf("o.PreserveUnmergeable must be true when o.KnockoutPrefix is specified")
	}

	//o.writeDebug("Source: %T :: Dest: %T", src, dest)

	if !o.MergeNilValues && src == nil {
		return dest, nil
	}
	if dest == nil && !o.PreserveUnmergeables {
		return src, nil
	}

	o.writeDebug("%#v", o)
	switch s := src.(type) {
	case map[string]interface{}:
		o.writeDebug("Hashes: %#v :: %#v", s, dest)
		switch d := dest.(type) {
		case map[string]interface{}:
			for sk, sv := range s {
				if _, ok := d[sk]; ok {
					o.writeDebug(" ==>merging: %#v => %#v :: %#v", sk, sv, d)
					r, err := deepMerge(sv, d[sk], o.copyWithIncreasedDebugIndent())
					if err != nil {
						return nil, err
					}
					d[sk] = r
				} else {
					o.writeDebug(" ==>copying over: %#v => %#v :: %#v", sk, sv, d)
					// dest[src_key] doesn't exist so we want to create and overwrite it (but we do this via deep_merge!)
					// note: we rescue here b/c some classes respond to "dup" but don't implement it (Numeric, TrueClass, FalseClass, NilClass among maybe others)
					// we dup src_value if possible because we're going to merge into it (since dest is empty)
					// TODO: duplicate values safely
					switch src_dup := sv.(type) {
					case []interface{}:
						if o.KeepArrayDuplicates {
							// note: in this case the merge will be additive, rather than a bounded set, so we can't simply merge src with itself
							// We need to merge src with an empty array
							if r, err := deepMerge(sv, make([]interface{}, 0), o.copyWithIncreasedDebugIndent()); err != nil {
								return nil, err
							} else {
								d[sk] = r
							}
						} else {
							r, err := deepMerge(sv, src_dup, o.copyWithIncreasedDebugIndent())
							if err != nil {
								return nil, err
							}
							d[sk] = r
						}
					default:
						r, err := deepMerge(sv, src_dup, o.copyWithIncreasedDebugIndent())
						if err != nil {
							return nil, err
						}
						d[sk] = r
					}
				}
			}
			return d, nil
		case []interface{}: // elseif kind_of?(Array)
			if o.ExtendExistingArrays {
				d = append(d, s)
				return d, nil
			}
			if !o.PreserveUnmergeables {
				return overwriteUnmergeables(s, d, o)
			}
			return d, nil
		default: // else dest isn't a hash, so we overwrite it completely (if permitted)
			if !o.PreserveUnmergeables {
				return overwriteUnmergeables(s, d, o)
			}
			// not allowed to overwrite an unmergeable so we keep the destination
			return dest, nil
		}
	case []interface{}:
		o.writeDebug("Arrays: %#v :: %#v", s, dest)
		if o.OverwriteArrays {
			o.writeDebug("> overwrite arrays")
			return src, nil
		} else {
			sourceAllStrings := sliceOfAll(s, isString)
			var hasNakedKnockoutPrefix = false
			if sourceAllStrings && arraySplitChar != "" && len(s) > 0 {
				ss := toStringSlice(s)
				o.writeDebug("split/join on source: %#v", src)
				ss = strings.Split(strings.Join(ss, arraySplitChar), arraySplitChar)
				for _, v := range ss {
					if o.KnockoutPrefix != nil && v == *o.KnockoutPrefix {
						hasNakedKnockoutPrefix = true
					}
				}
				if hasNakedKnockoutPrefix {
					ss = sliceWithout(ss, *o.KnockoutPrefix)
				}
				s = toAbstractSlice(ss)
			}
			switch d := dest.(type) {
			case []interface{}:
				// if there's a naked knockout_prefix in source, that means we are to truncate dest
				if o.KnockoutPrefix != nil && hasNakedKnockoutPrefix {
					d = make([]interface{}, 0)
					o.writeDebug("source has naked knockout prefix; truncating destination to: %#v", d)
				}

				destAllStrings := sliceOfAll(s, isString)
				if destAllStrings && arraySplitChar != "" && len(d) > 0 {
					dd := toStringSlice(d)
					o.writeDebug("split/join on dest: %#v", src)
					dd = strings.Split(strings.Join(dd, arraySplitChar), arraySplitChar)
					d = toAbstractSlice(dd)
				}
				if o.KnockoutPrefix != nil {
					// remove knockout prefix items from both source and dest
					knockoutPrefix := *o.KnockoutPrefix
					sWithoutKnockoutItems := make([]interface{}, 0)
					for _, koItem := range s {
						isKnockoutItem := false
						switch item := koItem.(type) {
						case string:
							if strings.HasPrefix(item, knockoutPrefix) {
								isKnockoutItem = true
								itemToKnockout := strings.TrimPrefix(item, knockoutPrefix)
								o.writeDebug("found %#v ==> knocking out: %#v", koItem, itemToKnockout)
								dd := toStringSlice(d)
								dd = sliceWithout(dd, item)
								dd = sliceWithout(dd, itemToKnockout)
								d = toAbstractSlice(dd)
							}
						}
						if !isKnockoutItem {
							sWithoutKnockoutItems = append(sWithoutKnockoutItems, koItem)
						}
					}
					s = sWithoutKnockoutItems
				}

				sourceAllHashes := allHashes(s)
				destAllHashes := allHashes(d)

				if o.MergeHashArrays && sourceAllHashes && destAllHashes {
					o.writeDebug("merge hashes in lists")
					list := make([]interface{}, 0)
					for i, dv := range d {
						if i < len(s) {
							sv := s[i]
							o.writeDebug("- index %d: %#v :: %#v", i, sv, dv)
							if r, err := deepMerge(sv, dv, o.copyWithIncreasedDebugIndent()); err != nil {
								return nil, err
							} else {
								list = append(list, r)
							}
						} else {
							// no source item at this index so preserve dest item
							list = append(list, dv)
						}
					}
					if len(s) > len(d) {
						// source has more items than dest; append the rest
						list = append(list, s[len(d):]...)
					}
					d = list
				} else if o.KeepArrayDuplicates {
					d = append(d, s...)
				} else {
					d = combineWithoutDuplicates(s, d, o)
				}
				if o.SortMergedArrays {
					o.writeDebug("HERE")
					if len(d) > 0 {
						switch d[0].(type) {
						case int:
							sort.Slice(d, func(i, j int) bool { return d[i].(int) < d[j].(int) })
						case int32:
							sort.Slice(d, func(i, j int) bool { return d[i].(int32) < d[j].(int32) })
						case int64:
							sort.Slice(d, func(i, j int) bool { return d[i].(int64) < d[j].(int64) })
						case float32:
							sort.Slice(d, func(i, j int) bool { return d[i].(float32) < d[j].(float32) })
						case float64:
							sort.Slice(d, func(i, j int) bool { return d[i].(float64) < d[j].(float64) })
						case string:
							sort.Slice(d, func(i, j int) bool { return strings.Compare(d[i].(string), d[j].(string)) < 0 })
						default:
							sort.Slice(d, func(i, j int) bool {
								return strings.Compare(fmt.Sprintf("%#v", d[i]), fmt.Sprintf("%#v", d[j])) < 0
							})
						}
					}
				}
				return d, nil
			default:
				return overwriteUnmergeables(s, d, o)
			}
		}
	default:
		o.writeDebug("Source is neither map nor slice; overwriting dest: %#v :: %#v", src, dest)
		switch d := dest.(type) {
		case []interface{}:
			if o.ExtendExistingArrays {
				d = append(d, s)
				return d, nil
			}
			return overwriteUnmergeables(s, d, o)
		default:
			return overwriteUnmergeables(s, d, o)
		}
	}
}

// [3, 4, 5] ==> [1, 2, 3] = [1, 2, 3, 4, 5]
func combineWithoutDuplicates(s []interface{}, d []interface{}, o *Config) []interface{} {
	for _, v := range s {
		if indexOf(d, v) < 0 {
			d = append(d, v)
		}
	}

	return d
}

func sliceOfAll(items []interface{}, testItem func(v interface{}) bool) bool {
	var trueForAll = true
	for _, v := range items {
		trueForAll = testItem(v) && trueForAll
	}
	return trueForAll
}

func isString(item interface{}) bool {
	switch item.(type) {
	case string:
		return true
	default:
		return false
	}
}

// allHashes returns true when all elements of the items slice are map[string]interface{}
func allHashes(items []interface{}) bool {
	return sliceOfAll(items, isHash)
}

func isHash(item interface{}) bool {
	switch item.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

func sliceWithout(ss []string, s string) []string {
	foundIndex := indexOfString(ss, s)
	if foundIndex < 0 {
		// item not found
		return ss
	}

	var result = make([]string, 0)
	for i, v := range ss {
		if i != foundIndex {
			result = append(result, v)
		}
	}

	return result
}

func indexOfString(ss []string, s string) int {
	for i, v := range ss {
		if strings.EqualFold(v, s) {
			return i
		}
	}
	return -1
}

func indexOf(ss []interface{}, s interface{}) int {
	// compare detailed string representations
	strS := fmt.Sprintf("%#v", s)
	for i, v := range ss {
		if fmt.Sprintf("%#v", v) == strS {
			return i
		}
	}
	return -1
}

// overwriteUnmergeables returns the result of writing src over top of dest using the configured options
func overwriteUnmergeables(src, dest interface{}, o *Config) (interface{}, error) {
	o.writeDebug("Overwrite: %#v :: %#v", src, dest)
	overwriteUnmergeable := !o.PreserveUnmergeables
	if o.KnockoutPrefix != nil && overwriteUnmergeable {
		knockoutPrefix := *o.KnockoutPrefix
		var srcTemp interface{}

		// use this flag instead of srcTemp to avoid comparing various types
		knockoutChangedSource := false

		// apply knockout prefix to source before overwriting dest
		switch s := src.(type) {
		case string:
			o.writeDebug("remove knockout string from source %#v before overwriting dest", s)
			srcTemp = strings.TrimLeft(s, knockoutPrefix)
			knockoutChangedSource = srcTemp != s
		case []interface{}:
			// remove all knockout elements before overwriting dest
			t := make([]interface{}, 0)
			for _, v := range s {
				switch vv := v.(type) {
				case string:
					if !strings.HasPrefix(vv, knockoutPrefix) {
						t = append(t, vv)
					}
					//knockoutChangedSource = knockoutChangedSource && (t[i] != vv)
				}
			}
			srcTemp = t
		default:
			srcTemp = s
		}

		o.writeDebug("comparing srcTemp %#v with original src %#v", srcTemp, src)
		if !knockoutChangedSource {
			// we didn't find a KnockoutPrefix so simply overwrite dest
			o.writeDebug("%#v -over -> %#v", src, dest)
			return srcTemp, nil
		} else {
			// found a KnockoutPrefix so delete dest
			o.writeDebug("knocking out %#v -over -> %#v", src, dest)
			return "", nil
		}
	} else if !o.PreserveUnmergeables {
		return src, nil
	}
	return dest, nil
}

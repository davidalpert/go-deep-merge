package DeepMerge

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
					o.writeDebug(" ==>merging over: %#v => %#v :: %#v", sk, sv, d)
					// dest[src_key] doesn't exist so we want to create and overwrite it (but we do this via deep_merge!)
					// note: we rescue here b/c some classes respond to "dup" but don't implement it (Numeric, TrueClass, FalseClass, NilClass among maybe others)
					// we dup src_value if possible because we're going to merge into it (since dest is empty)
					// TODO: duplicate values safely
					switch src_dup := sv.(type) {
					case []interface{}:
						if o.KeepArrayDuplicates {
							// note: in this case the merge will be additive, rather than a bounded set, so we can't simply merge src with itself
							// We need to merge src with an empty array
							return deepMerge(sv, make([]interface{}, 0), o.copyWithIncreasedDebugIndent())
						}
						r, err := deepMerge(sv, src_dup, o.copyWithIncreasedDebugIndent())
						if err != nil {
							return nil, err
						}
						d[sk] = r
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
			return d, nil
		default: // else dest isn't a hash, so we overwrite it completely (if permitted)
			if !o.PreserveUnmergeables {
				return overwriteUnmergeables(s, d, o)
			}
			// not allowed to overwrite an unmergeable so we keep the destination
			return dest, nil
		}
	case []string:
		o.writeDebug("Array<string>: %#v :: %#v", s, dest)
		if o.OverwriteArrays {
			o.writeDebug("> overwrite arrays")
			return src, nil
		} else {
			// if we are instructed, join/split any source arrays before processing
			if arraySplitChar != "" {
				o.writeDebug("split/join on source: %#v", src)
				s = strings.Split(strings.Join(s, arraySplitChar), arraySplitChar)
			}
			switch d := dest.(type) {
			case []string:
				// if we are instructed, join/split any dest arrays before processing
				if arraySplitChar != "" {
					o.writeDebug("split/join on dest: %#v", src)
					d = strings.Split(strings.Join(d, arraySplitChar), arraySplitChar)
				}
				if o.KnockoutPrefix != nil {
					// if there's a naked knockout_prefix in source, that means we are to truncate dest
					if indexOfString(s, *o.KnockoutPrefix) > -1 {
						return sliceWithout(s, *o.KnockoutPrefix), nil
					}
					// TODO: merge arrays
				}
			}
			return src, nil
		}
	case []interface{}:
		o.writeDebug("Arrays: %#v :: %#v", s, dest)
		if o.OverwriteArrays {
			o.writeDebug("> overwrite arrays")
			return src, nil
		} else {
			if o.KnockoutPrefix != nil {
				// treat as []string => []string ???
			}
			switch d := dest.(type) {
			case []interface{}:
				sourceAllHashes := allHashes(s)
				destAllHashes := allHashes(d)
				if sourceAllHashes && destAllHashes && o.MergeHashArrays {
					list := make([]interface{}, 0)
					for i, dv := range d {
						if i < len(s) {
							sv := s[i]
							if r, err := deepMerge(sv, dv, o.copyWithIncreasedDebugIndent()); err != nil {
								return nil, err
							} else {
								list = append(list, r)
							}
						}
					}
					if len(s) < len(d) {
						list = append(list, s[len(d)])
					}
					d = list
				} else if o.KeepArrayDuplicates {
					d = append(d, s...)
				} else {
					d = combineWithoutDuplicates(s, d, o)
				}
				if o.SortMergedArrays {
					sort.Slice(d, func(i, j int) bool {
						// TODO: how do we sort a slice of interfaces? cast each one?
						//return d[i] < d[j]
						return i < j
					})
				}
				return d, nil
			}
			return src, nil
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

// allHashes returns true when all elements of the items slice are map[string]interface{}
func allHashes(items []interface{}) bool {
	var trueForAll = true
	for _, v := range items {
		trueForAll = isHash(v) && trueForAll
	}
	return trueForAll
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
	for i, v := range ss {
		if v == s {
			return i
		}
	}
	return -1
}

// overwriteUnmergeables returns the result of writing src over top of dest using the configured options
func overwriteUnmergeables(src, dest interface{}, o *Config) (interface{}, error) {
	o.writeDebug("Others: %#v :: %#v", src, dest)
	if o.KnockoutPrefix != nil && !o.PreserveUnmergeables {
		var srcTemp interface{}

		// apply knockout prefix to source before overwriting dest
		switch s := src.(type) {
		case string:
			// remove knockout string from source before overwriting dest
			srcTemp = strings.TrimLeft(*o.KnockoutPrefix, s)
		case []interface{}:
			// remove all knockout elements before overwriting dest
			t := make([]interface{}, len(s))
			for i, v := range s {
				switch vv := v.(type) {
				case string:
					t[i] = strings.TrimLeft(*o.KnockoutPrefix, vv)
				}
			}
			srcTemp = t
		default:
			srcTemp = s
		}

		if srcTemp == src {
			// we didn't find a KnockoutPrefix so simply overwrite dest
			o.writeDebug("%#v -over -> %#v", src, dest)
			return srcTemp, nil
		} else {
			// found a KnockoutPrefix so delete dest
			o.writeDebug("%#v -over -> %#v", "", dest)
			return "", nil
		}
	} else if !o.PreserveUnmergeables {
		return src, nil
	}
	return dest, nil
}

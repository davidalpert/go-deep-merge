package DeepMerge

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestHashDeepMerge(t *testing.T) {
	var foo = make(map[string]interface{})
	foo["id"] = []int32{3, 4, 5}

	tests := []struct {
		name    string
		src     map[string]interface{}
		dest    map[string]interface{}
		opt     *Config
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "ko_deep_merge!",
			src:  map[string]interface{}{"id": []interface{}{3, 4, 5}},
			dest: map[string]interface{}{"id": []interface{}{1, 2, 3}},
			opt:  NewConfigDeeperMergeKO(),
			want: map[string]interface{}{"id": []interface{}{1, 2, 3, 4, 5}},
		},
		{
			name: "deep_merge!",
			src:  map[string]interface{}{"id": []interface{}{3, 4, 5}},
			dest: map[string]interface{}{"id": []interface{}{1, 2, 3}},
			opt:  NewConfigDeeperMergeBang(),
			want: map[string]interface{}{"id": []interface{}{1, 2, 3, 4, 5}},
		},
		{
			name: "deep_merge",
			src:  map[string]interface{}{"id": "xxx"},
			dest: map[string]interface{}{"id": []interface{}{1, 2, 3}},
			opt:  NewConfigDeeperMerge(),
			want: map[string]interface{}{"id": []interface{}{1, 2, 3}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tt.name), func(t *testing.T) {
			got, err := MergeWithOptions(tt.src, tt.dest, tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRubyHashToMap(t *testing.T) {
	tests := []struct {
		have    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			have: `{'id' => '2'}`,
			want: map[string]interface{}{"id": "2"},
		},
		//{
		//	have: `{"property" => {"bedroom_count" => {'2'=>3, "king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
		//	want: map[string]interface{}{
		//		"property": map[string]interface{}{
		//			"bathroom_count": []interface{}{"1"},
		//			"bedroom_count": map[string]interface{}{
		//				"2":        3,
		//				"king_bed": []interface{}{3},
		//			},
		//		}},
		//	//map[property:map[bathroom_count:[1] bedroom_count:map[2:3 king_bed:[3]]]]
		//},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := rubyHashToMap(tt.have)
			if (err != nil) != tt.wantErr {
				t.Errorf("rubyHashToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rubyHashToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var reRubySymbol = regexp.MustCompile(`(?m):("?[a-zA-Z$]([a-zA-Z0-9_])*([;_<>=?])?"?)`)

// rubyHashToMap takes a string in ruby hash rocket format (like the
// example from the deep_merge lib test cases) and converts it to a
// golang map[string]interface{} so that we can use ruby hash rocket
// literals as test input and expectations to achieve better parity
// and easier maintenance between this lib and the ruby deep_merge
// source
func rubyHashToMap(input string) (map[string]interface{}, error) {
	input = strings.ReplaceAll(input, "'", "\"")
	input = strings.ReplaceAll(input, "=>", ":")
	input = strings.ReplaceAll(input, "nil", "null")
	input = reRubySymbol.ReplaceAllString(input, "\"$1\"")

	return UnmarshallJSONToMap(input)
}

// UnmarshallJSONToMap converts string input in JSON format to a map
// for use in the Merge functions
func UnmarshallJSONToMap(input string) (map[string]interface{}, error) {
	if !strings.HasPrefix(input, "{") || !strings.HasSuffix(input, "}") {
		return nil, fmt.Errorf("this function only works on JSON objects (must start with '{' and end with '}' to serialize as a map")
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		return nil, fmt.Errorf("unmarshalling %#v: %v", input, err)
	} else {
		return result, err
	}
}

// TestDeepMerge contains merge tests (moving from basic to more complex)
func TestDeepMerge(t *testing.T) {
	var foo = make(map[string]interface{})
	foo["id"] = []int32{3, 4, 5}

	tests := []struct {
		name    string
		src     string
		dest    string
		opt     *Config
		want    string
		wantErr bool
	}{
		{
			name: "test merging an hash w/array into blank hash",
			src:  `{'id' => '2'}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().UnpackArraysWith(","),
			want: `{'id' => '2'}`,
		},
		{
			name: "test merging an hash w/array into blank hash",
			src:  `{'region' => {'id' => ['227', '2']}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().UnpackArraysWith(","),
			want: `{'region' => {'id' => ['227', '2']}}`,
		},
		{
			name: "merge from empty hash",
			src:  `{}`,
			dest: `{"property" => ["2","4"]}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => ["2","4"]}`,
		},
		{
			name: "merge to empty hash",
			src:  `{"property" => ["2","4"]}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => ["2","4"]}`,
		},
		{
			name: "simple string overwrite",
			src:  `{"name" => "value"}`,
			dest: `{"name" => "value1"}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"name" => "value"}`,
		},
		{
			name: "simple string overwrite of empty hash",
			src:  `{"name" => "value"}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"name" => "value"}`,
		},
		{
			name: "hashes holding array",
			src:  `{"property" => ["1","3"]}`,
			dest: `{"property" => ["2","4"]}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => ["2","4","1","3"]}`,
		},
		{
			name: "hashes holding array (overwrite)",
			src:  `{"property" => ["1","3"]}`,
			dest: `{"property" => ["2","4"]}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteArrays(true),
			want: `{"property" => ["1","3"]}`,
		},
		{
			name: "hashes holding array (sorted)",
			src:  `{"property" => ["1","3"]}`,
			dest: `{"property" => ["2","4"]}`,
			opt:  NewConfigDeeperMergeBang().WithSortMergedArrays(true),
			want: `{"property" => ["1","2","3","4"]}`,
		},
		{
			name: "hashes holding hashes holding arrays (array with duplicate elements is merged with dest then src",
			src:  `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => ["3", "2"], "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => ["3","2","1"], "bathroom_count" => ["2", "1", "4+"]}}`,
		},
		{
			name: "hash holding hash holding array v string (string is overwritten by array)",
			src:  `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "hash holding hash holding array v string (string is NOT overwritten by array)",
			src:  `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			want: `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "hash holding hash holding string v array (array is overwritten by string)",
			src:  `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "hash holding hash holding string v array (array does NOT overwrite string)",
			src:  `{"property" => {"bedroom_count" => "3", "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			want: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "hash holding hash holding hash v array (array is overwritten by hash)",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => 1}, "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => 1}, "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "hash holding hash holding hash v array (array is NOT overwritten by hash)",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => 1}, "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			want: `{"property" => {"bedroom_count" => ["1", "2"], "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "3 hash layers holding integers (integers are overwritten by source)",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => 1}, "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => 2, "queen_bed" => 4}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => 1}, "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int (arrays are merged)",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [3], "queen_bed" => [1]}, "bathroom_count" => ["1", "4+"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2,3], "queen_bed" => [4,1]}, "bathroom_count" => ["2","1","4+"]}}`,
		},
		{
			name: "1 hash overwriting 3 hash layers holding arrays of int",
			src:  `{"property" => "1"}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => "1"}`,
		},
		{
			name: "1 hash NOT overwriting 3 hash layers holding arrays of int",
			src:  `{"property" => "1"}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int (arrays are merged) but second hash's array is overwritten",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [3], "queen_bed" => [1]}, "bathroom_count" => "1"}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2,3], "queen_bed" => [4,1]}, "bathroom_count" => "1"}}`,
		},
		{
			name: "3 hash layers holding arrays of int (arrays are merged) but second hash's array is NOT overwritten",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [3], "queen_bed" => [1]}, "bathroom_count" => "1"}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2,3], "queen_bed" => [4,1]}, "bathroom_count" => ["2"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but one holds int. This one overwrites, but the rest merge",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => [1]}, "bathroom_count" => ["1"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => 3, "queen_bed" => [4,1]}, "bathroom_count" => ["2","1"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but source is incomplete.",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2,3], "queen_bed" => [4]}, "bathroom_count" => ["2","1"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but source is shorter and has new 2nd level ints.",
			src:  `{"property" => {"bedroom_count" => {"2"=>3, "king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"2"=>3, "king_bed" => [2,3], "queen_bed" => [4]}, "bathroom_count" => ["2","1"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but source is empty",
			src:  `{}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but dest is empty",
			src:  `{"property" => {"bedroom_count" => {"2"=>3, "king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"2"=>3, "king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but source includes a nil in the array",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [nil], "queen_bed" => [1, nil]}, "bathroom_count" => [nil, "1"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [2], "queen_bed" => [4]}, "bathroom_count" => ["2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [2,nil], "queen_bed" => [4, 1, nil]}, "bathroom_count" => ["2", nil, "1"]}}`,
		},
		{
			name: "3 hash layers holding arrays of int, but destination includes a nil in the array",
			src:  `{"property" => {"bedroom_count" => {"king_bed" => [3], "queen_bed" => [1]}, "bathroom_count" => ["1"]}}`,
			dest: `{"property" => {"bedroom_count" => {"king_bed" => [nil], "queen_bed" => [4, nil]}, "bathroom_count" => [nil,"2"]}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"property" => {"bedroom_count" => {"king_bed" => [nil, 3], "queen_bed" => [4, nil, 1]}, "bathroom_count" => [nil, "2", "1"]}}`,
		},
		{
			name: "if extend_existing_arrays == true && destination.kind_of?(Array) && source element is neither array nor hash, push source to destionation",
			src:  `{ "property" => "4" }`,
			dest: `{ "property" => ["1", "2", "3"] }`,
			opt:  NewConfigDeeperMergeBang().WithExtendExistingArrays(true),
			want: `{"property" => ["1", "2", "3", "4"]}`,
		},
		{
			name: "if extend_existing_arrays == true && destination.kind_of?(Array) && source.kind_of(Hash), push source to destination",
			src:  `{ "property" => {:number => "3"} }`,
			dest: `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:  NewConfigDeeperMergeBang().WithExtendExistingArrays(true),
			want: `{"property"=>[{:number=>"1"}, {:number=>"2"}, {:number=>"3"}]}`,
		},
		{
			// assert_raise(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst, {:knockout_prefix => ""})}
			name:    "knockout_prefix & overwrite unmergeable 1",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang().WithKnockout(""),
			wantErr: true,
		},
		{
			// assert_raise(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst, {:preserve_unmergeables => true, :knockout_prefix => ""})}
			name:    "knockout_prefix & overwrite unmergeable 2",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang().WithPreserveUnmergeables(true).WithKnockout(""),
			wantErr: true,
		},
		{
			// assert_raise(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst, {:preserve_unmergeables => true, :knockout_prefix => "--"})}
			name:    "knockout_prefix & overwrite unmergeable 3",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang().WithPreserveUnmergeables(true).WithKnockout("--"),
			wantErr: true,
		},
		{
			// assert_nothing_raised(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst, {:knockout_prefix => "--"})}
			name:    "knockout_prefix & overwrite unmergeable 4",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang().WithKnockout("--").EnableDebug(),
			wantErr: false,
		},
		{
			// assert_nothing_raised(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst)}
			name:    "knockout_prefix & overwrite unmergeable 5",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang(),
			wantErr: false,
		},
		{
			// assert_nothing_raised(DeepMerge::InvalidParameter) {DeepMerge::deep_merge!(hash_src, hash_dst, {:preserve_unmergeables => true})}
			name:    "knockout_prefix & overwrite unmergeable 6",
			src:     `{ "property" => {:number => "3"} }`,
			dest:    `{ "property" => [{:number => "1"}, {:number => "2"}] }`,
			opt:     NewConfigDeeperMergeBang().WithPreserveUnmergeables(true),
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tt.name), func(t *testing.T) {
			// Arrange
			s, err := rubyHashToMap(tt.src)
			if err != nil {
				t.Errorf("unmarshall source >>%s<< to map: %v", tt.src, err)
			}
			d, err := rubyHashToMap(tt.dest)
			if err != nil {
				t.Errorf("unmarshall dest >>%s<< to map: %v", tt.dest, err)
			}

			// Act
			got, err := MergeWithOptions(s, d, tt.opt)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != "" {
				w, err := rubyHashToMap(tt.want)
				if err != nil {
					t.Errorf("unmarshall expectation >>%s<< to map: %v", tt.want, err)
				}
				if !reflect.DeepEqual(got, w) {
					t.Errorf("Merge() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

package DeepMerge

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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
				assert.FailNow(t, "Merge()", "error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "Merge() got = %v, want %v", got, tt.want)
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
		{
			have: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}, "action"=>"browse", "task"=>"browse", "controller"=>"results", "property_order_by"=>"property_type.descr"}`,
			want: map[string]interface{}{"action": "browse", "controller": "results", "property_order_by": "property_type.descr", "region": map[string]interface{}{"id": "230", "ids": []interface{}{"227", "233", "324", "230", "230"}}, "task": "browse"},
		},
		{
			have: `{"property" => {"bedroom_count" => {'2'=>3, "king_bed" => [3]}, "bathroom_count" => ["1"]}}`,
			want: map[string]interface{}{"property": map[string]interface{}{"bathroom_count": []interface{}{"1"}, "bedroom_count": map[string]interface{}{"2": 3.0, "king_bed": []interface{}{3.0}}}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := rubyHashToMap(tt.have)
			if (err != nil) != tt.wantErr {
				assert.FailNow(t, "rubyHashToMap()", "error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "rubyHashToMap() got = %v, want %v", got, tt.want)
		})
	}
}

var reRubySymbol = regexp.MustCompile(`(?m):("?[a-zA-Z$]([a-zA-Z0-9_])*([;_<>?])?"?)`)

// rubyHashToMap takes a string in ruby hash rocket format (like the
// example from the deep_merge lib test cases) and converts it to a
// golang map[string]interface{} so that we can use ruby hash rocket
// literals as test input and expectations to achieve better parity
// and easier maintenance between this lib and the ruby deep_merge
// source
func rubyHashToMap(input string) (map[string]interface{}, error) {
	input = reRubySymbol.ReplaceAllString(input, "\"$1\"")
	input = strings.ReplaceAll(input, "'", "\"")
	input = strings.ReplaceAll(input, "=>", ":")
	input = strings.ReplaceAll(input, "nil", "null")

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
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'id' => '2'}`,
		},
		{
			name: "test merging an hash w/array into blank hash",
			src:  `{'region' => {'id' => ['227', '2']}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
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
			opt:     NewConfigDeeperMergeBang().WithKnockout("--"),
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
		{
			// wierd edge case where hash keys are actually arrays; I've converted them to string constants for this test
			name: "hash holding arrays of arrays",
			src:  `{"[\"1\", \"2\", \"3\"]" => ["1", "2"]}`,
			dest: `{"[\"4\", \"5\"]" => ["3"]}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{"[\"1\", \"2\", \"3\"]" => ["1", "2"], "[\"4\", \"5\"]" => ["3"]}`,
		},
		{
			name: "test merging of hash with blank hash, and make sure that source array split still functions",
			src:  `{'property' => {'bedroom_count' => ["1", "2,3"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX).WithUnpackArrays(","),
			want: `{'property' => {'bedroom_count' => ["1", "2", "3"]}}`,
		},
		{
			name: "test merging of hash with blank hash, and make sure that source array split does not function when turned off",
			src:  `{'property' => {'bedroom_count' => ["1", "2,3"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX),
			want: `{'property' => {'bedroom_count' => ["1", "2,3"]}}`,
		},
		{
			name: "test merging into a blank hash with overwrite_unmergeables turned on",
			src:  `{"action" => "browse", "controller" => "results"}`,
			dest: `{}`,
			// DeepMerge::deep_merge!(hash_src, hash_dst,{:overwrite_unmergeables =>  true, :knockout_prefix =>  FIELD_KNOCKOUT_PREFIX, :unpack_arrays =>  ","})
			opt:  NewConfigDeeperMergeBang().WithPreserveUnmergeables(false).WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX).WithUnpackArrays(","),
			want: `{"action" => "browse", "controller" => "results"}`,
		},
		{
			name: `special params/session style hash with knockout_merge elements in form src: ["1","2"] dest:["--1,--2", "3,4"]`,
			src:  `{"amenity"=>{"id"=>["--1,--2", "3,4"]}}`,
			dest: `{"amenity"=>{"id"=>["1", "2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX).WithUnpackArrays(","),
			// DeepMerge::deep_merge!(hash_params, hash_session, {:knockout_prefix => FIELD_KNOCKOUT_PREFIX, :unpack_arrays => ","})
			want: `{"amenity"=>{"id"=>["3","4"]}}`,
		},
		{
			name: `same as previous but without ko_split value, this merge should fail`,
			src:  `{"amenity"=>{"id"=>["--1,--2", "3,4"]}}`,
			dest: `{"amenity"=>{"id"=>["1", "2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX),
			want: `{"amenity"=>{"id"=>["1","2","3,4"]}}`,
		},
		{
			name: `special params/session style hash with knockout_merge elements in form src: ["1","2"] dest:["--1,--2", "3,4"]`,
			src:  `{"amenity"=>{"id"=>["--1,2", "3,4", "--5", "6"]}}`,
			dest: `{"amenity"=>{"id"=>["1", "2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX).WithUnpackArrays(","),
			want: `{"amenity"=>{"id"=>["2","3","4","6"]}}`,
		},
		{
			name: `special params/session style hash with knockout_merge elements in form src: ["--1,--2", "3,4", "--5", "6"] dest:["1,2", "3,4"]`,
			src:  `{"amenity"=>{"id"=>["--1,--2", "3,4", "--5", "6"]}}`,
			dest: `{"amenity"=>{"id"=>["1", "2", "3", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithKnockout(DEFAULT_FIELD_KNOCKOUT_PREFIX).WithUnpackArrays(","),
			want: `{"amenity"=>{"id"=>["3","4","6"]}}`,
		},
		{
			name: "overwrite_unmergeable_1",
			src:  `{"url_regions" => [], "region" => {"ids" => ["227,233"]}, "action" => "browse", "task" => "browse", "controller" => "results"}`,
			dest: `{"region" => {"ids" => ["227"]}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"url_regions" => [], "region" => {"ids" => ["227", "233"]}, "action" => "browse", "task" => "browse", "controller" => "results"}`,
		},
		{
			name: "overwrite_unmergeable_2",
			src:  `{"region"=>{"ids"=>["--","227"], "id"=>"230"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"ids"=>["227"], "id"=>"230"}}`,
		},
		{
			name: "overwrite_unmergeable_3",
			src:  `{"region"=>{"ids"=>["--","227", "232", "233"], "id"=>"232"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"ids"=>["227", "232", "233"], "id"=>"232"}}`,
		},
		{
			name: "overwrite_unmergeable_4",
			src:  `{"region"=>{"ids"=>["--,227,232,233"], "id"=>"232"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"ids"=>["227", "232", "233"], "id"=>"232"}}`,
		},
		{
			name: "overwrite_unmergeable_5",
			src:  `{"region"=>{"ids"=>["--,227,232","233"], "id"=>"232"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"ids"=>["227", "232", "233"], "id"=>"232"}}`,
		},
		{
			name: "overwrite_unmergeable_6",
			src:  `{"region"=>{"ids"=>["--,227"], "id"=>"230"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"ids"=>["227"], "id"=>"230"}}`,
		},
		{
			name: "overwrite_unmergeable_7",
			src:  `{"region"=>{"ids"=>["--,227"], "id"=>"230"}}`,
			dest: `{"region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}, "action"=>"browse", "task"=>"browse", "controller"=>"results", "property_order_by"=>"property_type.descr"}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(",").EnableDebug(),
			want: `{"region"=>{"ids"=>["227"], "id"=>"230"}, "action"=>"browse", "task"=>"browse",
		"controller"=>"results", "property_order_by"=>"property_type.descr"}`,
		},
		{
			name: "overwrite_unmergeable_8",
			src:  `{"query_uuid"=>"6386333d-389b-ab5c-8943-6f3a2aa914d7", "region"=>{"ids"=>["--,227"], "id"=>"230"}}`,
			dest: `{"query_uuid"=>"6386333d-389b-ab5c-8943-6f3a2aa914d7", "url_regions"=>[], "region"=>{"ids"=>["227", "233", "324", "230", "230"], "id"=>"230"}, "action"=>"browse", "task"=>"browse", "controller"=>"results", "property_order_by"=>"property_type.descr"}`,
			opt:  NewConfigDeeperMergeBang().WithOverwriteUnmergeables(true).WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"query_uuid" => "6386333d-389b-ab5c-8943-6f3a2aa914d7", "url_regions"=>[],
		"region"=>{"ids"=>["227"], "id"=>"230"}, "action"=>"browse", "task"=>"browse",
		"controller"=>"results", "property_order_by"=>"property_type.descr"}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tt.name), func(t *testing.T) {
			// Arrange
			s, err := rubyHashToMap(tt.src)
			assert.NoError(t, err, "unmarshall source >>%s<< to map: %v", tt.src, err)
			d, err := rubyHashToMap(tt.dest)
			assert.NoError(t, err, "unmarshall dest >>%s<< to map: %v", tt.dest, err)

			// Act
			got, err := MergeWithOptions(s, d, tt.opt)

			// Assert
			if (err != nil) != tt.wantErr {
				assert.FailNow(t, "Merge() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != "" {
				w, err := rubyHashToMap(tt.want)
				assert.NoError(t, err, "unmarshall expectation >>%s<< to map: %v", tt.want, err)
				assert.Equal(t, w, got, "MergeWithOptions()")
			}
		})
	}
}

// TestKnockoutPrefixes contains merge tests looking for correct behavior from
// specific real-world params/session merges using the custom modifiers built
// for param/session merges
func TestKnockoutPrefixes(t *testing.T) {
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
		// ko_split: nil
		{
			name: "typical params/session style hash with knockout_merge elements",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["1", "2", "3"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"property"=>{"bedroom_count"=>["2", "3"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["3"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"property"=>{"bedroom_count"=>["3","2"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"property"=>{"bedroom_count"=>["4","2","3"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["--1", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"property"=>{"bedroom_count"=>["4","2","3"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements",
			src:  `{"amenity"=>{"id"=>["--1", "--2", "3", "4"]}}`,
			dest: `{"amenity"=>{"id"=>["1", "2"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"amenity"=>{"id"=>["3","4"]}}`,
		},

		// ko_split: ","
		{
			name: "typical params/session style hash with knockout_merge elements, with ko_split: \",\"",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["1", "2", "3"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"property"=>{"bedroom_count"=>["2", "3"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements, with ko_split: \",\"",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["3"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"property"=>{"bedroom_count"=>["3","2"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements, with ko_split: \",\"",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"property"=>{"bedroom_count"=>["4","2","3"]}}`,
		},
		{
			name: "typical params/session style hash with knockout_merge elements, with ko_split: \",\"",
			src:  `{"property"=>{"bedroom_count"=>["--1", "2", "3"]}}`,
			dest: `{"property"=>{"bedroom_count"=>["--1", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"property"=>{"bedroom_count"=>["4","2","3"]}}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => "--"}`,
			dest: `{"amenity" => "1"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => ""}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--"]}`,
			dest: `{"amenity" => "1"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => []}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => "--"}`,
			dest: `{"amenity" => ["1"]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => ""}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--"]}`,
			dest: `{"amenity" => ["1"]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => []}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--"]}`,
			dest: `{"amenity" => "1"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => []}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--", "2"]}`,
			dest: `{'amenity' => ["1", "3", "7+"]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => ["2"]}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--", "2"]}`,
			dest: `{'amenity' => "5"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => ['2']}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => "--"}`,
			dest: `{"amenity"=>{"id"=>["1", "2", "3", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => ""}`,
		},
		{
			name: `knock out entire dest hash if "--" is passed for source`,
			src:  `{'amenity' => ["--"]}`,
			dest: `{"amenity"=>{"id"=>["1", "2", "3", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'amenity' => []}`,
		},
		{
			name: `knock out dest array if "--" is passed for source`,
			src:  `{"region" => {'ids' => "--"}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"]}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => {'ids' => ""}}`,
		},
		{
			name: `knock out dest array but leave other elements of hash intact`,
			src:  `{"region" => {'ids' => "--"}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => {'ids' => "", 'id'=>'11'}}`,
		},
		{
			name: `knock out entire tree of dest hash`,
			src:  `{"region" => "--"}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => ""}`,
		},
		{
			name: `knock out entire tree of dest hash - retaining array format`,
			src:  `{"region" => {'ids' => ["--"]}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => {'ids' => [], 'id'=>'11'}}`,
		},
		{
			name: `knock out entire tree of dest hash & replace with new content`,
			src:  `{"region" => {'ids' => ["2", "--", "6"]}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => {'ids' => ["2", "6"], 'id'=>'11'}}`,
		},
		{
			name: `knock out entire tree of dest hash & replace with new content`,
			src:  `{"region" => {'ids' => ["7", "--", "6"]}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' => {'ids' => ["7", "6"], 'id'=>'11'}}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tt.name), func(t *testing.T) {
			// Arrange
			s, err := rubyHashToMap(tt.src)
			assert.NoError(t, err, "unmarshall source >>%s<< to map: %v", tt.src, err)
			d, err := rubyHashToMap(tt.dest)
			assert.NoError(t, err, "unmarshall dest >>%s<< to map: %v", tt.dest, err)

			// Act
			got, err := MergeWithOptions(s, d, tt.opt)

			// Assert
			if (err != nil) != tt.wantErr {
				assert.FailNow(t, "Merge() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != "" {
				w, err := rubyHashToMap(tt.want)
				assert.NoError(t, err, "unmarshall expectation >>%s<< to map: %v", tt.want, err)
				assert.Equal(t, w, got, "Merge() got = %v, want %v", got, w)
			}
		})
	}
}

// TestEdges contains some tests marked "edge test" and "Example"
func TestEdges(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		dest    string
		opt     *Config
		want    string
		wantErr bool
	}{
		{
			name: `edge test: make sure that when we turn off knockout_prefix that all values are processed correctly`,
			src:  `{"region" => {'ids' => ["7", "--", "2", "6,8"]}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang().WithUnpackArrays(","),
			want: `{'region' => {'ids' => ["1", "2", "3", "4", "7", "--", "6", "8"], 'id'=>'11'}}`,
		},
		{
			name: `edge test 2: make sure that when we turn off source array split that all values are processed correctly`,
			src:  `{"region" => {'ids' => ["7", "3", "--", "6,8"]}}`,
			dest: `{"region"=>{"ids"=>["1", "2", "3", "4"], 'id'=>'11'}}`,
			opt:  NewConfigDeeperMergeBang(),
			want: `{'region' => {'ids' => ["1", "2", "3", "4", "7", "--", "6,8"], 'id'=>'11'}}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "1"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>"--1"}`,
			dest: `{"amenity"=>"1"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"amenity"=>""}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "2"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>"--1"}`,
			dest: `{"amenity"=>"2"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"amenity"=>""}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "1"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>["--1"]}`,
			dest: `{"amenity"=>"1"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"amenity"=>[]}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "1"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>["--1"]}`,
			dest: `{"amenity"=>["1"]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"amenity"=>[]}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "1"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>"--1"}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"amenity"=>""}`,
		},
		{
			name: `Example: src = {'key' => "--1"}, dst = {'key' => "1"} -> merges to {'key' => ""}`,
			src:  `{"amenity"=>"--1"}`,
			dest: `{"amenity"=>["1"]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"amenity"=>""}`,
		},
		{
			name: `are unmerged hashes passed unmodified w/out :unpack_arrays?`,
			src:  `{"amenity"=>{"id"=>["26,27"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix(),
			want: `{"amenity"=>{"id"=>["26,27"]}}`,
		},
		{
			name: `hash should be merged`,
			src:  `{"amenity"=>{"id"=>["26,27"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"amenity"=>{"id"=>["26","27"]}}`,
		},
		{
			name: `hashes with knockout values are suppressed`,
			src:  `{"amenity"=>{"id"=>["--26,--27,28"]}}`,
			dest: `{}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"amenity"=>{"id"=>["28"]}}`,
		},
		{
			name: `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'ids'=>['227','2','3','3']}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'ids'=>[]}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'ids'=>['227','2','3','3'], 'id' => '3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'ids'=>[], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'ids'=>['--']}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '2244', 'ids'=>[], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'ids'=>['--'], 'id' => '5'}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'ids'=>['--'], 'id' => '5'}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '2244', 'ids'=>[], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'ids'=>['--', '227'], 'id' => '5'}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'ids'=>['--', '227'], 'id' => '5'}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'muni_city_id' => '--', 'ids'=>'--', 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'muni_city_id' => '--', 'ids'=>'--', 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '', 'ids'=>'', 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'muni_city_id' => '--', 'ids'=>['--'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'muni_city_id' => '--', 'ids'=>['--'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '', 'ids'=>[], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{'region' =>{'muni_city_id' => '--', 'ids'=>['--','227'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			src:  `{'region' =>{'muni_city_id' => '--', 'ids'=>['--','227'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
			dest: `{'region' =>{'muni_city_id' => '2244', 'ids'=>['227','2','3','3'], 'id'=>'3'}, 'query_uuid' => 'zzz'}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{'region' =>{'muni_city_id' => '', 'ids'=>['227'], 'id'=>'5'}, 'query_uuid' => 'zzz'}`,
		},
		{
			name: `{"muni_city_id"=>"--", "id"=>""}`,
			src:  `{"muni_city_id"=>"--", "id"=>""}`,
			dest: `{"muni_city_id"=>"", "id"=>""}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"muni_city_id"=>"", "id"=>""}`,
		},
		{
			name: `{"region"=>{"muni_city_id"=>"--", "id"=>""}}`,
			src:  `{"region"=>{"muni_city_id"=>"--", "id"=>""}}`,
			dest: `{"region"=>{"muni_city_id"=>"", "id"=>""}}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"region"=>{"muni_city_id"=>"", "id"=>""}}`,
		},
		{
			name: `{"query_uuid"=>"a0dc3c84-ec7f-6756-bdb0-fff9157438ab", "url_regions"=>[], "region"=>{"muni_city_id"=>"--", "id"=>""}, "property"=>{"property_type_id"=>"", "search_rate_min"=>"", "search_rate_max"=>""}, "task"=>"search", "run_query"=>"Search"}`,
			src:  `{"query_uuid"=>"a0dc3c84-ec7f-6756-bdb0-fff9157438ab", "url_regions"=>[], "region"=>{"muni_city_id"=>"--", "id"=>""}, "property"=>{"property_type_id"=>"", "search_rate_min"=>"", "search_rate_max"=>""}, "task"=>"search", "run_query"=>"Search"}`,
			dest: `{"query_uuid"=>"a0dc3c84-ec7f-6756-bdb0-fff9157438ab", "url_regions"=>[], "region"=>{"muni_city_id"=>"", "id"=>""}, "property"=>{"property_type_id"=>"", "search_rate_min"=>"", "search_rate_max"=>""}, "task"=>"search", "run_query"=>"Search"}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","),
			want: `{"query_uuid"=>"a0dc3c84-ec7f-6756-bdb0-fff9157438ab", "url_regions"=>[], "region"=>{"muni_city_id"=>"", "id"=>""}, "property"=>{"property_type_id"=>"", "search_rate_min"=>"", "search_rate_max"=>""}, "task"=>"search", "run_query"=>"Search"}`,
		},
		{
			name: `hash of array of hashes`,
			src:  `{"item" => [{"1" => "3"}, {"2" => "4"}]}`,
			dest: `{"item" => [{"3" => "5"}]}`,
			opt:  NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().EnableDebug(),
			want: `{"item" => [{"3" => "5"}, {"1" => "3"}, {"2" => "4"}]}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tt.name), func(t *testing.T) {
			// Arrange
			s, err := rubyHashToMap(tt.src)
			assert.NoError(t, err, "unmarshall source >>%s<< to map: %v", tt.src, err)
			d, err := rubyHashToMap(tt.dest)
			assert.NoError(t, err, "unmarshall dest >>%s<< to map: %v", tt.dest, err)

			// Act
			got, err := MergeWithOptions(s, d, tt.opt)

			// Assert
			if (err != nil) != tt.wantErr {
				assert.FailNow(t, "Merge() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != "" {
				w, err := rubyHashToMap(tt.want)
				assert.NoError(t, err, "unmarshall expectation >>%s<< to map: %v", tt.want, err)
				assert.Equal(t, w, got, "Merge() got = %v, want %v", got, w)
			}

			// special case of a second merge
			if tt.name == "hash should be merged" {
				t.Run("second merge of same values should result in no change in output", func(t *testing.T) {
					//assert.FailNow(t, "here")

					hashParams, err := rubyHashToMap(`{"amenity"=>{"id"=>["26,27"]}}`)
					assert.NoError(t, err, tt.name+" part 2")
					//// DeepMerge::deep_merge!(hash_params, hash_session, {:knockout_prefix => FIELD_KNOCKOUT_PREFIX, :unpack_arrays => ","})
					got2, err := MergeWithOptions(hashParams, got, NewConfigDeeperMergeBang().WithDefaultKnockoutPrefix().WithUnpackArrays(","))
					assert.NoError(t, err, tt.name+" part 2 - merge")

					want2, err := rubyHashToMap(`{"amenity"=>{"id"=>["26","27"]}}`)
					assert.NoError(t, err, tt.name+" part 2 - want2")
					assert.Equal(t, want2, got2, "Merge() got = %v, want %v", got2, want2)
				})
			}
		})
	}
}

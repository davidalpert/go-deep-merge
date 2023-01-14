package DeepMerge

import (
	"fmt"
	"reflect"
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

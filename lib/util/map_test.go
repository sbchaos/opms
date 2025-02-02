package util_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/sbchaos/opms/lib/util"
)

func TestMapHelper(t *testing.T) {
	t.Run("MergeAnyMaps", func(t *testing.T) {
		t.Run("clones the map when only one parameter", func(t *testing.T) {
			orig := map[string]interface{}{
				"Key": "Value",
			}

			clone := util.MergeAnyMaps(orig)
			assert.NotNil(t, clone)
			assert.Equal(t, clone["Key"], orig["Key"])
			// Check if both map pointers are different
			assert.NotEqual(t, reflect.ValueOf(clone).Pointer(), reflect.ValueOf(orig).Pointer())
		})
		t.Run("merges when multiple maps", func(t *testing.T) {
			orig := map[string]interface{}{
				"Key1": "Value1",
				"Key2": "Value1",
			}
			orig2 := map[string]interface{}{
				"Key2": "Value2",
			}

			merged := util.MergeAnyMaps(orig, orig2)
			assert.NotNil(t, merged)
			assert.Len(t, merged, 2)
			assert.Equal(t, merged["Key1"], orig["Key1"])
			assert.Equal(t, merged["Key2"], orig2["Key2"])
		})
	})
	t.Run("MergeMaps", func(t *testing.T) {
		t.Run("merges string maps", func(t *testing.T) {
			mp1 := map[string]string{
				"Key": "3",
			}
			mp2 := map[string]string{
				"Key2": "4",
			}

			merged := util.MergeMaps(mp1, mp2)
			assert.NotNil(t, merged)
			assert.Equal(t, "3", merged["Key"])
			assert.Equal(t, "4", merged["Key2"])
		})
		t.Run("overrides the values in first map", func(t *testing.T) {
			mp1 := map[string]string{
				"Key": "3",
			}
			mp2 := map[string]string{
				"Key": "4",
			}

			merged := util.MergeMaps(mp1, mp2)
			assert.NotNil(t, merged)
			assert.Equal(t, "4", merged["Key"])
		})
	})
	t.Run("MapToList", func(t *testing.T) {
		t.Run("returns the values in map", func(t *testing.T) {
			mapping := map[string]string{"a": "b", "c": "d", "e": "f"}
			list := util.MapToList(mapping)
			assert.Len(t, list, 3)
			sort.Strings(list)
			assert.Equal(t, []string{"b", "d", "f"}, list)
		})
	})
	t.Run("ListToMap", func(t *testing.T) {
		t.Run("returns the list of strings as map", func(t *testing.T) {
			list := []string{"a", "b", "c", "d"}
			mapping := util.ListToMap(list)
			assert.Len(t, mapping, 4)
			assert.EqualValues(t, map[string]struct{}{"a": {}, "b": {}, "c": {}, "d": {}}, mapping)
		})
	})
	t.Run("AppendToMap", func(t *testing.T) {
		t.Run("appends data from string map", func(t *testing.T) {
			orig := map[string]interface{}{
				"Key": "Value1",
			}

			toAppend := map[string]string{
				"Key2": "Value2",
			}

			util.AppendToMap(orig, toAppend)
			assert.Len(t, orig, 2)
			assert.Equal(t, "Value2", orig["Key2"])
		})
	})
	t.Run("Contains", func(t *testing.T) {
		t.Run("returns false when map is nil", func(t *testing.T) {
			result := util.Contains[string, string](nil, "a")
			assert.False(t, result)
		})
		t.Run("returns false when some value not found", func(t *testing.T) {
			mp := map[string]string{"a": "b", "c": "d", "e": "f"}

			result := util.Contains(mp, "a", "c", "g")
			assert.False(t, result)
		})
		t.Run("returns true when all values found", func(t *testing.T) {
			mp := map[string]string{"a": "b", "c": "d", "e": "f"}

			result := util.Contains(mp, "a", "c", "e")
			assert.True(t, result)
		})
	})
	t.Run("ConfigAs", func(t *testing.T) {
		t.Run("returns false when map is nil", func(t *testing.T) {
			input, _ := structpb.NewStruct(map[string]any{
				"time": 10,
			})

			result := util.ConfigAs[float64](input.AsMap(), "time")
			assert.Equal(t, result, float64(10))
		})
	})
}

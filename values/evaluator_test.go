package values

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var lessTests = []struct {
	a, b     any
	expected bool
}{
	{nil, nil, false},
	{false, true, true},
	{false, false, false},
	{false, nil, false},
	{nil, false, false},
	{0, 1, true},
	{1, 0, false},
	{1, 1, false},
	{1, 2.1, true},
	{1.1, 2, true},
	{2.1, 1, false},
	{"a", "b", true},
	{"b", "a", false},
	{[]string{"a"}, []string{"a"}, false},
}

func TestLess(t *testing.T) {
	for i, test := range lessTests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			value := Less(test.a, test.b)
			require.Equalf(t, test.expected, value, "%#v < %#v", test.a, test.b)
		})
	}
}

func TestLength(t *testing.T) {
	require.Equal(t, 3, Length([]int{1, 2, 3}))
	require.Equal(t, 3, Length("abc"))
	require.Equal(t, 0, Length(map[string]int{"a": 1}))
}

func TestSort(t *testing.T) {
	array := []any{2, 1}
	Sort(array)
	require.Equal(t, []any{1, 2}, array)

	array = []any{"b", "a"}
	Sort(array)
	require.Equal(t, []any{"a", "b"}, array)

	array = []any{
		map[string]any{"key": 20},
		map[string]any{"key": 10},
		map[string]any{},
	}
	SortByProperty(array, "key", true)
	require.Nil(t, array[0].(map[string]any)["key"])
	require.Equal(t, 10, array[1].(map[string]any)["key"])
	require.Equal(t, 20, array[2].(map[string]any)["key"])
}

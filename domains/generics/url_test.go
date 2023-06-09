package generics

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURL_WithQuery(t *testing.T) {
	t.Run("Success/ShouldReturnCopy", func(t *testing.T) {
		src := URL{Host: "https://example.com", Path: "/path", Query: map[string][]string{"a": {"foo", "bar"}}}
		dst := src.WithQuery(map[string]interface{}{"a": "qux", "b": 3})

		require.Equal(t, map[string][]string{"a": {"foo", "bar"}}, src.Query)
		require.Equal(t, map[string][]string{"a": {"foo", "bar", "qux"}, "b": {"3"}}, dst.Query)
	})

	data := []struct {
		name     string
		src      URL
		data     map[string]interface{}
		expected URL
	}{
		{
			name: "Success",
			src:  URL{Host: "https://example.com", Path: "/path"},
			data: map[string]interface{}{"a": "foo", "b": 3},
			expected: URL{
				Host: "https://example.com",
				Path: "/path",
				Query: map[string][]string{
					"a": {"foo"},
					"b": {"3"},
				},
			},
		},
		{
			name: "Success/ShouldMergeQueries",
			src:  URL{Host: "https://example.com", Path: "/path", Query: map[string][]string{"a": {"foo", "bar"}, "c": {"baz"}}},
			data: map[string]interface{}{"a": "qux", "b": 3},
			expected: URL{
				Host: "https://example.com",
				Path: "/path",
				Query: map[string][]string{
					"a": {"foo", "bar", "qux"},
					"b": {"3"},
					"c": {"baz"},
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			require.Equal(t, d.expected, d.src.WithQuery(d.data))
		})
	}
}

func TestURL_String(t *testing.T) {
	data := []struct {
		name     string
		url      URL
		expected string
	}{
		{
			name:     "Success/ShouldReturnURL",
			url:      URL{Host: "https://example.com", Path: "/path"},
			expected: "https://example.com/path",
		},
		{
			name:     "Success/ShouldReturnURLWithQuery",
			url:      URL{Host: "https://example.com", Path: "/path", Query: map[string][]string{"a": {"foo", "bar"}}},
			expected: "https://example.com/path?a=foo,bar",
		},
		{
			name:     "Success/ShouldReturnURLWithMultipleQueries",
			url:      URL{Host: "https://example.com", Path: "/path", Query: map[string][]string{"a": {"foo", "bar"}, "c": {"baz"}}},
			expected: "https://example.com/path?a=foo,bar&c=baz",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			require.Equal(t, d.expected, d.url.String())
		})
	}
}

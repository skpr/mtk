package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	var testCases = []struct {
		comment string
		config  string
		want    Rules
	}{
		{
			"ignore",
			"test-data/ignore.yml",
			Rules{
				Ignore: []string{
					"ignore_this_table",
					"and_this_one",
				},
			},
		},
		{
			"nodata",
			"test-data/nodata.yml",
			Rules{
				NoData: []string{
					"table_with_structure_only_please",
					"yeah_do_this_one_too",
				},
			},
		},
		{
			"sanitize",
			"test-data/rewrite.yml",
			Rules{
				Rewrite: map[string]Rewrite{
					"accounts": map[string]string{
						"email": "concat(id, \"@sanitized\")",
						"password": "\"SANITIZED_PASSWORD\"",
					},
				},
			},
		},
		{
			"mixed",
			"test-data/mixed.yml",
			Rules{
				Ignore: []string{"foo"},
				NoData: []string{"bar"},
				Rewrite: map[string]Rewrite{
					"baz": map[string]string{
						"qux": "quux",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := Load(testCase.config)
		fmt.Println(actual.SanitizeMap())
		assert.Nil(t, err)
		assert.Equal(t, testCase.want, actual, testCase.comment)
	}
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	var testCases = []struct {
		comment string
		config  string
		want    File
	}{
		{
			"ignore",
			"test-data/ignore.yml",
			File{
				Ignore: []string{
					"ignore_this_table",
					"and_this_one",
				},
			},
		},
		{
			"nodata",
			"test-data/nodata.yml",
			File{
				NoData: []string{
					"table_with_structure_only_please",
					"yeah_do_this_one_too",
				},
			},
		},
		{
			"sanitize",
			"test-data/sanitize.yml",
			File{
				Sanitize: Sanitize{
					Tables: []Table{
						{
							"accounts",
							[]Field{
								{
									Name:  "email",
									Value: "SANITIZED_MAIL",
								},
								{
									Name:  "password",
									Value: "SANITIZED_PASSWORD",
								},
							},
						},
					},
				},
			},
		},
		{
			"mixed",
			"test-data/mixed.yml",
			File{
				Ignore: []string{"foo"},
				NoData: []string{"bar"},
				Sanitize: Sanitize{
					[]Table{
						Table{
							"baz",
							[]Field{
								{
									Name:  "qux",
									Value: "quux",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := Load(testCase.config)
		assert.Nil(t, err)
		assert.Equal(t, testCase.want, actual, testCase.comment)
	}
}

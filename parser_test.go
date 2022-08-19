package parser

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
	"testing"
)

func TestFindAttributeAndGetValue(t *testing.T) {
	attributes := Attributes{{Key: "key_test1", Val: "val_test2"}, {Key: "key_test2", Val: "val_test2"}, {Key: "key_test3", Val: "val_test3"}}

	type params struct {
		attributes Attributes
		key        string
	}

	type testCase struct {
		name   string
		params params
		want   string
	}

	tests := []testCase{
		{
			name: "Test key example 1",
			params: params{
				attributes: attributes,
				key:        "key_test2",
			},
			want: "val_test2",
		},
		{
			name: "Test key example 2",
			params: params{
				attributes: attributes,
				key:        "key_test3",
			},
			want: "val_test3",
		},
		{
			name: "Non-existent key test",
			params: params{
				attributes: attributes,
				key:        "fsafaf",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			received := FindAttributeAndGetValue(tt.params.attributes, func(attribute html.Attribute) bool {
				return attribute.Key == tt.params.key
			})

			assert.Equal(t, tt.want, received)
		})
	}

}

func TestParse(t *testing.T) {
	type params struct {
		filePath string
	}

	type testCase struct {
		name   string
		params params
		want   []Link
	}

	tests := []testCase{
		{
			name:   "Simple link node test",
			params: params{"testdata/example_1.html"},
			want: []Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
		},
		{
			name:   "Sibling links test",
			params: params{"testdata/example_2.html"},
			want: []Link{
				{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github!",
				},
			},
		},
		{
			name:   "Full HTML page",
			params: params{"testdata/example_3.html"},
			want: []Link{
				{
					Href: "#",
					Text: "Login",
				},
				{
					Href: "/lost",
					Text: "Lost? Need help?",
				},
				{
					Href: "https://twitter.com/marcusolsson",
					Text: "@marcusolsson",
				},
			},
		},
		{
			name:   "Link with comments inside",
			params: params{"testdata/example_4.html"},
			want: []Link{
				{
					Href: "/dog-cat",
					Text: "dog cat",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlReader, err := os.Open(tt.params.filePath)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			var actual []Link
			actual, err = Parse(htmlReader, atom.A)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			assert.Equal(t, tt.want, actual)
		})
	}
}

package search

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		minKeyworldLen int
		stopwords      []string
	}
	tests := []struct {
		name string
		args args
		want *Index
	}{
		{
			"simple search",
			args{3, []string{"the", "and"}},
			&Index{
				CurrentId:      1,
				DocMap:         make(map[int]Document),
				TagMap:         make(map[string][]int),
				KwMap:          make(map[string]map[int]int),
				stopWords:      []string{"the", "and"},
				minKeyworldLen: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.minKeyworldLen, tt.args.stopwords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Add(t *testing.T) {
	idx := New(3, []string{"the", "and"})
	err := idx.Add("http://example.com", "Example", "This is an example", []string{"example", "test"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(idx.DocMap))
	assert.Equal(t, 2, len(idx.TagMap))
	assert.Equal(t, 1, idx.TagMap["example"][0])
	n := titleKeywordValue + tagKeywordValue + 1
	assert.Equal(t, map[int]int{1: n}, idx.KwMap["example"])
}

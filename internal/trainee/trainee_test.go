package trainee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnglish(t *testing.T) {
	name := "Ben"
	language := "english"
	expected := "Hello Ben"
	actual, err := Hello(name, language)

	if err != nil {
		t.Errorf("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestAnotherEnglish(t *testing.T) {
	name := "Ben"
	language := "english"
	expected := "Hello Ben"
	actual, err := Hello(name, language)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestFibRecursive(t *testing.T) {
	expected1 := uint(55)
	actual1 := FibRecursive(10)

	expected2 := uint(55)
	actual2 := FibIterative(10)

	assert.Equal(t, expected1, actual1)
	assert.Equal(t, expected2, actual2)
}

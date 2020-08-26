package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
)

func TestPKCS7Padding(t *testing.T) {

	var tests = []struct {
		input     []byte
		blockSize int
		expected  []byte
	}{
		{[]byte{'g', 'o', 'l', 'a', 'n', 'g'}, 16, []byte{'g', 'o', 'l', 'a', 'n', 'g', 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}},
		{[]byte{'j', 'a', 'v', 'a'}, 3, []byte{'j', 'a', 'v', 'a', 2, 2}},
	}

	for _, test := range tests {
		output := pkcs7Padding(test.input, test.blockSize)
		assert.Equal(t, output, test.expected, "should be equal")
	}
}

func TestPKCS7UnPadding(t *testing.T) {

	var tests = []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{'g', 'o', 'l', 'a', 'n', 'g', 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}, []byte{'g', 'o', 'l', 'a', 'n', 'g'}},
		{[]byte{'j', 'a', 'v', 'a', 2, 2}, []byte{'j', 'a', 'v', 'a'}},
	}

	for _, test := range tests {
		output := pkcs7UnPadding(test.input)
		assert.Equal(t, output, test.expected, "should be equal")
	}
}

func TestSha1Sum(t *testing.T) {
	hash := sha1Sum([]byte{'g', 'o', 'l', 'a', 'n', 'g'}, []byte{'1', '2', '3', '4', '5'})
	assert.Equal(t, hex.EncodeToString(hash), "31fa5264168f0a276619c952712854a726456de2", "should be equal")
}

func TestSha256Sum(t *testing.T) {
	hash := sha256Sum([]byte{'g', 'o', 'l', 'a', 'n', 'g'})
	assert.Equal(t, hex.EncodeToString(hash), "d754ed9f64ac293b10268157f283ee23256fb32a4f8dedb25c8446ca5bcb0bb3", "should be equal")
}

func TestSha512Sum(t *testing.T) {
	hash := sha512Sum([]byte{'g', 'o', 'l', 'a', 'n', 'g'})
	assert.Equal(t, hex.EncodeToString(hash[:]), "df84c5d44709cfeb8a22c8cf006ac926c92c6823d37e112f2c68a22890e61615f97ad1d4eb1d3e043442063886b4ce2f15eaa73ea8ff769808fc76d47f607ec5", "should be equal")
}

func TestAppendBytes(t *testing.T) {
	b1 := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	b2 := []byte{'g', 'o', 'o', 'g', 'l', 'e'}
	b3 := []byte{1, 2, 3, 4}

	output := appendBytes(b1, b2, b3)
	assert.Equal(t, string(output), string(b1)+string(b2)+string(b3), "should be equal")
}

func TestIsOrderedMapEmpty(t *testing.T) {
	o := orderedmap.New()
	assert.Equal(t, isOrderedMapEmpty(o), true, "should be equal")

	o.Set("key", "value")
	assert.Equal(t, isOrderedMapEmpty(o), false, "should be equal")
}

func TestCloneOrderedMap(t *testing.T) {
	o := orderedmap.New()
	o.Set("key", "value")

	clone, _ := cloneOrderedMap(o)

	assert.NotSame(t, o, clone, "should not be equal")

	o.Set("key1", "value1")

	assert.Equal(t, len(clone.Keys()), 1, "should be equal")
}

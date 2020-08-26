package bridgeutil

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
)

func TestStringToOrderedMap(t *testing.T) {

	var tests = []string{
		`{"key":{"key1":"value1"}}`,
		`{"data":{"private_info":"79676feb56c7b8c222924d945ba3d7c73333c27b7bc94e8a76cbaa643db3722695d7b822aa3d62443f3bacbdb993b45ec9421769b15b97bd085c0fc21132de4c08a4626b28ddc40481e1563245b337ffb782113e364cc94e40348577eae4a714c9764e6c206439b1d86fa97c17f33164f2a2ca343dd1d5f9e7d2c68fbb8ed58d","transaction":{"originator_vasp":{"vasp_code":"VASPUSNY1","addrs":[{"address":"bnb1vynn9hamtqg9me7y6frja0rvfva9saprl55gl4","addr_extra_info":[]}]},"beneficiary_vasp":{"vasp_code":"VASPUSNY2","addrs":[{"address":"bnb1hj767k8nlf0jn6p3c3wvl0a66c4782a3f78d7e","addr_extra_info":[{"tag":"abc"}]}]},"currency_id":"sygna:0x80000090","amount":"4.51120135938784"},"data_dt":"2020-07-13T05:56:53.088Z","signature":"2f536f6fa60ad75be96269517cb05133b760c457e5becd1a276a4c920539dbb44ce49649b79d0c74c43c68df406df9824bb92fa4901c36503a67700a44d93eaa"},"callback":{"callback_url":"https://facb1c03d3dae42f07008d0c42979623.m.pipedream.net","signature":"be9000b96b5a86b971fe1818e23790beb33fc9d2b27d761ea70c067eb73adea06fd8aada3ec577f62e87b77ff18cb635bd48e1e33b677908b0bf92ea743c85b4"}}`,
	}

	for _, v := range tests {
		o := StringToOrderedMap(v)
		bOrderedMap, _ := json.Marshal(o)
		assert.Equal(t, string(bOrderedMap), v, "should be equal")
	}
}

func TestOrderedMapToString(t *testing.T) {
	childMap := orderedmap.New()
	childMap.Set("key1", "value1")

	o := orderedmap.New()
	o.Set("key", childMap)

	message, err := OrderedMapToString(o)
	assert.Nil(t, err)
	assert.Equal(t, message, `{"key":{"key1":"value1"}}`, "should be equal")

	o1 := orderedmap.New()
	o1.Set("key2", "value2")

	maps := make([]*orderedmap.OrderedMap, 2)
	maps[0] = o
	maps[1] = o1

	message, err = OrderedMapToString(maps...)
	assert.Nil(t, err)
	assert.Equal(t, message, `[{"key":{"key1":"value1"}},{"key2":"value2"}]`, "should be equal")
}

func TestCastArrayToOrderedMapArray(t *testing.T) {
	originalData := `{"vasp_data":[{"vasp_code":"AAAAAAAA798","vasp_name":"ASH","vasp_pubkey":"04629dac91cbe671b38b20822f03fe39252a0f93505111c330fbf531af91f3a05e439ec27c4e8ad0b705408bbe9f1e225beeb2b1a33b1b7a23a20040a8c95fca61"},{"vasp_code":"AABCASRR","vasp_name":"2489723983723987","vasp_pubkey":"047453b1a211cb1ad0fae92e9c6b8eccfc0935e6d70c15c9fd561929d50365b087cc8ee6609625d64f33fdda8591a1fd2eeb10ad41031a2fd7ea546bcc0f46ae73"},{"vasp_code":"ABCDKRZZ111","vasp_name":"ASDFGHJKL111111","vasp_pubkey":"22222222222222222222222"}],"signature":"d48333a7069e584449c396f4f340d865cb2ace07f058eb5d22315c76d443976e04d3d4716f6baf3ac9a729521c2ec7d066ed8c675c817c5394666917f1fd528c","timestamp":1596596057813}`
	o := StringToOrderedMap(originalData)

	vaspData, _ := o.Get("vasp_data")
	assert.Equal(t, reflect.TypeOf(vaspData).String(), "[]interface {}")

	vaspDataMaps := castArrayToOrderedMapArray(vaspData)
	assert.Equal(t, reflect.TypeOf(vaspDataMaps).String(), "[]*orderedmap.OrderedMap")

	bOrderedMap, _ := json.Marshal(o)
	assert.Equal(t, string(bOrderedMap), originalData, "should be equal")
}

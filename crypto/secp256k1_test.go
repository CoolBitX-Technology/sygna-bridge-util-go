package crypto

import (
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {

	o := orderedmap.New()
	o.UnmarshalJSON([]byte(`{"transfer_id":"b97903fd68fcff05cfe035482bc3cf7fd934505b4e0644e612087dca4bae37e4","txid":"6f721fba0d405df21fb27dd76cfe2b548907f3881c5625b9cfe624c15c3178ae"}`))

	originator := orderedmap.New()
	originator.Set("name", "Antoine Griezmann")
	originator.Set("date_of_birth", "1991-03-21")

	beneficiary := orderedmap.New()
	beneficiary.Set("name", "利昂內爾 梅西")

	o1 := orderedmap.New()
	o1.Set("originator", originator)
	o1.Set("beneficiary", beneficiary)

	o2 := orderedmap.New()
	o2.Set("username", "kunming")
	o2.Set("password", 1234)
	o2.Set("signature", "abcdef") //it would be replace empty before signing
	o2.Set("abc", 1.234)

	//sign data by javascript bridge util
	var tests = []struct {
		input    *orderedmap.OrderedMap
		expected string
	}{
		{o, "a599a99d018f544701e3ae1217f783581a23228d23a5fe18ff96e9fb6471d75127943bd791e3d69495a787cc0a689b4777c875f5302bf116ee88ac27f5562b2a"},
		{o1, "70be6318f31204c9fe28e0b30dabff02b2909105bd0a97d094d6ed5d497461077afa65b440afe37e238172a88b05f3240d0eabb09d008fa25ac0224a507c56b5"},
		{o2, "3bec3a43f9b5679647b0da870fb6c955488a47e7eb44f285491ebcf084ec69ca4170c8eb59f1f22002d5128c09affbadc8a56092c234dbfcedf916bf9dad17dc"},
	}

	for _, test := range tests {
		Sign(test.input, fakePrivateKey)
		signature, _ := test.input.Get("signature")
		assert.Equal(t, signature, test.expected, "should be equal")
	}
}

func TestVerify(t *testing.T) {

	o := orderedmap.New()
	o.UnmarshalJSON([]byte(`{"transfer_id":"b97903fd68fcff05cfe035482bc3cf7fd934505b4e0644e612087dca4bae37e4","txid":"6f721fba0d405df21fb27dd76cfe2b548907f3881c5625b9cfe624c15c3178ae","signature":"a599a99d018f544701e3ae1217f783581a23228d23a5fe18ff96e9fb6471d75127943bd791e3d69495a787cc0a689b4777c875f5302bf116ee88ac27f5562b2a"}`))

	originator := orderedmap.New()
	originator.Set("name", "Antoine Griezmann")
	originator.Set("date_of_birth", "1991-03-21")

	beneficiary := orderedmap.New()
	beneficiary.Set("name", "利昂內爾 梅西")

	o1 := orderedmap.New()
	o1.Set("originator", originator)
	o1.Set("beneficiary", beneficiary)
	o1.Set("signature", "70be6318f31204c9fe28e0b30dabff02b2909105bd0a97d094d6ed5d497461077afa65b440afe37e238172a88b05f3240d0eabb09d008fa25ac0224a507c56b5")

	o2 := orderedmap.New()
	o2.Set("username", "kunming")
	o2.Set("password", 1234)
	o2.Set("signature", "3bec3a43f9b5679647b0da870fb6c955488a47e7eb44f285491ebcf084ec69ca4170c8eb59f1f22002d5128c09affbadc8a56092c234dbfcedf916bf9dad17dc")
	o2.Set("abc", 1.234)

	//sequence is important
	o3 := orderedmap.New()
	o3.Set("username", "kunming")
	o3.Set("password", 1234)
	o3.Set("abc", 1.234)
	o3.Set("signature", "3bec3a43f9b5679647b0da870fb6c955488a47e7eb44f285491ebcf084ec69ca4170c8eb59f1f22002d5128c09affbadc8a56092c234dbfcedf916bf9dad17dc")

	var tests = []struct {
		input    *orderedmap.OrderedMap
		expected bool
	}{
		{o, true},
		{o1, true},
		{o2, true},
		{o3, false},
	}

	for _, test := range tests {
		valid, _ := Verify(test.input, fakePublicKey)
		assert.Equal(t, valid, test.expected, "should be equal")
	}
}

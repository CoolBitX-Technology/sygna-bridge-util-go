package crypto

import (
	"encoding/json"
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
)

const fakePrivateKey = "ba4523e5091939113423a709b5924708af30fc5a958ac71f48eb030b84494702"
const fakePublicKey = "04c1a0d4269ce2b0e1dab89e8defbfc9c0c780e6b769f1dba7cbc3531c8167ae7f0b49b1a36d574fd0cbb353f5d31152110daa541213cf0919c1be708a112163e3"

func TestEncrypt(t *testing.T) {

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
	o2.Set("signature", "abcdef")
	o2.Set("abc", 1.234)

	s3 := "abcdefghijk"

	var tests = []struct {
		input interface{}
	}{
		{o},
		{o1},
		{o2},
		{s3},
	}

	for _, test := range tests {
		var bInput []byte
		switch v := test.input.(type) {
		default:
			t.Fatalf("unexpected type %T", v)
		case *orderedmap.OrderedMap:
			bInput, _ = json.Marshal(test.input)
		case string:
			bInput = []byte(test.input.(string))
		}
		ciphertext, _ := Encrypt(bInput, fakePublicKey)

		decrypted, _ := Decrypt(ciphertext, fakePrivateKey)
		var bDecrypted []byte
		switch v := decrypted.(type) {
		default:
			t.Fatalf("unexpected type %T", v)
		case *orderedmap.OrderedMap:
			bDecrypted, _ = json.Marshal(decrypted)
		case string:
			bDecrypted = []byte(decrypted.(string))
		}
		assert.Equal(t, string(bInput), string(bDecrypted), "should be equal")
	}
}

func TestDecrypt(t *testing.T) {

	o := orderedmap.New()
	o.UnmarshalJSON([]byte(`{"transfer_id":"2535e84afdee114b090e0d6b65772b600c852618d50b38770f653d478490d876","permission_status":"ACCEPTED","signature":"01a62295e3e52a824c92427db26d13f21feb01a4fa3f99ead0de2baecbd3526d52597e27e3df2085bd191d6da1793770459e32f46c943ff18125af9dc409907f"}`))

	// {
	// 	currency_id: 'sygna:0x80000090',
	// 	currency_name: 'XRP',
	// 	currency_symbol: 'XRP',
	// 	is_active: true,
	// 	addr_extra_info: ['tag'],
	// }
	o1 := orderedmap.New()
	o1.Set("currency_id", "sygna:0x80000090")
	o1.Set("currency_name", "XRP")
	o1.Set("currency_symbol", "XRP")
	o1.Set("is_active", true)
	o1.Set("addr_extra_info", []string{"tag"})

	o2 := orderedmap.New()
	o2.Set("signature", "簽章")
	o2.Set("abc", 3.14)
	o2.Set("username", "kunming")
	o2.Set("password", 567890)

	s3 := "zxcvvbjgiyi5/喬丹"

	//encrypt data by javascript util
	var tests = []struct {
		input    string
		expected interface{}
	}{
		{"044fd221098ffae56ea3da9b72312b4d202290fb4d00eb6c021f3bee8a186868d71e95bb309cf8d353036d766ce8ae0b7582690260cc7623bbc62cae5a17cd3741f5b39864ccd7372ff80f397758f6a1b456cec2b66b71ce4c6063de2a0768c843788cf94c61b0242acff6332f80cf829505c2208e00ffb24962fd1dda78751f0312ab07614e5891e8d419db18785f73944bc736d85460d9d5f7fdc8c67d139581ee138a9b702af896f048cc746dd975f24456c4da87cdf16426724b85719ac6b1ebf20581464e4a1fe8d9229c7843d96969920f5b1696e5b300f681963bc90c39dedc95d4fd1d0483dcb40a5737e2435ae15059c559b5275dbde255333e0729f771a5d6abc226f31db8f88d68e25f82eeb53b72f5b652c652bf45d98d3adca102fabb5c830808b00b05dff5db4604b9f38a60ebc98e2034ce1ebf8743f4bc84eff1ee376e58da19acdac5c9889478205a89bc0aa921b0574280f943e7f4ec755d90281070", o},
		{"0446e24626005440d9f22914c04213dcfb55e8196c8e8045fe39544952e212f50bf4eabe0f84814a70db980f0c3af1fc552b48762cb607b34905fa01f258d0176f11e690d1393d555b0a1ed99d11afd908f3cf4b0f06e41052274b4533ca9598ff5fe20602aabf329e965aab86e7899e9ebd16a88f1b808c508c074f90d6ee1d26b9aa4f594e765abccacca456d71445c45d0b260742b618c203805b985e28f478744714010caf2dd02f55bb687420051d03ccfd39096493641cf4abb2ec412397ac6731d0c82424eb32bf0361cdf087044b0f8a67", o1},
		{"0485efc578145cef06674c8670e51af8d74bbd880afa55350ea1c61e4fbc816f0d45fcf11290a402a5bab5424efbe4b9e5118210fcdae4638983e501c2a2fa23bdbe044df6a112f72401d98aeb4df70aa5e84dfd3be6820e11c55259a2712365beb17b66a6d0770c61d4ed4419448d27d82e59b8673e4004e80874edc8a4a9c9a894809e9f0dedb8aa112f21f4de98ffe06036b618ab6270f9b8672efd06812277e6611cb3", o2},
		{"045831322455392ceeab6a0ca1898437aeabf6f72096767de90dc062d2b2ab91281baf9a898c72eb1248cd7c4a3781cef74b79779f6bdfae5952d3013d913a0fc9a4a7983f160247595ad097f8a80531bab06301a9330311d8dc71542747bb3ab594aa9d0b31816272184b31064b1f989548ef30fa", s3},
	}

	for _, test := range tests {
		decrypted, _ := Decrypt(test.input, fakePrivateKey)
		var bDecrypted []byte
		switch v := decrypted.(type) {
		default:
			t.Fatalf("unexpected type %T", v)
		case *orderedmap.OrderedMap:
			bDecrypted, _ = json.Marshal(decrypted)
		case string:
			bDecrypted = []byte(decrypted.(string))
		}

		var bExpected []byte
		switch v := test.expected.(type) {
		default:
			t.Fatalf("unexpected type %T", v)
		case *orderedmap.OrderedMap:
			bExpected, _ = json.Marshal(test.expected)
		case string:
			bExpected = []byte(test.expected.(string))
		}

		assert.Equal(t, string(bDecrypted), string(bExpected), "should be equal")
	}
}

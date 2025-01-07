package example

import (
	"log"

	bridgeutil "github.com/CoolBitX-Technology/sygna-bridge-util-go"
	"github.com/iancoleman/orderedmap"
)

const domain = bridgeutil.SygnaBridgeAPITestDomain

const (
	originatorAPIKey     = "{{originatorAPIKey}}"
	originatorPrivatekey = "{{originatorPrivatekey}}"
	originatorPublicKey  = "{{originatorPublicKey}}"
)

const (
	beneficiaryAPIKey     = "{{beneficiaryAPIKey}}"
	beneficiaryPrivatekey = "{{beneficiaryPrivatekey}}"
	beneficiaryPublicKey  = "{{beneficiaryPublicKey}}"
)

const sensitiveData = `
{
	"originator": {
		"originator_persons": [
			{
				"natural_person": {
					"name": {
						"name_identifiers": [
							{
								"primary_identifier": "Wu Xinli",
								"name_identifier_type": "LEGL"
							}
						]
					},
					"national_identification": {
						"national_identifier": "446005",
						"national_identifier_type": "RAID",
						"registration_authority": "RA000553"
					},
					"country_of_residence": "TZ"
				}
			}
		],
		"account_numbers": [
			"r3kmLJN5D28dHuH8vZNUZpMC43pEHpaocV"
		]
	},
	"beneficiary": {
		"beneficiary_persons": [
			{
				"legal_person": {
					"name": {
						"name_identifiers": [
							{
								"legal_person_name": "ABC Limited",
								"legal_person_name_identifier_type": "LEGL"
							}
						]
					}
				}
			}
		],
		"account_numbers": [
			"rAPERVgXZavGgiGv6xBgtiZurirW2yAmY"
		]
	}
}`

const OtherCDDInfo = "[National ID] Passport number 123456789"

func encryptAndDecrypt() {
	ciphertext, err := bridgeutil.EncryptString(sensitiveData, originatorPublicKey)
	if err != nil {
		panic(err)
	}
	log.Printf("plaintext encrypted: %v\n", ciphertext)

	plaintext, err := bridgeutil.Decrypt(ciphertext, originatorPrivatekey)
	if err != nil {
		panic(err)
	}
	strPlaintext, _ := bridgeutil.OrderedMapToString(plaintext.(*orderedmap.OrderedMap))
	log.Printf("plaintext decrypted: %v\n", strPlaintext)
}

func signAndVerify() {

	o := orderedmap.New()
	o.Set("transfer_id", "b97903fd68fcff05cfe035482bc3cf7fd934505b4e0644e612087dca4bae37e4")
	o.Set("txid", "6f721fba0d405df21fb27dd76cfe2b548907f3881c5625b9cfe624c15c3178ae")

	err := bridgeutil.Sign(o, originatorPrivatekey)
	if err != nil {
		panic(err)
	}
	strMap, _ := bridgeutil.OrderedMapToString(o)
	log.Printf("signed data: %v\n", strMap)

	valid, err := bridgeutil.Verify(o, originatorPublicKey)
	if err != nil {
		panic(err)
	}
	log.Printf("valid: %v\n", valid)
}

func getVASP() {
	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.GetVASP(true)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response...)
	log.Printf("GetVASP response: %v\n", strResponse)
}

func getVASPPublicKey() {
	targetVASPCode := "VASPUSNY1"

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.GetVASPPublicKey(targetVASPCode, true)
	if err != nil {
		panic(err)
	}
	log.Printf("GetVASPPublicKey response: %v\n", response)
}

func getStatus() {
	transferID := "077fdce779ce0d4eda296eb34759db6b361a85abac053bfa63cea411ef85bf44"

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.GetStatus(transferID)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("response: %v\n", strResponse)
}

func getCurrencies() {
	queryParam := orderedmap.New()
	queryParam.Set("currency_id", "sygna:0x80000090")

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.GetCurrencies(queryParam)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response...)
	log.Printf("GetCurrencies response: %v\n", strResponse)
}

func postPermissionRequest() {
	ciphertext, err := bridgeutil.EncryptString(sensitiveData, beneficiaryPublicKey)
	if err != nil {
		panic(err)
	}

	originatorAddr := orderedmap.New()
	originatorAddr.Set("address", "r3kmLJN5D28dHuH8vZNUZpMC43pEHpaocV")

	originatorAddrs := make([]*orderedmap.OrderedMap, 1)
	originatorAddrs[0] = originatorAddr

	originatorVASP := orderedmap.New()
	originatorVASP.Set("vasp_code", "VASPUSNY1")
	originatorVASP.Set("addrs", originatorAddrs)

	beneficiaryAddrInfo := orderedmap.New()
	beneficiaryAddrInfo.Set("tag", "abc")

	beneficiaryAddrInfos := make([]*orderedmap.OrderedMap, 1)
	beneficiaryAddrInfos[0] = beneficiaryAddrInfo

	beneficiaryAddr := orderedmap.New()
	beneficiaryAddr.Set("address", "rAPERVgXZavGgiGv6xBgtiZurirW2yAmY")
	beneficiaryAddr.Set("addr_extra_info", beneficiaryAddrInfos)

	beneficiaryAddrs := make([]*orderedmap.OrderedMap, 1)
	beneficiaryAddrs[0] = beneficiaryAddr

	beneficiaryVASP := orderedmap.New()
	beneficiaryVASP.Set("vasp_code", "VASPUSNY2")
	beneficiaryVASP.Set("addrs", beneficiaryAddrs)

	transaction := orderedmap.New()
	transaction.Set("originator_vasp", originatorVASP)
	transaction.Set("beneficiary_vasp", beneficiaryVASP)
	transaction.Set("currency_id", "sygna:0x80000090")
	transaction.Set("amount", "4.51120135938784")

	permissionRequestData := orderedmap.New()
	permissionRequestData.Set("private_info", ciphertext)
	permissionRequestData.Set("transaction", transaction)
	permissionRequestData.Set("data_dt", "2020-07-13T05:56:53.088Z")

	err = bridgeutil.Sign(permissionRequestData, originatorPrivatekey)
	if err != nil {
		panic(err)
	}

	callbackData := orderedmap.New()
	callbackData.Set("callback_url", "https://facb1c03d3dae42f07008d0c42979623.m.pipedream.net")
	err = bridgeutil.Sign(callbackData, originatorPrivatekey)
	if err != nil {
		panic(err)
	}

	postPermissionRequestData := orderedmap.New()
	postPermissionRequestData.Set("data", permissionRequestData)
	postPermissionRequestData.Set("callback", callbackData)

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostPermissionRequest(postPermissionRequestData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostPermissionRequest response: %v\n", strResponse)
}

func postPermission() {
	transferID := "e8867006137a94f13656198be8fa720cf5b822f70adf6fe71a25538d0b2c230e"
	postPermissionData := orderedmap.New()
	postPermissionData.Set("transfer_id", transferID)
	postPermissionData.Set("permission_status", bridgeutil.PermissionStatusAccepted)

	err := bridgeutil.Sign(postPermissionData, beneficiaryPrivatekey)
	if err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    beneficiaryAPIKey,
	}
	response, err := api.PostPermission(postPermissionData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostPermission response: %v\n", strResponse)
}

func postTransactionID() {
	transferID := "e8867006137a94f13656198be8fa720cf5b822f70adf6fe71a25538d0b2c230e"
	txID := "6f721fba0d405df21fb27dd76cfe2b548907f3881c5625b9cfe624c15c3178ae"
	postTxIDData := orderedmap.New()
	postTxIDData.Set("transfer_id", transferID)
	postTxIDData.Set("txid", txID)

	err := bridgeutil.Sign(postTxIDData, originatorPrivatekey)
	if err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostTransactionID(postTxIDData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostTxId response: %v\n", strResponse)
}

func postBeneficiaryEndpointURL() {
	postBeneficiaryEndpointURLData := orderedmap.New()
	postBeneficiaryEndpointURLData.Set("vasp_code", "VASPUSNY2")
	postBeneficiaryEndpointURLData.Set("callback_permission_request_url", "https://google.com")
	postBeneficiaryEndpointURLData.Set("callback_txid_url", "https://stackoverflow.com")
	postBeneficiaryEndpointURLData.Set("callback_validate_addr_url", "https://github.com")

	err := bridgeutil.Sign(postBeneficiaryEndpointURLData, beneficiaryPrivatekey)
	if err != nil {
		panic(err)
	}
	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    beneficiaryAPIKey,
	}
	response, err := api.PostBeneficiaryEndpointURL(postBeneficiaryEndpointURLData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("postBeneficiaryEndpointURL response: %v\n", strResponse)
}

func postRetry() {
	postRetryData := orderedmap.New()
	postRetryData.Set("vasp_code", "VASPUSNY2")

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    beneficiaryAPIKey,
	}
	response, err := api.PostRetry(postRetryData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("postRetry response: %v\n", strResponse)
}

func postWalletAddressFilter() {
	postWalletAddressFilterData := orderedmap.New()
	postWalletAddressFilterData.Set("currency_id", "sygna:0x80000000")
	postWalletAddressFilterData.Set("addrs", []string{
		"14YjfWmQGTqLBPkG26qG81MKnZwV8z7wh1",
		"bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
	})

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	ignoreKYT := false
	response, err := api.PostWalletAddressFilter(postWalletAddressFilterData, ignoreKYT)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response...)
	log.Printf("postWalletAddressFilter response: %v\n", strResponse)
}

func postServerStatus() {
	postServerStatusData := orderedmap.New()
	postServerStatusData.Set("vasp_code", "VASPUSNY2")
	postServerStatusData.Set("status", "maintaining")
	postServerStatusData.Set("started_at", 1724808400000)
	postServerStatusData.Set("ended_at", 1724808400000)

	if err := bridgeutil.Sign(postServerStatusData, beneficiaryPrivatekey); err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    beneficiaryAPIKey,
	}
	response, err := api.PostServerStatus(postServerStatusData)
	if err != nil {
		panic(err)
	}

	strResponse, err := bridgeutil.OrderedMapToString(response)
	if err != nil {
		panic(err)
	}
	log.Printf("postServerStatus response: %v\n", strResponse)
}

func postTransactionCDD() {
	ciphertext, err := bridgeutil.EncryptString(OtherCDDInfo, beneficiaryPublicKey)
	if err != nil {
		panic(err)
	}
	postTransactionCDDData := orderedmap.New()
	postTransactionCDDData.Set("transfer_id", "e8867006137a94f13656198be8fa720cf5b822f70adf6fe71a25538d0b2c230e")
	postTransactionCDDData.Set("other_cdd_info", ciphertext)

	err = bridgeutil.Sign(postTransactionCDDData, originatorPrivatekey)
	if err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostTransactionCDD(postTransactionCDDData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostTransactionCDD response: %v\n", strResponse)
}

func postTransactionCDDRequest() {
	requestCDDData := orderedmap.New()
	address := []string{"address_line", "country"}
	requestCDDData.Set("geographic_address", address)

	cddRequestBody := orderedmap.New()
	cddRequestBody.Set("transfer_id", "463412c611e53aefcd016a7efda1328e7a1067ab659a5a45a60c1c7023ceb133")
	cddRequestBody.Set("request_cdd_data", requestCDDData)
	if err := bridgeutil.Sign(cddRequestBody, beneficiaryPrivatekey); err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostTransactionCDDRequest(cddRequestBody)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostTransactionCDDRequest response: %v\n", strResponse)
}

func postVASPBeneficiaryCheckingRule() {
	checkingRuleData := orderedmap.New()

	// natural person
	naturalPerson := orderedmap.New()
	naturalPerson.Set("country_of_residence", true)
	naturalPerson.Set("customer_identification", false)

	dateAndPlaceOfBirth := orderedmap.New()
	dateAndPlaceOfBirth.Set("date_of_birth", true)
	dateAndPlaceOfBirth.Set("place_of_birth", true)

	naturalPerson.Set("date_and_place_of_birth", dateAndPlaceOfBirth)

	checkingRuleData.Set("natural_person", naturalPerson)

	// legal person
	legalPerson := orderedmap.New()
	legalPerson.Set("country_of_registration", false)
	legalPerson.Set("customer_identification", true)

	nameIdentifiers := orderedmap.New()
	nameIdentifiers.Set("legal_person_name_identifier_type", true)
	nameIdentifiers.Set("legal_person_name", true)

	legalPerson.Set("name_identifiers", nameIdentifiers)

	checkingRuleData.Set("legal_person", legalPerson)

	if err := bridgeutil.Sign(checkingRuleData, beneficiaryPrivatekey); err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostVASPBeneficiaryCheckingRule(checkingRuleData)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostVASPBeneficiaryCheckingRule response: %v\n", strResponse)
}

func postTransactionCancel() {
	cancelRequestBody := orderedmap.New()
	cancelRequestBody.Set("transfer_id", "463412c611e53aefcd016a7efda1328e7a1067ab659a5a45a60c1c7023ceb133")
	if err := bridgeutil.Sign(cancelRequestBody, originatorPrivatekey); err != nil {
		panic(err)
	}

	api := &bridgeutil.BridgeAPI{
		APIDomain: domain,
		APIKey:    originatorAPIKey,
	}
	response, err := api.PostTransactionCancel(cancelRequestBody)
	if err != nil {
		panic(err)
	}
	strResponse, _ := bridgeutil.OrderedMapToString(response)
	log.Printf("PostTransactionCancel response: %v\n", strResponse)
}

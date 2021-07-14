# Golang Sygna Bridge Util

This is a Golang library to help you build servers/services within Sygna Bridge Ecosystem. For more detail information, please see [Sygna Bridge](https://www.sygna.io/).

## Installation

```shell
go get github.com/CoolBitX-Technology/sygna-bridge-util-go
```

## Import

```golang

import (
	bridgeutil "github.com/CoolBitX-Technology/sygna-bridge-util-go"
)
```

## Crypto

Dealing with encrypting, decrypting, signing and verifying in Sygna Bridge.

### ECIES Encrypting an Decrypting

During the communication of VASPs, there are some private information that must be encrypted. We use ECIES(Elliptic Curve Integrated Encryption Scheme) to securely encrypt these private data so that they can only be accessed by the recipient.

We're using [IVMS101 (interVASP Messaging Standard)](https://intervasp.org/) as our private information format.

We also provide [IVMS101 Golang Utility](https://github.com/CoolBitX-Technology/sygna-bridge-ivms-utils/tree/master/golang) to construct data payload.

```golang
sensitiveData := `
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

privateInfo,err := bridgeutil.EncryptString(
  sensitiveData,
  recipientPubKey,
)

decryptedPrivateInfo,err := bridgeutil.Decrypt(
  privateInfo,
  recipientPrivateKey,
)
```

### Sign and Verify

In Sygna Bridge, we use secp256k1 ECDSA over sha256 of utf-8 json string to create signature on every API call. Since you need to provide the identical utf-8 string during verification, the order of key-value pair you put into the object is important.

The following example is the snippet of originator's signing process of `permissionRequest` API call. If you put the key `transaction` before `private_info` in the object, the verification will fail in the central server.

````golang
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

bridgeutil.Sign(permissionRequestData, originatorPrivateKey)

valid, err := bridgeutil.Verify(permissionRequestData, originatorPublicKey)


## API

API calls to communicate with Sygna Bridge server.

We use **baisc auth** with all the API calls. To simplify the process, we provide a API class to deal with authentication and post/ get request format.

```golang
api := &bridgeutil.BridgeAPI{
  APIDomain: domain,
  APIKey:    originatorAPIKey,
}
````

After you create the `BridgeAPI` struct, you can use it to make any API call to communicate with Sygna Bridge central server.

### Get VASP Information

```golang
// Get List of VASPs associated with public keys.
verify := true // set verify to true to verify the signature attached with api response automatically.
vasps,err := api.GetVASP(verify)

// Or call use GetVASPPublicKey() to directly get public key for a specific VASP.
publicKey := api.GetVASPPublicKey("VASPUSNY1", verify)
```

### For Originator

There are two API calls from **transaction originator** to Sygna Bridge Server defined in the protocol, which are `PostPermissionRequest` and `PostTransactionID`.

The full logic of originator would be like the following:

```golang
recipientPublicKey := api.GetVASPPublicKey("VASPUSNY1",verify)
privateIinfo := bridgeutil.Encrypt(
  // from example above
  sensitiveData,
  recipientPublicKey,
)

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


bridgeutil.Sign(
  permissionRequestData,
  senderPrivateKey,
)

callbackURL :"https://81f7d956.ngrok.io/api/v2/originator/transaction/premission"

callbackData := orderedmap.New()
callbackData.Set("callback_url", callbackURL)

bridgeutil.Sign(
  callbackData,
  senderPrivateKey,
)

postPermissionRequestData := orderedmap.New()
postPermissionRequestData.Set("data", permissionRequestData)
postPermissionRequestData.Set("callback", callbackData)

response, err := api.PostPermissionRequest(
  postPermissionRequestData,
)

transferID := response.Get("transfer_id")

// Broadcast your transaction to blockchain after got and api response at your api server.
txid := "1a0c9bef489a136f7e05671f7f7fada2b9d96ac9f44598e1bcaa4779ac564dcd"

// Inform Sygna Bridge that a specific transfer is successfully broadcasted to the blockchain.
txIDData := orderedmap.New()
txIDData.Set("transfer_id", transferID)
txIDData.Set("txid", txid)

bridge.Sign(txIDData, senderPrivateKey)
response, err := api.PostTransactionID(txIDData)
```

### For Beneficiary

There is only one api for Beneficiary VASP to call, which is `PostPermission`. After the beneficiary server confirm their legitimacy of a transfer request, they will sign `{ transfer_id, permission_status }` using `Sign()` function, and send the result with signature to Sygna Bridge Central Server.

```golang
permissionStatus := bridgeutil.PermissionStatusAccepted // or bridgeutil.PermissionStatusRejected

permissionData := orderedmap.New()
permissionData.Set("transfer_id", transferID)
permissionData.Set("permission_status", permissionStatus)

bridgeutil.Sign(
  permissionData,
  beneficiaryPrivateKey,
)
finalResult := api.PostPermission(permissionData)
```

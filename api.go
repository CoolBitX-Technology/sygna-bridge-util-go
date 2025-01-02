package bridgeutil

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iancoleman/orderedmap"
	"github.com/imroc/req"
)

const get = "GET"
const post = "POST"

// BridgeAPI is a convenient struct for using sygna API
type BridgeAPI struct {
	APIDomain string
	APIKey    string
	UserAgent string
}

func isHTTPStatusOK(statusCode int) bool {
	if statusCode >= 200 && statusCode < 300 {
		return true
	}
	return false
}

func parseResponse(r *req.Resp, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	resp := r.Response()
	statusCode := resp.StatusCode

	m := orderedmap.New()
	// if err is nil, the response would be json format
	err = r.ToJSON(m)
	if err == nil {
		if !isHTTPStatusOK(statusCode) {
			_, exist := m.Get("status")
			if !exist {
				m.Set("status", statusCode)
			}
			bMessage, _ := json.Marshal(m)
			return nil, errors.New(string(bMessage))
		}
		return m, nil
	}
	// if err is not nil, the response would be array json format
	var maps []*orderedmap.OrderedMap
	err = r.ToJSON(&maps)
	if err != nil {
		return nil, err
	}
	if !isHTTPStatusOK(statusCode) {
		bMessage, _ := json.Marshal(maps)
		return nil, errors.New(string(bMessage))
	}

	return maps, nil
}

func request(api *BridgeAPI, method, path string, v ...interface{}) (interface{}, error) {
	if api.UserAgent == "" {
		api.UserAgent = "util-go"
	}

	header := req.Header{
		"Content-type": "application/json;",
		"X-Api-Key":    api.APIKey,
		"User-Agent":   api.UserAgent,
	}

	options := make([]interface{}, len(v)+1)
	options[0] = header
	copy(options[1:], v)

	var r *req.Resp
	var err error

	url := api.APIDomain + path
	switch method {
	case get:
		r, err = req.Get(url, options...)
	case post:
		r, err = req.Post(url, options...)
	default:
		panic(errors.New("unsupported method"))
	}
	return parseResponse(r, err)
}

/*
GetVASP Get list of registered VASP associated with publicKey.
Set validate false to disable validating returned vasp list data.
Set isProdEnv true to use SygnaBridgeCentralPubkey to Verify data.

see https://developers.sygna.io/reference#bridgevasp-3
*/
func (api *BridgeAPI) GetVASP(validate bool, isProdEnv ...bool) ([]*orderedmap.OrderedMap, error) {
	response, err := request(api, get, "v2/bridge/vasp")
	if err != nil {
		return nil, err
	}
	vaspData, _ := response.(*orderedmap.OrderedMap).Get("vasp_data")
	mapVASPData := castArrayToOrderedMapArray(vaspData)

	if !validate {
		return mapVASPData, nil
	}

	isProd := false
	if len(isProdEnv) > 0 {
		isProd = isProdEnv[0]
	}

	publicKey := SygnaBridgeTestPubkey
	if isProd {
		publicKey = SygnaBridgeCentralPubkey
	}

	valid, err := Verify(response.(*orderedmap.OrderedMap), publicKey)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("get VASP info error: invalid signature")
	}

	return mapVASPData, nil
}

// GetVASPPublicKey A Wrapper function of GetVASP to return specific VASP's Public Key.
func (api *BridgeAPI) GetVASPPublicKey(targetVASPCode string, validate bool, isProdEnv ...bool) (string, error) {
	response, err := api.GetVASP(validate, isProdEnv...)

	if err != nil {
		return "", err
	}

	for _, v := range response {
		vaspCode, _ := v.Get("vasp_code")
		if vaspCode == targetVASPCode {
			publickKey, _ := v.Get("vasp_pubkey")
			return publickKey.(string), nil
		}
	}
	return "", errors.New("Invalid targetVASPCode")
}

/*
GetStatus Get detail of particular transaction permission request

see https://developers.sygna.io/reference#bridgestatus-3
*/
func (api *BridgeAPI) GetStatus(transferID string) (*orderedmap.OrderedMap, error) {
	param := req.Param{
		"transfer_id": transferID,
	}
	response, err := request(api, get, "v2/bridge/transaction/status", param)

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
GetCurrencies Get supported currencies

see https://developers.sygna.io/reference#bridgecurrencies
*/
func (api *BridgeAPI) GetCurrencies(queryParams *orderedmap.OrderedMap) ([]*orderedmap.OrderedMap, error) {
	param := req.Param{}

	if queryParams != nil {
		for _, k := range queryParams.Keys() {
			v, _ := queryParams.Get(k)
			param[k] = v
		}
	}
	response, err := request(api, get, "v2/bridge/transaction/currencies", param)

	if err != nil {
		return nil, err
	}

	supportedCoins, _ := response.(*orderedmap.OrderedMap).Get("supported_coins")
	mapSupportedCoins := castArrayToOrderedMapArray(supportedCoins)

	return mapSupportedCoins, nil
}

/*
PostBeneficiaryEndpointURL Revise beneficiary endpoint url

see https://developers.sygna.io/reference#bridgebeneficiaryendpointurl
*/
func (api *BridgeAPI) PostBeneficiaryEndpointURL(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/vasp/beneficiary-endpoint-url", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
PostPermissionRequest Should be called by the originator VASP to inform Sygna Bridge about the creation of a compliant transaction.

see https://developers.sygna.io/reference#bridgepermissionrequest-3
*/
func (api *BridgeAPI) PostPermissionRequest(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/permission-request", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
PostPermission Notify Sygna Bridge that you have confirmed specific permission Request

	from other VASP. Should be called by Beneficiary Server

see https://developers.sygna.io/reference#bridgepermission-3
*/
func (api *BridgeAPI) PostPermission(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/permission", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
PostTransactionID Send broadcasted transaction id to Sygna Bridge for purpose of storage.

see https://developers.sygna.io/reference#bridgetransactionid-3
*/
func (api *BridgeAPI) PostTransactionID(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/txid", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
PostRetry Retrieve the lost transfer requests

see https://developers.sygna.io/reference#bridgeretry-3
*/
func (api *BridgeAPI) PostRetry(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/retry", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

func (api *BridgeAPI) PostTransactionCDDRequest(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/cdd-request", req.BodyJSON(param))
	if err != nil {
		return nil, err
	}

	return response.(*orderedmap.OrderedMap), nil
}

func (api *BridgeAPI) PostTransactionCDD(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/transaction/cdd", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

/*
PostWalletAddressFilter Retrieve address information

You should active Blockchain Analytics to retrieve address information which is not belong to a SYGNA VASP.

see https://developers.sygna.io/reference#bridgewallet-address-filter
*/
func (api *BridgeAPI) PostWalletAddressFilter(param *orderedmap.OrderedMap, ignoreKYT ...bool) ([]*orderedmap.OrderedMap, error) {
	q := req.Param{}

	if len(ignoreKYT) > 0 {
		q["ignore_kyt"] = ignoreKYT[0]
	}
	response, err := request(api, post, "v2/bridge/wallet-address-filter", req.BodyJSON(param), q)
	if err != nil {
		return nil, err
	}

	return response.([]*orderedmap.OrderedMap), nil
}

// Get vasp details by vasp code
func (api *BridgeAPI) GetVASPDetails(vaspCode string, validate bool, isProdEnv ...bool) (*orderedmap.OrderedMap, error) {
	response, err := request(api, get, fmt.Sprintf("v2/bridge/vasp/detail/%s", vaspCode))
	if err != nil {
		return nil, err
	}
	vaspData, _ := response.(*orderedmap.OrderedMap).Get("vasp_data")
	VASPDataObject := castObjectToOrderedMapObject(vaspData)

	if !validate {
		return VASPDataObject, nil
	}

	isProd := false
	if len(isProdEnv) > 0 {
		isProd = isProdEnv[0]
	}

	publicKey := SygnaBridgeTestPubkey
	if isProd {
		publicKey = SygnaBridgeCentralPubkey
	}

	valid, err := Verify(response.(*orderedmap.OrderedMap), publicKey)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("get VASP info error: invalid signature")
	}

	return VASPDataObject, nil
}

// GetVASPUsage Get VASP usage by timestamp
func (api *BridgeAPI) GetVASPUsages(startAt, endAt int64, validate bool, isProdEnv ...bool) ([]*orderedmap.OrderedMap, error) {
	param := req.Param{
		"start_at": startAt,
		"end_at":   endAt,
	}
	response, err := request(api, get, "v2/bridge/vasp/usage", param)

	if err != nil {
		return nil, err
	}

	usageData, _ := response.(*orderedmap.OrderedMap).Get("data")
	usageDataObject := castArrayToOrderedMapArray(usageData)

	if !validate {
		return usageDataObject, nil
	}

	isProd := false
	if len(isProdEnv) > 0 {
		isProd = isProdEnv[0]
	}

	publicKey := SygnaBridgeTestPubkey
	if isProd {
		publicKey = SygnaBridgeCentralPubkey
	}

	valid, err := Verify(response.(*orderedmap.OrderedMap), publicKey)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("get VASP usages error: invalid signature")
	}

	return usageDataObject, nil
}

// PostServerStatus declares that the VASP’s server is currently in maintenance.
func (api *BridgeAPI) PostServerStatus(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/vasp/server-status", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

func (api *BridgeAPI) PostVASPBeneficiaryCheckingRule(param *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	response, err := request(api, post, "v2/bridge/vasp/beneficiary-checking-rule", req.BodyJSON(param))

	if err != nil {
		return nil, err
	}
	return response.(*orderedmap.OrderedMap), nil
}

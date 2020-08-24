package bridgeutil

//SygnaBridgeCentralPubkey public key for SygnaBridgeApiDomain
const SygnaBridgeCentralPubkey = "047b04ca933c0fccb7094af06bafb77e0fdd9264b45243cba0b72cd8f1bc8fc4e7454902d4bb6bad8ed4bc4dfae102858b6a7649e4febca0c5b266566aa4e59f12"

//SygnaBridgeTestPubkey public key for SygnaBridgeApiTestDomain
const SygnaBridgeTestPubkey = "04a6936f2bc43773cb4874980518b3f681c004464d167aebdc9e305e10d6fb6cdacb27a22812453e6c51ceabff5b1e2d2196d81a8d3e8e71e907948b01a7ea9ac8"

const (
	//SygnaBridgeAPIDomain production domain
	SygnaBridgeAPIDomain = "https://api.sygna.io/"
	//SygnaBridgeAPITestDomain test domain
	SygnaBridgeAPITestDomain = "https://test-api.sygna.io/"
)

const (
	//PermissionStatusAccepted accept transfer from originator vasp
	PermissionStatusAccepted = "ACCEPTED"
	//PermissionStatusRejected reject transfer from originator vasp
	PermissionStatusRejected = "REJECTED"
)

const (
	//RejectCodeBVRC001 When the originator VASP is going to send an unsupported currency to you.
	RejectCodeBVRC001 = "BVRC001"
	//RejectCodeBVRC002 When your service is under downtime or you are unable to reply with the request.
	RejectCodeBVRC002 = "BVRC002"
	//RejectCodeBVRC003 When your customer is not able to receive more transaction inflows.
	RejectCodeBVRC003 = "BVRC003"
	//RejectCodeBVRC004 When your customer fails your internal compliance check or the person is listed in your blacklist.
	RejectCodeBVRC004 = "BVRC004"
	//RejectCodeBVRC005 When private_info can not be decoded
	RejectCodeBVRC005 = "BVRC005"
	//RejectCodeBVRC006 When private_info can be decoded but the format is wrong
	RejectCodeBVRC006 = "BVRC006"
	//RejectCodeBVRC007 Beneficiary name is not matched with the name in the beneficiary VASP database.
	RejectCodeBVRC007 = "BVRC007"
	//RejectCodeBVRC999 When the reject reason is not included in the above options, please put your customized message in the reject_message.
	RejectCodeBVRC999 = "BVRC999"
)

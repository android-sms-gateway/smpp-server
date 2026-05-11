package sessions

type state int

const (
	StateOpen state = iota
	StateBoundTX
	StateBoundRX
	StateBoundTRX
)

// SMPP Error Codes - Standard SMPP v3.4 error codes.
const (
	ErrNoError               uint32 = 0x00000000 // ESME_ROK - No error
	ErrInvalidMsgLen         uint32 = 0x00000001 // ESME_RINVMSGLEN - Message Length is invalid
	ErrInvalidCmdLen         uint32 = 0x00000002 // ESME_RINVCMDLEN - Command Length is invalid (PDU length is too short or too large)
	ErrInvalidCmdID          uint32 = 0x00000003 // ESME_RINVCMDID - Invalid Command ID
	ErrInvalidBindStatus     uint32 = 0x00000004 // ESME_RINVBNDSTS - Incorrect BIND Status for given command (PDU sent in wrong session state)
	ErrAlreadyBound          uint32 = 0x00000005 // ESME_RALYBND - ESME Already in Bound State (bind request in already bound session)
	ErrInvalidPriority       uint32 = 0x00000006 // ESME_RINVPRTFLG - Invalid Priority Flag
	ErrInvalidRegDelivFlag   uint32 = 0x00000007 // ESME_RINVREGDLVFLG - Invalid Registered Delivery Flag
	ErrSystemError           uint32 = 0x00000008 // ESME_RSYSERR - System Error (MC system error, MC unavailable)
	ErrInvalidSrcAddr        uint32 = 0x0000000A // ESME_RINVSRCADR - Invalid Source Address
	ErrInvalidDstAddr        uint32 = 0x0000000B // ESME_RINVDSTADR - Invalid Destination Address
	ErrInvalidMsgID          uint32 = 0x0000000C // ESME_RINVMSGID - Message ID is invalid
	ErrBindFail              uint32 = 0x0000000D // ESME_RBINDFAIL - Bind Failed (generic bind failure)
	ErrInvalidPassword       uint32 = 0x0000000E // ESME_RINVPASWD - Invalid Password (password field length invalid)
	ErrInvalidSystemID       uint32 = 0x0000000F // ESME_RINVSYSID - Invalid System ID (system_id field length invalid)
	ErrCancelFail            uint32 = 0x00000011 // ESME_RCANCELFAIL - Cancel SM Failed
	ErrReplaceFail           uint32 = 0x00000013 // ESME_RREPLACEFAIL - Replace SM Failed
	ErrMsgQFull              uint32 = 0x00000014 // ESME_RMSGQFUL - Message Queue Full
	ErrInvalidServiceType    uint32 = 0x00000015 // ESME_RINVSERTYP - Invalid Service Type
	ErrInvalidNumDests       uint32 = 0x00000033 // ESME_RINVNUMDESTS - Invalid number of destinations
	ErrInvalidDLName         uint32 = 0x00000034 // ESME_RINVDLNAME - Invalid Distribution List name
	ErrInvalidDestFlag       uint32 = 0x00000040 // ESME_RINVDESTFLAG - Destination flag is invalid (submit_multi)
	ErrInvalidSubRep         uint32 = 0x00000042 // ESME_RINVSUBREP - Submit w/replace functionality not supported
	ErrInvalidESMClass       uint32 = 0x00000043 // ESME_RINVESMCLASS - Invalid esm_class field data
	ErrCantSubmitToDL        uint32 = 0x00000044 // ESME_RCNTSUBDL - Cannot Submit to Distribution List
	ErrSubmitFail            uint32 = 0x00000045 // ESME_RSUBMITFAIL - submit_sm, data_sm or submit_multi failed
	ErrInvalidSrcTON         uint32 = 0x00000048 // ESME_RINVSRCTON - Invalid Source address TON
	ErrInvalidSrcNPI         uint32 = 0x00000049 // ESME_RINVSRCNPI - Invalid Source address NPI
	ErrInvalidDstTON         uint32 = 0x00000050 // ESME_RINVDSTTON - Invalid Destination address TON
	ErrInvalidDstNPI         uint32 = 0x00000051 // ESME_RINVDSTNPI - Invalid Destination address NPI
	ErrInvalidSystemType     uint32 = 0x00000053 // ESME_RINVSYSTYP - Invalid system_type field
	ErrInvalidRepFlag        uint32 = 0x00000054 // ESME_RINVREPFLAG - Invalid replace_if_present flag
	ErrInvalidNumMsgs        uint32 = 0x00000055 // ESME_RINVNUMMSGS - Invalid number of messages
	ErrThrottled             uint32 = 0x00000058 // ESME_RTHROTTLED - Throttling error (ESME exceeded message limits)
	ErrInvalidSched          uint32 = 0x00000061 // ESME_RINVSCHED - Invalid Scheduled Delivery Time
	ErrInvalidExpiry         uint32 = 0x00000062 // ESME_RINVEXPIRY - Invalid message validity period (Expiry time)
	ErrInvalidDfltMsgID      uint32 = 0x00000063 // ESME_RINVDFTMSGID - Predefined Message ID is Invalid
	ErrRxTempAppError        uint32 = 0x00000064 // ESME_RX_T_APPN - ESME Receiver Temporary App Error Code
	ErrRxPermAppError        uint32 = 0x00000065 // ESME_RX_P_APPN - ESME Receiver Permanent App Error Code
	ErrRxRejectMsg           uint32 = 0x00000066 // ESME_RX_R_APPN - ESME Receiver Reject Message Error Code
	ErrQueryFail             uint32 = 0x00000067 // ESME_RQUERYFAIL - query_sm request failed
	ErrInvalidTLVStream      uint32 = 0x000000C0 // ESME_RINVTLVSTREAM - Error in optional part of PDU Body (TLV decoding error)
	ErrTLVNotAllowed         uint32 = 0x000000C1 // ESME_RTLVNOTALLWD - TLV not allowed
	ErrInvalidTLVLen         uint32 = 0x000000C2 // ESME_RINVTLVLEN - Invalid Parameter Length (TLV length invalid)
	ErrMissingTLV            uint32 = 0x000000C3 // ESME_RMISSINGTLV - Expected TLV missing
	ErrInvalidTLVVal         uint32 = 0x000000C4 // ESME_RINVTLVVAL - Invalid TLV Value
	ErrDeliveryFailure       uint32 = 0x000000FE // ESME_RDELIVERYFAILURE - Transaction Delivery Failure
	ErrUnknown               uint32 = 0x000000FF // ESME_RUNKNOWNERR - Unknown Error
	ErrServiceTypeUnauth     uint32 = 0x00000100 // ESME_RSERTYPUNAUTH - ESME Not authorised to use specified service_type
	ErrProhibited            uint32 = 0x00000101 // ESME_RPROHIBITED - ESME Prohibited from using specified operation
	ErrServiceTypeUnavail    uint32 = 0x00000102 // ESME_RSERTYPUNAVAIL - Specified service_type is unavailable
	ErrServiceTypeDenied     uint32 = 0x00000103 // ESME_RSERTYPDENIED - Specified service_type is denied
	ErrInvalidDCS            uint32 = 0x00000104 // ESME_RINVDCS - Invalid Data Coding Scheme
	ErrInvalidSrcAddrSubunit uint32 = 0x00000105 // ESME_RINVSRCADDRSUBUNIT - Source Address Subunit is invalid
	ErrInvalidDstAddrSubunit uint32 = 0x00000106 // ESME_RINVDSTADDRSUBUNIT - Destination Address Subunit is invalid
	ErrInvalidBCastFreqInt   uint32 = 0x00000107 // ESME_RINVBCASTFREQINT - Broadcast Frequency Interval is invalid
	ErrInvalidBCastAliasName uint32 = 0x00000108 // ESME_RINVBCASTALIAS_NAME - Broadcast Alias Name is invalid
	ErrInvalidBCastAreaFmt   uint32 = 0x00000109 // ESME_RINVBCASTAREAFMT - Broadcast Area Format is invalid
	ErrInvalidNumBCastAreas  uint32 = 0x0000010A // ESME_RINVNUMBCAST_AREAS - Number of Broadcast Areas is invalid
	ErrInvalidBCastCntType   uint32 = 0x0000010B // ESME_RINVBCASTCNTTYPE - Broadcast Content Type is invalid
	ErrInvalidBCastMsgClass  uint32 = 0x0000010C // ESME_RINVBCASTMSGCLASS - Broadcast Message Class is invalid
	ErrBCastFail             uint32 = 0x0000010D // ESME_RBCASTFAIL - broadcast_sm operation failed
	ErrBCastQueryFail        uint32 = 0x0000010E // ESME_RBCASTQUERYFAIL - query_broadcast_sm operation failed
	ErrBCastCancelFail       uint32 = 0x0000010F // ESME_RBCASTCANCELFAIL - cancel_broadcast_sm operation failed
	ErrInvalidBCastRep       uint32 = 0x00000110 // ESME_RINVBCAST_REP - Number of Repeated Broadcasts is invalid
	ErrInvalidBCastSvcGrp    uint32 = 0x00000111 // ESME_RINVBCASTSRVGRP - Broadcast Service Group is invalid
	ErrInvalidBCastChanInd   uint32 = 0x00000112 // ESME_RINVBCASTCHANIND - Broadcast Channel Indicator is invalid
)

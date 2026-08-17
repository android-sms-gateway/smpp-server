package sessions

import "github.com/android-sms-gateway/client-go/rest"

// bindStatusForPingError maps a gateway REST API ping failure to an SMPP
// bind error code. Auth rejections from the gateway (HTTP 4xx client errors)
// surface as ESME_RINVPASWD so ESME clients can distinguish bad credentials
// from upstream/network failures, which stay ESME_RBINDFAIL.
func bindStatusForPingError(err error) uint32 {
	if err == nil {
		return ErrNoError
	}
	if rest.IsClientError(err) {
		return ErrInvalidPassword
	}
	return ErrBindFail
}

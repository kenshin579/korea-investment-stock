package kis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		RtCode:  "1",
		MsgCode: "MCA00001",
		Message: "잘못된 요청",
		TrID:    "FHKST01010100",
	}
	assert.Contains(t, err.Error(), "MCA00001")
	assert.Contains(t, err.Error(), "잘못된 요청")
}

func TestAPIError_ErrorsAs(t *testing.T) {
	var err error = &APIError{RtCode: "1", MsgCode: "MCA00001", Message: "msg"}
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "MCA00001", apiErr.MsgCode)
}

func TestSentinelErrors(t *testing.T) {
	assert.Equal(t, "kis: token expired", ErrTokenExpired.Error())
	assert.Equal(t, "kis: rate limited", ErrRateLimited.Error())
	assert.Equal(t, "kis: resource not found", ErrNotFound.Error())
	assert.Equal(t, "kis: unauthorized", ErrUnauthorized.Error())
}

package knowledgesdk

import toolsdk "github.com/domainry/domainry-tools-sdk"

// Authority is the authenticated product identity projected into Knowledge.
// Identity owns authentication and grants; Knowledge owns resource decisions.
type Authority = toolsdk.Authority

// Error is the shared structured Knowledge service error.
type Error = toolsdk.Error

type ErrorClass = toolsdk.ErrorClass

const (
	ErrorClassBadRequest  = toolsdk.ErrorClassBadRequest
	ErrorClassForbidden   = toolsdk.ErrorClassForbidden
	ErrorClassNotFound    = toolsdk.ErrorClassNotFound
	ErrorClassConflict    = toolsdk.ErrorClassConflict
	ErrorClassRateLimited = toolsdk.ErrorClassRateLimited
	ErrorClassUnavailable = toolsdk.ErrorClassUnavailable
	ErrorClassInternal    = toolsdk.ErrorClassInternal
)

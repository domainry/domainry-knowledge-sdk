package knowledgesdk

import toolsdk "github.com/domainry/domainry-tools-sdk"

import "slices"

// Authority is the authenticated product identity projected into Knowledge.
// Identity owns authentication and grants; Knowledge owns resource decisions.
//
// Organization scope IDs are live Identity facts. Knowledge persists only a
// stable organization subject ID on a grant; it never copies the organization
// tree or treats these fields as a transferable access bundle.
type Authority struct {
	Known       bool     `json:"known"`
	RuntimeID   string   `json:"runtime_id"`
	WorkspaceID string   `json:"workspace_id"`
	UserID      string   `json:"user_id"`
	RoleKey     string   `json:"role_key,omitempty"`
	OrgID       string   `json:"org_id,omitempty"`
	OrgScopeIDs []string `json:"org_scope_ids,omitempty"`
	AccountType string   `json:"account_type,omitempty"`
}

func AuthorityFromTools(value toolsdk.Authority) Authority {
	return Authority{Known: value.Known, RuntimeID: value.RuntimeID, WorkspaceID: value.WorkspaceID, UserID: value.UserID, RoleKey: value.RoleKey}
}

func (value Authority) ToolsAuthority() toolsdk.Authority {
	return toolsdk.Authority{Known: value.Known, RuntimeID: value.RuntimeID, WorkspaceID: value.WorkspaceID, UserID: value.UserID, RoleKey: value.RoleKey}
}

func (value Authority) Equal(other Authority) bool {
	return value.Known == other.Known && value.RuntimeID == other.RuntimeID && value.WorkspaceID == other.WorkspaceID && value.UserID == other.UserID && value.RoleKey == other.RoleKey && value.OrgID == other.OrgID && slices.Equal(value.OrgScopeIDs, other.OrgScopeIDs) && value.AccountType == other.AccountType
}

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

package knowledgesdk

import (
	"fmt"
	"strings"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const AuthorizationOwner = "module:knowledge"

const (
	ActionLibrariesList    = "knowledge.libraries.list"
	ActionLibrariesCreate  = "knowledge.libraries.create"
	ActionLibrariesGet     = "knowledge.libraries.get"
	ActionLibrariesUpdate  = "knowledge.libraries.update"
	ActionLibrariesArchive = "knowledge.libraries.archive"
	ActionLibrariesRestore = "knowledge.libraries.restore"

	ActionLibraryMembersList   = "knowledge.library_members.list"
	ActionLibraryMembersSearch = "knowledge.library_members.search"
	ActionLibraryMembersSet    = "knowledge.library_members.set"
	ActionLibraryMembersRemove = "knowledge.library_members.remove"

	ActionPagesList        = "knowledge.pages.list"
	ActionPagesGet         = "knowledge.pages.get"
	ActionPagesCreate      = "knowledge.pages.create"
	ActionPagesUpdate      = "knowledge.pages.update"
	ActionPagesMove        = "knowledge.pages.move"
	ActionPagesReorder     = "knowledge.pages.reorder"
	ActionPagesCopy        = "knowledge.pages.copy"
	ActionPagesArchive     = "knowledge.pages.archive"
	ActionPagesRestore     = "knowledge.pages.restore"
	ActionPagesPublish     = "knowledge.pages.publish"
	ActionPagesMovePreview = "knowledge.pages.move_preview"

	ActionPageAccessGet = "knowledge.page_access.get"
	ActionPageAccessSet = "knowledge.page_access.set"

	ActionPageCollaborationConnect = "knowledge.page_collaboration.connect"
	ActionPageCollaborationSync    = "knowledge.page_collaboration.sync"
	ActionPageCollaborationClose   = "knowledge.page_collaboration.close"
	ActionPageCollaborationPublish = "knowledge.page_collaboration.publish"

	ActionPageVersionsList    = "knowledge.page_versions.list"
	ActionPageVersionsGet     = "knowledge.page_versions.get"
	ActionPageVersionsDiff    = "knowledge.page_versions.diff"
	ActionPageVersionsRestore = "knowledge.page_versions.restore"

	ActionPageCommentsList    = "knowledge.page_comments.list"
	ActionPageCommentsCreate  = "knowledge.page_comments.create"
	ActionPageCommentsUpdate  = "knowledge.page_comments.update"
	ActionPageCommentsResolve = "knowledge.page_comments.resolve"
	ActionPageCommentsDelete  = "knowledge.page_comments.delete"

	ActionPageAttachmentsList   = "knowledge.page_attachments.list"
	ActionPageAttachmentsSet    = "knowledge.page_attachments.set"
	ActionPageAttachmentsRemove = "knowledge.page_attachments.remove"

	ActionDocumentsList     = "knowledge.documents.list"
	ActionDocumentsGet      = "knowledge.documents.get"
	ActionDocumentsUpload   = "knowledge.documents.upload"
	ActionDocumentsDownload = "knowledge.documents.download"
	ActionDocumentsRetry    = "knowledge.documents.retry"
	ActionDocumentsReindex  = "knowledge.documents.reindex"
	ActionDocumentsDelete   = "knowledge.documents.delete"
	ActionDocumentsTransfer = "knowledge.documents.transfer"

	ActionSourcesList = "knowledge.sources.list"
	ActionSourcesBind = "knowledge.sources.bind"

	ActionAccessExplain     = "knowledge.access.explain"
	ActionPrincipalsSearch  = "knowledge.principals.search"
	ActionPrincipalsResolve = "knowledge.principals.resolve"

	ActionDocumentsProcess       = "knowledge.documents.process"
	ActionPageProjectionsProcess = "knowledge.page_projections.process"
)

type actionSpec struct {
	key, capability, capabilityLabel, method, path, label string
	effect                                                actioncontract.EffectClass
	risk                                                  actioncontract.RiskLevel
}

// AuthorizationActions is Knowledge's single source-owned executable
// manifest. HTTP routes, direct SDK calls and SaaS RPCs project the same Action
// keys and same-key Permission definitions.
func AuthorizationActions() ([]actioncontract.ActionDefinition, error) {
	specs := []actionSpec{
		{ActionLibrariesList, "knowledge.libraries", "Knowledge libraries", "GET", "/knowledge/libraries", "List knowledge libraries", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionLibrariesCreate, "knowledge.libraries", "Knowledge libraries", "POST", "/knowledge/libraries", "Create knowledge library", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionLibrariesGet, "knowledge.libraries", "Knowledge libraries", "GET", "/knowledge/libraries/{libraryID}", "Read knowledge library", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionLibrariesUpdate, "knowledge.libraries", "Knowledge libraries", "PATCH", "/knowledge/libraries/{libraryID}", "Update knowledge library", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionLibrariesArchive, "knowledge.libraries", "Knowledge libraries", "POST", "/knowledge/libraries/{libraryID}/archive", "Archive knowledge library", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionLibrariesRestore, "knowledge.libraries", "Knowledge libraries", "POST", "/knowledge/libraries/{libraryID}/restore", "Restore knowledge library", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionLibraryMembersList, "knowledge.library_members", "Knowledge library members", "GET", "/knowledge/libraries/{libraryID}/members", "List knowledge library members", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionLibraryMembersSearch, "knowledge.library_members", "Knowledge library members", "GET", "/knowledge/libraries/{libraryID}/members/search", "Search mentionable knowledge library members", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionLibraryMembersSet, "knowledge.library_members", "Knowledge library members", "PUT", "/knowledge/libraries/{libraryID}/members/{userID}", "Set knowledge library member", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionLibraryMembersRemove, "knowledge.library_members", "Knowledge library members", "DELETE", "/knowledge/libraries/{libraryID}/members/{userID}", "Remove knowledge library member", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesList, "knowledge.pages", "Knowledge pages", "GET", "/knowledge/libraries/{libraryID}/pages", "List knowledge pages", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesGet, "knowledge.pages", "Knowledge pages", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}", "Read knowledge page", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesCreate, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages", "Create knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesUpdate, "knowledge.pages", "Knowledge pages", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/draft", "Save knowledge page draft", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesMove, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/move", "Move knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesMovePreview, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/move/preview", "Preview knowledge page move access impact", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesReorder, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/reorder", "Reorder sibling knowledge pages", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesCopy, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/copy", "Copy knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesArchive, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/archive", "Archive knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesRestore, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/restore", "Restore knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesPublish, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/publish", "Publish knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageAccessGet, "knowledge.page_access", "Knowledge page access", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/access", "Read knowledge page access policy", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageAccessSet, "knowledge.page_access", "Knowledge page access", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/access", "Set knowledge page access policy", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageCollaborationConnect, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/connect", "Connect knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCollaborationSync, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/sync", "Synchronize knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCollaborationClose, "knowledge.page_collaboration", "Knowledge page collaboration", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/{sessionID}", "Close knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageCollaborationPublish, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/publish", "Publish knowledge collaboration snapshot", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageVersionsList, "knowledge.page_versions", "Knowledge page versions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions", "List knowledge page versions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsGet, "knowledge.page_versions", "Knowledge page versions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/{version}", "Read knowledge page version", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsDiff, "knowledge.page_versions", "Knowledge page versions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/diff", "Compare knowledge page versions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsRestore, "knowledge.page_versions", "Knowledge page versions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/{version}/restore", "Restore knowledge page version to draft", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsList, "knowledge.page_comments", "Knowledge page comments", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments", "List knowledge page comments", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageCommentsCreate, "knowledge.page_comments", "Knowledge page comments", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments", "Create knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageCommentsUpdate, "knowledge.page_comments", "Knowledge page comments", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}", "Update knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsResolve, "knowledge.page_comments", "Knowledge page comments", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}/resolve", "Resolve knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsDelete, "knowledge.page_comments", "Knowledge page comments", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}", "Delete knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageAttachmentsList, "knowledge.page_attachments", "Knowledge page attachments", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments", "List knowledge page attachments", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageAttachmentsSet, "knowledge.page_attachments", "Knowledge page attachments", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments/{documentID}", "Attach document to knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageAttachmentsRemove, "knowledge.page_attachments", "Knowledge page attachments", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments/{documentID}", "Detach document from knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionDocumentsList, "knowledge.documents", "Knowledge documents", "GET", "/knowledge/libraries/{libraryID}/documents", "List knowledge documents", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionDocumentsGet, "knowledge.documents", "Knowledge documents", "GET", "/knowledge/libraries/{libraryID}/documents/{documentID}", "Read knowledge document status", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionDocumentsUpload, "knowledge.documents", "Knowledge documents", "POST", "/knowledge/libraries/{libraryID}/documents", "Upload knowledge document", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionDocumentsDownload, "knowledge.documents", "Knowledge documents", "GET", "/knowledge/libraries/{libraryID}/documents/{documentID}/content", "Download knowledge document", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionDocumentsRetry, "knowledge.documents", "Knowledge documents", "POST", "/knowledge/libraries/{libraryID}/documents/{documentID}/retry", "Retry knowledge document processing", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionDocumentsReindex, "knowledge.documents", "Knowledge documents", "POST", "/knowledge/libraries/{libraryID}/documents/{documentID}/reindex", "Reindex knowledge document", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionDocumentsDelete, "knowledge.documents", "Knowledge documents", "DELETE", "/knowledge/libraries/{libraryID}/documents/{documentID}", "Delete knowledge document", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionDocumentsTransfer, "knowledge.documents", "Knowledge documents", "POST", "/knowledge/libraries/{libraryID}/documents/transfer", "Copy or move knowledge document", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionSourcesList, "knowledge.sources", "Knowledge sources", "GET", "/knowledge/libraries/{libraryID}/sources", "List knowledge sources", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionSourcesBind, "knowledge.sources", "Knowledge sources", "PUT", "/knowledge/libraries/{libraryID}/source", "Bind knowledge source", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionAccessExplain, "knowledge.access", "Knowledge access", "POST", "/knowledge/libraries/{libraryID}/access/explain", "Explain effective knowledge access", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPrincipalsSearch, "knowledge.principals", "Knowledge principals", "GET", "/knowledge/principals", "Search visible knowledge principals", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPrincipalsResolve, "knowledge.principals", "Knowledge principals", "POST", "/knowledge/principals/resolve", "Resolve visible knowledge principals", actioncontract.EffectRead, actioncontract.RiskMedium},
	}
	definitions := make([]actioncontract.ActionDefinition, 0, len(specs)+2)
	for _, spec := range specs {
		definition, err := normalizeRoleAction(spec)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}
	for _, definition := range []actioncontract.ActionDefinition{
		signedAction(ActionDocumentsProcess, "knowledge.operations", "Knowledge operations", "documents_process", "Process knowledge document work"),
		signedAction(ActionPageProjectionsProcess, "knowledge.operations", "Knowledge operations", "page_projections_process", "Process knowledge page projection work"),
	} {
		normalized, err := actioncontract.NormalizeDefinition(definition)
		if err != nil {
			return nil, fmt.Errorf("normalize Knowledge Action %q: %w", definition.Key, err)
		}
		definitions = append(definitions, normalized)
	}
	registry := actioncontract.NewRegistry()
	if err := registry.Register(definitions...); err != nil {
		return nil, fmt.Errorf("register Knowledge Actions: %w", err)
	}
	if err := registry.Freeze(); err != nil {
		return nil, fmt.Errorf("freeze Knowledge Actions: %w", err)
	}
	return registry.Definitions(), nil
}

// AuthorizationAction resolves one exact source-owned action. Callers use it
// when they must evaluate the same Permission projection used by HTTP and SaaS.
func AuthorizationAction(key string) (actioncontract.ActionDefinition, bool, error) {
	definitions, err := AuthorizationActions()
	if err != nil {
		return actioncontract.ActionDefinition{}, false, err
	}
	key = strings.TrimSpace(key)
	for _, definition := range definitions {
		if definition.Key == key {
			return actioncontract.CloneDefinition(definition), true, nil
		}
	}
	return actioncontract.ActionDefinition{}, false, nil
}

func normalizeRoleAction(spec actionSpec) (actioncontract.ActionDefinition, error) {
	separator := strings.LastIndex(spec.key, ".")
	idempotency := "not_applicable"
	auditClass := "owner_read_audit_policy"
	if spec.effect == actioncontract.EffectWrite {
		idempotency = "caller_key_required"
		auditClass = "mutation_audit_required"
	}
	definition := actioncontract.ActionDefinition{
		Key: spec.key, Owner: AuthorizationOwner, SourceKind: "module_http", CapabilityKey: spec.capability, CapabilityLabel: spec.capabilityLabel,
		OperationKey: spec.key[separator+1:], OperationLabel: spec.label, Label: spec.label,
		Exposures:     []actioncontract.Exposure{actioncontract.ExposurePublic, actioncontract.ExposureManagement},
		Authorization: actioncontract.Authorization{Strategy: actioncontract.AuthorizationAuthenticated},
		HTTP:          &actioncontract.HTTPBinding{Method: spec.method, RouteTemplate: spec.path},
		NonHTTP: []actioncontract.NonHTTPBinding{
			{Kind: "sdk", InvocationKey: spec.key},
			{Kind: "saas_rpc", InvocationKey: spec.key},
		},
		Permission: &actioncontract.PermissionDefinition{
			Key: spec.key, Owner: AuthorizationOwner, ResourceKey: spec.key[:separator], OperationKey: spec.key[separator+1:],
			Label: spec.label, Category: spec.capabilityLabel, GrantMode: actioncontract.PermissionGrantDirect, LifecycleStatus: actioncontract.LifecycleActive,
		},
		EffectClass: spec.effect, RiskLevel: spec.risk, IdempotencyDecision: idempotency,
		AuditClass: auditClass, AuditEvent: strings.ReplaceAll(spec.key, ".", "_"), LifecycleStatus: actioncontract.LifecycleActive,
	}
	normalized, err := actioncontract.NormalizeDefinition(definition)
	if err != nil {
		return actioncontract.ActionDefinition{}, fmt.Errorf("normalize Knowledge Action %q: %w", spec.key, err)
	}
	return normalized, nil
}

func signedAction(key, capability, capabilityLabel, operation, label string) actioncontract.ActionDefinition {
	return actioncontract.ActionDefinition{
		Key: key, Owner: AuthorizationOwner, SourceKind: "module", CapabilityKey: capability, CapabilityLabel: capabilityLabel,
		OperationKey: operation, OperationLabel: label, Label: label,
		Exposures:     []actioncontract.Exposure{actioncontract.ExposureOps},
		Authorization: actioncontract.Authorization{Strategy: actioncontract.AuthorizationSigned, PolicyKey: "knowledge.system_scope"},
		NonHTTP:       []actioncontract.NonHTTPBinding{{Kind: "worker", InvocationKey: key}},
		EffectClass:   actioncontract.EffectWrite, RiskLevel: actioncontract.RiskHigh,
		IdempotencyDecision: "system_operation_identity", AuditClass: "knowledge_system_operation", AuditEvent: strings.ReplaceAll(key, ".", "_"),
		LifecycleStatus: actioncontract.LifecycleActive,
	}
}

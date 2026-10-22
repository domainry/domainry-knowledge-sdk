package knowledgesdk

import (
	"fmt"
	"strings"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const AuthorizationOwner = "module:knowledge"

const (
	ActionOpenAPIGet = "knowledge.openapi.get"

	ActionLibrariesList    = "knowledge.libraries.list"
	ActionLibrariesCreate  = "knowledge.libraries.create"
	ActionLibrariesGet     = "knowledge.libraries.get"
	ActionLibrariesUpdate  = "knowledge.libraries.update"
	ActionLibrariesArchive = "knowledge.libraries.archive"
	ActionLibrariesRestore = "knowledge.libraries.restore"

	ActionLibraryMembersList               = "knowledge.library_members.list"
	ActionLibraryMembersSearch             = "knowledge.library_members.search"
	ActionLibraryMembersSet                = "knowledge.library_members.set"
	ActionLibraryMembersRemove             = "knowledge.library_members.remove"
	ActionLibraryOrganizationMembersList   = "knowledge.library_organization_members.list"
	ActionLibraryOrganizationMembersSet    = "knowledge.library_organization_members.set"
	ActionLibraryOrganizationMembersRemove = "knowledge.library_organization_members.remove"

	ActionPagesList                 = "knowledge.pages.list"
	ActionPagesGet                  = "knowledge.pages.get"
	ActionPagesCreate               = "knowledge.pages.create"
	ActionPageReferencesList        = "knowledge.page_references.list"
	ActionPageReferencesCreate      = "knowledge.page_references.create"
	ActionPagesUpdate               = "knowledge.pages.update"
	ActionPagesMove                 = "knowledge.pages.move"
	ActionPagesReorder              = "knowledge.pages.reorder"
	ActionPagesCopy                 = "knowledge.pages.copy"
	ActionPagesArchive              = "knowledge.pages.archive"
	ActionPagesRestore              = "knowledge.pages.restore"
	ActionPagesPublish              = "knowledge.pages.publish"
	ActionPagesMovePreview          = "knowledge.pages.move_preview"
	ActionPageNativeOperationsApply = "knowledge.page_native_operations.apply"
	ActionPageGridAutomationExecute = "knowledge.page_grid_automations.execute"

	ActionPageAccessGet        = "knowledge.page_access.get"
	ActionPageAccessSet        = "knowledge.page_access.set"
	ActionPageContentPolicyGet = "knowledge.page_content_policy.get"
	ActionPageContentPolicySet = "knowledge.page_content_policy.set"
	ActionPageContentCopy      = "knowledge.page_content.copy"
	ActionPageDownload         = "knowledge.page_content.download"
	ActionPagePrint            = "knowledge.page_content.print"
	ActionPageExport           = "knowledge.page_content.export"
	ActionPageDuplicate        = "knowledge.page_content.duplicate"

	ActionPageCollaborationConnect = "knowledge.page_collaboration.connect"
	ActionPageCollaborationSync    = "knowledge.page_collaboration.sync"
	ActionPageCollaborationEvents  = "knowledge.page_collaboration.events"
	ActionPageCollaborationClose   = "knowledge.page_collaboration.close"
	ActionPageCollaborationPublish = "knowledge.page_collaboration.publish"
	ActionPageLiveEvents           = "knowledge.page_live.events"
	ActionPagePersonalStateGet     = "knowledge.page_personal_state.get"
	ActionPagePersonalStateSet     = "knowledge.page_personal_state.set"
	ActionPageViewRecord           = "knowledge.page_views.record"
	ActionPageCollectionsList      = "knowledge.page_collections.list"
	ActionPageMetadataGet          = "knowledge.page_metadata.get"
	ActionPageMetadataSet          = "knowledge.page_metadata.set"
	ActionPageGovernanceSet        = "knowledge.page_governance.set"
	ActionPageGovernanceReview     = "knowledge.page_governance.review"
	ActionPageGovernanceAuditList  = "knowledge.page_governance_audit.list"
	ActionPageSearch               = "knowledge.page_search.search"
	ActionPageSearchEvents         = "knowledge.page_search.events"
	ActionPageIdentitiesSearch     = "knowledge.page_identities.search"
	ActionPageRelationshipGraphGet = "knowledge.page_relationship_graph.get"

	ActionPageVersionsList      = "knowledge.page_versions.list"
	ActionPageVersionsGet       = "knowledge.page_versions.get"
	ActionPageVersionsDiff      = "knowledge.page_versions.diff"
	ActionPageVersionsRestore   = "knowledge.page_versions.restore"
	ActionPageEditHistoryList   = "knowledge.page_edit_history.list"
	ActionPageSuggestionsList   = "knowledge.page_suggestions.list"
	ActionPageSuggestionsCreate = "knowledge.page_suggestions.create"
	ActionPageSuggestionsDecide = "knowledge.page_suggestions.decide"

	ActionPageCommentsList    = "knowledge.page_comments.list"
	ActionPageCommentsCreate  = "knowledge.page_comments.create"
	ActionPageCommentsUpdate  = "knowledge.page_comments.update"
	ActionPageCommentsResolve = "knowledge.page_comments.resolve"
	ActionPageCommentsDelete  = "knowledge.page_comments.delete"
	ActionPageCommentsReact   = "knowledge.page_comments.react"

	ActionPageAttachmentsList   = "knowledge.page_attachments.list"
	ActionPageAttachmentsSet    = "knowledge.page_attachments.set"
	ActionPageAttachmentsRemove = "knowledge.page_attachments.remove"

	ActionPageTasksList   = "knowledge.page_tasks.list"
	ActionPageTasksCreate = "knowledge.page_tasks.create"
	ActionPageTasksUpdate = "knowledge.page_tasks.update"

	ActionPagePollsList   = "knowledge.page_polls.list"
	ActionPagePollsCreate = "knowledge.page_polls.create"
	ActionPagePollsUpdate = "knowledge.page_polls.update"
	ActionPagePollsVote   = "knowledge.page_polls.vote"

	ActionPageFormSubmit          = "knowledge.page_forms.submit"
	ActionPageFormResponsesList   = "knowledge.page_forms.responses.list"
	ActionPageFormAnalyticsGet    = "knowledge.page_forms.analytics.get"
	ActionPageFormResponsesExport = "knowledge.page_forms.responses.export"

	ActionPageTemplatesList        = "knowledge.page_templates.list"
	ActionPageTemplatesCreate      = "knowledge.page_templates.create"
	ActionPageTemplatesUpdate      = "knowledge.page_templates.update"
	ActionPageTemplateVersionsList = "knowledge.page_template_versions.list"
	ActionPageTemplateVersionsGet  = "knowledge.page_template_versions.get"
	ActionPageTemplatesInstantiate = "knowledge.page_templates.instantiate"

	ActionWebhookSubscriptionsList   = "knowledge.webhook_subscriptions.list"
	ActionWebhookSubscriptionsCreate = "knowledge.webhook_subscriptions.create"
	ActionWebhookSubscriptionsUpdate = "knowledge.webhook_subscriptions.update"
	ActionWebhookEventsList          = "knowledge.webhook_events.list"
	ActionWebhookEventsReplay        = "knowledge.webhook_events.replay"
	ActionAgentReferencesResolve     = "knowledge.agent_references.resolve"
	ActionAIExecutionsCreate         = "knowledge.ai_executions.create"
	ActionAIExecutionsGet            = "knowledge.ai_executions.get"
	ActionAIExecutionsConfirm        = "knowledge.ai_executions.confirm"
	ActionApplicationCardsResolve    = "knowledge.application_cards.resolve"

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

	ActionAccessExplain                 = "knowledge.access.explain"
	ActionPrincipalsSearch              = "knowledge.principals.search"
	ActionPrincipalsResolve             = "knowledge.principals.resolve"
	ActionOrganizationPrincipalsResolve = "knowledge.organization_principals.resolve"
	ActionPageSharesList                = "knowledge.page_shares.list"
	ActionPageSharesCreate              = "knowledge.page_shares.create"
	ActionPageSharesRevoke              = "knowledge.page_shares.revoke"
	ActionPageSharesRedeem              = "knowledge.page_shares.redeem"
	ActionPageSharesRead                = "knowledge.page_shares.read"
	ActionPageShareContentCopy          = "knowledge.page_shares.content.copy"
	ActionPageShareContentPrint         = "knowledge.page_shares.content.print"
	ActionPageShareContentExport        = "knowledge.page_shares.content.export"

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
		{ActionOpenAPIGet, "knowledge.openapi", "Knowledge OpenAPI", "GET", "/knowledge/openapi.json", "Read authorized Knowledge OpenAPI document", actioncontract.EffectRead, actioncontract.RiskLow},
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
		{ActionLibraryOrganizationMembersList, "knowledge.library_members", "Knowledge library organization members", "GET", "/knowledge/libraries/{libraryID}/organization-members", "List knowledge library organization members", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionLibraryOrganizationMembersSet, "knowledge.library_members", "Knowledge library organization members", "PUT", "/knowledge/libraries/{libraryID}/organization-members/{organizationID}", "Set knowledge library organization member", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionLibraryOrganizationMembersRemove, "knowledge.library_members", "Knowledge library organization members", "DELETE", "/knowledge/libraries/{libraryID}/organization-members/{organizationID}", "Remove knowledge library organization member", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesList, "knowledge.pages", "Knowledge pages", "GET", "/knowledge/libraries/{libraryID}/pages", "List knowledge pages", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesGet, "knowledge.pages", "Knowledge pages", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}", "Read knowledge page", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesCreate, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages", "Create knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageReferencesList, "knowledge.page_references", "Knowledge page references", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/references", "List outgoing references or incoming backlinks", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageReferencesCreate, "knowledge.page_references", "Knowledge page references", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/references", "Create a child page and its reference block atomically", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesUpdate, "knowledge.pages", "Knowledge pages", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/draft", "Save knowledge page draft", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageNativeOperationsApply, "knowledge.page_native_operations", "Knowledge native editor operations", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/native-operations", "Apply conflict-aware native editor operations", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageGridAutomationExecute, "knowledge.page_grid_automations", "Knowledge grid automations", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/grid-automations/execute", "Execute one idempotent multidimensional-table automation", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesMove, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/move", "Move knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesMovePreview, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/move/preview", "Preview knowledge page move access impact", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagesReorder, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/reorder", "Reorder sibling knowledge pages", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesCopy, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/copy", "Copy knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagesArchive, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/archive", "Archive knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesRestore, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/restore", "Restore knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPagesPublish, "knowledge.pages", "Knowledge pages", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/publish", "Publish knowledge page", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageAccessGet, "knowledge.page_access", "Knowledge page access", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/access", "Read knowledge page access policy", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageAccessSet, "knowledge.page_access", "Knowledge page access", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/access", "Set knowledge page access policy", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageContentPolicyGet, "knowledge.page_content_policy", "Knowledge page content policy", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/content-policy", "Read knowledge page content policy", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageContentPolicySet, "knowledge.page_content_policy", "Knowledge page content policy", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/content-policy", "Set knowledge page content policy", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageContentCopy, "knowledge.page_content", "Knowledge page content", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/content-copy/authorize", "Authorize copying knowledge page content", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageDownload, "knowledge.page_content", "Knowledge page content", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/documents/{documentID}/download", "Download a file through its knowledge page", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPagePrint, "knowledge.page_content", "Knowledge page content", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/print/authorize", "Authorize printing a knowledge page", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageExport, "knowledge.page_content", "Knowledge page content", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/export/authorize", "Authorize exporting a knowledge page", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageDuplicate, "knowledge.page_content", "Knowledge page content", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/duplicate/authorize", "Authorize creating a derived copy of a knowledge page", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageSharesList, "knowledge.page_shares", "Knowledge page shares", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/shares", "List knowledge page shares", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageSharesCreate, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/shares", "Create a controlled knowledge page share", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageSharesRevoke, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/shares/{shareID}/revoke", "Revoke a knowledge page share", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageSharesRedeem, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/page-shares/{shareID}/redeem", "Redeem a controlled knowledge page share", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageSharesRead, "knowledge.page_shares", "Knowledge page shares", "GET", "/knowledge/page-shares/{shareID}/pages/{pageID}", "Read a redeemed published knowledge page share", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageShareContentCopy, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/page-shares/{shareID}/pages/{pageID}/content-copy/authorize", "Authorize copying content from a redeemed page share", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageShareContentPrint, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/page-shares/{shareID}/pages/{pageID}/print/authorize", "Authorize printing content from a redeemed page share", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageShareContentExport, "knowledge.page_shares", "Knowledge page shares", "POST", "/knowledge/page-shares/{shareID}/pages/{pageID}/export/authorize", "Authorize exporting content from a redeemed page share", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageCollaborationConnect, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/connect", "Connect knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCollaborationSync, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/sync", "Synchronize knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCollaborationEvents, "knowledge.page_collaboration", "Knowledge page collaboration", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/{sessionID}/events", "Stream knowledge page collaboration events", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageCollaborationClose, "knowledge.page_collaboration", "Knowledge page collaboration", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/{sessionID}", "Close knowledge page collaboration", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageCollaborationPublish, "knowledge.page_collaboration", "Knowledge page collaboration", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/collaboration/publish", "Publish knowledge collaboration snapshot", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageLiveEvents, "knowledge.page_live", "Knowledge page live updates", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/live-events", "Stream permission-safe knowledge page change signals", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagePersonalStateGet, "knowledge.page_navigation", "Knowledge page navigation", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/personal-state", "Read private page navigation state", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagePersonalStateSet, "knowledge.page_navigation", "Knowledge page navigation", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/personal-state", "Set favorite, subscription and shortcut state", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageViewRecord, "knowledge.page_navigation", "Knowledge page navigation", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/view", "Record an authorized page visit", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageCollectionsList, "knowledge.page_navigation", "Knowledge page navigation", "GET", "/knowledge/libraries/{libraryID}/page-collections", "List authorized page navigation collections", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageMetadataGet, "knowledge.page_metadata", "Knowledge page metadata", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/metadata", "Read page tags, owner, status and custom properties", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageMetadataSet, "knowledge.page_metadata", "Knowledge page metadata", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/metadata", "Set page aliases, relationships, tags, status and custom properties", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageGovernanceSet, "knowledge.page_governance", "Knowledge page governance", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/governance", "Set page owner, content administrators, retention and review policy", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionPageGovernanceReview, "knowledge.page_governance", "Knowledge page governance", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/governance/review", "Confirm a scheduled page governance review", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageGovernanceAuditList, "knowledge.page_governance_audit", "Knowledge page governance audit", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/governance/audit", "Query immutable page governance audit events", actioncontract.EffectRead, actioncontract.RiskHigh},
		{ActionPageSearch, "knowledge.page_search", "Knowledge page search", "POST", "/knowledge/libraries/{libraryID}/page-search", "Search authorized published page content and page metadata", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageSearchEvents, "knowledge.page_search", "Knowledge page search updates", "GET", "/knowledge/libraries/{libraryID}/page-search/events", "Stream permission-safe page search invalidation signals", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageIdentitiesSearch, "knowledge.page_discovery", "Knowledge page discovery", "GET", "/knowledge/libraries/{libraryID}/page-identities", "Search authorized page titles and aliases with disambiguating paths", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageRelationshipGraphGet, "knowledge.page_discovery", "Knowledge page discovery", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/relationship-graph", "Read the authorized bounded page relationship graph", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsList, "knowledge.page_versions", "Knowledge page versions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions", "List knowledge page versions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsGet, "knowledge.page_versions", "Knowledge page versions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/{version}", "Read knowledge page version", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsDiff, "knowledge.page_versions", "Knowledge page versions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/diff", "Compare knowledge page versions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageVersionsRestore, "knowledge.page_versions", "Knowledge page versions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/versions/{version}/restore", "Restore knowledge page version to draft", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageEditHistoryList, "knowledge.page_edit_history", "Knowledge page edit history", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/edit-history", "List continuous block-level page edit history", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageSuggestionsList, "knowledge.page_suggestions", "Knowledge page suggestions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/suggestions", "List knowledge page revision suggestions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageSuggestionsCreate, "knowledge.page_suggestions", "Knowledge page suggestions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/suggestions", "Create a knowledge page revision suggestion", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageSuggestionsDecide, "knowledge.page_suggestions", "Knowledge page suggestions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/suggestions/{suggestionID}/decision", "Accept, reject, or withdraw a revision suggestion", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsList, "knowledge.page_comments", "Knowledge page comments", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments", "List knowledge page comments", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageCommentsCreate, "knowledge.page_comments", "Knowledge page comments", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments", "Create knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageCommentsUpdate, "knowledge.page_comments", "Knowledge page comments", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}", "Update knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsResolve, "knowledge.page_comments", "Knowledge page comments", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}/resolve", "Resolve knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsDelete, "knowledge.page_comments", "Knowledge page comments", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}", "Delete knowledge page comment", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageCommentsReact, "knowledge.page_comments", "Knowledge page comments", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/comments/{commentID}/reaction", "Set knowledge page comment reaction", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageAttachmentsList, "knowledge.page_attachments", "Knowledge page attachments", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments", "List knowledge page attachments", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageAttachmentsSet, "knowledge.page_attachments", "Knowledge page attachments", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments/{documentID}", "Attach document to knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageAttachmentsRemove, "knowledge.page_attachments", "Knowledge page attachments", "DELETE", "/knowledge/libraries/{libraryID}/pages/{pageID}/attachments/{documentID}", "Detach document from knowledge page", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageTasksList, "knowledge.page_tasks", "Knowledge page tasks", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/tasks", "List knowledge page tasks", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageTasksCreate, "knowledge.page_tasks", "Knowledge page tasks", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/tasks", "Create knowledge page task", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageTasksUpdate, "knowledge.page_tasks", "Knowledge page tasks", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/tasks/{taskID}", "Update knowledge page task", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagePollsList, "knowledge.page_polls", "Knowledge page polls", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/polls", "List knowledge page polls", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPagePollsCreate, "knowledge.page_polls", "Knowledge page polls", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/polls", "Create knowledge page poll", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagePollsUpdate, "knowledge.page_polls", "Knowledge page polls", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/polls/{pollID}", "Update knowledge page poll", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPagePollsVote, "knowledge.page_polls", "Knowledge page polls", "PUT", "/knowledge/libraries/{libraryID}/pages/{pageID}/polls/{pollID}/ballot", "Vote in knowledge page poll", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageFormSubmit, "knowledge.page_forms", "Knowledge page forms", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/form/submissions", "Submit a published knowledge page form", actioncontract.EffectWrite, actioncontract.RiskLow},
		{ActionPageFormResponsesList, "knowledge.page_forms", "Knowledge page forms", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/form/submissions", "List knowledge page form responses", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageFormAnalyticsGet, "knowledge.page_forms", "Knowledge page forms", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/form/analytics", "Read knowledge page form analytics", actioncontract.EffectRead, actioncontract.RiskMedium},
		{ActionPageFormResponsesExport, "knowledge.page_forms", "Knowledge page forms", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/form/submissions/export", "Export knowledge page form responses", actioncontract.EffectRead, actioncontract.RiskHigh},
		{ActionPageTemplatesList, "knowledge.page_templates", "Knowledge page templates", "GET", "/knowledge/libraries/{libraryID}/templates", "List knowledge page templates", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageTemplatesCreate, "knowledge.page_templates", "Knowledge page templates", "POST", "/knowledge/libraries/{libraryID}/templates", "Save page subtree as a template", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageTemplatesUpdate, "knowledge.page_templates", "Knowledge page templates", "PUT", "/knowledge/libraries/{libraryID}/templates/{templateID}", "Create a knowledge page template version", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionPageTemplateVersionsList, "knowledge.page_templates", "Knowledge page templates", "GET", "/knowledge/libraries/{libraryID}/templates/{templateID}/versions", "List knowledge page template versions", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageTemplateVersionsGet, "knowledge.page_templates", "Knowledge page templates", "GET", "/knowledge/libraries/{libraryID}/templates/{templateID}/versions/{version}", "Read a knowledge page template version", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionPageTemplatesInstantiate, "knowledge.page_templates", "Knowledge page templates", "POST", "/knowledge/libraries/{libraryID}/templates/{templateID}/instantiate", "Create an independent page subtree from a template", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionWebhookSubscriptionsList, "knowledge.webhook_subscriptions", "Knowledge webhook subscriptions", "GET", "/knowledge/libraries/{libraryID}/webhook-subscriptions", "List Knowledge webhook subscriptions", actioncontract.EffectRead, actioncontract.RiskHigh},
		{ActionWebhookSubscriptionsCreate, "knowledge.webhook_subscriptions", "Knowledge webhook subscriptions", "POST", "/knowledge/libraries/{libraryID}/webhook-subscriptions", "Create a Knowledge webhook subscription", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionWebhookSubscriptionsUpdate, "knowledge.webhook_subscriptions", "Knowledge webhook subscriptions", "PUT", "/knowledge/libraries/{libraryID}/webhook-subscriptions/{subscriptionID}", "Update or revoke a Knowledge webhook subscription", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionWebhookEventsList, "knowledge.webhook_events", "Knowledge webhook events", "GET", "/knowledge/libraries/{libraryID}/webhook-events", "List durable Knowledge webhook events and deliveries", actioncontract.EffectRead, actioncontract.RiskHigh},
		{ActionWebhookEventsReplay, "knowledge.webhook_events", "Knowledge webhook events", "POST", "/knowledge/libraries/{libraryID}/webhook-events/{eventID}/replay", "Replay a Knowledge webhook event", actioncontract.EffectWrite, actioncontract.RiskHigh},
		{ActionAgentReferencesResolve, "knowledge.agent_references", "Knowledge Agent references", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/agent-references/resolve", "Resolve one immutable permission-checked Agent citation", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionAIExecutionsCreate, "knowledge.ai_executions", "Knowledge AI executions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions", "Create a source-pinned Knowledge AI execution", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionAIExecutionsGet, "knowledge.ai_executions", "Knowledge AI executions", "GET", "/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions/{executionID}", "Read one Knowledge AI execution and its sources", actioncontract.EffectRead, actioncontract.RiskLow},
		{ActionAIExecutionsConfirm, "knowledge.ai_executions", "Knowledge AI executions", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions/{executionID}/confirm", "Confirm generated text for insertion by the active editor", actioncontract.EffectWrite, actioncontract.RiskMedium},
		{ActionApplicationCardsResolve, "knowledge.application_cards", "Knowledge application cards", "POST", "/knowledge/libraries/{libraryID}/pages/{pageID}/application-cards/resolve", "Resolve one current authorized application card", actioncontract.EffectRead, actioncontract.RiskLow},
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
		{ActionOrganizationPrincipalsResolve, "knowledge.principals", "Knowledge organization principals", "POST", "/knowledge/organization-principals/resolve", "Resolve visible Identity organization principals", actioncontract.EffectRead, actioncontract.RiskMedium},
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

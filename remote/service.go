package remote

import (
	"context"
	"encoding/json"
	"time"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	knowledgepersistence "github.com/domainry/domainry-knowledge-sdk/persistence"
)

type serviceClient struct {
	client    *client
	options   knowledgecontract.Options
	runtimeID string
}

func (*serviceClient) Close()                  {}
func (*serviceClient) Start(context.Context)   {}
func (*serviceClient) WakeKnowledgeDocuments() {}

func (value *serviceClient) invoke(ctx context.Context, operation string, input any, output any) error {
	return value.client.invokeWithOptions(ctx, operation, input, output, value.options)
}

func (value *serviceClient) BindKnowledgeLibrarySource(ctx context.Context, id string, input knowledgesdk.KnowledgeLibrarySourceWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionSourcesBind, struct {
		ID        string                                   `json:"id"`
		Input     knowledgesdk.KnowledgeLibrarySourceWrite `json:"input"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageShares(ctx context.Context, libraryID, pageID, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSharePage, error) {
	var output knowledgesdk.KnowledgePageSharePage
	err := value.invoke(ctx, knowledgesdk.ActionPageSharesList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, after, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageShare(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageShareCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageShare, error) {
	var output knowledgesdk.KnowledgePageShare
	err := value.invoke(ctx, knowledgesdk.ActionPageSharesCreate, struct {
		LibraryID string                                `json:"library_id"`
		PageID    string                                `json:"page_id"`
		Input     knowledgesdk.KnowledgePageShareCreate `json:"input"`
		Authority knowledgesdk.Authority                `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) RevokeKnowledgePageShare(ctx context.Context, libraryID, pageID, shareID string, expectedRevision int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageShare, error) {
	var output knowledgesdk.KnowledgePageShare
	err := value.invoke(ctx, knowledgesdk.ActionPageSharesRevoke, struct {
		LibraryID        string                 `json:"library_id"`
		PageID           string                 `json:"page_id"`
		ShareID          string                 `json:"share_id"`
		ExpectedRevision int64                  `json:"expected_revision"`
		Authority        knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, shareID, expectedRevision, authority}, &output)
	return output, err
}

func (value *serviceClient) RedeemKnowledgePageShare(ctx context.Context, shareID string, input knowledgesdk.KnowledgePageShareRedeem, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageShare, error) {
	var output knowledgesdk.KnowledgePageShare
	err := value.invoke(ctx, knowledgesdk.ActionPageSharesRedeem, struct {
		ShareID   string                                `json:"share_id"`
		Input     knowledgesdk.KnowledgePageShareRedeem `json:"input"`
		Authority knowledgesdk.Authority                `json:"authority"`
	}{shareID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ReadKnowledgePageShare(ctx context.Context, shareID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSharedView, error) {
	var output knowledgesdk.KnowledgePageSharedView
	err := value.invoke(ctx, knowledgesdk.ActionPageSharesRead, struct {
		ShareID   string                 `json:"share_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{shareID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) AuthorizeKnowledgePageShareContentOperation(ctx context.Context, shareID, pageID, operation string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageContentDecision, error) {
	var output knowledgesdk.KnowledgePageContentDecision
	action := map[string]string{
		knowledgesdk.KnowledgePageContentOperationCopy:   knowledgesdk.ActionPageShareContentCopy,
		knowledgesdk.KnowledgePageContentOperationPrint:  knowledgesdk.ActionPageShareContentPrint,
		knowledgesdk.KnowledgePageContentOperationExport: knowledgesdk.ActionPageShareContentExport,
	}[operation]
	if action == "" {
		return output, &knowledgesdk.Error{Class: "bad_request", Code: "knowledge.page_content_operation_invalid"}
	}
	err := value.invoke(ctx, action, struct {
		ShareID   string                 `json:"share_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{shareID, pageID, authority}, &output)
	return output, err
}
func (value *serviceClient) CreateKnowledgeLibrary(ctx context.Context, input knowledgesdk.KnowledgeLibraryCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibrariesCreate, struct {
		Input     knowledgesdk.KnowledgeLibraryCreate `json:"input"`
		Authority knowledgesdk.Authority              `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value *serviceClient) DatasourceAccess(ctx context.Context, id, operation string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.datasource_access", struct {
		ID        string                 `json:"id"`
		Operation string                 `json:"operation"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, operation, authority}, &output)
	return output, err
}
func (value *serviceClient) DatasourceDefinitions(ctx context.Context, library string, authority knowledgesdk.Authority) ([]knowledgesdk.KnowledgeDatasourceDefinition, error) {
	var output []knowledgesdk.KnowledgeDatasourceDefinition
	err := value.invoke(ctx, "service.datasource_definitions", struct {
		Library   string                 `json:"library"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, authority}, &output)
	return output, err
}
func (value *serviceClient) DeleteKnowledgeDocument(ctx context.Context, library, id string, expected int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsDelete, struct {
		Library   string                 `json:"library"`
		ID        string                 `json:"id"`
		Expected  int64                  `json:"expected"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) RetryKnowledgeDocument(ctx context.Context, library, id string, expected int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsRetry, struct {
		Library   string                 `json:"library"`
		ID        string                 `json:"id"`
		Expected  int64                  `json:"expected"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) ReindexKnowledgeDocument(ctx context.Context, library, id string, expected int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsReindex, struct {
		Library   string                 `json:"library"`
		ID        string                 `json:"id"`
		Expected  int64                  `json:"expected"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) DocumentAccess(ctx context.Context, library, operation string, authority knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRepository, error) {
	if err := value.invoke(ctx, "service.document_access", struct {
		Library   string                 `json:"library"`
		Operation string                 `json:"operation"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, operation, authority}, nil); err != nil {
		return nil, err
	}
	return documentRepository{client: value.client}, nil
}
func (value *serviceClient) DocumentBinding(ctx context.Context, library string, authority knowledgesdk.Authority) (knowledgesdk.ManagedKnowledgeDocumentSource, error) {
	var output remoteSourceDescriptor
	err := value.invoke(ctx, "service.document_binding", struct {
		Library   string                 `json:"library"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, authority}, &output)
	if err != nil {
		return nil, err
	}
	return managedSource{client: value.client, descriptor: output}, nil
}
func (value *serviceClient) DownloadKnowledgeDocument(ctx context.Context, library, id string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentDownload, error) {
	request := knowledgecontract.SaaSDocumentDownloadRequest{Library: library, ID: id, Authority: authority}
	var document knowledgesdk.KnowledgeDocument
	grants, err := value.client.invokeWithGrants(ctx, knowledgecontract.SaaSDocumentDownloadPrepareOperation, request, &document, value.options)
	if err != nil {
		return knowledgesdk.KnowledgeDocumentDownload{}, err
	}
	request.Grants = grants
	header, err := encodeTransferHeader(request)
	if err != nil {
		return knowledgesdk.KnowledgeDocumentDownload{}, err
	}
	return knowledgesdk.KnowledgeDocumentDownload{
		Document: document,
		Content:  &rangeDocumentReader{client: value.client, ctx: ctx, document: document, transferHeader: header},
	}, nil
}
func (value *serviceClient) KnowledgeDocument(ctx context.Context, library, id string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsGet, struct {
		Library   string                 `json:"library"`
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, id, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePagePersonalState(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePersonalState, error) {
	var output knowledgesdk.KnowledgePagePersonalState
	err := value.invoke(ctx, knowledgesdk.ActionPagePersonalStateGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePagePersonalState(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePagePersonalStateWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePersonalState, error) {
	var output knowledgesdk.KnowledgePagePersonalState
	err := value.invoke(ctx, knowledgesdk.ActionPagePersonalStateSet, struct {
		LibraryID string                                       `json:"library_id"`
		PageID    string                                       `json:"page_id"`
		Input     knowledgesdk.KnowledgePagePersonalStateWrite `json:"input"`
		Authority knowledgesdk.Authority                       `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) RecordKnowledgePageView(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageViewRecord, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePersonalState, error) {
	var output knowledgesdk.KnowledgePagePersonalState
	err := value.invoke(ctx, knowledgesdk.ActionPageViewRecord, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePageViewRecord `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageCollection(ctx context.Context, libraryID string, kind knowledgesdk.KnowledgePageCollectionKind, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageCollectionPage, error) {
	var output knowledgesdk.KnowledgePageCollectionPage
	err := value.invoke(ctx, knowledgesdk.ActionPageCollectionsList, struct {
		LibraryID string                                   `json:"library_id"`
		Kind      knowledgesdk.KnowledgePageCollectionKind `json:"kind"`
		After     string                                   `json:"after"`
		Limit     int                                      `json:"limit"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, kind, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgePageMetadata(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageMetadata, error) {
	var output knowledgesdk.KnowledgePageMetadata
	err := value.invoke(ctx, knowledgesdk.ActionPageMetadataGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageMetadata(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageMetadataWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageMetadata, error) {
	var output knowledgesdk.KnowledgePageMetadata
	err := value.invoke(ctx, knowledgesdk.ActionPageMetadataSet, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		Input     knowledgesdk.KnowledgePageMetadataWrite `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageGovernance(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageGovernanceWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageMetadata, error) {
	var output knowledgesdk.KnowledgePageMetadata
	err := value.invoke(ctx, knowledgesdk.ActionPageGovernanceSet, struct {
		LibraryID string                                    `json:"library_id"`
		PageID    string                                    `json:"page_id"`
		Input     knowledgesdk.KnowledgePageGovernanceWrite `json:"input"`
		Authority knowledgesdk.Authority                    `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ReviewKnowledgePageGovernance(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageGovernanceReview, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageMetadata, error) {
	var output knowledgesdk.KnowledgePageMetadata
	err := value.invoke(ctx, knowledgesdk.ActionPageGovernanceReview, struct {
		LibraryID string                                     `json:"library_id"`
		PageID    string                                     `json:"page_id"`
		Input     knowledgesdk.KnowledgePageGovernanceReview `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageGovernanceAudit(ctx context.Context, input knowledgesdk.KnowledgeGovernanceAuditQuery, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeGovernanceAuditPage, error) {
	var output knowledgesdk.KnowledgeGovernanceAuditPage
	err := value.invoke(ctx, knowledgesdk.ActionPageGovernanceAuditList, struct {
		Input     knowledgesdk.KnowledgeGovernanceAuditQuery `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{input, authority}, &output)
	return output, err
}

func (value *serviceClient) SearchKnowledgePages(ctx context.Context, libraryID string, input knowledgesdk.KnowledgePageSearchInput, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSearchPage, error) {
	var output knowledgesdk.KnowledgePageSearchPage
	err := value.invoke(ctx, knowledgesdk.ActionPageSearch, struct {
		LibraryID string                                `json:"library_id"`
		Input     knowledgesdk.KnowledgePageSearchInput `json:"input"`
		Authority knowledgesdk.Authority                `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageSearchEvents(ctx context.Context, libraryID string, input knowledgesdk.KnowledgePageSearchEvents, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSearchEventBatch, error) {
	var out knowledgesdk.KnowledgePageSearchEventBatch
	err := value.invoke(ctx, knowledgesdk.ActionPageSearchEvents, struct {
		LibraryID string                                 `json:"library_id"`
		Input     knowledgesdk.KnowledgePageSearchEvents `json:"input"`
		Authority knowledgesdk.Authority                 `json:"authority"`
	}{LibraryID: libraryID, Input: input, Authority: authority}, &out)
	return out, err
}

func (value *serviceClient) SearchKnowledgePageIdentities(ctx context.Context, libraryID, query string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageIdentityPage, error) {
	var output knowledgesdk.KnowledgePageIdentityPage
	err := value.invoke(ctx, knowledgesdk.ActionPageIdentitiesSearch, struct {
		LibraryID string                 `json:"library_id"`
		Query     string                 `json:"query"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, query, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageRelationshipGraph(ctx context.Context, libraryID, pageID string, depth, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageRelationshipGraph, error) {
	var output knowledgesdk.KnowledgePageRelationshipGraph
	err := value.invoke(ctx, knowledgesdk.ActionPageRelationshipGraphGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Depth     int                    `json:"depth"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, depth, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeDocumentRecord(ctx context.Context, _ knowledgepersistence.KnowledgeDocumentRepository, library, id string, authority knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var output knowledgepersistence.KnowledgeDocumentRecord
	err := value.invoke(ctx, "service.document_record", struct {
		Library   string                 `json:"library"`
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, id, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeDocuments(ctx context.Context, library, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentPage, error) {
	var output knowledgesdk.KnowledgeDocumentPage
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsList, struct {
		Library   string                 `json:"library"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{library, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibraries(ctx context.Context, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryPage, error) {
	var output knowledgesdk.KnowledgeLibraryPage
	err := value.invoke(ctx, knowledgesdk.ActionLibrariesList, struct {
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibrary(ctx context.Context, id string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibrariesGet, struct {
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibraryMembers(ctx context.Context, id, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryMembers, error) {
	var output knowledgesdk.KnowledgeLibraryMembers
	err := value.invoke(ctx, knowledgesdk.ActionLibraryMembersList, struct {
		ID        string                 `json:"id"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibraryOrganizationMembers(ctx context.Context, id, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryOrganizationMembers, error) {
	var output knowledgesdk.KnowledgeLibraryOrganizationMembers
	err := value.invoke(ctx, knowledgesdk.ActionLibraryOrganizationMembersList, struct {
		ID        string                 `json:"id"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibrarySources(ctx context.Context, id, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrarySources, error) {
	var output knowledgesdk.KnowledgeLibrarySources
	err := value.invoke(ctx, knowledgesdk.ActionSourcesList, struct {
		ID        string                 `json:"id"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) LibraryAccess(ctx context.Context, operation string, authority knowledgesdk.Authority) (knowledgepersistence.KnowledgeLibraryRepository, error) {
	if err := value.invoke(ctx, "service.library_access", struct {
		Operation string                 `json:"operation"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{operation, authority}, nil); err != nil {
		return nil, err
	}
	return libraryRepository{service: value}, nil
}
func (value *serviceClient) LibraryAuthorize(ctx context.Context, operation string, item knowledgesdk.KnowledgeLibrary, authority knowledgesdk.Authority) error {
	return value.invoke(ctx, "service.library_authorize", struct {
		Operation string                        `json:"operation"`
		Item      knowledgesdk.KnowledgeLibrary `json:"item"`
		Authority knowledgesdk.Authority        `json:"authority"`
	}{operation, item, authority}, nil)
}
func (value *serviceClient) LibraryResult(ctx context.Context, item knowledgesdk.KnowledgeLibrary, inputErr error, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	if inputErr != nil {
		return knowledgesdk.KnowledgeLibrary{}, inputErr
	}
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_result", struct {
		Item      knowledgesdk.KnowledgeLibrary `json:"item"`
		Authority knowledgesdk.Authority        `json:"authority"`
	}{item, authority}, &output)
	return output, err
}
func (value *serviceClient) RemoveKnowledgeLibraryMember(ctx context.Context, id, user string, revision int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibraryMembersRemove, struct {
		ID        string                 `json:"id"`
		User      string                 `json:"user"`
		Revision  int64                  `json:"revision"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, user, revision, authority}, &output)
	return output, err
}
func (value *serviceClient) SetKnowledgeLibraryMember(ctx context.Context, id, user string, input knowledgesdk.KnowledgeLibraryMemberWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibraryMembersSet, struct {
		ID        string                                   `json:"id"`
		User      string                                   `json:"user"`
		Input     knowledgesdk.KnowledgeLibraryMemberWrite `json:"input"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{id, user, input, authority}, &output)
	return output, err
}
func (value *serviceClient) RemoveKnowledgeLibraryOrganizationMember(ctx context.Context, id, organizationID string, revision int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibraryOrganizationMembersRemove, struct {
		ID             string                 `json:"id"`
		OrganizationID string                 `json:"organization_id"`
		Revision       int64                  `json:"revision"`
		Authority      knowledgesdk.Authority `json:"authority"`
	}{id, organizationID, revision, authority}, &output)
	return output, err
}
func (value *serviceClient) SetKnowledgeLibraryOrganizationMember(ctx context.Context, id, organizationID string, input knowledgesdk.KnowledgeLibraryMemberWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibraryOrganizationMembersSet, struct {
		ID             string                                   `json:"id"`
		OrganizationID string                                   `json:"organization_id"`
		Input          knowledgesdk.KnowledgeLibraryMemberWrite `json:"input"`
		Authority      knowledgesdk.Authority                   `json:"authority"`
	}{id, organizationID, input, authority}, &output)
	return output, err
}
func (value *serviceClient) TransferKnowledgeDocument(ctx context.Context, library string, input knowledgesdk.KnowledgeDocumentTransfer, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	err := value.invoke(ctx, knowledgesdk.ActionDocumentsTransfer, struct {
		Library   string                                 `json:"library"`
		Input     knowledgesdk.KnowledgeDocumentTransfer `json:"input"`
		Authority knowledgesdk.Authority                 `json:"authority"`
	}{library, input, authority}, &output)
	return output, err
}
func (value *serviceClient) UpdateKnowledgeLibrary(ctx context.Context, id string, input knowledgesdk.KnowledgeLibraryUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	var output knowledgesdk.KnowledgeLibrary
	err := value.invoke(ctx, knowledgesdk.ActionLibrariesUpdate, struct {
		ID        string                              `json:"id"`
		Input     knowledgesdk.KnowledgeLibraryUpdate `json:"input"`
		Authority knowledgesdk.Authority              `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}
func (value *serviceClient) UploadKnowledgeDocument(ctx context.Context, library string, input knowledgesdk.KnowledgeDocumentUpload, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	return value.uploadKnowledgeDocument(ctx, knowledgecontract.SaaSDocumentUploadRequest{
		Operation: knowledgecontract.SaaSDocumentUploadOperation,
		Library:   library,
		Input:     input,
		Authority: authority,
	})
}

func (value *serviceClient) UploadKnowledgeDocumentChunk(ctx context.Context, library string, input knowledgesdk.KnowledgeDocumentUploadChunk, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentUploadProgress, error) {
	request := knowledgecontract.SaaSDocumentChunkUploadRequest{
		Operation: knowledgecontract.SaaSDocumentUploadChunkOperation,
		Library:   library, Input: input, Authority: authority,
	}
	var output knowledgesdk.KnowledgeDocumentUploadProgress
	grants, err := value.client.invokeWithGrants(ctx, knowledgecontract.SaaSDocumentUploadPrepareOperation, request, nil, value.options)
	if err != nil {
		return output, err
	}
	request.Grants = grants
	err = value.client.uploadStream(ctx, knowledgecontract.SaaSDocumentUploadPath, request, request.Input.Content, request.Input.ChunkBytes, &output)
	return output, err
}

func (value *serviceClient) UploadKnowledgeDocumentForSource(ctx context.Context, library string, input knowledgesdk.KnowledgeDocumentUpload, source knowledgesdk.KnowledgeDocumentSourceAccess, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error) {
	return value.uploadKnowledgeDocument(ctx, knowledgecontract.SaaSDocumentUploadRequest{
		Operation: knowledgecontract.SaaSDocumentUploadSourceOperation,
		Library:   library,
		Input:     input,
		Source:    &source,
		Authority: authority,
	})
}

func (value *serviceClient) uploadKnowledgeDocument(ctx context.Context, request knowledgecontract.SaaSDocumentUploadRequest) (knowledgesdk.KnowledgeDocument, error) {
	var output knowledgesdk.KnowledgeDocument
	grants, err := value.client.invokeWithGrants(ctx, knowledgecontract.SaaSDocumentUploadPrepareOperation, request, nil, value.options)
	if err != nil {
		return output, err
	}
	request.Grants = grants
	err = value.client.uploadStream(ctx, knowledgecontract.SaaSDocumentUploadPath, request, request.Input.Content, request.Input.Bytes, &output)
	return output, err
}

type libraryRepository struct{ service *serviceClient }

func (v libraryRepository) CreateKnowledgeLibrary(c context.Context, i knowledgesdk.KnowledgeLibraryCreate, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.CreateKnowledgeLibrary(c, i, a)
}
func (v libraryRepository) KnowledgeLibraries(c context.Context, x string, l int, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryPage, error) {
	return v.service.KnowledgeLibraries(c, x, l, a)
}
func (v libraryRepository) KnowledgeLibrary(c context.Context, id string, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.KnowledgeLibrary(c, id, a)
}
func (v libraryRepository) UpdateKnowledgeLibrary(c context.Context, id string, i knowledgesdk.KnowledgeLibraryUpdate, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.UpdateKnowledgeLibrary(c, id, i, a)
}
func (v libraryRepository) KnowledgeLibraryMembers(c context.Context, id, x string, l int, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryMembers, error) {
	return v.service.KnowledgeLibraryMembers(c, id, x, l, a)
}
func (v libraryRepository) SetKnowledgeLibraryMember(c context.Context, id, u string, i knowledgesdk.KnowledgeLibraryMemberWrite, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.SetKnowledgeLibraryMember(c, id, u, i, a)
}
func (v libraryRepository) RemoveKnowledgeLibraryMember(c context.Context, id, u string, r int64, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.RemoveKnowledgeLibraryMember(c, id, u, r, a)
}
func (v libraryRepository) KnowledgeLibraryOrganizationMembers(c context.Context, id, x string, l int, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryOrganizationMembers, error) {
	return v.service.KnowledgeLibraryOrganizationMembers(c, id, x, l, a)
}
func (v libraryRepository) SetKnowledgeLibraryOrganizationMember(c context.Context, id, organizationID string, i knowledgesdk.KnowledgeLibraryMemberWrite, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.SetKnowledgeLibraryOrganizationMember(c, id, organizationID, i, a)
}
func (v libraryRepository) RemoveKnowledgeLibraryOrganizationMember(c context.Context, id, organizationID string, r int64, a knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error) {
	return v.service.RemoveKnowledgeLibraryOrganizationMember(c, id, organizationID, r, a)
}

type documentRepository struct{ client *client }

func (v documentRepository) ActivateKnowledgeDocumentSource(c context.Context, s knowledgesdk.KnowledgeDocumentStorageScope, id string) error {
	return v.client.invoke(c, "documents.activate_source", struct {
		Scope knowledgesdk.KnowledgeDocumentStorageScope `json:"scope"`
		ID    string                                     `json:"id"`
	}{s, id}, nil)
}
func (v documentRepository) KnowledgeDocumentLibrarySource(c context.Context, id string, a knowledgesdk.Authority) (string, error) {
	var o string
	e := v.client.invoke(c, "documents.library_source", struct {
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeSourceManaged(c context.Context, id string) (bool, error) {
	var o bool
	e := v.client.invoke(c, "documents.source_managed", struct {
		ID string `json:"id"`
	}{id}, &o)
	return o, e
}
func (v documentRepository) ReserveKnowledgeDocument(c context.Context, i knowledgepersistence.KnowledgeDocumentReserve, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.reserve", struct {
		Input     knowledgepersistence.KnowledgeDocumentReserve `json:"input"`
		Authority knowledgesdk.Authority                        `json:"authority"`
	}{i, a}, &o)
	return o, e
}
func (v documentRepository) CommitKnowledgeDocumentContent(c context.Context, id string, r int64, ref string, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.commit_content", struct {
		ID        string                 `json:"id"`
		Revision  int64                  `json:"revision"`
		BodyRef   string                 `json:"body_ref"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, r, ref, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocumentRecord(c context.Context, id string, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.record", struct {
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocumentByRemoteID(c context.Context, l, s, id string, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.by_remote", struct {
		Library   string                 `json:"library"`
		Source    string                 `json:"source"`
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{l, s, id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocuments(c context.Context, l, x string, n int, a knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentPage, error) {
	var o knowledgesdk.KnowledgeDocumentPage
	e := v.client.invoke(c, "documents.list", struct {
		Library   string                 `json:"library"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{l, x, n, a}, &o)
	return o, e
}
func (v documentRepository) RequestKnowledgeDocumentDeletion(c context.Context, id string, r int64, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.request_delete", struct {
		ID        string                 `json:"id"`
		Revision  int64                  `json:"revision"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, r, a}, &o)
	return o, e
}
func (v documentRepository) RequestKnowledgeDocumentRetry(c context.Context, id string, r int64, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.request_retry", struct {
		ID        string                 `json:"id"`
		Revision  int64                  `json:"revision"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, r, a}, &o)
	return o, e
}
func (v documentRepository) RequestKnowledgeDocumentReindex(c context.Context, id string, r int64, a knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.request_reindex", struct {
		ID        string                 `json:"id"`
		Revision  int64                  `json:"revision"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, r, a}, &o)
	return o, e
}
func (v documentRepository) ClaimKnowledgeDocumentWork(c context.Context, owner, scope string, now time.Time, d time.Duration) (knowledgepersistence.KnowledgeDocumentLease, bool, error) {
	var o struct {
		Lease   knowledgepersistence.KnowledgeDocumentLease `json:"lease"`
		Claimed bool                                        `json:"claimed"`
	}
	e := v.client.invoke(c, "documents.claim_work", struct {
		Owner    string        `json:"owner"`
		Scope    string        `json:"scope"`
		Now      time.Time     `json:"now"`
		Duration time.Duration `json:"duration"`
	}{owner, scope, now, d}, &o)
	return o.Lease, o.Claimed, e
}
func (v documentRepository) KnowledgeDocumentWorkRecord(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) (knowledgepersistence.KnowledgeDocumentRecord, error) {
	var o knowledgepersistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.work_record", l, &o)
	return o, e
}
func (v documentRepository) StartKnowledgeDocumentPut(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) (knowledgepersistence.KnowledgeDocumentRecord, bool, error) {
	var o struct {
		Record  knowledgepersistence.KnowledgeDocumentRecord `json:"record"`
		Started bool                                         `json:"started"`
	}
	e := v.client.invoke(c, "documents.start_put", l, &o)
	return o.Record, o.Started, e
}
func (v documentRepository) StartKnowledgeDocumentReindexPut(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) (knowledgepersistence.KnowledgeDocumentRecord, bool, error) {
	var o struct {
		Record  knowledgepersistence.KnowledgeDocumentRecord `json:"record"`
		Started bool                                         `json:"started"`
	}
	e := v.client.invoke(c, "documents.start_reindex_put", l, &o)
	return o.Record, o.Started, e
}
func (v documentRepository) StartKnowledgeDocumentRetiredDelete(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) (knowledgepersistence.KnowledgeDocumentRecord, bool, error) {
	var o struct {
		Record  knowledgepersistence.KnowledgeDocumentRecord `json:"record"`
		Started bool                                         `json:"started"`
	}
	e := v.client.invoke(c, "documents.start_retired_delete", l, &o)
	return o.Record, o.Started, e
}
func (v documentRepository) StartKnowledgeDocumentScan(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) (knowledgepersistence.KnowledgeDocumentRecord, bool, error) {
	var o struct {
		Record  knowledgepersistence.KnowledgeDocumentRecord `json:"record"`
		Started bool                                         `json:"started"`
	}
	e := v.client.invoke(c, "documents.start_scan", l, &o)
	return o.Record, o.Started, e
}
func (v documentRepository) StartKnowledgeDocumentDelete(c context.Context, l knowledgepersistence.KnowledgeDocumentLease) error {
	return v.client.invoke(c, "documents.start_delete", l, nil)
}
func (v documentRepository) ApplyKnowledgeDocumentProgress(c context.Context, l knowledgepersistence.KnowledgeDocumentLease, p knowledgepersistence.KnowledgeDocumentProgress) error {
	return v.client.invoke(c, "documents.apply_progress", struct {
		Lease    knowledgepersistence.KnowledgeDocumentLease    `json:"lease"`
		Progress knowledgepersistence.KnowledgeDocumentProgress `json:"progress"`
	}{l, p}, nil)
}

type artifactRepository struct{ client *client }

func (value *serviceClient) KnowledgePages(ctx context.Context, libraryID, parentID, after string, limit int, includeArchived bool, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePage, error) {
	var output knowledgesdk.KnowledgePagePage
	err := value.invoke(ctx, knowledgesdk.ActionPagesList, struct {
		LibraryID       string                 `json:"library_id"`
		ParentID        string                 `json:"parent_id"`
		After           string                 `json:"after"`
		Limit           int                    `json:"limit"`
		IncludeArchived bool                   `json:"include_archived"`
		Authority       knowledgesdk.Authority `json:"authority"`
	}{libraryID, parentID, after, limit, includeArchived, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePage(ctx context.Context, libraryID, pageID string, includeDraft bool, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPagesGet, struct {
		LibraryID    string                 `json:"library_id"`
		PageID       string                 `json:"page_id"`
		IncludeDraft bool                   `json:"include_draft"`
		Authority    knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, includeDraft, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePage(ctx context.Context, libraryID string, input knowledgesdk.KnowledgePageCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPagesCreate, struct {
		LibraryID string                           `json:"library_id"`
		Input     knowledgesdk.KnowledgePageCreate `json:"input"`
		Authority knowledgesdk.Authority           `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageReferences(ctx context.Context, libraryID, pageID string, direction knowledgesdk.KnowledgePageReferenceDirection, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageReferencePage, error) {
	var output knowledgesdk.KnowledgePageReferencePage
	err := value.invoke(ctx, knowledgesdk.ActionPageReferencesList, struct {
		LibraryID string                                       `json:"library_id"`
		PageID    string                                       `json:"page_id"`
		Direction knowledgesdk.KnowledgePageReferenceDirection `json:"direction"`
		After     string                                       `json:"after"`
		Limit     int                                          `json:"limit"`
		Authority knowledgesdk.Authority                       `json:"authority"`
	}{libraryID, pageID, direction, after, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageReference(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageReferenceCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageReferenceCreateResult, error) {
	var output knowledgesdk.KnowledgePageReferenceCreateResult
	err := value.invoke(ctx, knowledgesdk.ActionPageReferencesCreate, struct {
		LibraryID string                                    `json:"library_id"`
		PageID    string                                    `json:"page_id"`
		Input     knowledgesdk.KnowledgePageReferenceCreate `json:"input"`
		Authority knowledgesdk.Authority                    `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SaveKnowledgePageDraft(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageDraftWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPagesUpdate, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePageDraftWrite `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ApplyKnowledgePageNativeOperations(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageNativeOperationBatch, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageNativeOperationResult, error) {
	var output knowledgesdk.KnowledgePageNativeOperationResult
	err := value.invoke(ctx, knowledgesdk.ActionPageNativeOperationsApply, struct {
		LibraryID string                                         `json:"library_id"`
		PageID    string                                         `json:"page_id"`
		Input     knowledgesdk.KnowledgePageNativeOperationBatch `json:"input"`
		Authority knowledgesdk.Authority                         `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ExecuteKnowledgeGridAutomation(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeGridAutomationExecute, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeGridAutomationResult, error) {
	var output knowledgesdk.KnowledgeGridAutomationResult
	err := value.invoke(ctx, knowledgesdk.ActionPageGridAutomationExecute, struct {
		LibraryID string                                      `json:"library_id"`
		PageID    string                                      `json:"page_id"`
		Input     knowledgesdk.KnowledgeGridAutomationExecute `json:"input"`
		Authority knowledgesdk.Authority                      `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) MoveKnowledgePage(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageMove, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePage, error) {
	var output knowledgesdk.KnowledgePage
	err := value.invoke(ctx, knowledgesdk.ActionPagesMove, struct {
		LibraryID string                         `json:"library_id"`
		PageID    string                         `json:"page_id"`
		Input     knowledgesdk.KnowledgePageMove `json:"input"`
		Authority knowledgesdk.Authority         `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) PreviewKnowledgePageMove(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageMovePreview, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageMoveImpact, error) {
	var output knowledgesdk.KnowledgePageMoveImpact
	err := value.invoke(ctx, knowledgesdk.ActionPagesMovePreview, struct {
		LibraryID string                                `json:"library_id"`
		PageID    string                                `json:"page_id"`
		Input     knowledgesdk.KnowledgePageMovePreview `json:"input"`
		Authority knowledgesdk.Authority                `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ReorderKnowledgePages(ctx context.Context, libraryID string, input knowledgesdk.KnowledgePageReorder, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePage, error) {
	var output knowledgesdk.KnowledgePagePage
	err := value.invoke(ctx, knowledgesdk.ActionPagesReorder, struct {
		LibraryID string                            `json:"library_id"`
		Input     knowledgesdk.KnowledgePageReorder `json:"input"`
		Authority knowledgesdk.Authority            `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) CopyKnowledgePage(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageCopy, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPagesCopy, struct {
		LibraryID string                         `json:"library_id"`
		PageID    string                         `json:"page_id"`
		Input     knowledgesdk.KnowledgePageCopy `json:"input"`
		Authority knowledgesdk.Authority         `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageArchived(ctx context.Context, libraryID, pageID string, archived bool, input knowledgesdk.KnowledgePageStateWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePage, error) {
	action := knowledgesdk.ActionPagesArchive
	if !archived {
		action = knowledgesdk.ActionPagesRestore
	}
	var output knowledgesdk.KnowledgePage
	err := value.invoke(ctx, action, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePageStateWrite `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) PublishKnowledgePage(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePagePublish, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPagesPublish, struct {
		LibraryID string                            `json:"library_id"`
		PageID    string                            `json:"page_id"`
		Input     knowledgesdk.KnowledgePagePublish `json:"input"`
		Authority knowledgesdk.Authority            `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageVersions(ctx context.Context, libraryID, pageID string, before int64, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageVersionPage, error) {
	var output knowledgesdk.KnowledgePageVersionPage
	err := value.invoke(ctx, knowledgesdk.ActionPageVersionsList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Before    int64                  `json:"before"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, before, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageEditHistory(ctx context.Context, libraryID, pageID string, before int64, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageEditHistoryPage, error) {
	var output knowledgesdk.KnowledgePageEditHistoryPage
	err := value.invoke(ctx, knowledgesdk.ActionPageEditHistoryList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Before    int64                  `json:"before"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, before, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageSuggestions(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSuggestionPage, error) {
	var output knowledgesdk.KnowledgePageSuggestionPage
	err := value.invoke(ctx, knowledgesdk.ActionPageSuggestionsList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageSuggestion(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageSuggestionCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSuggestion, error) {
	var output knowledgesdk.KnowledgePageSuggestion
	err := value.invoke(ctx, knowledgesdk.ActionPageSuggestionsCreate, struct {
		LibraryID string                                     `json:"library_id"`
		PageID    string                                     `json:"page_id"`
		Input     knowledgesdk.KnowledgePageSuggestionCreate `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) DecideKnowledgePageSuggestion(ctx context.Context, libraryID, pageID, suggestionID string, input knowledgesdk.KnowledgePageSuggestionDecision, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageSuggestion, error) {
	var output knowledgesdk.KnowledgePageSuggestion
	err := value.invoke(ctx, knowledgesdk.ActionPageSuggestionsDecide, struct {
		LibraryID    string                                       `json:"library_id"`
		PageID       string                                       `json:"page_id"`
		SuggestionID string                                       `json:"suggestion_id"`
		Input        knowledgesdk.KnowledgePageSuggestionDecision `json:"input"`
		Authority    knowledgesdk.Authority                       `json:"authority"`
	}{libraryID, pageID, suggestionID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageVersion(ctx context.Context, libraryID, pageID string, version int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageVersion, error) {
	var output knowledgesdk.KnowledgePageVersion
	err := value.invoke(ctx, knowledgesdk.ActionPageVersionsGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Version   int64                  `json:"version"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, version, authority}, &output)
	return output, err
}

func (value *serviceClient) DiffKnowledgePageVersions(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageVersionDiffInput, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageVersionDiff, error) {
	var output knowledgesdk.KnowledgePageVersionDiff
	err := value.invoke(ctx, knowledgesdk.ActionPageVersionsDiff, struct {
		LibraryID string                                     `json:"library_id"`
		PageID    string                                     `json:"page_id"`
		Input     knowledgesdk.KnowledgePageVersionDiffInput `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) RestoreKnowledgePageVersion(ctx context.Context, libraryID, pageID string, version int64, input knowledgesdk.KnowledgePageRestoreVersion, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPageVersionsRestore, struct {
		LibraryID string                                   `json:"library_id"`
		PageID    string                                   `json:"page_id"`
		Version   int64                                    `json:"version"`
		Input     knowledgesdk.KnowledgePageRestoreVersion `json:"input"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, pageID, version, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageComments(ctx context.Context, libraryID, pageID, after string, limit int, unresolvedOnly bool, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageCommentPage, error) {
	var output knowledgesdk.KnowledgePageCommentPage
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsList, struct {
		LibraryID      string                 `json:"library_id"`
		PageID         string                 `json:"page_id"`
		After          string                 `json:"after"`
		Limit          int                    `json:"limit"`
		UnresolvedOnly bool                   `json:"unresolved_only"`
		Authority      knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, after, limit, unresolvedOnly, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageComment(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageCommentCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageComment, error) {
	var output knowledgesdk.KnowledgePageComment
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsCreate, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		Input     knowledgesdk.KnowledgePageCommentCreate `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) UpdateKnowledgePageComment(ctx context.Context, libraryID, pageID, commentID string, input knowledgesdk.KnowledgePageCommentUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageComment, error) {
	var output knowledgesdk.KnowledgePageComment
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsUpdate, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		CommentID string                                  `json:"comment_id"`
		Input     knowledgesdk.KnowledgePageCommentUpdate `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, commentID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) ResolveKnowledgePageComment(ctx context.Context, libraryID, pageID, commentID string, input knowledgesdk.KnowledgePageCommentResolve, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageComment, error) {
	var output knowledgesdk.KnowledgePageComment
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsResolve, struct {
		LibraryID string                                   `json:"library_id"`
		PageID    string                                   `json:"page_id"`
		CommentID string                                   `json:"comment_id"`
		Input     knowledgesdk.KnowledgePageCommentResolve `json:"input"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, pageID, commentID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) DeleteKnowledgePageComment(ctx context.Context, libraryID, pageID, commentID string, input knowledgesdk.KnowledgePageCommentDelete, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageComment, error) {
	var output knowledgesdk.KnowledgePageComment
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsDelete, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		CommentID string                                  `json:"comment_id"`
		Input     knowledgesdk.KnowledgePageCommentDelete `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, commentID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageCommentReaction(ctx context.Context, libraryID, pageID, commentID string, input knowledgesdk.KnowledgePageCommentReactionWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageComment, error) {
	var output knowledgesdk.KnowledgePageComment
	err := value.invoke(ctx, knowledgesdk.ActionPageCommentsReact, struct {
		LibraryID string                                         `json:"library_id"`
		PageID    string                                         `json:"page_id"`
		CommentID string                                         `json:"comment_id"`
		Input     knowledgesdk.KnowledgePageCommentReactionWrite `json:"input"`
		Authority knowledgesdk.Authority                         `json:"authority"`
	}{libraryID, pageID, commentID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageAttachments(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageAttachmentPage, error) {
	var output knowledgesdk.KnowledgePageAttachmentPage
	err := value.invoke(ctx, knowledgesdk.ActionPageAttachmentsList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageAttachment(ctx context.Context, libraryID, pageID, documentID string, input knowledgesdk.KnowledgePageAttachmentWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageAttachment, error) {
	var output knowledgesdk.KnowledgePageAttachment
	err := value.invoke(ctx, knowledgesdk.ActionPageAttachmentsSet, struct {
		LibraryID  string                                    `json:"library_id"`
		PageID     string                                    `json:"page_id"`
		DocumentID string                                    `json:"document_id"`
		Input      knowledgesdk.KnowledgePageAttachmentWrite `json:"input"`
		Authority  knowledgesdk.Authority                    `json:"authority"`
	}{libraryID, pageID, documentID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) RemoveKnowledgePageAttachment(ctx context.Context, libraryID, pageID, documentID string, input knowledgesdk.KnowledgePageAttachmentRemove, authority knowledgesdk.Authority) error {
	return value.invoke(ctx, knowledgesdk.ActionPageAttachmentsRemove, struct {
		LibraryID  string                                     `json:"library_id"`
		PageID     string                                     `json:"page_id"`
		DocumentID string                                     `json:"document_id"`
		Input      knowledgesdk.KnowledgePageAttachmentRemove `json:"input"`
		Authority  knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, documentID, input, authority}, nil)
}

func (value *serviceClient) KnowledgePageTasks(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTaskPage, error) {
	var output knowledgesdk.KnowledgePageTaskPage
	err := value.invoke(ctx, knowledgesdk.ActionPageTasksList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageTask(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageTaskCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTask, error) {
	var output knowledgesdk.KnowledgePageTask
	err := value.invoke(ctx, knowledgesdk.ActionPageTasksCreate, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePageTaskCreate `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) UpdateKnowledgePageTask(ctx context.Context, libraryID, pageID, taskID string, input knowledgesdk.KnowledgePageTaskUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTask, error) {
	var output knowledgesdk.KnowledgePageTask
	err := value.invoke(ctx, knowledgesdk.ActionPageTasksUpdate, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		TaskID    string                               `json:"task_id"`
		Input     knowledgesdk.KnowledgePageTaskUpdate `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, taskID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePagePolls(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePollPage, error) {
	var output knowledgesdk.KnowledgePagePollPage
	err := value.invoke(ctx, knowledgesdk.ActionPagePollsList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePagePoll(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePagePollCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePollView, error) {
	var output knowledgesdk.KnowledgePagePollView
	err := value.invoke(ctx, knowledgesdk.ActionPagePollsCreate, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePagePollCreate `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) UpdateKnowledgePagePoll(ctx context.Context, libraryID, pageID, pollID string, input knowledgesdk.KnowledgePagePollUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePollView, error) {
	var output knowledgesdk.KnowledgePagePollView
	err := value.invoke(ctx, knowledgesdk.ActionPagePollsUpdate, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		PollID    string                               `json:"poll_id"`
		Input     knowledgesdk.KnowledgePagePollUpdate `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, pollID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) VoteKnowledgePagePoll(ctx context.Context, libraryID, pageID, pollID string, input knowledgesdk.KnowledgePagePollVote, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePagePollView, error) {
	var output knowledgesdk.KnowledgePagePollView
	err := value.invoke(ctx, knowledgesdk.ActionPagePollsVote, struct {
		LibraryID string                             `json:"library_id"`
		PageID    string                             `json:"page_id"`
		PollID    string                             `json:"poll_id"`
		Input     knowledgesdk.KnowledgePagePollVote `json:"input"`
		Authority knowledgesdk.Authority             `json:"authority"`
	}{libraryID, pageID, pollID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SubmitKnowledgePageForm(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeFormSubmit, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeFormSubmission, error) {
	var output knowledgesdk.KnowledgeFormSubmission
	err := value.invoke(ctx, knowledgesdk.ActionPageFormSubmit, struct {
		LibraryID string                           `json:"library_id"`
		PageID    string                           `json:"page_id"`
		Input     knowledgesdk.KnowledgeFormSubmit `json:"input"`
		Authority knowledgesdk.Authority           `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageFormSubmissions(ctx context.Context, libraryID, pageID, before string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeFormSubmissionPage, error) {
	var output knowledgesdk.KnowledgeFormSubmissionPage
	err := value.invoke(ctx, knowledgesdk.ActionPageFormResponsesList, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Before    string                 `json:"before"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, before, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageFormAnalytics(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeFormAnalytics, error) {
	var output knowledgesdk.KnowledgeFormAnalytics
	err := value.invoke(ctx, knowledgesdk.ActionPageFormAnalyticsGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) ExportKnowledgePageFormSubmissions(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeFormExport, error) {
	var output knowledgesdk.KnowledgeFormExport
	err := value.invoke(ctx, knowledgesdk.ActionPageFormResponsesExport, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageTemplates(ctx context.Context, libraryID, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplatePage, error) {
	var output knowledgesdk.KnowledgePageTemplatePage
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplatesList, struct {
		LibraryID string                 `json:"library_id"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, after, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgePageTemplate(ctx context.Context, libraryID string, input knowledgesdk.KnowledgePageTemplateCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplate, error) {
	var output knowledgesdk.KnowledgePageTemplate
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplatesCreate, struct {
		LibraryID string                                   `json:"library_id"`
		Input     knowledgesdk.KnowledgePageTemplateCreate `json:"input"`
		Authority knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) UpdateKnowledgePageTemplate(ctx context.Context, libraryID, templateID string, input knowledgesdk.KnowledgePageTemplateUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplate, error) {
	var output knowledgesdk.KnowledgePageTemplate
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplatesUpdate, struct {
		LibraryID  string                                   `json:"library_id"`
		TemplateID string                                   `json:"template_id"`
		Input      knowledgesdk.KnowledgePageTemplateUpdate `json:"input"`
		Authority  knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, templateID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageTemplateVersions(ctx context.Context, libraryID, templateID string, before int64, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplateVersionPage, error) {
	var output knowledgesdk.KnowledgePageTemplateVersionPage
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplateVersionsList, struct {
		LibraryID  string                 `json:"library_id"`
		TemplateID string                 `json:"template_id"`
		Before     int64                  `json:"before"`
		Limit      int                    `json:"limit"`
		Authority  knowledgesdk.Authority `json:"authority"`
	}{libraryID, templateID, before, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageTemplateVersion(ctx context.Context, libraryID, templateID string, version int64, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplateVersion, error) {
	var output knowledgesdk.KnowledgePageTemplateVersion
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplateVersionsGet, struct {
		LibraryID  string                 `json:"library_id"`
		TemplateID string                 `json:"template_id"`
		Version    int64                  `json:"version"`
		Authority  knowledgesdk.Authority `json:"authority"`
	}{libraryID, templateID, version, authority}, &output)
	return output, err
}

func (value *serviceClient) InstantiateKnowledgePageTemplate(ctx context.Context, libraryID, templateID string, input knowledgesdk.KnowledgePageTemplateInstantiate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageTemplateInstance, error) {
	var output knowledgesdk.KnowledgePageTemplateInstance
	err := value.invoke(ctx, knowledgesdk.ActionPageTemplatesInstantiate, struct {
		LibraryID  string                                        `json:"library_id"`
		TemplateID string                                        `json:"template_id"`
		Input      knowledgesdk.KnowledgePageTemplateInstantiate `json:"input"`
		Authority  knowledgesdk.Authority                        `json:"authority"`
	}{libraryID, templateID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SearchKnowledgePrincipals(ctx context.Context, query, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePrincipalPage, error) {
	var output knowledgesdk.KnowledgePrincipalPage
	err := value.invoke(ctx, knowledgesdk.ActionPrincipalsSearch, struct {
		Query     string                 `json:"query"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{query, after, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) ResolveKnowledgePrincipals(ctx context.Context, userIDs []string, authority knowledgesdk.Authority) ([]knowledgesdk.KnowledgePrincipal, error) {
	var output []knowledgesdk.KnowledgePrincipal
	err := value.invoke(ctx, knowledgesdk.ActionPrincipalsResolve, struct {
		UserIDs   []string               `json:"user_ids"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{userIDs, authority}, &output)
	return output, err
}

func (value *serviceClient) ResolveKnowledgeOrganizationPrincipals(ctx context.Context, organizationIDs []string, authority knowledgesdk.Authority) ([]knowledgesdk.KnowledgeOrganizationPrincipal, error) {
	var output []knowledgesdk.KnowledgeOrganizationPrincipal
	err := value.invoke(ctx, knowledgesdk.ActionOrganizationPrincipalsResolve, struct {
		OrganizationIDs []string               `json:"organization_ids"`
		Authority       knowledgesdk.Authority `json:"authority"`
	}{organizationIDs, authority}, &output)
	return output, err
}

func (value *serviceClient) SearchKnowledgeLibraryMembers(ctx context.Context, libraryID, query, after string, limit int, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePrincipalPage, error) {
	var output knowledgesdk.KnowledgePrincipalPage
	err := value.invoke(ctx, knowledgesdk.ActionLibraryMembersSearch, struct {
		LibraryID string                 `json:"library_id"`
		Query     string                 `json:"query"`
		After     string                 `json:"after"`
		Limit     int                    `json:"limit"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, query, after, limit, authority}, &output)
	return output, err
}

func (value *serviceClient) ExplainKnowledgeAccess(ctx context.Context, libraryID string, input knowledgesdk.KnowledgeAccessExplainRequest, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeAccessExplanation, error) {
	var output knowledgesdk.KnowledgeAccessExplanation
	err := value.invoke(ctx, knowledgesdk.ActionAccessExplain, struct {
		LibraryID string                                     `json:"library_id"`
		Input     knowledgesdk.KnowledgeAccessExplainRequest `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageAccess(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageAccessPolicy, error) {
	var output knowledgesdk.KnowledgePageAccessPolicy
	err := value.invoke(ctx, knowledgesdk.ActionPageAccessGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageAccess(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageAccessWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageAccessPolicy, error) {
	var output knowledgesdk.KnowledgePageAccessPolicy
	err := value.invoke(ctx, knowledgesdk.ActionPageAccessSet, struct {
		LibraryID string                                `json:"library_id"`
		PageID    string                                `json:"page_id"`
		Input     knowledgesdk.KnowledgePageAccessWrite `json:"input"`
		Authority knowledgesdk.Authority                `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageContentPolicy(ctx context.Context, libraryID, pageID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageContentPolicy, error) {
	var output knowledgesdk.KnowledgePageContentPolicy
	err := value.invoke(ctx, knowledgesdk.ActionPageContentPolicyGet, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) SetKnowledgePageContentPolicy(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageContentPolicyWrite, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageContentPolicy, error) {
	var output knowledgesdk.KnowledgePageContentPolicy
	err := value.invoke(ctx, knowledgesdk.ActionPageContentPolicySet, struct {
		LibraryID string                                       `json:"library_id"`
		PageID    string                                       `json:"page_id"`
		Input     knowledgesdk.KnowledgePageContentPolicyWrite `json:"input"`
		Authority knowledgesdk.Authority                       `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) AuthorizeKnowledgePageContentOperation(ctx context.Context, libraryID, pageID, operation string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageContentDecision, error) {
	var output knowledgesdk.KnowledgePageContentDecision
	action := map[string]string{
		knowledgesdk.KnowledgePageContentOperationCopy:      knowledgesdk.ActionPageContentCopy,
		knowledgesdk.KnowledgePageContentOperationPrint:     knowledgesdk.ActionPagePrint,
		knowledgesdk.KnowledgePageContentOperationExport:    knowledgesdk.ActionPageExport,
		knowledgesdk.KnowledgePageContentOperationDuplicate: knowledgesdk.ActionPageDuplicate,
	}[operation]
	if action == "" {
		return output, &knowledgesdk.Error{Class: "bad_request", Code: "knowledge.page_content_operation_invalid"}
	}
	err := value.invoke(ctx, action, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, authority}, &output)
	return output, err
}

func (value *serviceClient) DownloadKnowledgePageDocument(ctx context.Context, libraryID, pageID, documentID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentDownload, error) {
	request := knowledgecontract.SaaSDocumentDownloadRequest{Library: libraryID, ID: documentID, PageID: pageID, Authority: authority}
	var document knowledgesdk.KnowledgeDocument
	grants, err := value.client.invokeWithGrants(ctx, knowledgecontract.SaaSDocumentDownloadPrepareOperation, request, &document, value.options)
	if err != nil {
		return knowledgesdk.KnowledgeDocumentDownload{}, err
	}
	request.Grants = grants
	header, err := encodeTransferHeader(request)
	if err != nil {
		return knowledgesdk.KnowledgeDocumentDownload{}, err
	}
	return knowledgesdk.KnowledgeDocumentDownload{Document: document, Content: &rangeDocumentReader{client: value.client, ctx: ctx, document: document, transferHeader: header}}, nil
}

func (value *serviceClient) ConnectKnowledgeCollaboration(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeCollaborationConnect, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeCollaborationSession, error) {
	var output knowledgesdk.KnowledgeCollaborationSession
	err := value.invoke(ctx, knowledgesdk.ActionPageCollaborationConnect, struct {
		LibraryID string                                     `json:"library_id"`
		PageID    string                                     `json:"page_id"`
		Input     knowledgesdk.KnowledgeCollaborationConnect `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) SyncKnowledgeCollaboration(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeCollaborationSync, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeCollaborationSyncResult, error) {
	var output knowledgesdk.KnowledgeCollaborationSyncResult
	err := value.invoke(ctx, knowledgesdk.ActionPageCollaborationSync, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		Input     knowledgesdk.KnowledgeCollaborationSync `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgeCollaborationEvents(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeCollaborationEvents, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeCollaborationEventBatch, error) {
	var output knowledgesdk.KnowledgeCollaborationEventBatch
	err := value.invoke(ctx, knowledgesdk.ActionPageCollaborationEvents, struct {
		LibraryID string                                    `json:"library_id"`
		PageID    string                                    `json:"page_id"`
		Input     knowledgesdk.KnowledgeCollaborationEvents `json:"input"`
		Authority knowledgesdk.Authority                    `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgePageLiveEvents(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgePageLiveEvents, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageLiveEventBatch, error) {
	var output knowledgesdk.KnowledgePageLiveEventBatch
	err := value.invoke(ctx, knowledgesdk.ActionPageLiveEvents, struct {
		LibraryID string                               `json:"library_id"`
		PageID    string                               `json:"page_id"`
		Input     knowledgesdk.KnowledgePageLiveEvents `json:"input"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) CloseKnowledgeCollaboration(ctx context.Context, libraryID, pageID, sessionID string, authority knowledgesdk.Authority) error {
	return value.invoke(ctx, knowledgesdk.ActionPageCollaborationClose, struct {
		LibraryID string                 `json:"library_id"`
		PageID    string                 `json:"page_id"`
		SessionID string                 `json:"session_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, sessionID, authority}, nil)
}

func (value *serviceClient) PublishKnowledgeCollaboration(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeCollaborationPublish, authority knowledgesdk.Authority) (knowledgesdk.KnowledgePageView, error) {
	var output knowledgesdk.KnowledgePageView
	err := value.invoke(ctx, knowledgesdk.ActionPageCollaborationPublish, struct {
		LibraryID string                                     `json:"library_id"`
		PageID    string                                     `json:"page_id"`
		Input     knowledgesdk.KnowledgeCollaborationPublish `json:"input"`
		Authority knowledgesdk.Authority                     `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgeWebhookSubscriptions(ctx context.Context, libraryID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeWebhookSubscriptionPage, error) {
	var output knowledgesdk.KnowledgeWebhookSubscriptionPage
	err := value.invoke(ctx, knowledgesdk.ActionWebhookSubscriptionsList, struct {
		LibraryID string                 `json:"library_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgeWebhookSubscription(ctx context.Context, libraryID string, input knowledgesdk.KnowledgeWebhookSubscriptionCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeWebhookSubscription, error) {
	var output knowledgesdk.KnowledgeWebhookSubscription
	err := value.invoke(ctx, knowledgesdk.ActionWebhookSubscriptionsCreate, struct {
		LibraryID string                                          `json:"library_id"`
		Input     knowledgesdk.KnowledgeWebhookSubscriptionCreate `json:"input"`
		Authority knowledgesdk.Authority                          `json:"authority"`
	}{libraryID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) UpdateKnowledgeWebhookSubscription(ctx context.Context, libraryID, subscriptionID string, input knowledgesdk.KnowledgeWebhookSubscriptionUpdate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeWebhookSubscription, error) {
	var output knowledgesdk.KnowledgeWebhookSubscription
	err := value.invoke(ctx, knowledgesdk.ActionWebhookSubscriptionsUpdate, struct {
		LibraryID      string                                          `json:"library_id"`
		SubscriptionID string                                          `json:"subscription_id"`
		Input          knowledgesdk.KnowledgeWebhookSubscriptionUpdate `json:"input"`
		Authority      knowledgesdk.Authority                          `json:"authority"`
	}{libraryID, subscriptionID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgeWebhookEvents(ctx context.Context, libraryID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeWebhookEventPage, error) {
	var output knowledgesdk.KnowledgeWebhookEventPage
	err := value.invoke(ctx, knowledgesdk.ActionWebhookEventsList, struct {
		LibraryID string                 `json:"library_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, authority}, &output)
	return output, err
}

func (value *serviceClient) ReplayKnowledgeWebhookEvent(ctx context.Context, libraryID, eventID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeWebhookEventRecord, error) {
	var output knowledgesdk.KnowledgeWebhookEventRecord
	err := value.invoke(ctx, knowledgesdk.ActionWebhookEventsReplay, struct {
		LibraryID string                 `json:"library_id"`
		EventID   string                 `json:"event_id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{libraryID, eventID, authority}, &output)
	return output, err
}

func (value *serviceClient) ResolveKnowledgeAgentReference(ctx context.Context, reference knowledgesdk.KnowledgeAgentReference, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeAgentResolvedReference, error) {
	var output knowledgesdk.KnowledgeAgentResolvedReference
	err := value.invoke(ctx, knowledgesdk.ActionAgentReferencesResolve, struct {
		Reference knowledgesdk.KnowledgeAgentReference `json:"reference"`
		Authority knowledgesdk.Authority               `json:"authority"`
	}{reference, authority}, &output)
	return output, err
}

func (value *serviceClient) ResolveKnowledgeApplicationCard(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeApplicationCardResolve, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeApplicationCard, error) {
	var output knowledgesdk.KnowledgeApplicationCard
	err := value.invoke(ctx, knowledgesdk.ActionApplicationCardsResolve, struct {
		LibraryID string                                       `json:"library_id"`
		PageID    string                                       `json:"page_id"`
		Input     knowledgesdk.KnowledgeApplicationCardResolve `json:"input"`
		Authority knowledgesdk.Authority                       `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) CreateKnowledgeAIExecution(ctx context.Context, libraryID, pageID string, input knowledgesdk.KnowledgeAIExecutionCreate, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeAIExecution, error) {
	var output knowledgesdk.KnowledgeAIExecution
	err := value.invoke(ctx, knowledgesdk.ActionAIExecutionsCreate, struct {
		LibraryID string                                  `json:"library_id"`
		PageID    string                                  `json:"page_id"`
		Input     knowledgesdk.KnowledgeAIExecutionCreate `json:"input"`
		Authority knowledgesdk.Authority                  `json:"authority"`
	}{libraryID, pageID, input, authority}, &output)
	return output, err
}

func (value *serviceClient) KnowledgeAIExecution(ctx context.Context, libraryID, pageID, executionID string, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeAIExecution, error) {
	var output knowledgesdk.KnowledgeAIExecution
	err := value.invoke(ctx, knowledgesdk.ActionAIExecutionsGet, struct {
		LibraryID   string                 `json:"library_id"`
		PageID      string                 `json:"page_id"`
		ExecutionID string                 `json:"execution_id"`
		Authority   knowledgesdk.Authority `json:"authority"`
	}{libraryID, pageID, executionID, authority}, &output)
	return output, err
}

func (value *serviceClient) ConfirmKnowledgeAIExecution(ctx context.Context, libraryID, pageID, executionID string, input knowledgesdk.KnowledgeAIExecutionConfirm, authority knowledgesdk.Authority) (knowledgesdk.KnowledgeAIExecution, error) {
	var output knowledgesdk.KnowledgeAIExecution
	err := value.invoke(ctx, knowledgesdk.ActionAIExecutionsConfirm, struct {
		LibraryID   string                                   `json:"library_id"`
		PageID      string                                   `json:"page_id"`
		ExecutionID string                                   `json:"execution_id"`
		Input       knowledgesdk.KnowledgeAIExecutionConfirm `json:"input"`
		Authority   knowledgesdk.Authority                   `json:"authority"`
	}{libraryID, pageID, executionID, input, authority}, &output)
	return output, err
}

var _ knowledgecontract.Service = (*serviceClient)(nil)
var _ = json.Valid

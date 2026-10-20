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

var _ knowledgecontract.Service = (*serviceClient)(nil)
var _ = json.Valid

// Package remote implements the official Knowledge SaaS client. It owns HTTP
// transport only; Knowledge business behavior and persistence remain in the
// independently deployed Knowledge service.
package remote

import (
	"context"
	"fmt"
	"net/http"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

type Factory struct{ config Config }

func NewFactory(config Config) *Factory { return &Factory{config: config} }

func (factory *Factory) OpenSaaS(ctx context.Context, ref knowledgecontract.ApplicationRef) (saashost.Binding, error) {
	if err := ref.Validate(); err != nil {
		return nil, err
	}
	client, err := newClient(factory.config, ref.RuntimeID)
	if err != nil {
		return nil, err
	}
	var descriptor knowledgecontract.Descriptor
	if err := client.request(ctx, http.MethodGet, knowledgecontract.SaaSDiscoveryPath, nil, &descriptor); err != nil {
		return nil, err
	}
	if err := descriptor.Validate(); err != nil {
		return nil, remoteError("unavailable", "knowledge.protocol_incompatible", false, err)
	}
	if descriptor.Audience != ref.RuntimeID {
		return nil, remoteError("forbidden", "knowledge.runtime_mismatch", false, nil)
	}
	runtime := &runtimeClient{client: client}
	return &binding{descriptor: descriptor, client: client, runtime: runtime}, nil
}

type binding struct {
	descriptor knowledgecontract.Descriptor
	client     *client
	runtime    *runtimeClient
}

func (value *binding) Descriptor() knowledgecontract.Descriptor { return value.descriptor }
func (value *binding) Runtime() knowledgecontract.Runtime       { return value.runtime }
func (*binding) Close(context.Context) error                    { return nil }

type runtimeClient struct {
	client  *client
	options knowledgecontract.Options
}

func (value *runtimeClient) Validate(options *knowledgecontract.Options) error {
	if options == nil {
		return fmt.Errorf("Knowledge options are required")
	}
	return nil
}
func (value *runtimeClient) Prepare(_ string, options knowledgecontract.Options) (knowledgecontract.Knowledge, error) {
	value.options = options
	return remoteKnowledge{client: value.client, options: options}, nil
}
func (value *runtimeClient) Activate(runtimeID string, options knowledgecontract.Options) error {
	if runtimeID != value.client.runtimeID {
		return remoteError("forbidden", "knowledge.runtime_mismatch", false, nil)
	}
	value.options = options
	return value.client.invoke(context.Background(), "runtime.activate", struct{}{}, nil)
}
func (value *runtimeClient) NewService(runtimeID string, options knowledgecontract.Options) knowledgecontract.Service {
	value.options = options
	return &serviceClient{client: value.client, options: options, runtimeID: runtimeID}
}

type remoteKnowledge struct {
	client  *client
	options knowledgecontract.Options
}

func (value remoteKnowledge) SearchKnowledge(ctx context.Context, query string, authority knowledgesdk.Authority) (knowledgesdk.Result, error) {
	var output knowledgesdk.Result
	err := value.client.invokeWithOptions(ctx, "runtime.knowledge.search", struct {
		Query     string                 `json:"query"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{query, authority}, &output, value.options)
	return output, err
}
func (value remoteKnowledge) ReadKnowledge(ctx context.Context, id string, authority knowledgesdk.Authority) (knowledgesdk.Result, error) {
	var output knowledgesdk.Result
	err := value.client.invokeWithOptions(ctx, "runtime.knowledge.read", struct {
		ID        string                 `json:"id"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{id, authority}, &output, value.options)
	return output, err
}
func (value remoteKnowledge) RevalidateKnowledge(ctx context.Context, input knowledgesdk.Result, authority knowledgesdk.Authority) error {
	return value.client.invokeWithOptions(ctx, "runtime.knowledge.revalidate", struct {
		Input     knowledgesdk.Result    `json:"input"`
		Authority knowledgesdk.Authority `json:"authority"`
	}{input, authority}, nil, value.options)
}

var _ saashost.Factory = (*Factory)(nil)
var _ saashost.Binding = (*binding)(nil)
var _ knowledgecontract.Runtime = (*runtimeClient)(nil)

package knowledgesdk

import (
	"fmt"
	"strings"
)

// KnowledgeOpenAPIDocument is generated from the same source-owned Action
// manifest used by the executable HTTP adapter. It deliberately exposes only
// the supported integration surface instead of every browser-internal route.
func KnowledgeOpenAPIDocument() (map[string]any, error) {
	actions, err := AuthorizationActions()
	if err != nil {
		return nil, err
	}
	exposed := openAPIActionSet()
	paths := map[string]any{}
	projected := make(map[string]struct{}, len(exposed))
	for _, action := range actions {
		if _, ok := exposed[action.Key]; !ok {
			continue
		}
		if action.HTTP == nil {
			return nil, fmt.Errorf("OpenAPI Action %q has no HTTP binding", action.Key)
		}
		pathItem, _ := paths[action.HTTP.RouteTemplate].(map[string]any)
		if pathItem == nil {
			pathItem = map[string]any{}
			paths[action.HTTP.RouteTemplate] = pathItem
		}
		method := strings.ToLower(strings.TrimSpace(action.HTTP.Method))
		projected[action.Key] = struct{}{}
		pathItem[method] = map[string]any{
			"operationId":       strings.ReplaceAll(action.Key, ".", "_"),
			"summary":           action.Label,
			"tags":              []string{action.CapabilityLabel},
			"x-domainry-action": action.Key,
			"security": []any{
				map[string]any{"browserBearer": []string{}},
				map[string]any{"applicationServiceBearer": []string{action.Key}},
			},
			"responses": map[string]any{
				"200": map[string]any{"description": "Successful authorized response"},
				"400": map[string]any{"description": "Invalid request"},
				"401": map[string]any{"description": "Authentication required"},
				"403": map[string]any{"description": "Current authorization denied"},
				"409": map[string]any{"description": "Revision or content conflict"},
			},
		}
	}
	if len(projected) != len(exposed) {
		return nil, fmt.Errorf("Knowledge OpenAPI projection is incomplete: projected %d of %d Actions", len(projected), len(exposed))
	}
	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "Domainry Knowledge API",
			"version":     "1.0.0",
			"description": "Authorized API for page trees, published or explicitly requested draft content, comments, tasks, multidimensional tables, and search.",
		},
		"paths": paths,
		"components": map[string]any{"securitySchemes": map[string]any{
			"browserBearer": map[string]any{
				"type": "http", "scheme": "bearer", "bearerFormat": "IdentityAccessToken",
			},
			"applicationServiceBearer": map[string]any{
				"type": "http", "scheme": "bearer", "bearerFormat": "IdentityApplicationServiceToken",
				"description": "The token must carry the exact x-domainry-action grant and is re-authorized against current Knowledge permissions and ACLs.",
			},
		}},
	}, nil
}

func openAPIActionSet() map[string]struct{} {
	keys := []string{
		ActionOpenAPIGet,
		ActionPagesList, ActionPagesGet, ActionPagesCreate, ActionPagesUpdate, ActionPagesPublish,
		ActionPageReferencesList, ActionPageReferencesCreate,
		ActionPageCommentsList, ActionPageCommentsCreate, ActionPageCommentsUpdate, ActionPageCommentsResolve, ActionPageCommentsDelete, ActionPageCommentsReact,
		ActionPageTasksList, ActionPageTasksCreate, ActionPageTasksUpdate,
		ActionPageNativeOperationsApply, ActionPageGridAutomationExecute,
		ActionPageSearch,
		ActionAgentReferencesResolve,
		ActionAIExecutionsCreate, ActionAIExecutionsGet, ActionAIExecutionsConfirm,
		ActionApplicationCardsResolve,
	}
	set := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}

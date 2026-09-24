package contract

import (
	"context"

	sdk "github.com/domainry/domainry-agent-sdk"
)

// ArtifactExportReader reads an existing immutable export with current artifact
// read access. It never creates another export or grants artifact mutation rights.
type ArtifactExportReader interface {
	ReadArtifactExport(context.Context, string, sdk.ConversationAuthority) (sdk.ConversationArtifactDownload, error)
}

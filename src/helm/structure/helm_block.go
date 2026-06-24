package structure

import (
	"regexp"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

// HelmFramework is the framework name reported by Helm blocks.
const HelmFramework = "Helm"

// LabelsAttributeName is the YAML key under metadata that yor writes tags into (D3 — labels only).
const LabelsAttributeName = "labels"

// k8sLabelValueRegex implements the Kubernetes label value rule (D3.1):
// an empty value is allowed; otherwise it must be at most 63 chars and match the
// alphanumeric-bookended charset below. Key validity is handled separately.
var k8sLabelValueRegex = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9._-]*[A-Za-z0-9])?$`)

// k8sLabelKeyNameRegex validates the name portion of a label key (after an optional DNS prefix).
var k8sLabelKeyNameRegex = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9._-]*[A-Za-z0-9])?$`)

const k8sLabelValueMaxLength = 63
const k8sLabelKeyNameMaxLength = 63

// LabelSite is a single metadata.labels target within a document (D7 — top-level label scope).
// A document has exactly one site: the TOP-LEVEL metadata. Nested pod-template / jobTemplate /
// volumeClaimTemplate metadata and `spec.selector` metadata are NEVER produced — yor tags are
// written only to the resource's own metadata.labels (see D7).
type LabelSite struct {
	// MetadataLine is the 0-based line of this site's `metadata:` key.
	MetadataLine int
	// MetadataIndent is the indentation width of the `metadata:` key.
	MetadataIndent int
	// LabelsLine is the 0-based line of an existing `labels:` key under this metadata, or -1.
	LabelsLine int
	// Description identifies the site (e.g. "metadata", "spec.template.metadata") for logging/tests.
	Description string
}

// HelmBlock is a single taggable Kubernetes resource inside a Helm template. Usually that is a whole
// YAML document; for a Kubernetes List document (kind: List or any typed `*List`) it is instead a
// single entry of the wrapper's `items:` array (see helm_parser.go).
// Identity (D11) is file path + kind + document ordinal (+ item index for List entries), captured by
// Name/Type, DocOrdinal and ItemIndex.
type HelmBlock struct {
	structure.Block
	// DocOrdinal is the 0-based index of this document within its (possibly multi-doc) file.
	DocOrdinal int
	// ItemIndex is the 0-based index of this resource within a Kubernetes List wrapper's `items:`
	// array, or -1 when the block is a plain (non-List) document.
	ItemIndex int
	// Sites are the metadata.labels targets for this resource (D7 — the resource's own metadata).
	// For a List item this is the item's OWN nested metadata, never the wrapper's.
	Sites []LabelSite
}

// GetFramework implements IYamlBlock.
func (b *HelmBlock) GetFramework() string {
	return HelmFramework
}

// GetSeparator implements IBlock; Helm labels are written one-per-line.
func (b *HelmBlock) GetSeparator() string {
	return "\n"
}

// GetTagsLines returns the existing labels line range for this block.
func (b *HelmBlock) GetTagsLines() structure.Lines {
	return b.Block.TagLines
}

// UpdateTags merges new and existing tags onto the block (D3 — labels only, append-only,
// preserving an existing yor_trace via the shared MergeTags semantics, D11). Statically-invalid
// labels are filtered out before write by GetValidLabels (D3.1) — except values containing a Helm
// template directive (`{{ }}`), which are written verbatim (Strategy 2); the writer emits warnings
// for dropped labels so this method stays side-effect free.
func (b *HelmBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}
	_ = b.MergeTags()
}

// IsValidK8sLabel reports whether a tag is a valid Kubernetes label (D3.1).
// Both the key (optionally prefixed with a DNS subdomain) and the value must comply.
func IsValidK8sLabel(key string, value string) bool {
	return isValidK8sLabelKey(key) && isValidK8sLabelValue(value)
}

func isValidK8sLabelValue(value string) bool {
	if value == "" {
		// An empty label value is valid in Kubernetes.
		return true
	}
	if len(value) > k8sLabelValueMaxLength {
		return false
	}
	return k8sLabelValueRegex.MatchString(value)
}

func isValidK8sLabelKey(key string) bool {
	if key == "" {
		return false
	}
	name := key
	if idx := indexByte(key, '/'); idx >= 0 {
		prefix := key[:idx]
		name = key[idx+1:]
		if prefix == "" || len(prefix) > 253 {
			return false
		}
	}
	if name == "" || len(name) > k8sLabelKeyNameMaxLength {
		return false
	}
	return k8sLabelKeyNameRegex.MatchString(name)
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// tagsToWrite returns ONLY the yor-managed tags that should be written as labels — never the
// document's pre-existing foreign labels (e.g. `app`, `team`). yor tags are written to the
// top-level metadata.labels only (D7). It is the set of NewTags, except that an existing
// `yor_trace` is preserved verbatim instead of the freshly minted one (idempotency, D11).
func (b *HelmBlock) tagsToWrite() []tags.ITag {
	var existingTrace tags.ITag
	for _, t := range b.ExitingTags {
		if t.GetKey() == tags.YorTraceTagKey {
			existingTrace = t
		}
	}
	out := make([]tags.ITag, 0, len(b.NewTags))
	for _, t := range b.NewTags {
		if t.GetKey() == tags.YorTraceTagKey && existingTrace != nil {
			out = append(out, existingTrace)
			continue
		}
		out = append(out, t)
	}
	// If yor_trace exists in the file but was not (re)computed in NewTags, still preserve it.
	if existingTrace != nil {
		found := false
		for _, t := range out {
			if t.GetKey() == tags.YorTraceTagKey {
				found = true
				break
			}
		}
		if !found {
			out = append(out, existingTrace)
		}
	}
	return out
}

// GetValidLabels returns the yor-managed tags-to-write that pass Kubernetes label validation,
// plus the tags that were dropped so the caller can warn (D3.1).
//
// Exception (Strategy 2): a tag whose value is a Helm template directive (`{{ ... }}`) is written
// verbatim even though it is not a statically valid K8s label value. This lets yor_name mirror a
// templated metadata.name so the RENDERED label tracks the resource's real name; the rendered value
// may be invalid, which is an accepted product choice.
func (b *HelmBlock) GetValidLabels() (valid []tags.ITag, dropped []tags.ITag) {
	for _, t := range b.tagsToWrite() {
		if isHelmTemplatedValue(t.GetValue()) && isValidK8sLabelKey(t.GetKey()) {
			valid = append(valid, t)
			continue
		}
		if IsValidK8sLabel(t.GetKey(), t.GetValue()) {
			valid = append(valid, t)
		} else {
			dropped = append(dropped, t)
		}
	}
	return valid, dropped
}

package structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

const (
	chartFileName   = "Chart.yaml"
	valuesFileName  = "values.yaml"
	helpersFileName = "_helpers.tpl"
	notesFileName   = "NOTES.txt"
	templatesDir    = "templates"
	chartsDir       = "charts"
)

// HelmParser implements common.IParser for Helm charts.
// It never invokes Helm and never strict-unmarshals templates (D1); it scans text so that
// `{{ ... }}` directives are preserved verbatim.
type HelmParser struct {
	rootDir              string
	skippedByCommentList []string
}

// Name implements IParser.
func (p *HelmParser) Name() string {
	return HelmFramework
}

// Init implements IParser.
func (p *HelmParser) Init(rootDir string, _ map[string]string) {
	p.rootDir = rootDir
}

// Close implements IParser.
func (p *HelmParser) Close() {}

// GetSkippedDirs implements IParser; chart dependency dirs are skipped (D9).
func (p *HelmParser) GetSkippedDirs() []string {
	return []string{chartsDir}
}

// GetSupportedFileExtensions implements IParser.
func (p *HelmParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension}
}

// GetSkipResourcesByComment implements IParser (D12).
func (p *HelmParser) GetSkipResourcesByComment() []string {
	return p.skippedByCommentList
}

// ValidFile implements IParser. A file is "Helm" if it lives under a templates/ directory AND a
// Chart.yaml exists at the chart root (D9). Chart.yaml apiVersion is intentionally NOT inspected —
// detection is version-tolerant (accepts any apiVersion). Helm-internal files are excluded.
func (p *HelmParser) ValidFile(filePath string) bool {
	base := filepath.Base(filePath)
	if base == valuesFileName || base == helpersFileName || base == notesFileName || base == chartFileName {
		return false
	}
	if strings.HasSuffix(base, ".tpl") {
		return false
	}
	return isUnderTemplatesWithChart(filePath)
}

// isUnderTemplatesWithChart reports whether filePath is inside a templates/ directory that has a
// sibling Chart.yaml at the chart root (the parent of templates/).
func isUnderTemplatesWithChart(filePath string) bool {
	abs, err := filepath.Abs(filePath)
	if err != nil {
		abs = filePath
	}
	dir := filepath.Dir(abs)
	for {
		if filepath.Base(dir) == templatesDir {
			chartRoot := filepath.Dir(dir)
			if fileExists(filepath.Join(chartRoot, chartFileName)) {
				return true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ParseFile implements IParser. It splits the file into `---`-separated documents (D2) and, for
// each document, locates every metadata.labels site (top-level + nested) via text scanning (D4/D7).
// A document is taggable iff at least one literal `metadata:` block is locatable. `{{ }}` lines are
// never interpreted.
func (p *HelmParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	if !p.ValidFile(filePath) {
		return nil, nil
	}
	// #nosec G304 - file is from user
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSuffix(string(contentBytes), "\n"), "\n")

	parsedBlocks := make([]structure.IBlock, 0)
	docStarts := documentStartLines(lines)
	for ordinal, start := range docStarts {
		end := len(lines) - 1
		if ordinal+1 < len(docStarts) {
			end = docStarts[ordinal+1] - 1
		}
		blocks := p.parseDocument(filePath, lines, start, end, ordinal)
		for _, block := range blocks {
			parsedBlocks = append(parsedBlocks, block)
		}
	}
	return parsedBlocks, nil
}

// documentStartLines returns the 0-based start line of each YAML document in the file.
func documentStartLines(lines []string) []int {
	starts := []int{0}
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if strings.TrimSpace(line) == "---" {
			starts = append(starts, i+1)
		}
	}
	filtered := make([]int, 0, len(starts))
	for idx, s := range starts {
		end := len(lines) - 1
		if idx+1 < len(starts) {
			end = starts[idx+1] - 2
		}
		if hasContent(lines, s, end) {
			filtered = append(filtered, s)
		}
	}
	if len(filtered) == 0 {
		return []int{0}
	}
	return filtered
}

func hasContent(lines []string, start, end int) bool {
	for i := start; i <= end && i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" || t == "---" || strings.HasPrefix(t, "#") {
			continue
		}
		return true
	}
	return false
}

// parseDocument builds the taggable HelmBlock(s) for a single YAML document, or nil if the document
// has no locatable metadata block (D4). A normal document yields a single block; a Kubernetes List
// document (kind: List or any typed `*List`) yields one block PER `items:` entry, each tagged on its
// own nested metadata. The List wrapper itself is never tagged. Skipped documents are warned.
func (p *HelmParser) parseDocument(filePath string, lines []string, start, end, ordinal int) []*HelmBlock {
	kind := findKind(lines, start, end)

	if isListDocument(lines, start, end, kind) {
		return p.parseListDocument(filePath, lines, start, end, ordinal, kind)
	}

	sites := findLabelSites(lines, start, end)

	skip := isSkipped(lines, start, end)

	if len(sites) == 0 {
		// metadata not locatable: templated include, conditional, or malformed → SKIP + warn (D4).
		logger.Warning(fmt.Sprintf(
			"skipping Helm document with no locatable metadata in %s (kind=%s, doc=%d)",
			filePath, displayKind(kind), ordinal))
		return nil
	}

	// Existing tags are read from the TOP-LEVEL metadata.labels (identity/idempotency baseline, D11).
	var existingTags []tags.ITag
	for _, s := range sites {
		if s.Description == "metadata" {
			existingTags = extractExistingLabels(lines, s.LabelsLine, end)
			break
		}
	}

	// yor_name is derived from the document's TOP-LEVEL metadata.name (this Name feeds
	// GetResourceName() -> the yor_name tagger in code2cloud/yor_name.go). The same yor_name is
	// intentionally written to every label site of this document (top-level + nested); uniqueness
	// is required ACROSS documents, which metadata.name naturally provides.
	resourceName := resourceNameForDocument(lines, sites, kind, ordinal)

	block := &HelmBlock{
		Block: structure.Block{
			FilePath:          filePath,
			ExitingTags:       existingTags,
			RawBlock:          nil,
			IsTaggable:        true,
			TagsAttributeName: LabelsAttributeName,
			// Block line numbers are exposed 1-based to match every other parser
			// (Terraform/CFN/Serverless) and the git tag group / git blame, which
			// are all 1-based. Internally start/end remain 0-based slice indices.
			Lines: structure.Lines{Start: start + 1, End: end + 1},
			Name:  resourceName,
			Type:  displayKind(kind),
		},
		DocOrdinal: ordinal,
		ItemIndex:  -1,
		Sites:      sites,
	}

	if skip {
		p.skippedByCommentList = append(p.skippedByCommentList, block.GetResourceID())
	}
	return []*HelmBlock{block}
}

// parseListDocument builds one HelmBlock per entry of a Kubernetes List wrapper's `items:` array.
// Each item is treated as its own taggable resource: its yor_name is derived from the item's OWN
// metadata.name (reusing the literal / templated / fallback rules) and tags are written to the
// item's OWN nested metadata.labels. The List wrapper's own metadata is NEVER a label site (D7).
func (p *HelmParser) parseListDocument(filePath string, lines []string, start, end, ordinal int, kind string) []*HelmBlock {
	skip := isSkipped(lines, start, end)
	listName := topLevelMetadataNameInRange(lines, start, end)

	itemRanges := listItemRanges(lines, start, end)
	if len(itemRanges) == 0 {
		logger.Warning(fmt.Sprintf(
			"skipping Helm List document with no locatable items in %s (kind=%s, doc=%d)",
			filePath, displayKind(kind), ordinal))
		return nil
	}

	blocks := make([]*HelmBlock, 0, len(itemRanges))
	for itemIndex, r := range itemRanges {
		itemKind := findKindInRange(lines, r.start, r.end)
		sites := findItemLabelSites(lines, r.start, r.end, r.indent)
		if len(sites) == 0 {
			// The item has no locatable metadata: templated include or malformed → SKIP + warn.
			logger.Warning(fmt.Sprintf(
				"skipping Helm List item with no locatable metadata in %s (kind=%s, doc=%d, item=%d)",
				filePath, displayKind(itemKind), ordinal, itemIndex))
			continue
		}

		existingTags := extractExistingLabels(lines, sites[0].LabelsLine, r.end)

		resourceName := resourceNameForItem(lines, sites, itemKind, listName, ordinal, itemIndex)

		block := &HelmBlock{
			Block: structure.Block{
				FilePath:          filePath,
				ExitingTags:       existingTags,
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: LabelsAttributeName,
				// 1-based block lines (see parseDocument) to align with git blame.
				Lines: structure.Lines{Start: r.start + 1, End: r.end + 1},
				Name:  resourceName,
				Type:  displayKind(itemKind),
			},
			DocOrdinal: ordinal,
			ItemIndex:  itemIndex,
			Sites:      sites,
		}

		if skip {
			p.skippedByCommentList = append(p.skippedByCommentList, block.GetResourceID())
		}
		blocks = append(blocks, block)
	}
	return blocks
}

// isListDocument reports whether a document should be handled as a Kubernetes List wrapper. BOTH
// signals must hold: the kind ends in `List` (the reserved collection-kind marker) AND a literal
// top-level `items:` sequence is present. Requiring both avoids two failure modes:
//   - A non-List document that happens to have a top-level `items:` field (e.g. some CRDs) is NOT
//     misread as a List (it lacks a `*List` kind).
//   - A single resource whose kind merely ends in `list` (e.g. a hypothetical `kind: Blocklist`)
//     but has no top-level `items:` falls through to the normal single-resource path and is still
//     tagged correctly, instead of being routed to List handling and silently dropped.
//
// This strengthens — and never violates — the LOCKED SPEC suffix rule: it only ever treats a
// `*List`-named document as a List, and only when its items are actually present.
func isListDocument(lines []string, start, end int, kind string) bool {
	return isListKind(kind) && hasTopLevelItems(lines, start, end)
}

// hasTopLevelItems reports whether a literal top-level `items:` key (indentation 0) exists within
// the document range. The scan is purely indentation/text based and ignores `{{ }}` template lines.
func hasTopLevelItems(lines []string, start, end int) bool {
	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		if indentationWidth(raw) != 0 {
			continue
		}
		if trimmed == "items:" {
			return true
		}
	}
	return false
}

// isListKind reports whether a Kubernetes kind denotes a List wrapper: the generic `List` or any
// typed `*List` variant (e.g. ConfigMapList, DeploymentList). The suffix `List` is matched
// case-insensitively, mirroring the LOCKED SPEC.
func isListKind(kind string) bool {
	if kind == "" {
		return false
	}
	return strings.HasSuffix(strings.ToLower(kind), "list")
}

// findKind returns the value of the top-level `kind:` key within the document range, or "".
func findKind(lines []string, start, end int) string {
	return findKindInRange(lines, start, end)
}

// findKindInRange returns the value of the shallowest-indentation `kind:` key within [start,end],
// or "". For a top-level document the shallowest indent is 0; for a List item it is the item's
// content indent. Only the first `kind:` at the minimum observed indent is returned.
func findKindInRange(lines []string, start, end int) string {
	minIndent := -1
	for i := start; i <= end && i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		indent := indentationWidth(lines[i])
		if strings.HasPrefix(trimmed, "- ") {
			indent += 2
		}
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		return ""
	}
	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		content := trimmed
		indent := indentationWidth(raw)
		if strings.HasPrefix(trimmed, "- ") {
			content = strings.TrimSpace(trimmed[2:])
			indent += 2
		}
		if indent != minIndent {
			continue
		}
		if strings.HasPrefix(content, "kind:") {
			return strings.TrimSpace(strings.TrimPrefix(content, "kind:"))
		}
	}
	return ""
}

// itemRange describes one `items:` entry of a List document: its line range [start,end] and the
// indentation of the item's content keys (one level deeper than the `- ` bullet).
type itemRange struct {
	start  int
	end    int
	indent int
}

// listItemRanges locates the `items:` sequence within a List document range and returns one
// itemRange per top-level entry. Each entry begins at its `- ` bullet line and runs until the next
// bullet at the same indent or the end of the sequence. The scan is purely indentation/text based.
func listItemRanges(lines []string, start, end int) []itemRange {
	itemsLine := -1
	itemsIndent := 0
	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		if indentationWidth(raw) != 0 {
			continue
		}
		if trimmed == "items:" {
			itemsLine = i
			itemsIndent = indentationWidth(raw)
			break
		}
	}
	if itemsLine < 0 {
		return nil
	}

	var ranges []itemRange
	bulletIndent := -1
	for i := itemsLine + 1; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := indentationWidth(raw)
		isBullet := strings.HasPrefix(trimmed, "- ") || trimmed == "-"
		// YAML allows block-sequence bullets to sit at the SAME indentation as their parent
		// `items:` key (the canonical `kubectl ... -o yaml` List style), so a bullet at exactly
		// itemsIndent is still part of the sequence. The sequence ends only when a NON-bullet line
		// returns to itemsIndent-or-shallower (a sibling key of items:), or a bullet appears
		// shallower than itemsIndent.
		if (isBullet && indent < itemsIndent) || (!isBullet && indent <= itemsIndent) {
			// Left the items: sequence.
			break
		}
		if isBullet {
			if bulletIndent < 0 {
				bulletIndent = indent
			}
			if indent != bulletIndent {
				// Deeper bullets belong to a nested sequence inside the current item.
				continue
			}
			// Close the previous range.
			if len(ranges) > 0 {
				ranges[len(ranges)-1].end = i - 1
			}
			ranges = append(ranges, itemRange{start: i, end: end, indent: indent + 2})
		}
	}
	return ranges
}

func displayKind(kind string) string {
	if kind == "" {
		return "Unknown"
	}
	return kind
}

// resourceNameForDocument derives the block Name (which feeds GetResourceName() -> the yor_name
// tag, code2cloud/yor_name.go) from the document's TOP-LEVEL metadata.name.
//
// yor_name is written as a Kubernetes LABEL value. Cases:
//   - Templated name (contains a `{{ }}` directive) -> MIRROR it verbatim (Strategy 2). The Helm
//     directive is preserved (D1) so the RENDERED yor_name tracks the resource's real, rendered
//     metadata.name. This is an intentional product choice: the statically-written value is not a
//     valid K8s label and the rendered value MAY be invalid (e.g. >63 chars), which is accepted.
//   - Literal name that is a valid label value       -> use it verbatim.
//   - Absent name, or a literal that is not a valid label value -> fall back to a deterministic
//     `<lowercased-kind>-<ordinal>` (e.g. "configmap-0"): a valid label, unique per document within
//     a file via the ordinal, and stable across runs for idempotency (D11).
func resourceNameForDocument(lines []string, sites []LabelSite, kind string, ordinal int) string {
	name := topLevelMetadataName(lines, sites)
	if name != "" && (isHelmTemplatedValue(name) || isValidK8sLabelValue(name)) {
		return name
	}
	return fmt.Sprintf("%s-%d", strings.ToLower(displayKind(kind)), ordinal)
}

// resourceNameForItem derives the yor_name for a single Kubernetes List item, reusing the literal /
// templated rules of resourceNameForDocument against the item's OWN metadata.name. When the item has
// no usable name the fallback is deterministic and collision-safe (required for idempotency, D11):
//
//   - With a List-wrapper top-level metadata.name (listName != ""): <list-name>-<kind>-<doc>-<item>
//     (e.g. "my-list-configmap-0-1").
//   - Without a wrapper name: <kind>-<doc>-<item> (e.g. "configmap-0-1").
//
// `doc` is the document ordinal within the file and `item` is the item index within the list, so two
// unnamed items can never collide.
func resourceNameForItem(lines []string, sites []LabelSite, kind, listName string, ordinal, item int) string {
	name := metadataNameAtSite(lines, sites)
	if name != "" && (isHelmTemplatedValue(name) || isValidK8sLabelValue(name)) {
		return name
	}
	loweredKind := strings.ToLower(displayKind(kind))
	if listName != "" {
		return fmt.Sprintf("%s-%s-%d-%d", listName, loweredKind, ordinal, item)
	}
	return fmt.Sprintf("%s-%d-%d", loweredKind, ordinal, item)
}

// metadataNameAtSite returns the value of the `name:` key directly under the metadata block of the
// first site (an item's own metadata), or "" if absent. Quotes are stripped; templated values are
// returned verbatim so the caller can detect them.
func metadataNameAtSite(lines []string, sites []LabelSite) string {
	if len(sites) == 0 {
		return ""
	}
	metaLine := sites[0].MetadataLine
	metaIndent := sites[0].MetadataIndent
	childIndent := metaIndent + 2
	for i := metaLine + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := indentationWidth(line)
		if indent <= metaIndent {
			return ""
		}
		if indent != childIndent {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
			return strings.Trim(value, `"'`)
		}
	}
	return ""
}

// isHelmTemplatedValue reports whether a value contains a Helm/Go template directive (`{{ ... }}`).
// Such a value cannot be statically validated as a Kubernetes label, but is mirrored verbatim into
// yor_name (Strategy 2) and exempted from static label validation before write.
func isHelmTemplatedValue(value string) bool {
	return strings.Contains(value, "{{")
}

// topLevelMetadataName returns the value of the `name:` key directly under the TOP-LEVEL
// `metadata:` block (the site whose Description is "metadata"), or "" if absent. The value's
// surrounding quotes are stripped; templated values are returned verbatim so the caller can detect
// and reject them.
func topLevelMetadataName(lines []string, sites []LabelSite) string {
	metaLine := -1
	metaIndent := 0
	for _, s := range sites {
		if s.Description == "metadata" {
			metaLine = s.MetadataLine
			metaIndent = s.MetadataIndent
			break
		}
	}
	if metaLine < 0 {
		return ""
	}
	childIndent := metaIndent + 2
	for i := metaLine + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := indentationWidth(line)
		if indent <= metaIndent {
			// Left the metadata block without finding name.
			return ""
		}
		if indent != childIndent {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
			return strings.Trim(value, `"'`)
		}
	}
	return ""
}

// topLevelMetadataNameInRange returns the value of the `name:` key directly under the document's
// TOP-LEVEL `metadata:` block (indent 0) within [start,end], or "". It is used to read a List
// wrapper's own metadata.name for the item fallback-naming prefix. The wrapper's metadata is read
// ONLY for the prefix; it is never itself a tag site (D7). Quotes are stripped.
func topLevelMetadataNameInRange(lines []string, start, end int) string {
	metaLine := -1
	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		if indentationWidth(raw) != 0 {
			continue
		}
		if trimmed == "metadata:" {
			metaLine = i
			break
		}
	}
	if metaLine < 0 {
		return ""
	}
	for i := metaLine + 1; i <= end && i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := indentationWidth(line)
		if indent == 0 {
			return ""
		}
		if indent != 2 {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
			return strings.Trim(value, `"'`)
		}
	}
	return ""
}

// findItemLabelSites locates a single Kubernetes List item's OWN metadata.labels target. The item's
// content keys sit at itemIndent; its own `metadata:` key is the first `metadata:` at exactly
// itemIndent within [start,end]. Deeper nested metadata (e.g. a pod-template's metadata inside the
// item) is intentionally NEVER returned, mirroring the top-level-only label scope (D7). The scan is
// purely indentation/text based.
func findItemLabelSites(lines []string, start, end, itemIndent int) []LabelSite {
	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "{{") {
			continue
		}
		// Normalize the first line of the item (the "- key: ..." bullet) to its content indent.
		indent := indentationWidth(raw)
		content := trimmed
		if strings.HasPrefix(trimmed, "- ") {
			content = strings.TrimSpace(trimmed[2:])
			indent += 2
		} else if trimmed == "-" {
			continue
		}
		if indent != itemIndent {
			continue
		}
		if content == "metadata:" || strings.HasPrefix(content, "metadata:") && yamlKeyName(content) == "metadata" {
			labelsLine := findLabelsLineForMetadata(lines, i, end, itemIndent)
			return []LabelSite{{
				MetadataLine:   i,
				MetadataIndent: itemIndent,
				LabelsLine:     labelsLine,
				Description:    "items.metadata",
			}}
		}
	}
	return nil
}

// pathFrame tracks one ancestor key on the indentation stack while scanning a document.
type pathFrame struct {
	indent int
	key    string
}

// findLabelSites scans a single document and returns the single approved metadata.labels target:
// the TOP-LEVEL metadata only (D7 — yor tags are written exclusively to the resource's own
// metadata.labels). Nested metadata (spec.template, spec.jobTemplate, spec.volumeClaimTemplates[*])
// and any `spec.selector` metadata are intentionally NEVER returned, so the pod/job/PVC templates
// and the immutable selector are left untouched. `{{ }}` lines are ignored. The scan is purely
// indentation/text based — it never unmarshals YAML.
func findLabelSites(lines []string, start, end int) []LabelSite {
	var sites []LabelSite
	var stack []pathFrame

	for i := start; i <= end && i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Skip pure template-control lines so they don't corrupt the path stack.
		if strings.HasPrefix(trimmed, "{{") {
			continue
		}

		indent := indentationWidth(raw)

		// A list item ("- ...") keeps its parent's path; treat the content after "- " as a key line
		// at a deeper indent so nested metadata inside list items is still found.
		isListItem := strings.HasPrefix(trimmed, "- ")
		keyPart := trimmed
		keyIndent := indent
		if isListItem {
			keyPart = strings.TrimSpace(trimmed[2:])
			keyIndent = indent + 2
		}

		// Pop frames that are no longer ancestors of this line.
		for len(stack) > 0 && keyIndent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}

		key := yamlKeyName(keyPart)
		if key == "" {
			continue
		}

		if key == "metadata" {
			if isApprovedMetadataPath(stack) {
				labelsLine := findLabelsLineForMetadata(lines, i, end, keyIndent)
				sites = append(sites, LabelSite{
					MetadataLine:   i,
					MetadataIndent: keyIndent,
					LabelsLine:     labelsLine,
					Description:    describePath(stack),
				})
			}
			// Do not push metadata onto the stack as a navigable ancestor for further sites.
			continue
		}

		// Push this key as an ancestor for subsequent deeper lines.
		stack = append(stack, pathFrame{indent: keyIndent, key: key})
	}
	return sites
}

// yamlKeyName extracts the key name from a "key: value" or "key:" line, or "" if the line is not a
// mapping key (e.g. a scalar list item value).
func yamlKeyName(s string) string {
	idx := strings.Index(s, ":")
	if idx <= 0 {
		return ""
	}
	key := strings.TrimSpace(s[:idx])
	// Reject keys containing spaces/template noise.
	if key == "" || strings.ContainsAny(key, " \t{}") {
		return ""
	}
	return key
}

// isApprovedMetadataPath reports whether a `metadata` key sitting under the given ancestor stack is
// an approved label site (D7). Only the TOP-LEVEL metadata (empty ancestor stack) is approved:
// yor tags are written exclusively to the resource's own metadata.labels. Nested metadata
// (spec.template, spec.jobTemplate, spec.volumeClaimTemplates[*]) and anything under a
// `spec.selector` are never tagged.
func isApprovedMetadataPath(stack []pathFrame) bool {
	keys := make([]string, 0, len(stack))
	for _, f := range stack {
		keys = append(keys, f.key)
	}

	// Only the top-level metadata (no ancestor keys) is an approved label site.
	return joinKeys(keys) == ""
}

func describePath(stack []pathFrame) string {
	keys := make([]string, 0, len(stack)+1)
	for _, f := range stack {
		keys = append(keys, f.key)
	}
	if len(keys) == 0 {
		return "metadata"
	}
	return joinKeys(keys) + ".metadata"
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ".")
}

// findLabelsLineForMetadata locates an existing `labels:` key directly under a metadata block whose
// key is at metaIndent, or -1.
func findLabelsLineForMetadata(lines []string, metaLine, end, metaIndent int) int {
	childIndent := metaIndent + 2
	for i := metaLine + 1; i <= end && i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := indentationWidth(line)
		if indent <= metaIndent {
			return -1
		}
		if indent == childIndent && strings.TrimSpace(line) == "labels:" {
			return i
		}
	}
	return -1
}

// extractExistingLabels reads existing label key/value pairs from a labels block (append-only
// baseline + idempotency, D3/D11).
func extractExistingLabels(lines []string, labelsLine, end int) []tags.ITag {
	if labelsLine < 0 {
		return nil
	}
	labelsIndent := indentationWidth(lines[labelsLine])
	var existing []tags.ITag
	for i := labelsLine + 1; i <= end && i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := indentationWidth(line)
		if indent <= labelsIndent {
			break
		}
		kv := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		value = strings.Trim(value, `"'`)
		existing = append(existing, &tags.Tag{Key: key, Value: value})
	}
	return existing
}

// isSkipped reports whether a #yor:skip / #yor:skipall comment (case-insensitive, D12) appears as a
// comment line within the document, at or above its first non-comment content line.
func isSkipped(lines []string, start, end int) bool {
	for i := start; i <= end && i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" || t == "---" {
			continue
		}
		if strings.HasPrefix(t, "#") {
			upper := strings.ToUpper(t)
			if upper == "#YOR:SKIP" || upper == "#YOR:SKIPALL" {
				return true
			}
			continue
		}
		break
	}
	return false
}

// WriteFile implements IParser. It validates via a temp-file round-trip before overwriting user
// files (D13): write to a temp file, re-parse it, and only then write the destination.
func (p *HelmParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	tempFile, err := os.CreateTemp(filepath.Dir(readFilePath), "temp.*.yaml")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	if err = WriteHelmFile(readFilePath, blocks, tempFile.Name()); err != nil {
		return err
	}
	if _, err = p.parseForValidation(tempFile.Name()); err != nil {
		return fmt.Errorf("editing file %v resulted in a malformed template, please open a github issue with the relevant details", readFilePath)
	}
	return WriteHelmFile(readFilePath, blocks, writeFilePath)
}

// parseForValidation re-parses a written temp file without re-applying detection (the temp file is
// outside a real chart layout). It only confirms the document structure is still scannable.
func (p *HelmParser) parseForValidation(filePath string) ([]structure.IBlock, error) {
	// #nosec G304 - file is from user
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSuffix(string(contentBytes), "\n"), "\n")
	parsedBlocks := make([]structure.IBlock, 0)
	docStarts := documentStartLines(lines)
	for ordinal, start := range docStarts {
		end := len(lines) - 1
		if ordinal+1 < len(docStarts) {
			end = docStarts[ordinal+1] - 1
		}
		sites := findLabelSites(lines, start, end)
		if len(sites) == 0 {
			continue
		}
		parsedBlocks = append(parsedBlocks, &HelmBlock{
			Block: structure.Block{
				FilePath:   filePath,
				IsTaggable: true,
				// 1-based block lines (see parseDocument) to align with git blame.
				Lines: structure.Lines{Start: start + 1, End: end + 1},
			},
			DocOrdinal: ordinal,
			Sites:      sites,
		})
	}
	return parsedBlocks, nil
}

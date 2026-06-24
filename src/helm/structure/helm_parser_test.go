package structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
)

const basicChartTemplates = "../../../tests/helm/resources/basic_chart/templates"
const customChartTemplates = "../../../tests/helm/resources/custom_resource_chart/templates"

// traceForOrdinal yields a deterministic yor_trace per document ordinal so golden files are stable.
func traceForOrdinal(ordinal int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", ordinal+1)
}

// traceForItem yields a deterministic yor_trace per (document ordinal, list item index) so golden
// files for Kubernetes List documents are stable and each item gets a distinct trace.
func traceForItem(ordinal, item int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-0000%04d%04d", ordinal+1, item+1)
}

// tagBlocks injects deterministic yor_name and yor_trace into each taggable, non-skipped block,
// mirroring what the runner's tag groups would do. yor_name is taken from the block's real
// GetResourceName() (which the parser derives from the document's metadata.name, with a
// deterministic fallback for templated/absent names) — NOT a hand-rolled lowercased kind — so the
// goldens faithfully represent production behavior. yor_trace stays one value per document ordinal.
func tagBlocks(blocks []structure.IBlock, skipped []string) {
	skippedSet := map[string]bool{}
	for _, s := range skipped {
		skippedSet[s] = true
	}
	for _, b := range blocks {
		hb := b.(*HelmBlock)
		if skippedSet[hb.GetResourceID()] {
			continue
		}
		trace := traceForOrdinal(hb.DocOrdinal)
		if hb.ItemIndex >= 0 {
			// List items share a DocOrdinal but must each receive a distinct, deterministic trace.
			trace = traceForItem(hb.DocOrdinal, hb.ItemIndex)
		}
		hb.AddNewTags([]tags.ITag{
			&tags.Tag{Key: tags.YorNameTagKey, Value: hb.GetResourceName()},
			&tags.Tag{Key: tags.YorTraceTagKey, Value: trace},
		})
	}
}

// runGolden parses inputName, tags blocks, writes, and compares against expectedName.
func runGolden(t *testing.T, dir, inputName, expectedName string) {
	t.Helper()
	p := &HelmParser{}
	p.Init(dir, nil)
	inputPath := filepath.Join(dir, inputName)
	blocks, err := p.ParseFile(inputPath)
	assert.NoError(t, err)

	tagBlocks(blocks, p.GetSkipResourcesByComment())

	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()

	err = p.WriteFile(inputPath, blocks, out.Name())
	assert.NoError(t, err)

	got, err := os.ReadFile(out.Name())
	assert.NoError(t, err)
	want, err := os.ReadFile(filepath.Join(dir, expectedName))
	assert.NoError(t, err)
	assert.Equal(t, string(want), string(got))
}

func TestHelmParser_Detection(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	assert.True(t, p.ValidFile(filepath.Join(basicChartTemplates, "deployment.yaml")), "template with Chart.yaml at root is Helm (D9)")
	assert.False(t, p.ValidFile(filepath.Join(basicChartTemplates, "_helpers.tpl")), "_helpers.tpl excluded")
	assert.False(t, p.ValidFile(filepath.Join(basicChartTemplates, "NOTES.txt")), "NOTES.txt excluded")
	assert.False(t, p.ValidFile("../../../tests/helm/resources/basic_chart/values.yaml"), "values.yaml excluded")
	assert.False(t, p.ValidFile("../../../tests/serverless/resources/no_tags/serverless.yml"), "non-Helm yaml not detected")
}

func TestHelmParser_Detection_VersionTolerant(t *testing.T) {
	p := &HelmParser{}
	p.Init(customChartTemplates, nil)
	// custom_resource_chart uses Chart.yaml apiVersion: v1 — detection must still accept it (D9).
	assert.True(t, p.ValidFile(filepath.Join(customChartTemplates, "rollout.yaml")))
}

func TestHelmParser_StructuralTaggability(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// metadata with no labels -> taggable (D4: create labels).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "deployment.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.True(t, blocks[0].IsBlockTaggable())
	assert.Equal(t, "Deployment", blocks[0].GetResourceType())

	// metadata not locatable (templated) -> skipped, no blocks (D4).
	tBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "templated_meta.yaml"))
	assert.NoError(t, err)
	assert.Len(t, tBlocks, 0)

	// kind: List wrapper whose items each have their OWN nested metadata -> one taggable block per
	// item (the wrapper itself is never tagged). no_metadata.yaml has a single ConfigMap item.
	nBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "no_metadata.yaml"))
	assert.NoError(t, err)
	assert.Len(t, nBlocks, 1)
	assert.True(t, nBlocks[0].IsBlockTaggable())
	assert.Equal(t, "ConfigMap", nBlocks[0].GetResourceType())
	assert.Equal(t, "nested-cm", nBlocks[0].GetResourceName())
}

func TestHelmParser_CustomResourceTaggable(t *testing.T) {
	p := &HelmParser{}
	p.Init(customChartTemplates, nil)
	blocks, err := p.ParseFile(filepath.Join(customChartTemplates, "rollout.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1, "custom resource with locatable metadata is taggable (D4, no allowlist)")
	assert.Equal(t, "Rollout", blocks[0].GetResourceType())
}

func TestHelmParser_YorName_FromMetadataName(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// Literal metadata.name -> yor_name equals it verbatim.
	depBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "deployment.yaml"))
	assert.NoError(t, err)
	assert.Len(t, depBlocks, 1)
	assert.Equal(t, "basic-chart-web", depBlocks[0].GetResourceName(),
		"yor_name is derived from the document's literal top-level metadata.name")
	assert.True(t, IsValidK8sLabel(tags.YorNameTagKey, depBlocks[0].GetResourceName()),
		"a literal metadata.name yor_name must be a valid K8s label value")

	// Multi-doc literal names -> yor_name differs across documents.
	mdBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "multi_doc.yaml"))
	assert.NoError(t, err)
	assert.Len(t, mdBlocks, 2)
	assert.Equal(t, "basic-chart-svc", mdBlocks[0].GetResourceName())
	assert.Equal(t, "basic-chart-cm", mdBlocks[1].GetResourceName())
	assert.NotEqual(t, mdBlocks[0].GetResourceName(), mdBlocks[1].GetResourceName(),
		"yor_name is unique across documents")
}

func TestHelmParser_YorName_TemplatedNamePreserved(t *testing.T) {
	p := &HelmParser{}
	p.Init(customChartTemplates, nil)

	// rollout.yaml has a TEMPLATED metadata.name ({{ .Release.Name }}-rollout). By design
	// (Strategy 2) yor_name mirrors metadata.name VERBATIM, including the Helm template directive,
	// so the rendered yor_name tracks the real resource name. The template expression is preserved
	// exactly (D1) even though, statically, it is not a valid K8s label value.
	blocks, err := p.ParseFile(filepath.Join(customChartTemplates, "rollout.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.Equal(t, "{{ .Release.Name }}-rollout", blocks[0].GetResourceName(),
		"templated metadata.name is mirrored verbatim into yor_name (Strategy 2)")
	assert.Contains(t, blocks[0].GetResourceName(), "{{",
		"the Helm template directive is preserved in yor_name")
}

func TestHelmParser_YorName_AbsentNameFallback(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// cronjob.yaml has a literal name; assert the fallback path indirectly via a name we know is
	// literal stays literal (guards against the templated change leaking into literal handling).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "cronjob.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.Equal(t, "basic-chart-cron", blocks[0].GetResourceName(),
		"a literal metadata.name is still used verbatim and is unaffected by templated-name support")
	assert.True(t, IsValidK8sLabel(tags.YorNameTagKey, blocks[0].GetResourceName()))
}

func TestHelmParser_Idempotency_TemplatedYorName_StandardKind(t *testing.T) {
	// Re-tagging an already-tagged STANDARD-kind file whose yor_name is a Helm template containing
	// embedded quotes ({{ include "..." . }}) must be a no-op (D11). Append-only is keyed on the
	// label key, so the existing yor_name is preserved verbatim and the file is unchanged.
	dir := basicChartTemplates
	p := &HelmParser{}
	p.Init(dir, nil)

	tagged := filepath.Join(dir, "templated_name_tagged.yaml")
	blocks, err := p.ParseFile(tagged)
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)

	tagBlocks(blocks, nil)
	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(tagged, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	want, _ := os.ReadFile(tagged)
	assert.Equal(t, string(want), string(got),
		"re-run on a tagged standard-kind file with a templated yor_name yields no diff (D11)")
}

func TestHelmParser_YorName_TemplatedName_StandardKind(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A STANDARD Kubernetes kind (Deployment, not a CRD) with a templated metadata.name must also
	// have its templated name mirrored verbatim into yor_name (Strategy 2 is kind-agnostic).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "templated_name.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.Equal(t, "Deployment", blocks[0].GetResourceType())
	assert.Equal(t, `{{ include "basic-chart.fullname" . }}`, blocks[0].GetResourceName(),
		"a templated metadata.name on a standard kind is mirrored verbatim into yor_name")
	assert.Contains(t, blocks[0].GetResourceName(), "{{",
		"the Helm template directive is preserved in yor_name")
}

func TestHelmParser_Golden_TemplatedName_StandardKind(t *testing.T) {
	runGolden(t, basicChartTemplates, "templated_name.yaml", "templated_name_tagged.yaml")
}

func TestHelmParser_Golden_CreateLabels(t *testing.T) {
	runGolden(t, basicChartTemplates, "deployment.yaml", "deployment_tagged.yaml")
}

func TestHelmParser_Golden_AppendLabels(t *testing.T) {
	runGolden(t, basicChartTemplates, "existing_labels.yaml", "existing_labels_tagged.yaml")
}

func TestHelmParser_Golden_CronJob(t *testing.T) {
	runGolden(t, basicChartTemplates, "cronjob.yaml", "cronjob_tagged.yaml")
}

func TestHelmParser_Idempotency_TemplatedYorName(t *testing.T) {
	// Re-tagging an already-tagged file whose yor_name is a Helm template ({{ }}) must be a no-op
	// (D11). This guards Strategy 2: the BARE templated value is read back, recognized as already
	// present (append-only is keyed on the label key), and the template directive is preserved
	// verbatim (D1).
	dir := customChartTemplates
	p := &HelmParser{}
	p.Init(dir, nil)

	tagged := filepath.Join(dir, "rollout_tagged.yaml")
	blocks, err := p.ParseFile(tagged)
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.Equal(t, "{{ .Release.Name }}-rollout", blocks[0].GetResourceName(),
		"templated yor_name is read back verbatim from the tagged file")

	tagBlocks(blocks, nil)
	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(tagged, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	want, _ := os.ReadFile(tagged)
	assert.Equal(t, string(want), string(got),
		"re-run on a tagged file with a templated yor_name yields no diff (D11)")
}

func TestHelmParser_TopLevelOnlyLabelSite(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// Deployment: ONLY the top-level metadata is a label site (D7). The pod-template metadata and
	// the immutable selector are never tagged.
	depBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "deployment.yaml"))
	assert.NoError(t, err)
	assert.Len(t, depBlocks, 1)
	depSites := siteDescriptions(depBlocks[0].(*HelmBlock))
	assert.Equal(t, []string{"metadata"}, depSites, "only the top-level metadata is tagged")

	// CronJob: ONLY top-level metadata — jobTemplate / pod-template metadata are not sites.
	cronBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "cronjob.yaml"))
	assert.NoError(t, err)
	assert.Len(t, cronBlocks, 1)
	cronSites := siteDescriptions(cronBlocks[0].(*HelmBlock))
	assert.Equal(t, []string{"metadata"}, cronSites)

	// StatefulSet (existing_labels): ONLY top-level metadata — pod-template and
	// volumeClaimTemplates metadata are not sites, and selector is never a site.
	stsBlocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "existing_labels.yaml"))
	assert.NoError(t, err)
	assert.Len(t, stsBlocks, 1)
	stsSites := siteDescriptions(stsBlocks[0].(*HelmBlock))
	assert.Equal(t, []string{"metadata"}, stsSites)
}

func siteDescriptions(b *HelmBlock) []string {
	out := make([]string, 0, len(b.Sites))
	for _, s := range b.Sites {
		out = append(out, s.Description)
	}
	return out
}

func TestHelmParser_Golden_MultiDoc(t *testing.T) {
	runGolden(t, basicChartTemplates, "multi_doc.yaml", "multi_doc_tagged.yaml")
}

func TestHelmParser_Golden_CustomResource(t *testing.T) {
	runGolden(t, customChartTemplates, "rollout.yaml", "rollout_tagged.yaml")
}

func TestHelmParser_Golden_SkipUnchanged(t *testing.T) {
	// skip.yaml carries #yor:skip; the block is collected as skipped and receives no tags.
	runGolden(t, basicChartTemplates, "skip.yaml", "skip_tagged.yaml")
}

func TestHelmParser_TemplatePreservation(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)
	inputPath := filepath.Join(basicChartTemplates, "deployment.yaml")
	blocks, err := p.ParseFile(inputPath)
	assert.NoError(t, err)
	tagBlocks(blocks, nil)

	out, err := os.CreateTemp(basicChartTemplates, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(inputPath, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	assert.Contains(t, string(got), `{{ .Values.replicaCount }}`, "template directives preserved verbatim (D1)")
}

func TestHelmParser_SkipComment(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)
	_, err := p.ParseFile(filepath.Join(basicChartTemplates, "skip.yaml"))
	assert.NoError(t, err)
	assert.Len(t, p.GetSkipResourcesByComment(), 1, "#yor:skip recorded (D12)")
}

func TestHelmParser_Idempotency(t *testing.T) {
	// Tagging an already-tagged file produces no further changes (D11).
	dir := basicChartTemplates
	p := &HelmParser{}
	p.Init(dir, nil)

	first := filepath.Join(dir, "deployment_tagged.yaml")
	blocks, err := p.ParseFile(first)
	assert.NoError(t, err)
	// Existing yor_trace should be detected from the tagged file.
	assert.Len(t, blocks, 1)
	existing := blocks[0].GetExistingTags()
	var hasTrace bool
	for _, tg := range existing {
		if tg.GetKey() == tags.YorTraceTagKey {
			hasTrace = true
		}
	}
	assert.True(t, hasTrace, "existing yor_trace read back from tagged file")

	// Re-tag with a different trace; output must keep the original trace value.
	tagBlocks(blocks, nil)
	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(first, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	want, _ := os.ReadFile(first)
	assert.Equal(t, string(want), string(got), "re-run on tagged file yields no diff (D11)")
}

func TestHelmParser_SafetyGuard_DoesNotCorruptOriginal(t *testing.T) {
	// Write to a temp dir copy and ensure the original input file is never mutated by WriteFile.
	dir := basicChartTemplates
	p := &HelmParser{}
	p.Init(dir, nil)
	inputPath := filepath.Join(dir, "deployment.yaml")
	original, _ := os.ReadFile(inputPath)

	blocks, err := p.ParseFile(inputPath)
	assert.NoError(t, err)
	tagBlocks(blocks, nil)

	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(inputPath, blocks, out.Name()))

	after, _ := os.ReadFile(inputPath)
	assert.Equal(t, string(original), string(after), "original input file unchanged when writing elsewhere (D13)")
}

// --- Kubernetes List kind support ---

func TestHelmParser_List_PerItemBlocks(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A generic kind: List with two named items yields ONE block per item, each derived from the
	// item's OWN metadata.name. The List wrapper itself is never a block.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2, "one block per list item, never the wrapper")

	b0 := blocks[0].(*HelmBlock)
	b1 := blocks[1].(*HelmBlock)
	assert.Equal(t, "app-config", b0.GetResourceName())
	assert.Equal(t, "feature-flags", b1.GetResourceName())
	assert.Equal(t, "ConfigMap", b0.GetResourceType())
	assert.Equal(t, 0, b0.ItemIndex, "first item index")
	assert.Equal(t, 1, b1.ItemIndex, "second item index")

	// Each item's label site is its OWN nested metadata, not top-level.
	assert.Len(t, b0.Sites, 1)
	assert.True(t, b0.Sites[0].MetadataIndent > 0, "item metadata is nested deeper than top level")
}

func TestHelmParser_List_TypedListVariant(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A typed *List (ConfigMapList) is handled exactly like the generic List (case-insensitive
	// suffix match on "List").
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "typed_list.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, "typed-one", blocks[0].GetResourceName())
	assert.Equal(t, "typed-two", blocks[1].GetResourceName())
}

func TestHelmParser_List_FallbackNaming_WithPrefix(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// Unnamed items + a List wrapper WITH a top-level metadata.name -> <list-name>-<kind>-<doc>-<item>.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_unnamed_prefixed.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, "my-list-configmap-0-0", blocks[0].GetResourceName())
	assert.Equal(t, "my-list-configmap-0-1", blocks[1].GetResourceName())
	assert.NotEqual(t, blocks[0].GetResourceName(), blocks[1].GetResourceName(),
		"unnamed items must get collision-safe distinct names")
}

func TestHelmParser_List_FallbackNaming_NoPrefix(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// Unnamed items + a List wrapper WITHOUT a top-level metadata.name -> <kind>-<doc>-<item>.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_unnamed.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, "configmap-0-0", blocks[0].GetResourceName())
	assert.Equal(t, "configmap-0-1", blocks[1].GetResourceName())
}

func TestHelmParser_Golden_List(t *testing.T) {
	runGolden(t, basicChartTemplates, "list.yaml", "list_tagged.yaml")
}

func TestHelmParser_Golden_TypedList(t *testing.T) {
	runGolden(t, basicChartTemplates, "typed_list.yaml", "typed_list_tagged.yaml")
}

func TestHelmParser_Golden_List_UnnamedPrefixed(t *testing.T) {
	runGolden(t, basicChartTemplates, "list_unnamed_prefixed.yaml", "list_unnamed_prefixed_tagged.yaml")
}

func TestHelmParser_Golden_List_Unnamed(t *testing.T) {
	runGolden(t, basicChartTemplates, "list_unnamed.yaml", "list_unnamed_tagged.yaml")
}

func TestHelmParser_List_FlowBulletsAtItemsIndent(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// Regression: YAML allows block-sequence bullets at the SAME indentation as the parent
	// `items:` key (the canonical `kubectl ... -o yaml` List style). Such items must still be
	// enumerated and tagged per item — previously the whole List was dropped as "no locatable
	// items", leaving every item untagged.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_flow.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2, "items whose bullets sit at the items: indent are still enumerated")
	assert.Equal(t, "flow-one", blocks[0].GetResourceName())
	assert.Equal(t, "flow-two", blocks[1].GetResourceName())
	assert.Equal(t, 0, blocks[0].(*HelmBlock).ItemIndex)
	assert.Equal(t, 1, blocks[1].(*HelmBlock).ItemIndex)
	assert.Equal(t, []string{"items.metadata"}, siteDescriptions(blocks[0].(*HelmBlock)),
		"each item is tagged on its own nested metadata, not the wrapper")
}

func TestHelmParser_Golden_List_FlowBullets(t *testing.T) {
	// Golden round-trip for the zero-indent (items:-level) bullet style: each item tagged on its
	// own metadata.labels; wrapper untouched.
	runGolden(t, basicChartTemplates, "list_flow.yaml", "list_flow_tagged.yaml")
}

func TestHelmParser_List_Idempotency(t *testing.T) {
	// Re-tagging an already-tagged List file is a no-op and preserves each item's existing
	// yor_trace (D11).
	dir := basicChartTemplates
	p := &HelmParser{}
	p.Init(dir, nil)

	tagged := filepath.Join(dir, "list_tagged.yaml")
	blocks, err := p.ParseFile(tagged)
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)

	tagBlocks(blocks, nil)
	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(tagged, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	want, _ := os.ReadFile(tagged)
	assert.Equal(t, string(want), string(got),
		"re-run on a tagged List file yields no diff and preserves each item's yor_trace (D11)")
}

func TestHelmParser_IsListDocument_RequiresBothSignals(t *testing.T) {
	// Detection requires BOTH a `*List` kind AND a top-level items: sequence.
	listDoc := strings.Split("apiVersion: v1\nkind: ConfigMapList\nitems:\n  - kind: ConfigMap\n    metadata:\n      name: a", "\n")
	assert.True(t, isListDocument(listDoc, 0, len(listDoc)-1, "ConfigMapList"),
		"a *List kind WITH top-level items: is a List document")

	// Case-insensitive on the suffix.
	genericLower := strings.Split("kind: list\nitems:\n  - kind: ConfigMap", "\n")
	assert.True(t, isListDocument(genericLower, 0, len(genericLower)-1, "list"),
		"the List suffix is matched case-insensitively")

	// A *List-named SINGLE resource with NO top-level items: is NOT a List (handled normally).
	singleton := strings.Split("kind: Blocklist\nmetadata:\n  name: deny-all\nspec:\n  rules:\n    - deny", "\n")
	assert.False(t, isListDocument(singleton, 0, len(singleton)-1, "Blocklist"),
		"a *List-named single resource without top-level items: is NOT a List document")

	// A non-*List document that happens to have a top-level items: is NOT a List.
	itemsButNotList := strings.Split("kind: ConfigMap\nmetadata:\n  name: cm\nitems:\n  - a\n  - b", "\n")
	assert.False(t, isListDocument(itemsButNotList, 0, len(itemsButNotList)-1, "ConfigMap"),
		"a non-*List kind with a top-level items: field is NOT a List document")

	// An items: nested under spec (not top-level) does not satisfy the items: signal.
	nestedItems := strings.Split("kind: WidgetList\nspec:\n  items:\n    - a", "\n")
	assert.False(t, isListDocument(nestedItems, 0, len(nestedItems)-1, "WidgetList"),
		"a non-top-level items: does not make a document a List")
}

func TestHelmParser_List_ListLikeSingletonTaggedNormally(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// kind: Blocklist ends in 'list' but is a single resource (no top-level items:). It must be
	// tagged via the NORMAL single-resource path (one block, top-level metadata), not dropped.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "listlike_singleton.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1, "a *List-named single resource is one normal block, not List items")
	b := blocks[0].(*HelmBlock)
	assert.Equal(t, "Blocklist", b.GetResourceType())
	assert.Equal(t, "deny-all", b.GetResourceName())
	assert.Equal(t, -1, b.ItemIndex, "single resource has no item index")
	assert.Equal(t, []string{"metadata"}, siteDescriptions(b), "tagged on its own top-level metadata")
}

func TestHelmParser_Golden_ListLikeSingleton(t *testing.T) {
	runGolden(t, basicChartTemplates, "listlike_singleton.yaml", "listlike_singleton_tagged.yaml")
}

func TestHelmParser_List_TemplatedItemName(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A `{{ }}` directive in an ITEM's own metadata.name is mirrored verbatim into yor_name
	// (Strategy 2), exactly like the single-resource templated-name path.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_templated_item.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	name := blocks[0].GetResourceName()
	assert.Equal(t, "{{ .Release.Name }}-cm", name)
	assert.Contains(t, name, "{{", "templated item name preserves the directive verbatim")
}

func TestHelmParser_List_MultiDocOrdinal(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A List that follows a normal document (`---`) sits at document ordinal 1, so its unnamed
	// items' fallback names embed <doc>=1: configmap-1-0 / configmap-1-1. This exercises the
	// non-zero document-ordinal component of the fallback-naming format.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_multidoc.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 3, "leading single doc + two list items")

	// Block 0 is the leading single ConfigMap (doc 0).
	assert.Equal(t, "leading-cm", blocks[0].GetResourceName())
	assert.Equal(t, -1, blocks[0].(*HelmBlock).ItemIndex)
	assert.Equal(t, 0, blocks[0].(*HelmBlock).DocOrdinal)

	// Blocks 1 & 2 are the List items at doc ordinal 1 with non-zero <doc> in their fallback names.
	assert.Equal(t, "configmap-1-0", blocks[1].GetResourceName())
	assert.Equal(t, "configmap-1-1", blocks[2].GetResourceName())
	assert.Equal(t, 1, blocks[1].(*HelmBlock).DocOrdinal)
	assert.Equal(t, 1, blocks[2].(*HelmBlock).DocOrdinal)
}

func TestHelmParser_Golden_List_MultiDocOrdinal(t *testing.T) {
	runGolden(t, basicChartTemplates, "list_multidoc.yaml", "list_multidoc_tagged.yaml")
}

func TestHelmParser_List_MixedNaming(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// One named item and one unnamed item under a named wrapper: the named item keeps its literal
	// name; the unnamed item gets the with-prefix fallback at its own item index.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_mixed_naming.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, "named-cm", blocks[0].GetResourceName(), "named item keeps its literal name")
	assert.Equal(t, "my-list-configmap-0-1", blocks[1].GetResourceName(),
		"unnamed item falls back at its own item index")
}

func TestHelmParser_List_HeterogeneousKinds(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A List with differing item kinds derives each item's fallback name from its OWN kind.
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_heterogeneous.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, "ConfigMap", blocks[0].GetResourceType())
	assert.Equal(t, "Service", blocks[1].GetResourceType())
	assert.Equal(t, "configmap-0-0", blocks[0].GetResourceName())
	assert.Equal(t, "service-0-1", blocks[1].GetResourceName(),
		"the second item's fallback uses its OWN kind (service), not the first item's")
}

func TestHelmParser_Golden_List_AppendItemLabels(t *testing.T) {
	// An item that already has metadata.labels keeps them and yor appends yor_name/yor_trace under
	// the SAME labels block (append-only baseline, D3).
	runGolden(t, basicChartTemplates, "list_append_item.yaml", "list_append_item_tagged.yaml")
}

func TestHelmParser_Golden_List_SkipUnchanged(t *testing.T) {
	// A #yor:skip comment on a List document skips every item, leaving the file byte-for-byte
	// identical.
	runGolden(t, basicChartTemplates, "list_skip.yaml", "list_skip_tagged.yaml")
}

func TestHelmParser_List_EmptyItemsNoBlocks(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A List with an empty items: sequence has no taggable entries → zero blocks (degenerate case).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_empty.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 0, "an empty List yields no taggable blocks")
}

func TestHelmParser_List_TemplatedItemMetadataSkipped(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A List item whose metadata is supplied via a template include has no locatable literal
	// metadata block → that item is skipped (mirroring the top-level templated_meta.yaml negative).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_templated_item_metadata.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 0, "an item with only templated metadata is not statically taggable")
}

// reRunIsNoOp parses an already-tagged file, re-tags, writes to a temp path, and asserts the output
// is byte-for-byte identical to the input (idempotency, D11).
func reRunIsNoOp(t *testing.T, dir, taggedName string) {
	t.Helper()
	p := &HelmParser{}
	p.Init(dir, nil)
	tagged := filepath.Join(dir, taggedName)
	blocks, err := p.ParseFile(tagged)
	assert.NoError(t, err)

	tagBlocks(blocks, p.GetSkipResourcesByComment())
	out, err := os.CreateTemp(dir, "out.*.yaml")
	assert.NoError(t, err)
	defer func() { _ = os.Remove(out.Name()) }()
	assert.NoError(t, p.WriteFile(tagged, blocks, out.Name()))

	got, _ := os.ReadFile(out.Name())
	want, _ := os.ReadFile(tagged)
	assert.Equal(t, string(want), string(got),
		"re-run on %s must be a no-op (idempotent), preserving each item's tags", taggedName)
}

func TestHelmParser_List_ItemNestedMetadataNotTagged(t *testing.T) {
	p := &HelmParser{}
	p.Init(basicChartTemplates, nil)

	// A Deployment item inside a List carries nested spec.template.metadata and spec.selector. Only
	// the item's OWN top-level metadata is a label site; the nested pod-template metadata and the
	// selector must NEVER be tagged (per-item analogue of TestHelmParser_TopLevelOnlyLabelSite).
	blocks, err := p.ParseFile(filepath.Join(basicChartTemplates, "list_nested_item.yaml"))
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)
	b := blocks[0].(*HelmBlock)
	assert.Equal(t, "Deployment", b.GetResourceType())
	assert.Equal(t, "web", b.GetResourceName())
	assert.Len(t, b.Sites, 1, "exactly one label site: the item's own metadata")
	// The single site must be the item's own metadata, NOT the deeper pod-template metadata.
	assert.Equal(t, []string{"items.metadata"}, siteDescriptions(b),
		"the only site is the item's own metadata, never nested spec.template.metadata")
}

func TestHelmParser_Golden_List_ItemNestedMetadataNotTagged(t *testing.T) {
	// Golden confirms tags land ONLY on the item's own metadata; nested spec.template.metadata and
	// the selector are byte-for-byte preserved.
	runGolden(t, basicChartTemplates, "list_nested_item.yaml", "list_nested_item_tagged.yaml")
}

func TestHelmParser_List_Idempotency_FallbackNames(t *testing.T) {
	// Re-running on List files whose items carry deterministic FALLBACK names must be a no-op: the
	// <kind>-<doc>-<item> (and <list-name>-<kind>-<doc>-<item>) names must reproduce identically on
	// re-parse, with each item's yor_trace preserved.
	reRunIsNoOp(t, basicChartTemplates, "list_unnamed_tagged.yaml")
	reRunIsNoOp(t, basicChartTemplates, "list_unnamed_prefixed_tagged.yaml")
}

func TestHelmParser_List_Idempotency_TemplatedItemName(t *testing.T) {
	// Re-running on a List item whose yor_name is a BARE templated directive is a no-op; the
	// `{{ }}` directive round-trips verbatim and the item's yor_trace is preserved.
	reRunIsNoOp(t, basicChartTemplates, "list_templated_item_tagged.yaml")
}

package structure

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

// WriteHelmFile is the dedicated Helm writer (D2). It is intentionally NOT the shared
// WriteYAMLFile: it edits the file as text so that Go template directives (`{{ ... }}`) are
// preserved verbatim (D1), it supports multiple `---`-separated documents per file (D2), and it
// writes into the single located metadata.labels site of each document — the TOP-LEVEL metadata
// only (D7) — append-only. Nested pod-template / jobTemplate / volumeClaimTemplate metadata and the
// immutable `spec.selector` are never written.
func WriteHelmFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	// #nosec G304 - file is from user
	originalBytes, err := os.ReadFile(readFilePath)
	if err != nil {
		return err
	}
	hadTrailingNewline := strings.HasSuffix(string(originalBytes), "\n")
	content := strings.TrimSuffix(string(originalBytes), "\n")
	lines := strings.Split(content, "\n")

	type insertion struct {
		// atLine is the 0-based index AFTER which the new lines are inserted.
		atLine int
		lines  []string
	}
	insertions := make([]insertion, 0)

	for _, b := range blocks {
		hb, ok := b.(*HelmBlock)
		if !ok || !hb.IsBlockTaggable() {
			continue
		}
		hb.UpdateTags()
		validLabels, dropped := hb.GetValidLabels()
		for _, d := range dropped {
			logger.Warning(fmt.Sprintf(
				"skipping invalid Kubernetes label for %s (kind=%s, doc=%d): %s=%q",
				readFilePath, hb.GetResourceType(), hb.DocOrdinal, d.GetKey(), d.GetValue()))
		}
		if len(validLabels) == 0 {
			continue
		}

		// Apply to the located label site(s) (D7 — top-level metadata only). Append-only is
		// enforced PER SITE: a key already present in that site's own labels block is never
		// overwritten, so re-runs are no-ops (D11).
		for _, site := range hb.Sites {
			if site.LabelsLine >= 0 {
				existingByKey := existingLabelKeysAtSite(lines, site.LabelsLine)
				labelIndent := indentationWidth(lines[site.LabelsLine])
				childIndent := strings.Repeat(" ", labelIndent+2)
				newLines := make([]string, 0)
				for _, t := range sortTags(validLabels) {
					if _, exists := existingByKey[t.GetKey()]; exists {
						continue
					}
					newLines = append(newLines, fmt.Sprintf("%s%s: %s", childIndent, t.GetKey(), quoteLabelValue(t.GetValue())))
				}
				if len(newLines) > 0 {
					insertions = append(insertions, insertion{atLine: lastLabelChildLine(lines, site.LabelsLine, labelIndent), lines: newLines})
				}
			} else {
				// Create a labels: sub-block under this metadata: (D4 — missing labels is not a skip).
				metaIndent := site.MetadataIndent
				labelsIndent := strings.Repeat(" ", metaIndent+2)
				childIndent := strings.Repeat(" ", metaIndent+4)
				newLines := []string{fmt.Sprintf("%slabels:", labelsIndent)}
				for _, t := range sortTags(validLabels) {
					newLines = append(newLines, fmt.Sprintf("%s%s: %s", childIndent, t.GetKey(), quoteLabelValue(t.GetValue())))
				}
				insertions = append(insertions, insertion{atLine: site.MetadataLine, lines: newLines})
			}
		}
	}

	// Apply insertions bottom-up so earlier indices remain valid.
	sort.Slice(insertions, func(i, j int) bool {
		return insertions[i].atLine > insertions[j].atLine
	})
	for _, ins := range insertions {
		at := ins.atLine + 1
		if at > len(lines) {
			at = len(lines)
		}
		updated := make([]string, 0, len(lines)+len(ins.lines))
		updated = append(updated, lines[:at]...)
		updated = append(updated, ins.lines...)
		updated = append(updated, lines[at:]...)
		lines = updated
	}

	out := strings.Join(lines, "\n")
	if hadTrailingNewline {
		out += "\n"
	}
	// #nosec G306 - match existing writers' permissions
	return os.WriteFile(writeFilePath, []byte(out), 0644)
}

// existingLabelKeysAtSite returns the set of label keys already present directly under the labels:
// block at labelsLine, so writing is append-only per site.
func existingLabelKeysAtSite(lines []string, labelsLine int) map[string]struct{} {
	keys := map[string]struct{}{}
	if labelsLine < 0 {
		return keys
	}
	labelsIndent := indentationWidth(lines[labelsLine])
	for i := labelsLine + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		if indentationWidth(line) <= labelsIndent {
			break
		}
		kv := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(kv) == 0 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		if key != "" {
			keys[key] = struct{}{}
		}
	}
	return keys
}

// lastLabelChildLine returns the 0-based index of the last line belonging to the labels: block,
// i.e. the line after which new label children should be inserted.
func lastLabelChildLine(lines []string, labelsLine int, labelsIndent int) int {
	last := labelsLine
	for i := labelsLine + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		if indentationWidth(line) <= labelsIndent {
			break
		}
		last = i
	}
	return last
}

func indentationWidth(line string) int {
	width := 0
	for _, r := range line {
		if r == ' ' {
			width++
		} else {
			break
		}
	}
	return width
}

// quoteLabelValue renders a label value for writing.
//
// A Helm-templated value (`{{ ... }}`) is emitted BARE — byte-for-byte identical to the source
// (e.g. metadata.name) — so the written yor_name matches the resource's name exactly, with no
// surrounding quotes or escaping (Strategy 2). The temp-file validation round-trip in WriteFile
// guards against the rare bare template that would not be well-formed YAML (e.g. one containing a
// top-level ": ").
//
// Every other value is wrapped in double quotes to keep YAML well-formed for values that could
// otherwise be misread (e.g. starting with a digit). yor label values are already validated.
func quoteLabelValue(value string) string {
	if isHelmTemplatedValue(value) {
		return value
	}
	return fmt.Sprintf("%q", value)
}

func sortTags(in []tags.ITag) []tags.ITag {
	out := make([]tags.ITag, len(in))
	copy(out, in)
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetKey() < out[j].GetKey()
	})
	return out
}

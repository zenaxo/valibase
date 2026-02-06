package gen

import (
	"fmt"
	"strings"

	"github.com/zenaxo/valibase/internal/utils"
)

func writeCollectionsSegment(w *tsw, colls []string) {
	w.WL("")
	w.WL("// All available PocketBase collections as a const map")
	w.WL("export const Collections = {")
	w.Indent()
	for _, c := range colls {
		w.WL(fmt.Sprintf("%s: '%s',", utils.ToPascalCase(c), c))
	}
	w.Dedent()
	w.WL("} as const;")
	w.WL("export type CollectionKey = keyof typeof Collections;")
	w.WL("export type CollectionName = (typeof Collections)[CollectionKey];")
}

func writeCollectionSection(w *tsw, c collectionRecord, byID map[string]collectionRecord) {
	n := nameParts(c.Name)
	sectionComment(w, n.collectionName)

	// Gather field snippets
	viewFields := make([]string, 0, len(c.Fields))
	inputFields := make([]string, 0, len(c.Fields))
	expandFields := make([]string, 0, len(c.Fields))

	for _, f := range c.Fields {
		if shouldSkipField(f) {
			continue
		}

		view, input := emitField(f)
		viewFields = append(viewFields, view)
		inputFields = append(inputFields, input)

		if expandField, ok := expandFieldSnippet(f, byID); ok {
			expandFields = append(expandFields, expandField)
		}
	}

	w.W(collectionFieldsSchema(n.collectionName, strings.Join(viewFields, ",\n\t")))
	w.W(collectionInputSchema(n.collectionName, strings.Join(inputFields, ",\n\t")))

	w.W(fmt.Sprintf(`
export type %sFields = v.InferOutput<typeof %sResponse>;
`, n.pascalSingular, n.lowerCamelSingular))

	writeExportType(w, n.collectionName, n.pascalSingular, expandFields)

	w.W(createUpdateExports(n.collectionName, c.Type))
}

func collectionFieldsSchema(collectionName, content string) string {
	n := nameParts(collectionName)

	var b strings.Builder
	b.WriteString("\n/**\n")
	fmt.Fprintf(&b, "* Raw field schema for %q\n", n.collectionName)
	b.WriteString("*/\n")
	fmt.Fprintf(&b, "export const %sResponse = %s;\n",
		n.lowerCamelSingular,
		v.Object(fmt.Sprintf(`{
	...systemFieldsSchema('%s').entries,
	%s
}`, n.collectionName, content)),
	)

	return b.String()
}

func collectionInputSchema(collectionName, content string) string {
	n := nameParts(collectionName)

	var b strings.Builder
	b.WriteString("\n/*\n")
	fmt.Fprintf(&b, "* Input schema for creating/updating %q\n", n.collectionName)
	b.WriteString("*/\n")
	fmt.Fprintf(&b, "export const %sInput = %s;\n",
		n.lowerCamelSingular,
		v.Object(fmt.Sprintf(`{
	%s
}`, content)),
	)

	return b.String()
}

func collectionExpandType(collectionName, content string) string {
	n := nameParts(collectionName)

	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return fmt.Sprintf("\n// No expand relations defined for %s\nexport type %sExpand = {};\n",
			n.pascalSingular,
			n.pascalSingular,
		)
	}

	var b strings.Builder
	b.WriteString("\n/**\n")
	fmt.Fprintf(&b, "* Relations that can be expanded when loading %q\n", n.collectionName)
	b.WriteString("*/\n")
	fmt.Fprintf(&b, "export type %sExpand = {\n\t%s\n}\n",
		n.pascalSingular,
		content,
	)
	return b.String()
}

func relationExpandTSType(f fieldSchema, byID map[string]collectionRecord) (string, bool) {
	if f.Type != FieldRelation || f.RelationCollectionID == nil {
		return "", false
	}

	targetColl, ok := byID[*f.RelationCollectionID]
	if !ok {
		return "", false
	}

	target := utils.ToPascalCase(utils.ToSingular(targetColl.Name))
	if f.MaxSelect != nil && *f.MaxSelect == 1 {
		return target, true
	}
	return target + "[]", true
}

func expandFieldSnippet(f fieldSchema, byID map[string]collectionRecord) (string, bool) {
	tsType, ok := relationExpandTSType(f, byID)
	if !ok {
		return "", false
	}
	return sanitizeFieldName(f.Name) + "?: " + tsType, true
}

func writeExportType(w *tsw, collectionName, pascalSingular string, expandFields []string) {
	parts := []string{pascalSingular + "Fields"}

	if len(expandFields) > 0 {
		expandBody := strings.Join(expandFields, ";\n\t")
		w.W(collectionExpandType(collectionName, expandBody))
		parts = append(parts, `Expand<Partial<`+pascalSingular+`Expand>>`)
	}

	w.W(fmt.Sprintf(`
export type %s = %s;
`, pascalSingular, strings.Join(parts, " & ")))
}

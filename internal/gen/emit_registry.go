package gen

import (
	"fmt"

	"github.com/zenaxo/valibase/internal/utils"
)

func createUpdateExports(collectionName string, cType collectionType) string {
	pascalSingular := utils.ToPascalCase(utils.ToSingular(collectionName))
	lowerCamelSingular := utils.ToLowerCamelCase(utils.ToSingular(collectionName))

	createFn := "createBaseSchema"
	updateFn := "updateBaseSchema"
	if cType == CollectionAuth {
		createFn = "createAuthSchema"
		updateFn = "updateAuthSchema"
	}

	return fmt.Sprintf(`
/**
 * Create/Update schemas and their inferred input types for "%[1]s" records.
 */
export const create%[1]sSchema = %[3]s(%[2]sInput);
export const update%[1]sSchema = %[4]s(%[2]sInput);

// Inferred input types from the above schemas
export type Create%[1]sInput = v.InferOutput<typeof create%[1]sSchema>;
export type Update%[1]sInput = v.InferOutput<typeof update%[1]sSchema>;
`, pascalSingular, lowerCamelSingular, createFn, updateFn)
}

func writeRegistry(w *tsw, collectionNames []string) {
	w.WL("")
	w.WL("// Central registry of all generated collection schemas")
	w.WL("export const registry = {")
	w.Indent()

	for _, coll := range collectionNames {
		lowerCamelSingular := utils.ToLowerCamelCase(utils.ToSingular(coll))
		pascalSingular := utils.ToPascalCase(utils.ToSingular(coll))

		w.WL(fmt.Sprintf("// Schemas for the %q collection", coll))
		w.WL(fmt.Sprintf("%s: {", coll))
		w.Indent()
		w.WL(fmt.Sprintf("response: %sResponse,", lowerCamelSingular))
		w.WL(fmt.Sprintf("create: create%sSchema,", pascalSingular))
		w.WL(fmt.Sprintf("update: update%sSchema", pascalSingular))
		w.Dedent()
		w.WL("},")
		w.WL("")
	}

	w.Dedent()
	w.WL("} as const;")
	w.WL("")
	w.WL("export type CollectionsMap = typeof registry;")
	w.WL("export type CollectionNameKey = keyof CollectionsMap;")
	w.WL("")
	w.WL("// Helper type map: collection name -> strongly typed record")
	w.WL("export type ResponseTypes = {")
	w.Indent()

	for _, coll := range collectionNames {
		pascalSingular := utils.ToPascalCase(utils.ToSingular(coll))
		w.WL(fmt.Sprintf("%s: %s;", coll, pascalSingular))
	}

	w.Dedent()
	w.WL("};")
}

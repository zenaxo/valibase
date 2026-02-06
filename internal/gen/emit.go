package gen

import (
	"fmt"
	"strings"
)

var fieldsToIgnore = map[string]struct{}{
	"created":      {},
	"updated":      {},
	"collectionId": {},
	"id":           {},
}

func shouldSkipField(f fieldSchema) bool {
	if _, ok := fieldsToIgnore[f.Name]; ok {
		return true
	}
	return f.Hidden
}

// GenerateTS generates the full TypeScript output for the provided collections.
func GenerateTS(collections []collectionRecord) string {
	w := &tsw{}

	collectionNames := make([]string, 0, len(collections))
	byID := make(map[string]collectionRecord, len(collections))

	for _, c := range collections {
		collectionNames = append(collectionNames, c.Name)
		byID[c.ID] = c
	}

	w.W(imports())
	writeCollectionsSegment(w, collectionNames)
	w.W(typeHelpers())

	for _, c := range collections {
		writeCollectionSection(w, c, byID)
	}

	writeRegistry(w, collectionNames)
	w.W(tail())

	return w.String()
}

func sectionComment(w *tsw, name string) {
	sep := strings.Repeat("=", 90)
	w.WL("")
	w.WL(fmt.Sprintf("/*%s", sep))
	w.WL(fmt.Sprintf("%s COLLECTION", strings.ToUpper(name)))
	w.WL(fmt.Sprintf("%s*/", sep))
}

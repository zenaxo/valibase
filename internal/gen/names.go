package gen

import "github.com/zenaxo/valibase/internal/utils"

type collNameParts struct {
	collectionName     string
	singular           string
	lowerCamelSingular string
	pascalSingular     string
}

func nameParts(collectionName string) collNameParts {
	singular := utils.ToSingular(collectionName)
	return collNameParts{
		collectionName:     collectionName,
		singular:           singular,
		lowerCamelSingular: utils.ToLowerCamelCase(singular),
		pascalSingular:     utils.ToPascalCase(singular),
	}
}

func sanitizeFieldName(name string) string {
	return utils.SanitizeFieldName(name)
}

package gen

import (
	"fmt"

	"github.com/zenaxo/valibase/internal/valibot"
)

var v = valibot.V

func emitField(f fieldSchema) (view, input string) {
	prefix := fmt.Sprintf("%s: ", sanitizeFieldName(f.Name))
	viewSchema, inputSchema := fieldSchemas(f)
	return prefix + viewSchema, prefix + inputSchema
}

func fieldSchemas(f fieldSchema) (view, input string) {
	switch f.Type {
	case FieldBool:
		return boolField(f.Required)
	case FieldAutoDate:
		return autoDateFieldsSchemas(f.Required)
	case FieldDate:
		return dateFieldsSchemas(f.Required)
	case FieldEditor:
		return editorFieldsSchemas(f.Required)
	case FieldEmail:
		return emailFieldsSchema(f.Required)
	case FieldFile:
		return fileFieldSchemas(f)
	case FieldGeoPoint:
		return geoPointFieldSchemas(f.Required)
	case FieldJSON:
		return jsonFieldSchemas(f.Required)
	case FieldNumber:
		return numberFieldSchemas(f)
	case FieldRelation:
		return relationFieldSchemas(f)
	case FieldSelect:
		return selectFieldSchemas(f)
	case FieldText:
		return textFieldSchemas(f)
	case FieldURL:
		return urlFieldSchemas(f)
	default:
		return "v.any()", "v.any()"
	}
}

func pbTextOptional(required bool, base string) (view, input string) {
	if required {
		return base, base
	}
	return v.OptionalTextResponse(base), v.Optional(base)
}

func pbOptionalArray(required bool, base string) string {
	if required {
		return base
	}
	return v.Optional(base, "[]")
}

func diffField(required bool, viewSchema, inputSchema string) (view, input string) {
	if required {
		return viewSchema, inputSchema
	}
	return v.Optional(viewSchema), v.Optional(inputSchema)
}

func jsonFieldSchemas(required bool) (view, input string) {
	return pbTextOptional(required, v.JSON())
}

func editorFieldsSchemas(required bool) (view, input string) {
	return pbTextOptional(required, v.Editor())
}

func dateFieldsSchemas(required bool) (view, input string) {
	return pbTextOptional(required, v.IsoDate())
}

func autoDateFieldsSchemas(required bool) (view, input string) {
	return pbTextOptional(required, v.IsoAutoDate())
}

func boolField(required bool) (view, input string) {
	base := v.Boolean()
	if required {
		return base, v.Literal("true")
	}
	return v.Optional(base), v.Optional(base)
}

func emailFieldsSchema(required bool) (view, input string) {
	base := v.EmailSchema()
	view = v.OptionalTextResponse(base)
	if required {
		input = base
	} else {
		input = v.Optional(base)
	}
	return view, input
}

func geoPointFieldSchemas(required bool) (view, input string) {
	base := v.GeoPoint()
	if required {
		return base, base
	}
	return v.Optional(base), v.Optional(base)
}

func relationFieldSchemas(f fieldSchema) (view, input string) {
	maxSel := f.MaxSelect
	isMany := maxSel == nil || *maxSel != 1

	if isMany {
		base := v.Pipe(v.Array(v.String()), v.Brand("RelationMultiple"))
		return pbOptionalArray(f.Required, base), pbOptionalArray(f.Required, base)
	}

	base := v.Pipe(
		v.String(),
		v.Length(15),
		v.Brand("Relation"),
	)
	return pbTextOptional(f.Required, base)
}

func fileFieldSchemas(f fieldSchema) (view, input string) {
	maxSel := f.MaxSelect
	isMany := maxSel == nil || *maxSel != 1

	mods := []string{v.File}
	if f.MimeTypes != nil {
		mods = append(mods, v.MimeTypes(*f.MimeTypes))
	}
	if f.MaxSize != nil {
		mods = append(mods, v.MaxSize(*f.MaxSize))
	}
	perFileInput := v.Pipe(mods...)

	if isMany {
		base := v.Array(perFileInput)
		input = pbOptionalArray(f.Required, base)
	} else {
		if f.Required {
			input = perFileInput
		} else {
			input = v.Optional(perFileInput)
		}
	}

	if isMany {
		base := v.Array(v.FileName)
		view = pbOptionalArray(f.Required, base)
	} else {
		if f.Required {
			view = v.File
		} else {
			view = v.OptionalTextResponse(v.FileName)
		}
	}

	return view, input
}

func urlFieldSchemas(f fieldSchema) (view, input string) {
	if f.OnlyDomains != nil {
		inputBase := v.OnlyDomains(*f.OnlyDomains)
		if f.Required {
			return v.URLSchema(), inputBase
		}
		return v.OptionalTextResponse(v.URLSchema()), v.Optional(inputBase)
	}

	if f.ExceptDomains != nil {
		inputBase := v.ExceptDomains(*f.ExceptDomains)
		if f.Required {
			return v.URLSchema(), inputBase
		}
		return v.OptionalTextResponse(v.URLSchema()), v.Optional(inputBase)
	}

	return pbTextOptional(f.Required, v.URLSchema())
}

func selectFieldSchemas(f fieldSchema) (view, input string) {
	viewBase := v.Array(v.String())
	maxSel := f.MaxSelect
	enum := v.StringEnum(f.Values)

	if maxSel != nil && *maxSel == 1 {
		return diffField(f.Required, viewBase, enum)
	}

	baseInput := v.Array(enum)

	if maxSel == nil && !f.Required {
		return diffField(f.Required, viewBase, baseInput)
	}

	switch {
	case f.Required && maxSel != nil:
		input = v.Pipe(baseInput, v.MinLength(1), v.MaxLength(*maxSel))
	case f.Required && maxSel == nil:
		input = v.Pipe(baseInput, v.MinLength(1))
	case !f.Required && maxSel != nil:
		input = v.Pipe(baseInput, v.MaxLength(*maxSel))
	}

	return diffField(f.Required, viewBase, input)
}

func textFieldSchemas(f fieldSchema) (view, input string) {
	required := f.Required
	min, max, pattern := f.Min, f.Max, f.Pattern

	if min == nil && max == nil && pattern == nil {
		return pbTextOptional(required, v.String())
	}

	if required {
		view = v.Pipe(v.String())
	} else {
		view = v.OptionalTextResponse(v.String())
	}

	mods := []string{v.String()}

	if (min != nil && max != nil) && (*min == *max) {
		mods = append(mods, v.Length(*min))
	} else {
		if min != nil {
			mods = append(mods, v.MinLength(*min))
		}
		if max != nil {
			mods = append(mods, v.MaxLength(*max))
		}
	}

	if pattern != nil {
		mods = append(mods, v.Pattern(*pattern))
	}

	input = v.Pipe(mods...)
	if !required {
		input = v.Optional(input)
	}

	return view, input
}

func numberFieldSchemas(f fieldSchema) (view, input string) {
	base := v.Number()
	if f.Required {
		view = base
	} else {
		view = v.Optional(base)
	}

	mods := []string{}
	if f.NoDecimals {
		mods = append(mods, v.Integer())
	}

	min := f.MinValue
	max := f.MaxValue

	switch {
	case min != nil && max != nil:
		if *min == *max {
			mods = append(mods, v.Value(*min))
		} else {
			mods = append(mods, v.MinValue(*min), v.MaxValue(*max))
		}
	case min != nil:
		mods = append(mods, v.MinValue(*min))
	case max != nil:
		mods = append(mods, v.MaxValue(*max))
	}

	args := append([]string{v.Number()}, mods...)
	inputBase := v.Pipe(args...)

	if f.Required {
		return view, inputBase
	}
	return view, v.Optional(inputBase)
}

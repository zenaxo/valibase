package gen

import "github.com/pocketbase/pocketbase/core"

func BuildCollections(dbColls []*core.Collection) []collectionRecord {
	collections := make([]collectionRecord, 0, len(dbColls))

	for _, c := range dbColls {
		if c == nil {
			continue
		}

		col := collectionRecord{
			ID:         c.Id,
			Type:       toCollectionType(c),
			Name:       c.Name,
			System:     c.System,
			ListRule:   c.ListRule,
			ViewRule:   c.ViewRule,
			CreateRule: c.CreateRule,
			UpdateRule: c.UpdateRule,
			DeleteRule: c.DeleteRule,
		}

		fields := c.Fields
		col.Fields = make([]fieldSchema, 0, len(fields))

		for _, f := range fields {
			if f == nil {
				continue
			}
			col.Fields = append(col.Fields, toFieldSchema(f))
		}

		collections = append(collections, col)
	}

	return collections
}

func toCollectionType(c *core.Collection) collectionType {
	switch {
	case c.IsAuth():
		return CollectionAuth
	case c.IsView():
		return CollectionView
	default:
		return CollectionBase
	}
}

func toFieldSchema(f core.Field) fieldSchema {
	out := fieldSchema{
		ID:     f.GetId(),
		Name:   f.GetName(),
		Type:   fieldType(f.Type()),
		System: f.GetSystem(),
		Hidden: f.GetHidden(),
		Unique: false, // derive from indexes if needed
	}

	switch tf := f.(type) {
	case *core.TextField:
		out.Required = tf.Required
		if tf.Min != 0 {
			out.Min = ptrInt(tf.Min)
		}
		if tf.Max != 0 {
			out.Max = ptrInt(tf.Max)
		}
		if tf.Pattern != "" {
			out.Pattern = ptrString(tf.Pattern)
		}

	case *core.NumberField:
		out.Required = tf.Required
		out.NoDecimals = tf.OnlyInt
		out.MinValue = tf.Min
		out.MaxValue = tf.Max

	case *core.BoolField:
		out.Required = tf.Required

	case *core.EmailField:
		out.Required = tf.Required

	case *core.URLField:
		out.Required = tf.Required
		if len(tf.OnlyDomains) > 0 {
			out.OnlyDomains = ptrStringSlice(tf.OnlyDomains)
		}
		if len(tf.ExceptDomains) > 0 {
			out.ExceptDomains = ptrStringSlice(tf.ExceptDomains)
		}

	case *core.DateField:
		out.Required = tf.Required

	case *core.JSONField:
		out.Required = tf.Required

	case *core.EditorField:
		out.Required = tf.Required

	case *core.GeoPointField:
		out.Required = tf.Required

	case *core.SelectField:
		out.Required = tf.Required
		if len(tf.Values) > 0 {
			out.Values = append([]string(nil), tf.Values...)
		}
		if tf.MaxSelect != 0 {
			out.MaxSelect = ptrInt(tf.MaxSelect)
		}

	case *core.RelationField:
		out.Required = tf.Required
		if tf.CollectionId != "" {
			out.RelationCollectionID = ptrString(tf.CollectionId)
		}
		if tf.MaxSelect != 0 {
			out.MaxSelect = ptrInt(tf.MaxSelect)
		}

	case *core.FileField:
		out.Required = tf.Required
		if tf.MaxSelect != 0 {
			out.MaxSelect = ptrInt(tf.MaxSelect)
		}
		if tf.MaxSize != 0 {
			out.MaxSize = ptrInt64(tf.MaxSize)
		}
		if len(tf.MimeTypes) > 0 {
			out.MimeTypes = ptrStringSlice(tf.MimeTypes)
		}

	case *core.AutodateField:
		// system-managed; keep defaults
	default:
		// unknown field; keep defaults
	}

	// preserve prior implied-required behavior
	if !out.Required && out.Min != nil && *out.Min > 0 {
		out.Required = true
	}

	return out
}

func ptrInt(v int) *int          { return &v }
func ptrInt64(v int64) *int64    { return &v }
func ptrString(v string) *string { return &v }

func ptrStringSlice(v []string) *[]string {
	cp := append([]string(nil), v...)
	return &cp
}

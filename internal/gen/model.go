package gen

type recordOptions struct {
	// selection / array constraints
	MaxSelect *int

	// text / collection constraints
	Min     *int
	Max     *int
	Pattern *string
	Values  []string

	// numeric constraints
	MinValue   *float64
	MaxValue   *float64
	NoDecimals bool

	// file constraints
	MaxSize   *int64
	MimeTypes *[]string

	// url/domain constraints
	ExceptDomains *[]string
	OnlyDomains   *[]string
}

type fieldType string

const (
	FieldText     fieldType = "text"
	FieldFile     fieldType = "file"
	FieldNumber   fieldType = "number"
	FieldBool     fieldType = "bool"
	FieldEmail    fieldType = "email"
	FieldURL      fieldType = "url"
	FieldDate     fieldType = "date"
	FieldAutoDate fieldType = "autodate"
	FieldSelect   fieldType = "select"
	FieldJSON     fieldType = "json"
	FieldRelation fieldType = "relation"
	FieldEditor   fieldType = "editor"
	FieldGeoPoint fieldType = "geoPoint"
)

type fieldSchema struct {
	ID       string
	Name     string
	Type     fieldType
	System   bool
	Required bool
	Unique   bool
	Hidden   bool

	RelationCollectionID *string

	recordOptions
}

type collectionType string

const (
	CollectionBase collectionType = "base"
	CollectionAuth collectionType = "auth"
	CollectionView collectionType = "view"
)

type collectionRecord struct {
	ID     string
	Type   collectionType
	Name   string
	System bool

	Fields []fieldSchema

	ListRule   *string
	ViewRule   *string
	CreateRule *string
	UpdateRule *string
	DeleteRule *string
}

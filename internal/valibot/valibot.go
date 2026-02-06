package valibot

import (
	"fmt"
	"strings"

	"github.com/zenaxo/valibase/internal/utils"
)

type valibot struct {
	FileName string // fileNameSchema
	File     string // fileSchema
	Relation string // recordIdSchema
	Password string // passwordSchema
}

var V = valibot{
	FileName: "fileNameSchema",
	File:     "fileSchema",
	Relation: "recordIdSchema",
	Password: "passwordSchema",
}

type callProps struct {
	fn      string
	arg     any
	message string
}

func callWithMessage(props *callProps) string {
	if props.message != "" {
		return fmt.Sprintf(`v.%s(%v, "%s")`, props.fn, props.arg, props.message)
	}
	return fmt.Sprintf("v.%s(%v)", props.fn, props.arg)
}

/*
Helper to reduce repetition and boilerplate

Example:
"minLength", 10 ->

	v.minLength(10)
*/
func call(fn string, arg any) string {
	return fmt.Sprintf("v.%s(%v)", fn, arg)
}

/*
Creates a string schema

Example:

	v.string()
*/
func (v valibot) String() string {
	return call("string", "")
}

/*
Creates a number schema

Example:

	v.number()
*/
func (v valibot) Number() string {
	return call("number", "")
}

/*
Creates a number schema

Example:

	v.number()
*/
func (v valibot) Boolean() string {
	return call("boolean", "")
}

/*
Creates a URL schema

Example:

	v.url()
*/
func (v valibot) URL() string {
	return call("url", "")
}

/*
Creates a complete URL schema

Example:

	urlSchema

	->

	v.pipe(
		v.string(),
		v.nonEmpty(),
		v.url('The url is badly formatted'),
		v.brand('URL')
	)
*/
func (v valibot) URLSchema() string {
	return "urlSchema"
}

/*
Creates an email schema with a custom error message

Example:

	v.email('The email is badly formatted')
*/
func (v valibot) Email() string {
	return call("email", "'The url is badly formatted'")
}

/*
Creates an object schema with the provided shape

Example:
"{ id: v.string() }" ->

	v.object({ id: v.string() })
*/
func (v valibot) Object(shape string) string {
	return call("object", shape)
}

/*
Creates an array schema containing schemas

Example:
todoSchema ->

	v.array(todoSchema)
*/
func (v valibot) Array(schemas ...string) string {
	return call("array", createModsString(schemas...))
}

/*
Creates a union schema from multiple schemas

Example:

	v.union([aSchema, bSchema])
*/
func (v valibot) Union(schemas ...string) string {
	return call("union", fmt.Sprintf("[%s]", createModsString(schemas...)))
}

/*
Creates a literal schema from one or more literal values

Example:

	v.literal('a')
*/
func (v valibot) Literal(schemas ...string) string {
	return call("literal", createModsString(schemas...))
}

/*
Creates a picklist schema from possible values

Example:

	v.picklist('a', 'b', 'c')
*/
func (v valibot) Picklist(values ...string) string {
	return call("picklist", createModsString(values...))
}

/*
Wraps one or more schemas in an optional schema

Example:
todoSchema ->

	v.optional(todoSchema)
*/
func (v valibot) Optional(schemas ...string) string {
	return call("optional", createModsString(schemas...))
}

/*
Creates a pipe schema containing schemas

Example:

	v.string(), v.minLength(10) -> v.pipe(v.string(), v.minLength(10))
*/
func (v valibot) Pipe(schemas ...string) string {
	return call("pipe", createModsString(schemas...))
}

/*
Creates a transform schema with a transformation function

Example:

	v.transform((value) => value.trim())
*/
func (v valibot) Transform(transformation string) string {
	return call("transform", transformation)
}

/*
Creates a lazy schema for recursively defined schemas

Example:

	v.lazy(() => userSchema)
*/
func (v valibot) Lazy(schemas ...string) string {
	return call("lazy", createModsString(schemas...))
}

/*
Creates a check schema using a custom validation function

Example:

	v.check((input) => input.length > 0)
*/
func (v valibot) Check(fn string) string {
	return call("check", fn)
}

/*
Creates a branded schema for nominal typing

Example:

	v.brand('UserId')
*/
func (v valibot) Brand(name string) string {
	return call("brand", fmt.Sprintf("'%s'", name))
}

/*
Creates a nonEmpty schema
Ensures the value is not empty

Example:

	v.nonEmpty()
*/
func (v valibot) NonEmpty() string {
	return call("nonEmpty", "")
}

/*
Creates a minLength schema of n length
Mostly used in combination with string schemas

Example:
10 ->

	v.minLength(10)
*/
func (v valibot) MinLength(n int) string {
	props := callProps{
		fn:      "minLength",
		arg:     n,
		message: fmt.Sprintf("Input must be at least %v characters", n),
	}
	return callWithMessage(&props)
}

/*
Creates a maxLength schema of n length
Mostly used in combination with string schemas

Example:
10 ->

	v.maxLength(10)
*/
func (v valibot) MaxLength(n int) string {
	props := callProps{
		fn:      "maxLength",
		arg:     n,
		message: fmt.Sprintf("Input must be at most %v characters", n),
	}
	return callWithMessage(&props)
}

/*
Creates an EXACT length schema of n length
Mostly used in combination with string schemas

Example:
10 ->

	v.length(10)
*/
func (v valibot) Length(n int) string {
	props := callProps{
		fn:      "length",
		arg:     n,
		message: fmt.Sprintf("Input must be exactly %v", n),
	}
	return callWithMessage(&props)
}

/*
Creates a minValue schema of n value
Used in combination with numeric schemas

Example:
10 ->

	v.minValue(10)
*/
func (v valibot) MinValue(n float64) string {
	props := callProps{
		fn:      "minValue",
		arg:     n,
		message: fmt.Sprintf("Input must be greater than %v", n-1),
	}
	return callWithMessage(&props)
}

/*
Creates a maxValue schema of n value
Used in combination with numeric schemas

Example:
10 ->

	v.maxValue(10)
*/
func (v valibot) MaxValue(n float64) string {
	props := callProps{
		fn:      "maxValue",
		arg:     n,
		message: fmt.Sprintf("Input must be lower than %v", n+1),
	}
	return callWithMessage(&props)
}

/*
Creates a value schema of n value
This means that the numeric value must be exactly n
Used in combination with numeric schemas

Example:
10 ->

	v.value(10)
*/
func (v valibot) Value(n float64) string {
	return call("value", n)
}

/*
Creates a regex schema containing pattern p

Example:

	v.regex(new RegExp(`^[\w][\w\.\-]*$`))
*/
func (v valibot) Pattern(p string) string {
	return call("regex", "/"+p+"/, 'Invalid format'")
}

/*
Helper schema that converts optional text to undefined and handles "" as optional

# Should only be used with string modifiers

Example:

	optionalTextResponse(schema)

	->

	v.pipe(
		v.union([v.literal(''), schema]),
		v.transform((input) => (input !== '' ? input : undefined))
	);
*/
func (v valibot) OptionalTextResponse(schemas ...string) string {
	var s string
	if len(schemas) == 0 {
		s = v.String()
	} else {
		s = createModsString(schemas...)
	}
	return "optionalTextResponse(" + s + ")"
}

/*
Creates a stringEnum schema with opts options

Example:
hello, world ->

	stringEnum("hello", "world")
*/
func (v valibot) StringEnum(opts []string) string {
	return "stringEnum(" + utils.ToQuotedStringArray(opts) + ")"
}

/*
Creates an onlyDomains schema with domains

Example:
facebook.com, instagram.com ->

	onlyDomains("facebook.com", "instagram.com")
*/
func (v valibot) OnlyDomains(domains []string) string {
	return "onlyDomains(" + utils.ToQuotedStringArray(domains) + ")"
}

/*
Creates an exceptDomains schema with domains

Example:
facebook.com, instagram.com ->

	exceptDomains("facebook.com", "instagram.com")
*/
func (v valibot) ExceptDomains(domains []string) string {
	return "exceptDomains(" + utils.ToQuotedStringArray(domains) + ")"
}

/*
withExpand returns a withExpand field

Example:

	// mock schema (todoExpandSchema)
	withExpand(todoExpandSchema)
*/
func (v valibot) WithExpand(schemas ...string) string {
	return "withExpand(" + createModsString(schemas...) + ")"
}

/*
SystemFields returns a systemFields field for a collection

Example:

	// mock collection (todos)
	systemFieldsSchema(Collections.Todos)
*/
func (v valibot) SystemFields(collection string) string {
	return "systemFieldsSchema(Collections." + utils.ToPascalCase(collection) + ")"
}

/*
MimeTypes returns a mimeType schema with a user friendly error message

Example:

	// mock types
	v.mimeTypes([]string{"image/jpeg", "image/png"})
*/
func (v valibot) MimeTypes(types []string) string {
	quoted := utils.ToQuotedStringArray(types)

	var capitalized []string
	for _, t := range types {
		fileType := strings.SplitN(t, "/", 2)
		if len(fileType) == 2 {
			capitalized = append(capitalized, strings.ToUpper(fileType[1]))
		} else {
			capitalized = append(capitalized, strings.ToUpper(t))
		}
	}

	errorMsg := "'Please select one of the following file types: " +
		strings.Join(capitalized, " or ") + "'"

	return call("mimeType", "["+quoted+"], "+errorMsg)
}

/*
MaxSize returns a maxSize schema for a file with a human readable label

Example:
10485760 (10 MB) ->

	v.maxSize(10 * 1024 * 1024, "Please select a file smaller than 10MB")
*/
func (v valibot) MaxSize(size int64) string {
	expr, label := utils.SizeExpression(size)
	return call("maxSize", expr+", 'Please select a file smaller than "+label+"'")
}

/*
Def:

	export const isoDateStringSchema = v.pipe(v.string(), v.isoTimestamp(), v.brand('Date'));
*/
func (v valibot) IsoDate() string {
	return "isoDateStringSchema"
}

/*
Def:

	export const isoAutoDateStringSchema = v.pipe(v.string(), v.isoTimestamp(), v.brand('AutoDate'));
*/
func (v valibot) IsoAutoDate() string {
	return "isoAutoDateStringSchema"
}

/*
Def:

	export const jsonSchema = v.pipe(v.string(), v.brand('JSON'));
*/
func (v valibot) JSON() string {
	return "jsonSchema"
}

/*
Def:

	export const editorSchema = v.pipe(v.string(), v.brand('Editor'));
*/
func (v valibot) Editor() string {
	return "editorSchema"
}

/*
Def:

	export const geoPointSchema = v.pipe(
		v.object({
			lon: v.number(),
			lat: v.number()
		}),
		v.brand('GeoPoint')
	);
*/
func (v valibot) GeoPoint() string {
	return "geoPointSchema"
}

func createModsString(mods ...string) string {
	var newMods []string
	for _, m := range mods {
		if m != "" {
			newMods = append(newMods, m)
		}
	}
	return strings.Join(newMods, ", ")
}

/*
Def:

	export const emailSchema = v.pipe(v.string(), v.email(), v.brand('Email'));
*/
func (v valibot) EmailSchema() string {
	return "emailSchema"
}

func (v valibot) Integer() string {
	return call("integer", `"Only integers are allowed."`)
}

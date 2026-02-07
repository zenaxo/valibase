# Valibase — Valibot schema generator for PocketBase

Valibase is a lightweight Go package that generates complete Valibot schemas and
TypeScript types from your PocketBase database schema.

It inspects your collections and fields at runtime and produces strongly typed
frontend-ready validation schemas. Valibase is intended to run during development
and integrates naturally with PocketBase hooks so your types stay in sync as your
database evolves.

---

## Features

- Generates Valibot input, create, and update schemas
- Typed PocketBase client wrappers
- Automatic regeneration when collections change
- Idiomatic Go API
- Minimal public surface: one exported entry point

---

## Requirements

- Go 1.22+
- PocketBase v0.22+
- Valibot in your frontend project

---

## Getting started

### 1. Install

```bash
go get github.com/zenaxo/valibase
```
2. Generate schemas from PocketBase

Integrate Valibase into your PocketBase hooks:
```go
package main

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/zenaxo/valibase/generator"
)

func generate(app core.App, outPath string) {
	if err := generator.GenerateTypes(app, outPath); err != nil {
		app.Logger().Error("failed to generate types", "error", err)
	}
}

func main() {
	outPath := os.Getenv("VALIBASE_OUT_PATH")
	isDev := os.Getenv("APP_ENV") == "development"

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if isDev {
			generate(se.App, outPath)

			// Regenerate after collection updates.
			app.OnCollectionAfterUpdateSuccess().BindFunc(func(e *core.CollectionEvent) error {
				generate(e.App, outPath)
				return e.Next()
			})
		}

		return se.Next()
	})
}
```
3. Use the generated types in TypeScript
```ts
import PocketBase from 'pocketbase'
import { Collections, type TypedPocketBase, type User } from '../database/database'

const pb = new PocketBase('http://localhost:8090') as TypedPocketBase

export async function getUsers(): Promise<User[]> {
  return pb.collection(Collections.Users).getFullList()
}
```

## Generated schema
For a full example of the generated output, see:

- Input schemas contain only fields and types.
- Create schemas include all fields and validation rules.
- Update schemas are partial versions of create schemas.

## Supported fields
Text
- MinLength
- MaxLength
- Length
- Pattern (regex)
- Optional / empty

Number
- MinValue
- MaxValue
- Value
- NoDecimals (`v.integer()`)

Date
- AutoDate fields (branded)
- Date fields

File
- MimeTypes
- MaxSize
- Single / multiple
Other
- Email
- URL
  - OnlyDomains
  - ExceptDomains
- GeoPoint
- JSON (branded text)
- Editor (branded text)
- Relation (branded text)
Select

## Status
Valibase is on early development and the API may change.
Feedback and contributions are welcom.

## License
MIT

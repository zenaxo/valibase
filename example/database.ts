/**
 *
 * This file was automatically @generated and should not be modified
 *
 * To use the client, import createTypedPocketBase
 *
 */

import type PocketBase from "pocketbase";
import type { RecordService } from "pocketbase";

import * as v from "valibot";

// All available PocketBase collections as a const map
export const Collections = {
  Mfas: "_mfas",
  Otps: "_otps",
  ExternalAuths: "_externalAuths",
  AuthOrigins: "_authOrigins",
  Superusers: "_superusers",
  Users: "users",
  Todos: "todos",
} as const;
export type CollectionKey = keyof typeof Collections;
export type CollectionName = (typeof Collections)[CollectionKey];
export const collectionIdSchema = v.pipe(
  v.string(),
  v.length(15),
  v.brand("CollectionId"),
);
export const recordIdSchema = v.pipe(
  v.string(),
  v.length(15),
  v.brand("RecordId"),
);

export const isoDateStringSchema = v.pipe(
  v.string(),
  v.isoTimestamp(),
  v.brand("Date"),
);
export const isoAutoDateStringSchema = v.pipe(
  v.string(),
  v.isoTimestamp(),
  v.brand("AutoDate"),
);

// Basic primitives
export const emailSchema = v.pipe(
  v.string(),
  v.email("Please enter a valid email address"),
  v.brand("Email"),
);
export const fileNameSchema = v.pipe(v.string(), v.brand("FileName"));
export const fileNameArraySchema = v.array(fileNameSchema);
export const fileSchema = v.pipe(v.file(), v.brand("File"));
export const fileArraySchema = v.array(fileSchema);

export const geoPointSchema = v.pipe(
  v.object({
    lon: v.number(),
    lat: v.number(),
  }),
  v.brand("GeoPoint"),
);

export const editorSchema = v.pipe(v.string(), v.brand("Editor"));
export const jsonSchema = v.pipe(v.string(), v.brand("JSON"));
export const urlSchema = v.pipe(
  v.string(),
  v.nonEmpty(),
  v.url("The url is badly formatted."),
  v.brand("URL"),
);

// PocketBase returns undefined fields as an empty string, this handles this issue and converts to undefined
const optionalTextResponse = <
  I extends string,
  O extends string | undefined,
  E extends v.BaseIssue<unknown>,
>(
  schema: v.BaseSchema<I, O, E>,
) =>
  v.pipe(
    v.union([v.literal(""), schema]),
    v.transform((input) => (input !== "" ? input : undefined)),
  );

// Restrict URLs to a fixed allow-list of hostnames
export const onlyDomains = <T extends readonly [string, ...string[]]>(
  ...domains: T
) =>
  v.pipe(
    v.string(),
    v.nonEmpty(),
    v.url("The url is badly formatted"),
    v.brand("OnlyDomains"),
    v.check(
      (input) => {
        let hostname: string;
        try {
          hostname = new URL(input).hostname;
        } catch {
          return false;
        }
        return domains.some(
          (d) => hostname === d || hostname.endsWith(`.${d}`),
        );
      },
      `The URL must be one of: ${domains.join(", ")}`,
    ),
  );

// Forbid URLs that match a blocked list of hostnames
export const exceptDomains = <T extends readonly [string, ...string[]]>(
  ...domains: T
) =>
  v.pipe(
    v.string(),
    v.nonEmpty(),
    v.url("The url is badly formatted"),
    v.brand("ExceptDomains"),
    v.check(
      (input) => {
        let hostname: string;
        try {
          hostname = new URL(input).hostname;
        } catch {
          return false;
        }
        return !domains.some(
          (d) => hostname === d || hostname.endsWith(`.${d}`),
        );
      },
      `The URL must not be any one of: ${domains.join(", ")}`,
    ),
  );

export type CollectionId = v.InferOutput<typeof collectionIdSchema>;
export type RecordId = v.InferOutput<typeof recordIdSchema>;
export type IsoAutoDate = v.InferOutput<typeof isoAutoDateStringSchema>;
export type IsoDate = v.InferOutput<typeof isoDateStringSchema>;
export type Email = v.InferOutput<typeof emailSchema>;
export type FileName = v.InferOutput<typeof fileNameSchema>;
export type FileNameArray = v.InferOutput<typeof fileNameArraySchema>;
export type File = v.InferOutput<typeof fileSchema>;
export type FileArray = v.InferOutput<typeof fileArraySchema>;
export type GeoPoint = v.InferOutput<typeof geoPointSchema>;
export type Editor = v.InferOutput<typeof editorSchema>;
export type JSON = v.InferOutput<typeof jsonSchema>;
export type URL = v.InferOutput<typeof urlSchema>;

export type Expand<E extends object> = {
  expand?: E;
};
export type OAuth2Providers<P extends Array<string>> = {
  oauth2Providers: P;
};

/* =========================================
 * Generic helpers
 * =======================================*/

// Wraps schema in v.optional, keeps type inference intact
export const optional = <I, O, E extends v.BaseIssue<unknown>>(
  schema: v.BaseSchema<I, O, E>,
) => v.optional(schema);

// Tiny helper for string enums using a picklist
export const stringEnum = <T extends readonly [string, ...string[]]>(
  ...values: T
) => v.picklist(values);

/* =========================================
 * System fields + auth/common helpers
 * =======================================*/

const systemFieldsSchema = <N extends CollectionName>(name: N) =>
  v.pipe(
    v.object({
      id: recordIdSchema,
      collectionId: collectionIdSchema,
      collectionName: v.literal(name),
      created: isoAutoDateStringSchema,
      updated: isoDateStringSchema,
    }),
    v.brand("SystemFields"),
  );

// Basic password and password-related schemas
export const passwordSchema = v.pipe(
  v.string(),
  v.minLength(8),
  v.brand("Password"),
);
export type Password = v.InferOutput<typeof passwordSchema>;

// Schema used when creating a password + confirmation pair
export const passwordConfirmSchema = v.pipe(
  v.object({
    password: passwordSchema,
    passwordConfirm: v.string(),
  }),
  v.forward(
    v.check((i) => i.password === i.passwordConfirm, "Passwords do not match"),
    ["passwordConfirm"],
  ),
);

// Schema used when updating a password (optional fields + consistency checks)
export const newPasswordSchema = v.pipe(
  v.object({
    password: v.optional(passwordSchema),
    passwordConfirm: v.optional(v.string()),
    oldPassword: v.optional(v.string()),
  }),
  v.forward(
    v.check(
      (i) => !i.password || !!i.passwordConfirm,
      "Please confirm your new password",
    ),
    ["passwordConfirm"],
  ),
  v.forward(
    v.check(
      (i) => !i.password || i.password === i.passwordConfirm,
      "Passwords do not match",
    ),
    ["passwordConfirm"],
  ),
  v.forward(
    v.check(
      (i) => !i.password || !!i.oldPassword,
      "Old password is required to change password",
    ),
    ["oldPassword"],
  ),
);

// Base schema helpers used by all collections
export const createBaseSchema = <
  TEntries extends v.ObjectEntries,
  TMessage extends v.ErrorMessage<v.ObjectIssue> | undefined,
>(
  fields: v.ObjectSchema<TEntries, TMessage>,
) => fields;

export const updateBaseSchema = <
  TEntries extends v.ObjectEntries,
  TMessage extends v.ErrorMessage<v.ObjectIssue> | undefined,
>(
  fields: v.ObjectSchema<TEntries, TMessage>,
) => v.partial(fields);

// Auth-aware schema helpers that compose base fields with auth schemas
export const createAuthSchema = <
  TEntries extends v.ObjectEntries,
  TMessage extends v.ErrorMessage<v.ObjectIssue> | undefined,
>(
  schema: v.ObjectSchema<TEntries, TMessage>,
) => v.intersect([schema, passwordConfirmSchema]);

export const updateAuthSchema = <
  TEntries extends v.ObjectEntries,
  TMessage extends v.ErrorMessage<v.ObjectIssue> | undefined,
>(
  schema: v.ObjectSchema<TEntries, TMessage>,
) => v.intersect([schema, newPasswordSchema]);

/*==========================================================================================
_MFAS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "_mfas"
 */
export const mfaResponse = v.object({
  ...systemFieldsSchema("_mfas").entries,
  collectionRef: v.string(),
  recordRef: v.string(),
  method: v.string(),
});

/*
 * Input schema for creating/updating "_mfas"
 */
export const mfaInput = v.object({
  collectionRef: v.string(),
  recordRef: v.string(),
  method: v.string(),
});

export type MfaFields = v.InferOutput<typeof mfaResponse>;

export type Mfa = MfaFields;

/**
 * Create/Update schemas and their inferred input types for "Mfa" records.
 */
export const createMfaSchema = createBaseSchema(mfaInput);
export const updateMfaSchema = updateBaseSchema(mfaInput);

// Inferred input types from the above schemas
export type CreateMfaInput = v.InferOutput<typeof createMfaSchema>;
export type UpdateMfaInput = v.InferOutput<typeof updateMfaSchema>;

/*==========================================================================================
_OTPS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "_otps"
 */
export const otpResponse = v.object({
  ...systemFieldsSchema("_otps").entries,
  collectionRef: v.string(),
  recordRef: v.string(),
});

/*
 * Input schema for creating/updating "_otps"
 */
export const otpInput = v.object({
  collectionRef: v.string(),
  recordRef: v.string(),
});

export type OtpFields = v.InferOutput<typeof otpResponse>;

export type Otp = OtpFields;

/**
 * Create/Update schemas and their inferred input types for "Otp" records.
 */
export const createOtpSchema = createBaseSchema(otpInput);
export const updateOtpSchema = updateBaseSchema(otpInput);

// Inferred input types from the above schemas
export type CreateOtpInput = v.InferOutput<typeof createOtpSchema>;
export type UpdateOtpInput = v.InferOutput<typeof updateOtpSchema>;

/*==========================================================================================
_EXTERNALAUTHS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "_externalAuths"
 */
export const externalAuthResponse = v.object({
  ...systemFieldsSchema("_externalAuths").entries,
  collectionRef: v.string(),
  recordRef: v.string(),
  provider: v.string(),
  providerId: v.string(),
});

/*
 * Input schema for creating/updating "_externalAuths"
 */
export const externalAuthInput = v.object({
  collectionRef: v.string(),
  recordRef: v.string(),
  provider: v.string(),
  providerId: v.string(),
});

export type ExternalAuthFields = v.InferOutput<typeof externalAuthResponse>;

export type ExternalAuth = ExternalAuthFields;

/**
 * Create/Update schemas and their inferred input types for "ExternalAuth" records.
 */
export const createExternalAuthSchema = createBaseSchema(externalAuthInput);
export const updateExternalAuthSchema = updateBaseSchema(externalAuthInput);

// Inferred input types from the above schemas
export type CreateExternalAuthInput = v.InferOutput<
  typeof createExternalAuthSchema
>;
export type UpdateExternalAuthInput = v.InferOutput<
  typeof updateExternalAuthSchema
>;

/*==========================================================================================
_AUTHORIGINS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "_authOrigins"
 */
export const authOriginResponse = v.object({
  ...systemFieldsSchema("_authOrigins").entries,
  collectionRef: v.string(),
  recordRef: v.string(),
  fingerprint: v.string(),
});

/*
 * Input schema for creating/updating "_authOrigins"
 */
export const authOriginInput = v.object({
  collectionRef: v.string(),
  recordRef: v.string(),
  fingerprint: v.string(),
});

export type AuthOriginFields = v.InferOutput<typeof authOriginResponse>;

export type AuthOrigin = AuthOriginFields;

/**
 * Create/Update schemas and their inferred input types for "AuthOrigin" records.
 */
export const createAuthOriginSchema = createBaseSchema(authOriginInput);
export const updateAuthOriginSchema = updateBaseSchema(authOriginInput);

// Inferred input types from the above schemas
export type CreateAuthOriginInput = v.InferOutput<
  typeof createAuthOriginSchema
>;
export type UpdateAuthOriginInput = v.InferOutput<
  typeof updateAuthOriginSchema
>;

/*==========================================================================================
_SUPERUSERS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "_superusers"
 */
export const superuserResponse = v.object({
  ...systemFieldsSchema("_superusers").entries,
  email: optionalTextResponse(emailSchema),
  emailVisibility: v.optional(v.boolean()),
  verified: v.optional(v.boolean()),
});

/*
 * Input schema for creating/updating "_superusers"
 */
export const superuserInput = v.object({
  email: emailSchema,
  emailVisibility: v.optional(v.boolean()),
  verified: v.optional(v.boolean()),
});

export type SuperuserFields = v.InferOutput<typeof superuserResponse>;

export type Superuser = SuperuserFields;

/**
 * Create/Update schemas and their inferred input types for "Superuser" records.
 */
export const createSuperuserSchema = createAuthSchema(superuserInput);
export const updateSuperuserSchema = updateAuthSchema(superuserInput);

// Inferred input types from the above schemas
export type CreateSuperuserInput = v.InferOutput<typeof createSuperuserSchema>;
export type UpdateSuperuserInput = v.InferOutput<typeof updateSuperuserSchema>;

/*==========================================================================================
USERS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "users"
 */
export const userResponse = v.object({
  ...systemFieldsSchema("users").entries,
  avatar: optionalTextResponse(fileNameSchema),
  name: optionalTextResponse(v.string()),
  username: v.pipe(v.string()),
  tasks: v.optional(
    v.pipe(v.array(v.string()), v.brand("RelationMultiple")),
    [],
  ),
  languages: v.optional(v.array(v.string())),
  address: v.optional(geoPointSchema),
  email: optionalTextResponse(emailSchema),
  emailVisibility: v.optional(v.boolean()),
  verified: v.optional(v.boolean()),
});

/*
 * Input schema for creating/updating "users"
 */
export const userInput = v.object({
  avatar: v.optional(
    v.pipe(
      fileSchema,
      v.mimeType(
        ["image/jpeg", "image/png", "image/svg+xml", "image/gif", "image/webp"],
        "Please select one of the following file types: JPEG or PNG or SVG+XML or GIF or WEBP",
      ),
    ),
  ),
  name: v.optional(
    v.pipe(
      v.string(),
      v.maxLength(255, "Input must be at most 255 characters"),
    ),
  ),
  username: v.pipe(
    v.string(),
    v.minLength(3, "Input must be at least 3 characters"),
    v.maxLength(80, "Input must be at most 80 characters"),
    v.regex(/^[\w][\w\.\-]*$/, "Invalid format"),
  ),
  tasks: v.optional(
    v.pipe(v.array(v.string()), v.brand("RelationMultiple")),
    [],
  ),
  languages: v.optional(stringEnum("Swedish", "English", "Spanish")),
  address: v.optional(geoPointSchema),
  email: emailSchema,
  emailVisibility: v.optional(v.boolean()),
  verified: v.optional(v.boolean()),
});

export type UserFields = v.InferOutput<typeof userResponse>;

/**
 * Relations that can be expanded when loading "users"
 */
export type UserExpand = {
  tasks?: Todo[];
};

export type User = UserFields & Expand<Partial<UserExpand>>;

/**
 * Create/Update schemas and their inferred input types for "User" records.
 */
export const createUserSchema = createAuthSchema(userInput);
export const updateUserSchema = updateAuthSchema(userInput);

// Inferred input types from the above schemas
export type CreateUserInput = v.InferOutput<typeof createUserSchema>;
export type UpdateUserInput = v.InferOutput<typeof updateUserSchema>;

/*==========================================================================================
TODOS COLLECTION
==========================================================================================*/

/**
 * Raw field schema for "todos"
 */
export const todoResponse = v.object({
  ...systemFieldsSchema("todos").entries,
  name: v.pipe(v.string()),
  description: v.pipe(v.string()),
  link: optionalTextResponse(urlSchema),
  state: v.optional(v.array(v.string())),
  completed: optionalTextResponse(isoDateStringSchema),
});

/*
 * Input schema for creating/updating "todos"
 */
export const todoInput = v.object({
  name: v.pipe(
    v.string(),
    v.minLength(3, "Input must be at least 3 characters"),
    v.maxLength(120, "Input must be at most 120 characters"),
  ),
  description: v.pipe(
    v.string(),
    v.minLength(3, "Input must be at least 3 characters"),
    v.maxLength(150, "Input must be at most 150 characters"),
  ),
  link: v.optional(onlyDomains("mytodos.com")),
  state: v.optional(stringEnum("Due", "Completed")),
  completed: v.optional(isoDateStringSchema),
});

export type TodoFields = v.InferOutput<typeof todoResponse>;

export type Todo = TodoFields;

/**
 * Create/Update schemas and their inferred input types for "Todo" records.
 */
export const createTodoSchema = createBaseSchema(todoInput);
export const updateTodoSchema = updateBaseSchema(todoInput);

// Inferred input types from the above schemas
export type CreateTodoInput = v.InferOutput<typeof createTodoSchema>;
export type UpdateTodoInput = v.InferOutput<typeof updateTodoSchema>;

// Central registry of all generated collection schemas
export const registry = {
  // Schemas for the "_mfas" collection
  _mfas: {
    response: mfaResponse,
    create: createMfaSchema,
    update: updateMfaSchema,
  },

  // Schemas for the "_otps" collection
  _otps: {
    response: otpResponse,
    create: createOtpSchema,
    update: updateOtpSchema,
  },

  // Schemas for the "_externalAuths" collection
  _externalAuths: {
    response: externalAuthResponse,
    create: createExternalAuthSchema,
    update: updateExternalAuthSchema,
  },

  // Schemas for the "_authOrigins" collection
  _authOrigins: {
    response: authOriginResponse,
    create: createAuthOriginSchema,
    update: updateAuthOriginSchema,
  },

  // Schemas for the "_superusers" collection
  _superusers: {
    response: superuserResponse,
    create: createSuperuserSchema,
    update: updateSuperuserSchema,
  },

  // Schemas for the "users" collection
  users: {
    response: userResponse,
    create: createUserSchema,
    update: updateUserSchema,
  },

  // Schemas for the "todos" collection
  todos: {
    response: todoResponse,
    create: createTodoSchema,
    update: updateTodoSchema,
  },
} as const;

export type CollectionsMap = typeof registry;
export type CollectionNameKey = keyof CollectionsMap;

// Helper type map: collection name -> strongly typed record
export type ResponseTypes = {
  _mfas: Mfa;
  _otps: Otp;
  _externalAuths: ExternalAuth;
  _authOrigins: AuthOrigin;
  _superusers: Superuser;
  users: User;
  todos: Todo;
};
// Schema helpers
type CreateSchemaOf<N extends CollectionNameKey> = CollectionsMap[N]["create"];
type UpdateSchemaOf<N extends CollectionNameKey> = CollectionsMap[N]["update"];

// Record type for a given collection name key
export type RecordOf<N extends CollectionNameKey> = ResponseTypes[N];

// Type of the payload for create operations
export type Create<N extends CollectionNameKey> = v.InferOutput<
  CreateSchemaOf<N>
>;

// Type of the payload for update operations
export type Update<N extends CollectionNameKey> = v.InferOutput<
  UpdateSchemaOf<N>
>;

/**
 * # TypedPocketBase
 * - Automatic schema generation
 * - Run time validation
 * - Includes validation through valibot
 * ### Usage:
 *
 * 		import type { TypedPocketBase } from '.../path-to-database.ts'
 *
 * 		const pb = new PocketBase(PUBLIC_PB) as TypedPocketBase
 *
 *		// Returns User
 * 		const users = pb.collection('users').getOne()
 *
 */
export type TypedPocketBase = {
  collection<T extends CollectionNameKey>(
    idOrName: T,
  ): RecordService<ResponseTypes[T]>;
} & PocketBase;

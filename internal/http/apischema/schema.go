package apischema

type Schema struct {
	Type                 string
	Format               string
	Description          string
	Enum                 []string
	Properties           []Property
	Required             []string
	Items                *Schema
	AdditionalProperties *Schema
}

type Property struct {
	Name     string
	Schema   Schema
	Required bool
}

func F(name string, schema Schema) Property {
	return Property{Name: name, Schema: schema}
}

func RF(name string, schema Schema) Property {
	return Property{Name: name, Schema: schema, Required: true}
}

func Str(desc string) Schema {
	return Schema{Type: "string", Description: desc}
}

func StrEnum(desc string, values ...string) Schema {
	return Schema{Type: "string", Description: desc, Enum: values}
}

func DateTime(desc string) Schema {
	return Schema{Type: "string", Format: "date-time", Description: desc}
}

func Int(desc string) Schema {
	return Schema{Type: "integer", Format: "int32", Description: desc}
}

func Int64(desc string) Schema {
	return Schema{Type: "integer", Format: "int64", Description: desc}
}

func Num(desc string) Schema {
	return Schema{Type: "number", Format: "double", Description: desc}
}

func Bool(desc string) Schema {
	return Schema{Type: "boolean", Description: desc}
}

func Any(desc string) Schema {
	return Schema{Description: desc}
}

func Obj(desc string, props ...Property) Schema {
	schema := Schema{Type: "object", Description: desc, Properties: props}
	for _, prop := range props {
		if prop.Required {
			schema.Required = append(schema.Required, prop.Name)
		}
	}
	return schema
}

func Arr(desc string, items Schema) Schema {
	return Schema{Type: "array", Description: desc, Items: &items}
}

func Map(desc string, values Schema) Schema {
	return Schema{Type: "object", Description: desc, AdditionalProperties: &values}
}

type Param struct {
	Name        string
	Description string
	Type        string
	Required    bool
}

func PathParam(name, desc string) Param {
	return Param{Name: name, Description: desc, Type: "string", Required: true}
}

func QueryParam(name, desc string) Param {
	return Param{Name: name, Description: desc, Type: "string"}
}

func QueryParamInt(name, desc string) Param {
	return Param{Name: name, Description: desc, Type: "integer"}
}

type Body struct {
	Content string
	Schema  Schema
}

func JSON(schema Schema) *Body {
	return &Body{Content: "application/json", Schema: schema}
}

func Form(schema Schema) *Body {
	return &Body{Content: "multipart/form-data", Schema: schema}
}

type Operation struct {
	Method     string
	Path       string
	Summary    string
	Tags       []string
	Auth       string
	PathParams []Param
	Query      []Param
	Request    *Body
	Response   Schema
	Raw        bool
	Content    string
	Code       int
	Errors     []int
}

const (
	ContentJSON   = "application/json"
	ContentSSE    = "text/event-stream"
	ContentBinary = "application/octet-stream"
	ContentText   = "text/plain"
)

func (op Operation) ResponseContent() string {
	if op.Content != "" {
		return op.Content
	}
	return ContentJSON
}

func (op Operation) SuccessCode() int {
	if op.Code != 0 {
		return op.Code
	}
	return 200
}

func (op Operation) Authenticated() bool {
	return op.Auth != ""
}

const (
	AuthUser     = "authenticated"
	AuthAdmin    = "instance admin"
	AuthOrgOwner = "organization owner"
	AuthOrgAdmin = "organization admin"
	AuthOrg      = "organization member"
)

func Scope(scope string) string {
	return "scope " + scope
}

func ProjectRole(role string) string {
	return "project " + role
}

func ServiceRole(role string) string {
	if role == "" {
		return "service member"
	}
	return "service " + role
}

func All() []Operation {
	ops := []Operation{}
	ops = append(ops, authOperations()...)
	ops = append(ops, organizationOperations()...)
	ops = append(ops, projectOperations()...)
	ops = append(ops, appOperations()...)
	ops = append(ops, serviceOperations()...)
	ops = append(ops, deploymentOperations()...)
	ops = append(ops, databaseOperations()...)
	ops = append(ops, backupOperations()...)
	ops = append(ops, serverOperations()...)
	ops = append(ops, migrationOperations()...)
	ops = append(ops, analyticsOperations()...)
	ops = append(ops, clusterOperations()...)
	ops = append(ops, systemOperations()...)
	return ops
}

func Find(method, path string) (Operation, bool) {
	for _, op := range All() {
		if op.Method == method && op.Path == path {
			return op, true
		}
	}
	return Operation{}, false
}

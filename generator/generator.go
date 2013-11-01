package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
)

// Generates a parser for a single Go AST type.
func GenerateEncoder(typeSpec *ast.TypeSpec, w io.Writer) error {
	var b bytes.Buffer
	name := typeSpec.Name.Name

	// Ignore type non-struct specs.
	typ, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	// Retrieve all relevant fields.
	fields := make([]*ast.Field, 0)
	for _, f := range typ.Fields.List {
		if ident, ok := f.Type.(*ast.Ident); ok {
			switch ident.Name {
			case "string", "int":
				fields = append(fields, f)
			}
		}
	}

	// Generate the encoder.
	fmt.Fprintf(&b, "type %sJSONEncoder struct {", name)
	fmt.Fprintln(&b, "w io.Writer")
	fmt.Fprintln(&b, "}")

	// Generate the constructor.
	fmt.Fprintf(&b, "func New%sJSONEncoder(w io.Writer) *%sJSONEncoder {\n", name, name)
	fmt.Fprintf(&b, "return &%sJSONEncoder{w: w}\n", name)
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b, "")

	fmt.Fprintf(&b, "func (e *%sJSONEncoder) Encode(v *%s) error {\n", name, name)
	for _, f := range fields {
		if ident, ok := f.Type.(*ast.Ident); ok {
			name := f.Names[0]

			switch ident.Name {
            case "string":
                fmt.Fprintf(&b, "\tif err := encoding.WriteString(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "int":
                fmt.Fprintf(&b, "\tif err := encoding.WriteInt(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "uint":
                fmt.Fprintf(&b, "\tif err := encoding.WriteUint(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "float32":
                fmt.Fprintf(&b, "\tif err := encoding.WriteFloat32(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "float64":
                fmt.Fprintf(&b, "\tif err := encoding.WriteFloat64(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "bool":
                fmt.Fprintf(&b, "\tif err := encoding.WriteBool(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
			}
		}
	}
	fmt.Fprintf(&b, "\treturn nil\n")
	fmt.Fprintf(&b, "}\n")

	// Write to formatted output stream.
	if _, err := b.WriteTo(w); err != nil {
		return err
	}

	return nil
}

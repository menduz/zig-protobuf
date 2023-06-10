package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	. "google.golang.org/protobuf/reflect/protoreflect"
	"os"
	"path"
	"strings"
)

type APIContext struct {
	File    protogen.File
	Imports []Import
}

type Import struct {
	Path         string
	PathAbsolute string
}

func canonical(absolutePath string) string {
	canonical := strings.ReplaceAll(absolutePath, "/", "_")
	canonical = strings.ReplaceAll(canonical, "\\", "_")
	canonical = strings.ReplaceAll(canonical, ".", "_")
	return canonical
}

func fieldName(field Name) Name {
	if field == "packed" || field == "type" || field == "null" || field == "error" {
		return "@\"" + field + "\""
	}
	return field
}

func getFieldKindName(field *protogen.Field) (string, error) {
	var prefix = ""

	switch field.Desc.Kind() {
	case MessageKind, StringKind, BytesKind:
		prefix = "?"
	default:
		if field.Desc.HasOptionalKeyword() || field.Desc.IsList() {
			prefix = "?"
		}
	}

	if field.Desc.IsMap() {
		return "", fmt.Errorf("Maps are not supported, field type: %s", field.Desc.Name)
	}

	if field.Desc.IsList() {
		prefix = "ArrayList("
	}

	switch field.Desc.Kind() {
	case Sint32Kind, Sfixed32Kind, Int32Kind:
		prefix += "i32"
	case Uint32Kind, Fixed32Kind:
		prefix += "u32"
	case Sint64Kind, Sfixed64Kind, Int64Kind:
		prefix += "i64"
	case Uint64Kind, Fixed64Kind:
		prefix += "u64"
	case BoolKind:
		prefix += "bool"
	case DoubleKind:
		prefix += "f64"
	case FloatKind:
		prefix += "f32"
	case StringKind, BytesKind: // TODO: validate if repeated strings and bytes are supported
		prefix += "[]const u8"
	case MessageKind:
		if field.Message.Desc.ParentFile().FullName() != field.Desc.ParentFile().FullName() {
			prefix += canonical(string(field.Message.Desc.ParentFile().Path()))
			prefix += "."
		}
		prefix += string(field.Message.Desc.Name())
	case EnumKind:
		prefix += string(field.Enum.Desc.Name())
	default:
		return "", fmt.Errorf("unmanaged field type in getFieldKindName %s", field.Desc.Kind())
	}

	if field.Desc.IsList() {
		prefix += ")"
	}

	return prefix, nil
}

func getFieldDescriptor(field *protogen.Field) (string, error) {

	typeName, err := getFieldKindName(field)

	if err != nil {
		return "", err
	}

	if field.Desc.IsList() {

		kindOfList := ".List"

		if field.Desc.IsPacked() {
			kindOfList = ".PackedList"
		}

		switch field.Desc.Kind() {
		case StringKind, BytesKind:
			return fmt.Sprintf("fd(%d, .{ %s = .String }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		case Sfixed64Kind, Sfixed32Kind, Fixed32Kind, Fixed64Kind, DoubleKind, FloatKind, Int64Kind:
			return fmt.Sprintf("fd(%d, .{ %s = .FixedInt }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		case Sint32Kind, Sint64Kind:
			return fmt.Sprintf("fd(%d, .{ %s = .{ .Varint = .ZigZagOptimized } }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		case Uint32Kind, Uint64Kind, BoolKind, Int32Kind:
			return fmt.Sprintf("fd(%d, .{ %s = .{ .Varint = .Simple } }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		case MessageKind:
			return fmt.Sprintf("fd(%d, .{ %s = .SubMessage }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		case EnumKind:
			return fmt.Sprintf("fd(%d, .{ %s = .{ .Varint = .Simple } }, %s)", field.Desc.Number(), kindOfList, typeName), nil
		default:
			return "", fmt.Errorf("unmanaged field type in  getFieldDescriptor 1 %s", field.Desc.Kind())
		}
	} else {
		switch field.Desc.Kind() {
		case Sfixed64Kind, Sfixed32Kind, Fixed32Kind, Fixed64Kind, DoubleKind, FloatKind, Int64Kind:
			return fmt.Sprintf("fd(%d, .FixedInt, %s)", field.Desc.Number(), typeName), nil
		case Sint32Kind, Sint64Kind:
			return fmt.Sprintf("fd(%d, .{ .Varint = .ZigZagOptimized }, %s)", field.Desc.Number(), typeName), nil
		case Uint32Kind, Uint64Kind, BoolKind, Int32Kind:
			return fmt.Sprintf("fd(%d, .{ .Varint = .Simple }, %s)", field.Desc.Number(), typeName), nil
		case StringKind, BytesKind:
			return fmt.Sprintf("fd(%d, .String, %s)", field.Desc.Number(), typeName), nil
		case MessageKind:
			return fmt.Sprintf("fd(%d, .{ .SubMessage = {} }, %s)", field.Desc.Number(), typeName), nil
		case EnumKind:
			return fmt.Sprintf("fd(%d, .{ .Varint = .Simple }, %s)", field.Desc.Number(), typeName), nil
		default:
			return "", fmt.Errorf("unmanaged field type in  getFieldDescriptor 2 %s", field.Desc.Kind())
		}
	}
}

func generateFieldDescriptor(field *protogen.Field, g *protogen.GeneratedFile) error {
	if fieldDesc, err := getFieldDescriptor(field); err != nil {
		return err
	} else {
		g.P("        .", fieldName(field.Desc.Name()), " = ", fieldDesc, ",")
	}
	return nil
}

func generateFieldDef(field *protogen.Field, g *protogen.GeneratedFile) error {
	if fieldKindName, err := getFieldKindName(field); err != nil {
		return err
	} else {
		// if field.Desc.Kind() == MessageKind && field.Desc.Message().ParentFile() != field.Desc.ParentFile() {
		// 	g.P("    //", field.Desc.Message().ParentFile())
		// }

		for _, c := range strings.Split(field.Comments.Leading.String(), "\n") {
			if strings.TrimSpace(c) != "" {
				g.P(c)
			}
		}
		g.P("    ", fieldName(field.Desc.Name()), ": ", fieldKindName, ",")
	}
	return nil
}

func generateMessages(g *protogen.GeneratedFile, messages []*protogen.Message) error {
	for _, m := range messages {
		msgName := m.Desc.Name()

		for _, c := range strings.Split(m.Comments.Leading.String(), "\n") {
			if strings.TrimSpace(c) != "" {
				g.P(c)
			}
		}

		g.P("pub const ", msgName, " = struct {")

		// field definitions
		for _, field := range m.Fields {
			if err := generateFieldDef(field, g); err != nil {
				return err
			}
		}

		generateEnums(g, m.Enums)
		generateMessages(g, m.Messages)

		g.P()
		// field descriptors
		g.P("    pub const _desc_table = .{")
		for _, field := range m.Fields {
			if err := generateFieldDescriptor(field, g); err != nil {
				return err
			}
		}
		g.P("    };")
		g.P()

		g.P("    pub fn encode(self: ", msgName, ", allocator: Allocator) ![]u8 {")
		g.P("        return pb_encode(self, allocator);")
		g.P("    }")
		g.P()

		g.P("    pub fn decode(input: []const u8, allocator: Allocator) !", msgName, " {")
		g.P("        return pb_decode(", msgName, ", input, allocator);")
		g.P("    }")
		g.P()

		g.P("    pub fn init(allocator: Allocator) ", msgName, " {")
		g.P("        return pb_init(", msgName, ", allocator);")
		g.P("    }")
		g.P()

		g.P("    pub fn deinit(self: ", msgName, ") void {")
		g.P("        return pb_deinit(self);")
		g.P("    }")

		g.P("};")
		g.P("")
	}
	return nil
}

func generateEnums(g *protogen.GeneratedFile, enums []*protogen.Enum) {
	for _, m := range enums {
		msgName := m.Desc.Name()

		for _, c := range strings.Split(m.Comments.Leading.String(), "\n") {
			if strings.TrimSpace(c) != "" {
				g.P(c)
			}
		}
		g.P("pub const ", msgName, " = enum(i32) {") // TODO: type
		for _, f := range m.Values {
			for _, c := range strings.Split(f.Comments.Leading.String(), "\n") {
				if strings.TrimSpace(c) != "" {
					g.P(c)
				}
			}
			g.P("    ", fieldName(f.Desc.Name()), " = ", f.Desc.Number(), ",")
		}
		g.P("    _,")
		g.P("};")
		g.P("")
	}
}

func generateFile(p *protogen.Plugin, ctx *APIContext) error {
	f := ctx.File

	// Skip generating file if there is no message.
	if len(f.Messages) == 0 {
		return nil
	}

	fullPath := f.Proto.GetName()
	filename := ""
	if ext := path.Ext(fullPath); ext == ".proto" {
		base := path.Base(fullPath)
		filename = base[:len(base)-len(path.Ext(base))]
	}
	filename += ".pb.zig"
	filename = path.Join(path.Dir(fullPath), filename)

	g := p.NewGeneratedFile(filename, "")

	g.P("// Code generated by protoc-gen-zig-go")
	g.P()
	g.P("const std = @import(\"std\");")
	g.P("const mem = std.mem;")
	g.P("const Allocator = mem.Allocator;")
	g.P("const ArrayList = std.ArrayList;")
	g.P()
	g.P("const protobuf = @import(\"protobuf\");")
	g.P("const FieldDescriptor = protobuf.FieldDescriptor;")
	g.P("const pb_decode = protobuf.pb_decode;")
	g.P("const pb_encode = protobuf.pb_encode;")
	g.P("const pb_deinit = protobuf.pb_deinit;")
	g.P("const pb_init = protobuf.pb_init;")
	g.P("const fd = protobuf.fd;")
	g.P()

	for _, dep := range ctx.Imports {
		varname := canonical(dep.PathAbsolute)

		g.P("const ", varname, " = @import(\"", dep.Path, "\");")
	}

	generateEnums(g, f.Enums)

	generateMessages(g, f.Messages)

	return nil
}

func (ctx *APIContext) ApplyImports(f *protogen.File) {
	var deps []Import

	for _, dep := range f.Proto.Dependency {
		if dep == "google/protobuf/timestamp.proto" {
			continue
		}
		importPath := path.Dir(dep)
		sourceDir := path.Dir(f.Proto.GetName() + ".pb.zig")
		sourceComponents := strings.Split(sourceDir, fmt.Sprintf("%c", os.PathSeparator))
		distanceFromRoot := len(sourceComponents)
		for _, pathComponent := range sourceComponents {
			if strings.HasPrefix(importPath, pathComponent) {
				importPath = strings.TrimPrefix(importPath, pathComponent)
				distanceFromRoot--
			}
		}
		deps = append(deps, Import{fullImportPath(zigPbFilename(dep), importPath, distanceFromRoot), dep})
	}
	ctx.Imports = deps
}

func zigPbFilename(name string) string {
	if ext := path.Ext(name); ext == ".proto" {
		base := path.Base(name)
		name = base[:len(base)-len(path.Ext(base))]
	}

	name += ".pb.zig"
	return name
}

func fullImportPath(fileName string, importPath string, distanceFromRoot int) string {
	fullPath := fileName
	fullPath = path.Join(importPath, fullPath)
	if distanceFromRoot > 0 {
		for i := 0; i < distanceFromRoot; i++ {
			fullPath = path.Join("..", fullPath)
		}
	}

	return fullPath
}

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}

			ctx := APIContext{File: *file}

			ctx.ApplyImports(file)

			if err := generateFile(plugin, &ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

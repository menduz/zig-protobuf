const warn = @import("std").debug.warn;
const std = @import("std");
const plugin = @import("google/protobuf/compiler/plugin.pb.zig");
const descriptor = @import("google/protobuf/descriptor.pb.zig");
const mem = std.mem;

const allocator = std.heap.page_allocator;

pub fn main() !void {
    const stdin = &std.io.getStdIn();
    const stdout = &std.io.getStdOut();

    // Read the contents
    const buffer_size = 1024 * 1024 * 10;
    const file_buffer = try stdin.readToEndAlloc(allocator, buffer_size);
    defer allocator.free(file_buffer);

    // plugin
    var request: plugin.CodeGeneratorRequest = try plugin.CodeGeneratorRequest.decode(file_buffer, allocator);

    std.debug.print("Parameter: {s}\n", .{request.parameter orelse "<empty>"});
    std.debug.print("Files to generate:\n", .{});
    for (request.file_to_generate.items) |a| {
        std.debug.print("  {s}\n", .{a});
    }

    var response = plugin.CodeGeneratorResponse.init(allocator);

    std.debug.print("Files:\n", .{});
    for (request.proto_file.items) |proto| {
        const t: descriptor.FileDescriptorProto = proto;
        std.debug.print("  {?s}\n", .{t.name});
        for (t.dependency.items) |dep| {
            std.debug.print("    depends on: {?s}\n", .{dep});
        }

        try response.file.append(try outputFile(request, proto));
    }

    const r = try response.encode(allocator);
    _ = try stdout.write(r);
}

fn outputFile(ctx: plugin.CodeGeneratorRequest, file: descriptor.FileDescriptorProto) !plugin.CodeGeneratorResponse.File {
    _ = ctx;

    var ret = plugin.CodeGeneratorResponse.File.init(allocator);

    ret.name = try std.fmt.allocPrint(allocator, "{?s}.zig", .{file.name});
    ret.content = "// test";

    return ret;
}

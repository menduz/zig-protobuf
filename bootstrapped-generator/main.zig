const warn = @import("std").debug.warn;
const std = @import("std");
const plugin = @import("google/protobuf/compiler/plugin.pb.zig");
const mem = std.mem;

pub fn main() !void {
    const allocator = std.heap.page_allocator;
    const stdin = &std.io.getStdIn();

    // Read the contents
    const buffer_size = 2000;
    const file_buffer = try stdin.readToEndAlloc(allocator, buffer_size);
    defer allocator.free(file_buffer);

    // plugin
    var request = plugin.CodeGeneratorRequest.decode(file_buffer, allocator);
    std.debug.print("{any}", .{request});
}

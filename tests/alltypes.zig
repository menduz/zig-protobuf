const std = @import("std");
const testing = std.testing;

const protobuf = @import("protobuf");
const tests = @import("./generated/tests.pb.zig");
const longName = @import("./generated/some/really/long/name/which/does/not/really/make/any/sense/but/sometimes/we/still/see/stuff/like/this.pb.zig");

pub fn printAllDecoded(input: []const u8) !void {
    var iterator = protobuf.WireDecoderIterator{ .input = input };
    std.debug.print("Decoding: {s}\n", .{std.fmt.fmtSliceHexUpper(input)});
    while (try iterator.next()) |extracted_data| {
        std.debug.print("  {any}\n", .{extracted_data});
    }
}

test "long name" {
    // - this test allocates an object only. used to instruct zig to try to compile the file
    // - it also ensures that SubMessage deinit() works
    var demo = longName.WouldYouParseThisForMePlease.init(testing.allocator);
    demo.field = .{ .field = "asd" };
    defer demo.deinit();

    const obtained = try demo.encode(testing.allocator);
    defer testing.allocator.free(obtained);
}

test "packed int32_list" {
    var demo = tests.Packed.init(testing.allocator);
    try demo.int32_list.append(0x01);
    try demo.int32_list.append(0x02);
    try demo.int32_list.append(0x03);
    try demo.int32_list.append(0x04);
    defer demo.deinit();

    const obtained = try demo.encode(testing.allocator);
    defer testing.allocator.free(obtained);

    try testing.expectEqualSlices(u8, &[_]u8{ 0x0A, 0x04, 0x01, 0x02, 0x03, 0x04 }, obtained);

    const decoded = try tests.Packed.decode(obtained, testing.allocator);
    defer decoded.deinit();
    try testing.expectEqualSlices(i32, demo.int32_list.items, decoded.int32_list.items);

    // TODO: cross test against Packed type
}

test "unpacked int32_list" {
    var demo = tests.UnPacked.init(testing.allocator);
    try demo.int32_list.append(0x01);
    try demo.int32_list.append(0x02);
    try demo.int32_list.append(0x03);
    try demo.int32_list.append(0x04);
    defer demo.deinit();

    const obtained = try demo.encode(testing.allocator);
    defer testing.allocator.free(obtained);

    try testing.expectEqualSlices(u8, &[_]u8{ 0x08, 0x01, 0x08, 0x02, 0x08, 0x03, 0x08, 0x04 }, obtained);

    const decoded = try tests.UnPacked.decode(obtained, testing.allocator);
    defer decoded.deinit();
    try testing.expectEqualSlices(i32, demo.int32_list.items, decoded.int32_list.items);

    // TODO: cross test against Packed type
}

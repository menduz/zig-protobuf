// Code generated by protoc-gen-zig
 ///! package some.really.long.name.which.does.not.really.make.any.sense.but.sometimes.we.still.see.stuff.like.this
const std = @import("std");
const mem = std.mem;
const Allocator = mem.Allocator;
const ArrayList = std.ArrayList;

const protobuf = @import("protobuf");
const pb_decode = protobuf.pb_decode;
const pb_encode = protobuf.pb_encode;
const pb_deinit = protobuf.pb_deinit;
const pb_init = protobuf.pb_init;
const fd = protobuf.fd;

pub const WouldYouParseThisForMePlease = struct {
    field: ?Test,

    pub const _desc_table = .{
        .field = fd(1, .{ .SubMessage = {} }),
    };


    pub fn encode(self: WouldYouParseThisForMePlease, allocator: Allocator) ![]u8 {
        return pb_encode(self, allocator);
    }
    pub fn decode(input: []const u8, allocator: Allocator) !WouldYouParseThisForMePlease {
        return pb_decode(WouldYouParseThisForMePlease, input, allocator);
    }
    pub fn init(allocator: Allocator) WouldYouParseThisForMePlease {
        return pb_init(WouldYouParseThisForMePlease, allocator);
    }
    pub fn deinit(self: WouldYouParseThisForMePlease) void {
        return pb_deinit(self);
    }
};

pub const Test = struct {
    field: ?[]const u8,

    pub const _desc_table = .{
        .field = fd(1, .String),
    };


    pub fn encode(self: Test, allocator: Allocator) ![]u8 {
        return pb_encode(self, allocator);
    }
    pub fn decode(input: []const u8, allocator: Allocator) !Test {
        return pb_decode(Test, input, allocator);
    }
    pub fn init(allocator: Allocator) Test {
        return pb_init(Test, allocator);
    }
    pub fn deinit(self: Test) void {
        return pb_deinit(self);
    }
};

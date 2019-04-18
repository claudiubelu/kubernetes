package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
    "go_test",
)

go_binary(
    name = "consume-memory",
    embed = [":go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["consume_memory.go"],
    importpath = "k8s.io/kubernetes/test/images/resource-consumer/consume-memory",
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["consume_memory_test.go"],
    embed = [":go_default_library"],
)

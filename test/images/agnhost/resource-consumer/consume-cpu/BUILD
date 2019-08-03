package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "consume_cpu.go",
        "consume_cpu_windows.go",
        "main.go",
    ],
    importpath = "k8s.io/kubernetes/test/images/agnhost/resource-consumer/consume-cpu",
    deps = [
        "//vendor/github.com/spf13/cobra:go_default_library",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//vendor/bitbucket.org/bertimus9/systemstat:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//vendor/golang.org/x/sys/windows:go_default_library",
        ],
        "//conditions:default": [],
    }),
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

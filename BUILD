load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix ledstripinterface
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "ledstripinterface",
    visibility = ["//visibility:private"],
    deps = [
        "//lib:go_default_library",
        "//pb:go_default_library",
    ],
)

go_binary(
    name = "ledstripinterface",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

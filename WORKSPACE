#  Copyright (c) 2023 Uber Technologies, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6b65cb7917b4d1709f9410ffe00ecf3e160edf674b78c54a894471320862184f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.0/rules_go-v0.39.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "ecba0f04f96b4960a5b250c8e8eeec42281035970aa8852dda73098274d14a1d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

http_archive(
    name = "tree_sitter",
    build_file = "//third_party/tree_sitter:BUILD.system",
    sha256 = "b355e968ec2d0241bbd96748e00a9038f83968f85d822ecb9940cbe4c42e182e",
    strip_prefix = "tree-sitter-0.20.7/lib",
    urls = [
        "https://github.com/tree-sitter/tree-sitter/archive/refs/tags/v0.20.7.tar.gz",
    ],
)

http_archive(
    name = "tree_sitter_go",
    build_file = "//third_party/tree_sitter/go:BUILD.system",
    sha256 = "ba7ca9571c7515d64782eae492130057449232d8c13c7b38d65b36c450dbc72f",
    strip_prefix = "tree-sitter-go-0.19.1/src",
    urls = [
        "https://github.com/tree-sitter/tree-sitter-go/archive/refs/tags/v0.19.1.tar.gz",
    ],
)

http_archive(
    name = "tree_sitter_java",
    build_file = "//third_party/tree_sitter/java:BUILD.system",
    sha256 = "d6e425ecec92f73fecd2c1f701b6f46858af2380db16b220c9dfd5075b049863",
    strip_prefix = "tree-sitter-java-0.19.1/src",
    urls = [
        "https://github.com/tree-sitter/tree-sitter-java/archive/refs/tags/v0.19.1.tar.gz",
    ],
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_yoheimuta_go_protoparser_v4",
    importpath = "github.com/yoheimuta/go-protoparser/v4",
    sum = "h1:80LGfVM25sCoNDD08hv9O0ShQMjoTrIE76j5ON+gq3U=",
    version = "v4.7.0",
)

go_repository(
    name = "org_uber_go_thriftrw",
    importpath = "go.uber.org/thriftrw",
    sum = "h1:pRuFLzbGvTcnYwGSjizWRHlbJUzGhu84sRiL1h1kUd8=",
    version = "v1.29.2",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
    version = "v3.0.1",
)

go_repository(
    name = "com_github_xwb1989_sqlparser",
    importpath = "github.com/xwb1989/sqlparser",
    sum = "h1:zzrxE1FKn5ryBNl9eKOeqQ58Y/Qpo3Q9QNxKHX5uzzQ=",
    version = "v0.0.0-20180606152119-120387863bf2",
)

go_repository(
    name = "org_uber_go_multierr",
    importpath = "go.uber.org/multierr",
    sum = "h1:S0h4aNzvfcFsC3dRF1jLoaov7oRaKqRGC/pUEJ2yvPQ=",
    version = "v1.10.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:+h33VjcLVPDHtOdpUCuF+7gSuG3yGIftsP1YvFihtJ8=",
    version = "v1.8.2",
)

go_repository(
    name = "net_starlark_go",
    build_extra_args = ["-exclude=vendor"],
    build_file_generation = "on",
    build_file_proto_mode = "disable",
    importpath = "go.starlark.net",
    sum = "h1:TckyegD9V0g1yvueWiMQ4BfEGeieq7yKq05oNVzXgJw=",
    version = "v0.0.0-20230118143110-ddd531cdb2da",
)

go_rules_dependencies()

go_register_toolchains(version = "1.20.1")

gazelle_dependencies()

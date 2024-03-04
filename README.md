# Last Diff Analyzer

This repo contains code for the Last Diff Analyzer: a multi-language tool for checking semantic equivalence for code.

## Description

Last Diff Analyzer is a multi-language platform-agnostic automated tool for checking semantic equivalence of code pieces. It can be integrated in code hosting platforms to:

* automatically approve the most recent changes (last diffs) to improve development velocity;
* skip expensive CI tests for semantically equivalent code changes;
* and more...

Most importantly, Last Diff Analyzer provides a highly-extensible multi-language static analysis framework that unifies the language structures of multiple languages while keeping unique language constructs. On top of this foundation, we have currently implemented support for Golang and Java in this repo, while providing easy extensions to other languages such as Kotlin, Swift, or TypeScript.

More technical details can be found in the paper.

## Integration
> **Warning** 
> We do not provide integration code in this repository, the following lists the logic for integrating with any code hosting platforms that support getting diff files and the actual source code between or at specific commits.

Last Diff Analyzer is platform-agnostic, meaning that you can integrate it with your preferred code hosting platforms. 

All it requires is the source directories and diff files for base and last code revisions to be checked.

For example, if current `main` branch is on commit A, and the developers have created two code revisions B and C that
branch off `main`. You will need to give the following inputs to Last Diff Analyzer to get the result:

* source directories for changed files on commit B
* diff file containing changes from A...B
* source directories for changed files on commit C
* diff file containing changes from A...C

The following code shows an implementation for programatically invoking Last Diff Analyzer:

```golang
// The path to the diff files for A...B (baseDiff) and A...C (lastDiff).
baseDiff, lastDiff := ..., ...
// The path to source directories containing changed files on B (baseDir) and C (lastDir).
baseDir, lastDir := ..., ...

// Setup the analyzer(s) with a configurable feature set.
analyzer := Analyzer{
    BaseDiff: baseDiff,
    LastDiff: lastDiff,
}
// Add core analyzer (currently supporting Golang and Java) with different features on.
analyzer.SubAnalyzers = append(a.SubAnalyzers, &core.Analyzer{RenamingOn: true, LoggingOn: true})

// Add basic approvers (only supporting comment changes) for other languages we support.
subAnalyzers := [...]Analyzer{
    &bazel.Analyzer{AnalyzableFileName: "BUILD.bazel.test"},
    &yaml.Analyzer{},
    &gomod.Analyzer{},
    &sql.Analyzer{},
    &protobuf.Analyzer{},
    &starlark.Analyzer{},
    &thrift.Analyzer{}
}
analyzer.SubAnalyzers = append(analyzer.SubAnalyzers, subAnalyzers...)

// Run the analyzer
status, err := a.Run(baseDir, lastDir)
// status will be analyzer.Approve, analyzer.Reject, or analyzer.Failure.
```

You can also check `analyzer/analyzer_test.go` for more details.

## Run

As discussed, we do not provide actual integration code in this repository. However, a comprehensive test suites (`analyzer/testdata`) are bundled to illustrate the features of Last Diff Analyzer. To run Last Diff Analyzer against the test suites, simply run

```bash
> bazel test //analyzer/...
```

Note that you need to have bazel installed since this project is built and managed by bazel.

Alternatively, you can run bazel using the official bazel docker container with directory mounting:
```bash
> docker run --rm -v $(pwd):/workspace -w /workspace gcr.io/bazel-public/bazel:6.3.2 test //analyzer/...
```

See the [official documentation](https://bazel.build/install/docker-container) on more details for running the official bazel docker container. 

## Citing this work

You are encouraged to cite the following paper if you use this tool in academic research:

```bibtex
@inproceedings{wang2023last,
  title={Last Diff Analyzer: Multi-language Automated Approver for Behavior-Preserving Code Revisions},
  author={Wang, Yuxin and Welc, Adam and Clapp, Lazaro and Chen, Lingchao},
  booktitle={Proceedings of the 31st ACM Joint European Software Engineering Conference and Symposium on the Foundations of Software Engineering},
  pages={1693--1704},
  year={2023}
}
```

## Support

For research purposes only.  This project is released "as-is" with no support expected.  Updates may be posted, but are not guaranteed and no timeline is available.

## License

This project is copyright 2023 Uber Technologies, Inc., and licensed under Apache 2.0.

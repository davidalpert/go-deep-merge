<!-- PROJECT SHIELDS -->
<!--
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![License: GPL v3][license-shield]][license-url]
<!-- [![Issues][issues-shield]][issues-url] -->
<!-- [![Forks][forks-shield]][forks-url] -->
<!-- ![GitHub Contributors][contributors-shield] -->
<!-- ![GitHub Contributors Image][contributors-image-url] -->

<!-- PROJECT LOGO -->
<br />
<h1 align="center">goconfig</h1>

<p align="center">
  A tool for managing, merging, and applying config files.
  <br />
  <a href="./README.md">README</a>
  ·
  <a href="./CONTRIBUTING.md">CONTRIBUTING</a>
  ·
  <a href="./CHANGELOG.md"><strong>CHANGELOG</strong></a>
  <br />
  <!-- <a href="https://github.com/davidalpert/goconfig">View Demo</a>
  · -->
  <a href="https://github.com/davidalpert/goconfig/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/goconfig/issues">Request Feature</a>
</p>

## Changelog


<a name="v1.0.0"></a>
## [v1.0.0] - 2023-01-21
### Bug Fixes
- deep merging hash arrays
- edge cases
- overwrite unmergeables introduced an orphan empty string
- knockout handling
- port additional tests and fix edge cases
- port additional test cases and fix some edge cases

### Build
- download vale styles only if missing
- give up on doctor for release builds
- use correct syntax for continue-on-error
- only use continue_on_error where it makes sense
- allow release builds to opt out of regenerating version detail
- don't regenerate version detail for the same version
- don't run bundler in a release build
- don't install goreleaser using doctor in github_actions
- ensure preview-release-notes has a version to use
- ensure we can run doctor in release builds
- install all go tools before ruby tools
- fix workflow syntax
- can't find go 1.16
- add github actions
- automate release tasks
- newer version_gen.go that accepts a VERSION env var
- run unit tests and rebuild before running features
- prevent common build tasks from running more than once
- add a task to install ruby gems
- update to latest version of some tools
- fix tools script self-reference
- add ruby support to start building out feature specs
- add a ci task
- install goconvey as a testing tool
- **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0

### Code Refactoring
- collect providers in a folder
- make 'go config get' generic
- pull aws paramstore client out to new package and simplify the internal API
- extract common cfgset merge logic
- migrate DeepMerge package into a v1 folder and add a basic CLI app

### Docs
- add contributing docs

### Features
- sync merged configs to AWS Parameter Store
- get keys from aws recursive
- goconfig aws get <key>
- goconfig sync folder
- goconfig merge files <src_path> <dest_path>

### Fixup
- with go.mod

### Test Coverage
- add infrastructure for cucumber/aruba features
- port knockout tests (all passing)
- add goconvey


<a name="v0.0.1"></a>
## v0.0.1 - 2023-01-19

[Unreleased]: https://github.com/davidalpert/paper-moon/compare/v1.0.0...HEAD
[v1.0.0]: https://github.com/davidalpert/paper-moon/compare/v0.0.1...v1.0.0

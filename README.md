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
  A tool for managing, merging, and applying configuration files.
  <br />
  <a href="./README.md"><strong>README</strong></a>
  ·
  <a href="./CHANGELOG.md">CHANGELOG</a>
  <br />
  <!-- <a href="https://github.com/davidalpert/goconfig">View Demo</a>
  · -->
  <a href="https://github.com/davidalpert/goconfig/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/goconfig/issues">Request Feature</a>
</p>

<details open="open">
  <summary><h2 style="display: inline-block">Table of contents</h2></summary>

- [About the project](#about-the-project)
  - [Built with](#built-with)
- [Getting started](#getting-started)
  - [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Local development](#local-development)
  - [Prerequisites](#prerequisites)
  - [Taskfile targets](#taskfile-targets)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

</details>

<!-- ABOUT THE PROJECT -->
## About the project

<!-- [![Paper-Moon Screen Shot][product-screenshot]](https://example.com) -->

### Built with

* [Golang 1.16](https://golang.org/)

<!-- GETTING STARTED -->
## Getting started

To get a local copy up and running follow these simple steps.

### Installation

1. Clone the repository
   ```sh
   git clone https://github.com/davidalpert/goconfig.git
   ```

<!-- USAGE EXAMPLES -->
## Usage

Run the `goconfig` binary with no arguments to show command-line help.

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/davidalpert/goconfig/issues) for a list of proposed features and known issues.

<!-- CONTRIBUTING -->
## Local development

### Prerequisites

`goconfig` builds and ships as a single-file binary with no prerequisites to make installation and updates easy.

* [golang](https://golang.org/doc/manage-install)
  * with a working go installation:
    ```
    go install golang.org/dl/go1.16@latest
    go1.16 download
    ```
* this repository includes a `./.tools/doctor.sh` script which validates your local environment and installs or helps you install missing dependencies
* [Taskfile](https://taskfile.dev/) a task runner

### Taskfile targets

This repository includes a `Taskfile` for help running common tasks.

Run `task` with no arguments to list the available targets:
```
$ task
goconfig v0.0.0

task: Available tasks for this project:
* autotest:       run tests continuously using goconvey's test UI
* build:          build
* gen:            run code-generation
* help:           list targets
* test:           run tests
```

<!-- CONTRIBUTING -->
## Contributing

See the [CONTRIBUTING](CONTRIBUTING.md) guide for local development setup and contribution guidelines.

1. Fork the Project
2. Create your Feature Branch
    ```
    git checkout -b feature/AmazingFeature
    ```
3. Commit your Changes
    ```
    git commit -m 'Add some AmazingFeature'
    ```
4. Push to the Branch
    ```
    git push origin feature/AmazingFeature
    ```
5. Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the GPU v3 License. See [LICENSE](LICENSE) for more information.

<!-- CONTACT -->
## Contact

David Alpert - [@davidalpert](https://twitter.com/davidalpert)

Project Link: [https://github.com/davidalpert/goconfig](https://github.com/davidalpert/goconfig)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/davidalpert/goconfig
[contributors-image-url]: https://contrib.rocks/image?repo=davidalpert/goconfig
[forks-shield]: https://img.shields.io/github/forks/davidalpert/goconfig
[forks-url]: https://github.com/davidalpert/goconfig/network/members
[issues-shield]: https://img.shields.io/github/issues/davidalpert/goconfig
[issues-url]: https://github.com/davidalpert/goconfig/issues
[license-shield]: https://img.shields.io/badge/License-GPLv3-blue.svg
[license-url]: https://www.gnu.org/licenses/gpl-3.0
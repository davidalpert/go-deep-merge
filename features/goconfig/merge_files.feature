Feature: merge files <src_path> <dest_path>

  this command can be used to merge two files together and the
  result is written to standard out and formatted based on the
  output flags (defaulting to yaml)

  NOTE: the keys in yaml output are sorted

  Background:
    Given I have installed "goconfig" locally into the path
    And a file named "config/app1/default.yaml" with:
      """
      log_level: WARN
      env: REQUIRED
      auth:
        username: default_user
        password: default_pass
      """
    And a file named "config/app1/dev.yaml" with:
      """
      log_level: DEBUG
      env: dev
      auth:
        password: dev_pass
      """

  # @announce-stdout
  Scenario: yaml merge to yaml output sorts keys
    When I successfully run `goconfig merge files config/app1/dev.yaml config/app1/default.yaml -o yaml`
    Then the stdout should contain:
      """
      auth:
        password: dev_pass
        username: default_user
      env: dev
      log_level: DEBUG
      """

  Scenario: yaml merge to json output sorts keys
    When I successfully run `goconfig merge files config/app1/dev.yaml config/app1/default.yaml -o json`
    Then the stdout should contain:
      """
      {
        "auth": {
          "password": "dev_pass",
          "username": "default_user"
        },
        "env": "dev",
        "log_level": "DEBUG"
      }
      """

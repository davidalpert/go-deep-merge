Feature: sync folder

  this command can be used to apply a convention-based merge
  pattern and then sync differences to a folder

  Background:
    Given I have installed "goconfig" locally into the path
    And I use a fixture named "simple-configuration"

  # @announce-stdout
  Scenario: yaml merge to yaml output sorts keys
    When I successfully run `goconfig sync folder config --out-folder out`
    Then the file "out/app1/prd.yaml" should contain:
      """
      env: prd
      log_level: WARN
      region: unknown
      """
    And the file "out/app1/stg.yaml" should contain:
      """
      env: stg
      log_level: WARN
      region: unknown
      """
    And the file "out/app1/dev.yaml" should contain:
      """
      env: dev
      log_level: DEBUG
      region: unknown
      """
    And the file "out/app1/dev.us-east-1.yaml" should contain:
      """
      env: dev
      log_level: DEBUG
      region: us-east-1
      """

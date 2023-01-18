Feature: usage
  Background:
    Given I have installed "goconfig" locally into the path

  # @announce-stdout @announce-stderr
  Scenario: help
    When I run `goconfig`
    Then the stdout should show usage
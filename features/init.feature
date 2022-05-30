Feature: Application Bootstrap

    Background:
        Given working directory is temporary

    Scenario: Find all transaction in range with invalid credentials provider
        Given there is a file ".n26/config.toml" with content:
        """
        [n26]
        credentials = 'invalid'
        device = 'ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c'
        """

        When I run command "transactions -d"

        Then console output is:
        """
        panic: could not build credentials provider option: unsupported credentials prov
        ider
        """

    Scenario: Find all transaction in range with invalid format
        When I run command "transactions -d --format invalid"

        Then console output is:
        """
        panic: unknown output format
        """

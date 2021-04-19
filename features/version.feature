Feature: Version

    Scenario: Short info
        Given I run command "version"

        Then console output matches:
        """
        dev \(rev: ; go[0-9.]+; [^)]+\)
        """

    Scenario: Full info
        Given I run command "version -f"

        Then console output matches:
        """
        dev \(rev: ; go[0-9.]+; [^)]+\)

        build user:.*
        build date:.*

        dependencies:.*
        """

Feature: Configure

    Scenario: No config file and do not use keychain
        Given I see a confirm prompt "Do you want to save your credentials to system keychain? (y/N)", I answer no

        When I run command "config"

        Then console output is:
        """
        ? Do you want to save your credentials to system keychain? No

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = ""
          device = "<uuid>"

        """

    Scenario: No config file and use keychain
        Given I see a confirm prompt "Do you want to save your credentials to system keychain? (y/N)", I answer yes
        And I see a password prompt "Enter username (input is hidden, leave it empty if no change) >", I answer "user@example.org"
        And I see a password prompt "Enter password (input is hidden, leave it empty if no change) >", I answer "123456"

        When I run command "config"

        Then console output is:
        """
        ? Do you want to save your credentials to system keychain? Yes
        ? Enter username (input is hidden, leave it empty if no change) > **************
        **
        ? Enter password (input is hidden, leave it empty if no change) > ******

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = "keychain"
          device = "<uuid>"

        """

        And keychain has username "user@example.org" and password "123456"

    Scenario: Have config file and change device id
        Given I create a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        And I see a confirm prompt "Do you want to generate a new device id? (y/N)", I answer yes
        And I see a confirm prompt "Do you want to save your credentials to system keychain? (Y/n)", I answer no

        When I run command "config"

        Then console output is:
        """
        ? Do you want to generate a new device id? Yes
        ? Do you want to save your credentials to system keychain? No

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = ""
          device = "<uuid>"

        """

        And configured device is not "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        And keychain has no credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

    Scenario: Have config file and no change device id
        Given I create a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I see a confirm prompt "Do you want to generate a new device id? (y/N)", I answer no
        And I see a confirm prompt "Do you want to save your credentials to system keychain? (Y/n)", I answer no

        When I run command "config"

        Then console output is:
        """
        ? Do you want to generate a new device id? No
        ? Do you want to save your credentials to system keychain? No

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = ""
          device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

        """

    Scenario: Have config file and start using keychain
        Given I create a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = ""
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I see a confirm prompt "Do you want to generate a new device id? (y/N)", I answer no
        And I see a confirm prompt "Do you want to save your credentials to system keychain? (y/N)", I answer yes
        And I see a password prompt "Enter username (input is hidden, leave it empty if no change) >", I answer "user@example.org"
        And I see a password prompt "Enter password (input is hidden, leave it empty if no change) >", I answer "123456"

        When I run command "config"

        Then console output is:
        """
        ? Do you want to generate a new device id? No
        ? Do you want to save your credentials to system keychain? Yes
        ? Enter username (input is hidden, leave it empty if no change) > **************
        **
        ? Enter password (input is hidden, leave it empty if no change) > ******

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = "keychain"
          device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

        """

        And keychain has username "user@example.org" and password "123456"

    Scenario: Have config file and stop using keychain
        Given I create a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        And I see a confirm prompt "Do you want to generate a new device id? (y/N)", I answer no
        And I see a confirm prompt "Do you want to save your credentials to system keychain? (Y/n)", I answer no

        When I run command "config"

        Then console output is:
        """
        ? Do you want to generate a new device id? No
        ? Do you want to save your credentials to system keychain? No

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = ""
          device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

        """

        And keychain has no credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

    Scenario: Have config file and do not change credentials
        Given I create a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        And I see a confirm prompt "Do you want to generate a new device id? (y/N)", I answer no
        And I see a confirm prompt "Do you want to save your credentials to system keychain? (Y/n)", I answer yes
        And I see a password prompt "Enter username (input is hidden, leave it empty if no change) >", I answer "user@example.org"
        And I see a password prompt "Enter password (input is hidden, leave it empty if no change) >", I answer ""

        When I run command "config"

        Then console output is:
        """
        ? Do you want to generate a new device id? No
        ? Do you want to save your credentials to system keychain? Yes
        ? Enter username (input is hidden, leave it empty if no change) > **************
        **
        ? Enter password (input is hidden, leave it empty if no change) >

        saved
        """

        And there is a file ".n26/config.toml" with content:
        """

        [n26]
          credentials = "keychain"
          device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"

        """

        And keychain has username "user@example.org" and password "123456"

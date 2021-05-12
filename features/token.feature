Feature: Token Reuse

    Background:
        Given working directory is temporary

    Scenario: Find all transaction in range without username
        Given I see a password prompt "Enter username (input is hidden) >", I interrupt

        When I run command "transactions -v"

        Then console output is:
        """
        ? Enter username (input is hidden) >
        could not find transactions: could not get token: missing username
        """

    Scenario: Find all transaction in range without username (debug)
        Given I see a password prompt "Enter username (input is hidden) >", I interrupt

        When I run command "transactions -d"

        Then console output is:
        """
        ? Enter username (input is hidden) >
        panic: could not find transactions: could not get token: missing username
        """

    Scenario: Find all transaction in range without password
        Given I see a password prompt "Enter username (input is hidden) >", I answer "user@example.org"
        Given I see a password prompt "Enter password (input is hidden) >", I interrupt

        When I run command "transactions -v"

        Then console output is:
        """
        ? Enter username (input is hidden) > ****************
        ? Enter password (input is hidden) >
        could not find transactions: could not get token: missing password
        """

    Scenario: Find all transaction in range and reuse the token
        Given there is a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain
        And n26 receives a success login request with username "user@example.org", password "123456" and device id "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        And n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        # Reuse the token, there is no extra requests to login or refresh token.
        Given n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain

    Scenario: Find all transaction in range and refresh the token
        Given there is a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain
        And n26 receives a success login request with username "user@example.org", password "123456" and device id "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        And n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        # Token is expired so I can use the refresh token to get a new one after 30 minutes.
        Given I freeze the clock
        And I add 30m to the clock
        And n26 receives a success refresh token request
        And n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain

    Scenario: Find all transaction in range and relogin
        Given there is a file ".n26/config.toml" with content:
        """
        [n26]
            credentials = "keychain"
            device = "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        """
        And I create a credentials "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain with content:
        """
        {"username":"user@example.org", "password":"123456"}
        """
        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain
        And n26 receives a success login request with username "user@example.org", password "123456" and device id "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        And n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        # Both token and refresh toke are expired after one day, so I have to get a new one.
        Given I freeze the clock
        And I add 24h to the clock
        And n26 receives a success login request with username "user@example.org", password "123456" and device id "ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c"
        And n26 receives a request to find all transactions in between "2020-01-02T03:04:05Z" and "2020-02-03T04:05:06Z" and responses:
        """
        [
            {
                "id": "801d35f4-f550-446a-974a-0d5dc2c1f55d",
                "userId": "7e3f710b-349d-4203-9c5d-cfbc716e1b8e",
                "type": "CT",
                "amount": 10,
                "currencyCode": "EUR",
                "visibleTS": 1617631557000,
                "partnerBic": "NTSBDEB1XXX",
                "partnerName": "Jane Doe",
                "accountId": "98f0afa3-e906-493a-a37f-afe29c7f9f2e",
                "partnerIban": "DEXX1001100126XXXXXXXX",
                "category": "micro-v2-income",
                "cardId": "f2252c42-c188-4b43-ab68-131024782b3d",
                "userCertified": 1617545157000,
                "pending": false,
                "transactionNature": "NORMAL",
                "createdTS": 1617541557000,
                "smartLinkId": "fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c",
                "smartContactId": "3edce485-6853-40bf-aa08-309c2eb3e7db",
                "linkId": "6f06f5fb-074d-4242-b280-db2af2fe6405",
                "confirmed": 1617545157000
            }
        ]
        """

        When I run command "transactions --from 2020-01-02T03:04:05Z --to 2020-02-03T04:05:06Z --format json"

        Then console output is:
        """
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        [{"id":"801d35f4-f550-446a-974a-0d5dc2c1f55d","userId":"7e3f710b-349d-4203-9c5d-
        cfbc716e1b8e","type":"CT","amount":10,"currencyCode":"EUR","visibleTS":161763155
        7000,"partnerBic":"NTSBDEB1XXX","partnerName":"Jane Doe","accountId":"98f0afa3-e
        906-493a-a37f-afe29c7f9f2e","partnerIban":"DEXX1001100126XXXXXXXX","category":"m
        icro-v2-income","cardId":"f2252c42-c188-4b43-ab68-131024782b3d","userCertified":
        1617545157000,"pending":false,"transactionNature":"NORMAL","createdTS":161754155
        7000,"smartLinkId":"fcdec3cb-47b2-4ca3-b98d-b326e1cc5a0c","smartContactId":"3edc
        e485-6853-40bf-aa08-309c2eb3e7db","linkId":"6f06f5fb-074d-4242-b280-db2af2fe6405
        ","confirmed":1617545157000}]
        """

        Then I delete token "user@example.org:ed24ad1f-94a4-4ac6-a097-f2bc54f58f0c" in keychain

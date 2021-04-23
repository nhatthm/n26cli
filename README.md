# N26 CLI

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/n26cli)](https://github.com/nhatthm/n26cli/releases/latest)
[![Build Status](https://github.com/nhatthm/n26cli/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/n26cli/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/n26cli/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/n26cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/n26cli)](https://goreportcard.com/report/github.com/nhatthm/n26cli)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/n26cli)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnhatthm%2Fn26cli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnhatthm%2Fn26cli?ref=badge_shield)

An awesome tool for managing your N26 account from the terminal

## Prerequisites

- `Go >= 1.16`
- Mac or Linux (dbus) (see [more](https://github.com/zalando/go-keyring#dependencies))

## Install

TBD

## Configuration

You don't need to configure the tool before using it. However, N26 requires a Device ID (in UUID format) to use its APIs. If you don't provide a specific Device
ID, you may be blocked with "Too Many Login Attempts" error. Therefore, it's strongly recommended configuring beforehand.

Just run `n26 config` and follow the prompt, it will create a new Device ID automatically for you and persist to
`~/.n26/config.toml`. This config file will be automatically loaded whenever you run a command.

If you wish to remember your username (and/or password), just say `yes` while being asked and type it in.

```
$ n26 config
? Do you want to generate a new device id? No
? Do you want to save your credentials to system keychain? Yes
? Enter username (input is hidden, leave it empty if no change) >
? Enter password (input is hidden, leave it empty if no change) >

saved
```

The given username and password will be persisted to your system keychain/keyring, securely. If you don't want to save your password, just leave in empty. The
tool will ask for it whenever you run a command. (see [Authentication](#Authentication))

## Usage

```
Usage:
  n26 [command]

Available Commands:
  config       configure
  help         Help about any command
  transactions show all transactions in a time period
  version      show version

Flags:
  -c, --config string   configuration file (default "~/.n26/config.toml")
  -d, --debug           debug output
  -h, --help            help for n26
  -v, --verbose         verbose output

Use "n26 [command] --help" for more information about a command.
```

## Authentication

The tool will look for your credentials following this order, top goes first:

- System Keychain/Keyring (if you decide to use it).
- `--username` and `--password` from the command arguments.
- `N26_USERNAME` and `N26_PASSWORD` from env vars.
- If username or password is missing, the tool will ask for it.

After successfully login with your password, you will be asked to confirm the login on your device, and you will have 1 minute to do it before getting timed
out.

If you do it on time, N26 will grant an access token, and the tool will persist it to your system keyring.

## Security & GDPR

The tool has `0` tracking and will NEVER track your usage. Therefore, I will never know if you have any issues with it unless you create one in
[the board](https://github.com/nhatthm/n26cli/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc).

The tool does NOT and will NOT share your email and credentials with any other 3rd parties.

The tool will save your access token to your system keyring, and also your credentials (if you choose to save it, see the [Configuration](#Configuration)
and [Authentication](#Authentication) section).

The tool NEVER shows your email and credentials on the screen, in the logs or verbose / debug output.

## Development

The tool uses these libraries for working with N26 APIs and Authentication.

- https://github.com/nhatthm/n26api
- https://github.com/nhatthm/n26keychain
- https://github.com/nhatthm/n26prompt

## Tests

The tool is tested by unit tests and integration tests with mocked API server (
see [`nhatthm/n26api/testkit`](https://github.com/nhatthm/n26api#integration-test))

## Donation

If you like this tool, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnhatthm%2Fn26cli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnhatthm%2Fn26cli?ref=badge_large)
# OpenSAPS

OpenSAPS stands for "Open Slack APi Server". This is an open-source implementation of Slack API server that can be used to integrate applications into each other using Slack API.

Initially this project was created for integrating Gitlab and Gitea into Matrix, because there was no good incoming webhooks support. But it can be used for anything that provides Slack Webhooks support.

Join [#opensaps:pztrn.name](https://matrix.to/#/#opensaps:pztrn.name) Matrix room for help and chat!

## Installation

```
go get -u -v -d gitlab.com/pztrn/opensaps
go install -v gitlab.com/pztrn/opensaps
```

Or drop into [tags section](https://gitlab.com/pztrn/opensaps/tags) to grab a precompiled binary!

## Configuration

Take a look at ``opensaps.example.yaml`` for configuration example and into [docs section](/doc/configuration.md) for configuration file fields description.

## Usage

The only parameter OpenSAPS binary accepts is a configuration file path. Do it like:

```bash
opensaps -config /path/to/opensaps.yaml
```

## About hooks and parsers

While configuring a webhook in your application, please, set username exactly same as one of parsers in ``parsers`` directory! Otherwise parser "default" will be used, which will just concatenate text and attachments into one message!

Also note - that nickname will be ignored while sending message to pushers. Nickname under which messages will appear depends on your account's configuration.

## Known to work good software

[There is a list of software](/doc/software_that_works_good.md) that known to work fine with OpenSAPS. Check it out!

----

## IMPORTANT NOTICE

This project isn't affiliated nor developed by Slack itself.

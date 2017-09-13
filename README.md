# OpenSAPS

OpenSAPS stands for "Open Slack APi Server". This is an open-source
implementation of Slack API server that can be used to integrate
applications into each other using Slack API.

Initially this project was created for integrating Gitlab and Gitea
into Matrix, because there was no good incoming webhooks support.

# Installation

```
go get -u -v -d lab.pztrn.name/pztrn/opensaps
go install -v lab.pztrn.name/pztrn/opensaps
```

# Configuration

Take a look at ``opensaps.example.yaml`` for configuration example.
Right now there is no documentation about configuration file, but it
will appear in future.

# Usage

The only parameter OpenSAPS binary accepts is a configuration file
path. Do it like:

```
opensaps -config /path/to/opensaps.yaml
```

There is some documentation available - check out ``doc`` directory!

# About hooks and parsers

While configuring a webhook in your application, please, set username
exactly same as one of parsers in ``parsers`` directory! Otherwise parser
"default" will be used, which will just concatenate text and attachments
into one message!

Also note - that nickname will be ignored while sending message to
pushers.

----

# IMPORTANT NOTICE

This project isn't affiliated nor developed by Slack itself.

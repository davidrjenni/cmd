# taskfs

taskfs - Google tasks as 9P file server

## Documentation

[![GoDoc](https://godoc.org/github.com/davidrjenni/cmd/taskfs?status.svg)](https://godoc.org/github.com/davidrjenni/cmd/taskfs)

## Installation

```
% go get github.com/davidrjenni/cmd/taskfs
```

You need to set the environment variables `TASKFS_CLIENT_ID` and `TASKFS_CLIENT_SECRET` of the OAuth2 configuration.
For further information see: https://developers.google.com/google-apps/tasks/auth.

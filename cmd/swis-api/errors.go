package main

import "errors"

var errMissingSecretOrToken = errors.New("missing ROOT_TOKEN or DUMP_TOKEN env vars")

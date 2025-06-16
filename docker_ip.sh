#!/bin/sh
export HOSTNAME=$(hostname -i)
exec "$@"
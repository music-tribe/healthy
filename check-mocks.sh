#!/bin/bash -e

make mocks && git diff --exit-code && echo $?
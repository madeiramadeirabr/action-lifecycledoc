#!/bin/sh -l

lifecycledoc --outputFormat github-action --titlePrefix $1 $2 >> $GITHUB_OUTPUT
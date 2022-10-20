#!/bin/sh -l

if [ $1 = 'github-action-json' ]; then
    lifecycledoc --outputFormat github-action-json --titlePrefix $2 $3 >> $GITHUB_OUTPUT
else
    lifecycledoc --outputFormat $1 --titlePrefix $2 $3 >> $GITHUB_STEP_SUMMARY
fi
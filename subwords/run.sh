#!/usr/bin/env bash

# Wordlists available from http://wordlist.aspell.net/12dicts/
# Using 2of12inf.txt for results.

set -e

TESTWORD=${3:-stripe}
ARGS=${4:-}

case "$2" in
  test)
    # Print the score for a given word:
    #   run.sh {input} test foo
    go run *.go --testword $TESTWORD -v < ${1}
    ;;
  all)
    # Print top 10000 words, sorted by score:
    #   run.sh {input} all unused -count=10000
    go run *.go < ${1} ${ARGS} | tee ${1}.output.txt
    ;;
  *)
    echo $"Usage: $0 {input dictionary} {test|all} {test word} {extra args}"
    exit 1
esac

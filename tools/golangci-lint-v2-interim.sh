#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# This script is a bridge to use golangci-lint v2 with the v1 --out-format flag.
# It allows JetBrains IDEs to remain usable while we migrate to v2.
main ()
{
    tmpfile="$(mktemp)"
    trap 'rm -f "${tmpfile}"' EXIT

    local -a args=()
    local prev_arg=""

    for arg in "$@"
    do
        if [ "${prev_arg}" = "--out-format" ]
        then
            args+=("${tmpfile}")
            prev_arg=""
            continue
        fi

        if [ "${arg}" = "--out-format" ]
        then
            args+=("--output.json.path")
        else
            args+=("${arg}")
        fi

        prev_arg="${arg}"
    done

    golangci-lint "${args[@]}"
    cat "${tmpfile}"
}

main "$@"

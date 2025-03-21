#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
    echo "must be run from repository root"
    exit 255
fi

source ./scripts/lib.sh

API_ROOT="./api"

function dirs {
    find . -type f -name '*.proto' -exec dirname {} \; | xargs -n1 basename | sort -u
}

function pb_files() {
    local pb_files
    pb_files=$(find . -type f -name "*.proto")
    echo "$pb_files"
}

function gen_for_modules() {
    local go_out="internal/common/genproto"
    if [ -d "$go_out" ]; then
        log_warning "found existing $go_out, cleaning all files under $go_out"
        run rm -rf "$go_out"
    fi

    for dir in $(dirs); do
        local service="${dir:0:${#dir}-2}"
        local pb_file="${service}.proto"

        if [ -d "$go_out/$dir" ]; then
            log_warning "cleaning all files under $go_out/$dir"
            run rm -rf "$go_out/$dir"/*
        else
            run mkdir -p "$go_out/$dir"
        fi
        log_info "generating code for $service"

        run protoc \
            -I="/usr/local/include/" \
            -I="${API_ROOT}" \
            "--go_out=internal/common/genproto" --go_opt=paths=source_relative \
            --go-grpc_opt=require_unimplemented_servers=false \
            "--go-grpc_out=internal/common/genproto" --go-grpc_opt=paths=source_relative \
            "${API_ROOT}/${dir}/$pb_file"
    done
    log_success "protoc done!"
}


echo "directories containing protos to be built:" $(dirs)
echo "found pd_files:" $(pb_files)
gen_for_modules



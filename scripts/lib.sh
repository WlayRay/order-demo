#!/usr/bin/env bash

# 定义颜色变量
COLOR_BOLD="\033[1m"
COLOR_RED="\033[31m"
COLOR_ORANGE="\033[33m"
COLOR_LIGHTCYAN="\033[36m"
COLOR_GREEN="\033[32m"
COLOR_WHITE="\033[37m"
COLOR_NONE="\033[0m"

# 定义 run 函数
run() {
    echo "Running: $*"
    "$@"
    if [ $? -ne 0 ]; then
        echo "Command failed: $*"
        exit 1
    fi
}

function log_error() {
    >&2 echo -n -e "${COLOR_BOLD}${COLOR_RED}"
    >&2 echo "$@"
    >&2 echo -n -e "${COLOR_NONE}"
}

function log_warning() {
    >&2 echo -n -e "${COLOR_ORANGE}"
    >&2 echo "$@"
    >&2 echo -n -e "${COLOR_NONE}"
}

function log_callout() {
    >&2 echo -n -e "${COLOR_LIGHTCYAN}"
    >&2 echo "$@"
    >&2 echo -n -e "${COLOR_NONE}"
}

function log_info() {
    >&2 echo -n -e "${COLOR_WHITE}"
    >&2 echo "$@"
    >&2 echo -n -e "${COLOR_NONE}"
}

function log_success() {
    >&2 echo -n -e "${COLOR_GREEN}"
    >&2 echo "$@"
    >&2 echo -n -e "${COLOR_NONE}"
}

function prepare_dir() {
    local dir="$1"
    if [ -d "dir" ]; then
        log_warning "Directory $dir already exists. Delete all files under it"
        run find "$dir" -mindepth 1 -delete
    else
        log_callout "Directory $dir does not exist. Creating it."
        run mkdir -p "$dir"
    fi
}
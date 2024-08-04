#!/usr/bin/bash

ENV_FILE=".env";

export_envs() {
    for line in "${lines[@]}"; do
        printf "export %s\n" $line;
        export $line;
    done
}

clean_envs() {
    for line in "${lines[@]}"; do
        pair=(${line//=/ })
        printf "unset %s\n" ${pair[0]};
        unset ${pair[0]};
    done
}

read_file() {
    IFS=$'\n' read -d '' -r -a lines < ${ENV_FILE};
}

run_cmd_with_envs() {
    read_file;
    export_envs;
    exec_cmd "$@";
    clean_envs;
}

exec_cmd() {
    local args="$@";
    printf "executing command: %s\n" "$args";
    $args;
}

main() {
    if [ -f ${ENV_FILE} ]; then
        printf "env file found\n";
        run_cmd_with_envs "$@";
    else
        printf "env file not found\n";
        exec_cmd "$@";
    fi
}

main "$@";

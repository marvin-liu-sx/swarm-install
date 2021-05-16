#!/usr/bin/env bash

move_account(){
    dir_num=("1" "2")
    docker-compose down

    for dir in ${dir_num[@]}
    do
       mkdir -p "/dev/docker/bee_bee-$dir/_data"
       sudo cp -rv "/var/lib/docker/volumes/bee_bee-${dir}/_data/keys" "/dev/docker/bee_bee-${dir}/_data"
       sudo cp -rv "/var/lib/docker/volumes/bee_bee-${dir}/_data/statestore" "/dev/docker/bee_bee-${dir}/_data"
       docker volume rm bee_bee-$dir
    done
}
move_account

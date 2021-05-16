#!/usr/bin/env bash

move_account(){
    docker-compose down

    for dir in {1..20}
    do
       echo move data for $dir ...
       mkdir -p "/data/docker/bee_bee-$dir/_data"
       chmod -R 755 "/data/docker/bee_bee-$dir/_data"
       sudo cp -rv "/var/lib/docker/volumes/bee_bee-${dir}/_data/keys" "/data/docker/bee_bee-${dir}/_data"
       sudo cp -rv "/var/lib/docker/volumes/bee_bee-${dir}/_data/statestore" "/data/docker/bee_bee-${dir}/_data"
       docker volume rm bee_bee-$dir
       echo move data done for $dir !!!
    done
}
move_account
#!/usr/bin/env bash

before_install() {
    rm -rf /home/ubuntu/user-service
}

after_install() {
    echo "Pass this step"
}

application_start() {
    cd /home/ubuntu/user-service && /home/ubuntu/user-service/user-service --action=setup-service
}

application_stop() {
    echo "Pass this step"
}

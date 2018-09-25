#!/bin/bash

# cluster22 
ssh-keyscan  172.22.0.{10..13}  > $TRAVIS_BUILD_DIR/known_hosts

# cluster25
ssh-keyscan -p 25 172.25.0.{10..13}  >> $TRAVIS_BUILD_DIR/known_hosts


cat $TRAVIS_BUILD_DIR/known_hosts


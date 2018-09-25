#!/bin/bash

# cluster22 
ssh-keyscan -v 172.22.0.{10..13}  > $TRAVIS_BUILD_DIR/known_hosts

# cluster25
ssh-keyscan -vp 25 172.25.0.{10..13}  >> $TRAVIS_BUILD_DIR/known_hosts


cat $TRAVIS_BUILD_DIR/known_hosts




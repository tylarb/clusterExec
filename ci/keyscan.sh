#!/bin/bash

# cluster22 
ssh-keyscan 172.22.0.{10..13}  >> ./known_hosts

# cluster25
ssh-keyscan 172.25.0.{10..13}  >> ./known_hosts




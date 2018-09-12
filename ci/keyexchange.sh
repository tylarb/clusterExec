#!/bin/bash


ssh-keyscan -H localhost >> ./known_hosts
ssh-keyscan -H localhost -p 25 >> ./known_hosts




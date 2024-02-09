#!/bin/bash
HOME=$(readlink -f ~)
echo $HOME > home.txt

./setup.sh
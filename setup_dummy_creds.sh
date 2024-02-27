#!/bin/bash

DIR="$HOME/.aws"
echo $DIR
if [ ! -d $DIR ]; then
    echo "Create $DIR"
    mkdir $DIR
fi


cat << EOF >> $DIR/credentials
[default]
aws_access_key_id = AKIADUMMY
aws_secret_access_key = DUMMY
EOF

cat << EOF >> $DIR/config
[default]
region = ap-northeast-1
EOF
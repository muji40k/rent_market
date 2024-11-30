#! /bin/bash

out=

for file in $@;
do
    if [ -f $file ]; then
        out+="--env-file $file "
    fi
done

echo $out


#!/bin/bash

if [ $# -eq 0 ]
then
    echo "Usage:"
    echo "$0 version"
    exit 1
fi
sed -i 's/new_ver/'$1'/g' zettaship-git/DEBIAN/control
mkdir -p zettaship-git/etc
mkdir -p zettaship-git/usr/bin
git clone https://github.com/s-kostyaev/zettaship.git
go get github.com/BurntSushi/toml
go get github.com/zazab/zhash
go get github.com/olekukonko/tablewriter
go get github.com/op/go-logging
cd zettaship
git checkout dev
cd zfs
go build
cp zfs ../../zettaship-git/usr/bin/
cp zettaship.toml ../../zettaship-git/etc/
cd ../..
dpkg -b zettaship-git zettaship-git_$1_amd64.deb
sed -i 's/'$1'/new_ver/g' zettaship-git/DEBIAN/control

#!/bin/bash
branch=dev

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
cd zettaship
git checkout $branch
cd zfs
deps=`go list -f '{{join .Deps "\n"}}' |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`
for dep in $deps; do go get $dep; done
go build
cp -f zfs ../../zettaship-git/usr/bin/
cp -f zettaship.toml ../../zettaship-git/etc/
cd ../..
dpkg -b zettaship-git zettaship-git_$1_amd64.deb
sed -i 's/'$1'/new_ver/g' zettaship-git/DEBIAN/control

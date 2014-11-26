# Maintainer:  <s-kostyaev@ngs>
pkgname=zettaship-git
pkgver=0.4
pkgrel=1
pkgdesc="client for using zfs in lxc container"
arch=('i686' 'x86_64')
url="https://github.com/s-kostyaev/zettaship"
license=('unknown')
depends=('git')
makedepends=('go')
backup=('etc/zettaship.toml')
branch='dev'
source=("${pkgname}::git+https://github.com/s-kostyaev/zettaship#branch=${branch}")
md5sums=('SKIP')
noextract=()
build() {
  go get github.com/BurntSushi/toml
  go get github.com/op/go-logging
  go get github.com/olekukonko/tablewriter
  go get github.com/zazab/zhash
  cd ${srcdir}/${pkgname}/zfs
  go build -o zfs
}

package() {
  install -D -m 755 ${srcdir}/${pkgname}/zfs/zfs ${pkgdir}/usr/bin/zfs
  install -D -m 644 ${srcdir}/${pkgname}/zfs/zettaship.toml ${pkgdir}/etc/zettaship.toml
}

#!/bin/bash

# Deploy our built packages to packagecloud.

REPO=dev
[ "$SOURCE_BRANCH" != "master" ] || REPO=big
echo "deploy source branch: $SOURCE_BRANCH"

source .metadata.sh
# deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}.arm64.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}.amd64.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}.armhf.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}.i386.deb
# rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}.arm64.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}.amd64.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}.armhf.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}.i386.rpm

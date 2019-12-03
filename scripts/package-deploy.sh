#!/bin/bash

# Deploy our built packages to packagecloud.

REPO=dev
[ "$SOURCE_BRANCH" != "master" ] || REPO=big
echo "deploy source branch: $SOURCE_BRANCH"

source .metadata.sh
# deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}-${ITERATION}.arm64.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}-${ITERATION}.amd64.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}-${ITERATION}.armhf.deb
package_cloud push golift/${REPO}/debian/stretch release/unifi-poller_${VERSION}-${ITERATION}.i386.deb
# rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}-${ITERATION}.arm64.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}-${ITERATION}.amd64.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}-${ITERATION}.armhf.rpm
package_cloud push golift/${REPO}/el/5 release/unifi-poller-${VERSION}-${ITERATION}.i386.rpm

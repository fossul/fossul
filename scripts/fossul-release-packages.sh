#!/bin/bash

PLUGIN_DIR="/home/ktenzer/plugins"

echo "Packaging configs"
# Tarball configs
cd $GOPATH/src/github.com/fossul/fossul/src/cli/configs
tar cf $GOBIN/default-configs.tar default
if [ $? != 0 ]; then exit 1; fi
gzip -f $GOBIN/default-configs.tar
if [ $? != 0 ]; then exit 1; fi

echo "Packaging services"
# Tarball server service
cd $GOBIN
tar cf server-service.tar server fossul-server-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f server-service.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball app service
cd $GOBIN
tar cf app-service.tar app fossul-app-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f app-service.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball storage service
cd $GOBIN
tar cf storage-service.tar storage fossul-storage-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f storage-service.tar
if [ $? != 0 ]; then exit 1; fi

echo "Packaging plugins"
# Tarball app plugins
cd $PLUGIN_DIR
tar cf plugins-app.tar app
if [ $? != 0 ]; then exit 1; fi
gzip -f plugins-app.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball storage plugins
cd $PLUGIN_DIR
tar cf $PLUGIN_DIR/plugins-storage.tar storage
if [ $? != 0 ]; then exit 1; fi
gzip -f $PLUGIN_DIR/plugins-storage.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball storage plugins
cd $PLUGIN_DIR
tar cf $PLUGIN_DIR/plugins-archive.tar archive
if [ $? != 0 ]; then exit 1; fi
gzip -f $PLUGIN_DIR/plugins-archive.tar
if [ $? != 0 ]; then exit 1; fi

echo "Completed"

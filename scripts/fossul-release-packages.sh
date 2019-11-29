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
tar cf $GOBIN/server-service.tar $GOBIN/server $GOBIN/fossul-server-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f $GOBIN/server-service.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball app service
tar cf $GOBIN/app-service.tar $GOBIN/app $GOBIN/fossul-app-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f $GOBIN/app-service.tar
if [ $? != 0 ]; then exit 1; fi

# Tarball storage service
tar cf $GOBIN/storage-service.tar $GOBIN/storage $GOBIN/fossul-storage-startup.sh
if [ $? != 0 ]; then exit 1; fi
gzip -f $GOBIN/storage-service.tar
if [ $? != 0 ]; then exit 1; fi

echo "Packaging plugins"
# Tarball app plugins
cd $PLUGIN_DIR
tar cf $PLUGIN_DIR/plugins-app.tar app
if [ $? != 0 ]; then exit 1; fi
gzip -f $PLUGIN_DIR/plugins-app.tar
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

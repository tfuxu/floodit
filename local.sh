#!/usr/bin/bash

read -p "Do you want to perform a clean build? [N/y] " answer

if [[ "$answer" == "y" ]]; then
    rm -r builddir
fi

meson setup builddir
meson configure builddir -Dprefix="$(pwd)/builddir" -Dbuildtype=debug

ninja -C builddir install
meson devenv -C builddir ./bin/floodit

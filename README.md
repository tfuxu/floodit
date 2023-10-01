# Flood It!
This repository showcases an example of compiling Gotk4 projects using Meson build system.

## How to build:

### Flatpak Builder:

#### Prerequisites:

- Flatpak Builder `flatpak-builder`
- GNOME SDK runtime `org.gnome.Sdk//45`
- GNOME Platform runtime `org.gnome.Platform//45`

Install required runtimes:
```shell
flatpak install org.gnome.Sdk//45 org.gnome.Platform//45
```

#### Building Instructions:

##### User installation
```shell
git clone https://github.com/tfuxu/flood_it.git
cd flood_it
flatpak-builder --install --user --force-clean repo/ build-aux/flatpak/io.github.tfuxu.flood_it.json
```

##### System installation
```shell
git clone https://github.com/tfuxu/flood_it.git
cd flood_it
flatpak-builder --install --system --force-clean repo/ build-aux/flatpak/io.github.tfuxu.flood_it.json
```

### Meson Build System:

#### Prerequisites:

The following packages are required to build this project:

- Golang >= 1.18 `go`
- Gtk4 `gtk4`
- Meson `meson`
- Ninja `ninja-build`

#### Building Instructions:

##### Global installation

```shell
git clone https://github.com/tfuxu/flood_it.git
cd flood_it
meson setup builddir
meson configure builddir -Dprefix=/usr/local
ninja -C builddir install
```

##### Local build (for testing and development purposes)

```shell
git clone https://github.com/tfuxu/flood_it.git
cd flood_it
meson setup builddir
meson configure builddir -Dprefix="$(pwd)/builddir" -Dbuildtype=debug
ninja -C builddir install
meson devenv -C builddir ./bin/flood_it
```

> **Note** 
> During testing and development, as a convenience, you can use the [`local_run.sh`](./local_run.sh) script to quickly rebuild local builds.

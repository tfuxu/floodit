project('flood_it',
          version: '1.0.0',
    meson_version: '>= 0.64.0',
  default_options: [ 'warning_level=2',
                     'werror=false',
                   ],
)

# Import modules
gnome = import('gnome')
i18n = import('i18n')

# Project information
PROJECT_URL = 'https://github.com/tfuxu/flood_it'
BUGTRACKER_URL = 'https://github.com/tfuxu/flood_it/issues'
HELP_URL = 'https://github.com/tfuxu/flood_it/discussions'
#TRANSLATE_URL = 'https://hosted.weblate.org/projects/tfuxu/flood_it'

# Constants
PROJECT_RDNN_NAME = 'io.github.tfuxu.flood_it'
ROOT_PATH = '/io/github/tfuxu/flood_it'
PKGDATA_DIR = join_paths(get_option('prefix'), get_option('datadir'), meson.project_name())
MODULE_DIR = join_paths(PKGDATA_DIR, 'flood_it')

git_bin = find_program('git', required: false)

# Set APPLICATION_ID and VERSION_SUFFIX
if get_option('buildtype') == 'debug'
  if git_bin.found()
    VCS_TAG = run_command('git', 'rev-parse', '--short', 'HEAD', check: true).stdout().strip()
  else
    VCS_TAG = ''
  endif
  if VCS_TAG == ''
    VERSION_SUFFIX = '-devel'
  else
    VERSION_SUFFIX = '-@0@'.format(VCS_TAG)
  endif
  APPLICATION_ID = '@0@.Devel'.format(PROJECT_RDNN_NAME)
else
  VERSION_SUFFIX = ''
  APPLICATION_ID = PROJECT_RDNN_NAME
endif

# Required dependencies
dependency('glib-2.0')
dependency('gtk4', version: '>= 4.10.0')

go_root = run_command('printenv', 'GOROOT', check: false)

# Find Go installation
if go_root.returncode() != 0
  go_bin = find_program('go')
else
  go_bin = find_program(join_paths(go_root.stdout().strip(), 'bin', 'go'))
endif

module_sources = files(
  'go.mod',
  'go.sum'
)

# Install configuration data
conf = configuration_data()
conf.set('APP_ID', APPLICATION_ID)
conf.set('ROOT_PATH', ROOT_PATH)
conf.set('PKGDATA_DIR', PKGDATA_DIR)
conf.set('DATA_DIR', join_paths(get_option('prefix'), get_option('datadir')))
conf.set('LOCALE_DIR', join_paths(get_option('prefix'), get_option('localedir')))
conf.set('VERSION', meson.project_version() + VERSION_SUFFIX)
conf.set('BUILD_TYPE', get_option('buildtype'))
conf.set('SCHEMAS_DIR', PKGDATA_DIR)
conf.set('SOURCE_DIR', meson.current_source_dir())
conf.set('BUILD_DIR', meson.current_build_dir())

# Subdirs
subdir('data')
subdir('po')
subdir('src')

flood_it_build_path = join_paths(meson.current_source_dir(), 'src')
#sources = [module_sources, go_sources]

# Compile project
go_build = custom_target(
  'go-build',
  build_by_default: true,
  build_always_stale: true,
  #input: module_sources,
  output: meson.project_name(),
  console: true,
  install: true,
  install_dir: get_option('bindir'),
  command: [go_bin, 'build', '-v', '-o', '@OUTPUT@', flood_it_build_path]
)

# Execute post-installation GTK/GNOME scripts
gnome.post_install(
  glib_compile_schemas: true,
  gtk_update_icon_cache: true,
  update_desktop_database: true,
)

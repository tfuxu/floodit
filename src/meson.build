# Generate constants module
configure_file(
  input: 'constants/constants.go.in',
  output: 'constants.go',
  configuration: configuration_data({
    'APP_ID': APPLICATION_ID,
    'ROOT_PATH': ROOT_PATH,
    'RELEASE_VER': meson.project_version(),
    'VERSION': meson.project_version() + VERSION_SUFFIX,
    'BUILD_TYPE': get_option('buildtype'),
    'PROJECT_URL': PROJECT_URL,
    'BUGTRACKER_URL': BUGTRACKER_URL,
    'HELP_URL': HELP_URL,
    #'TRANSLATE_URL': TRANSLATE_URL,
    'DATA_DIR': conf.get('DATA_DIR'),
    'PKGDATA_DIR': PKGDATA_DIR,
    'LOCALE_DIR': conf.get('LOCALE_DIR'),
  }),
)

# Copy the constants.go output to the source directory
run_command(
  'cp',
  join_paths(meson.project_build_root(), 'src', 'constants.go'),
  join_paths(meson.project_source_root(), 'src', 'constants', 'constants.go'),
  check: true
)

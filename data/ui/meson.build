blueprints = custom_target('blueprints',
  input: files(
    'app_preferences_dialog.blp',
    'color_keyboard.blp',
    'custom_mode_page.blp',
    'difficulty_page.blp',
    'game_page.blp',
    'game_rules_dialog.blp',
    'main_window.blp',
    'result_page.blp',
    'shortcuts_dialog.blp',
    'starting_view.blp'
  ),
  output: '.',
  command: [find_program('blueprint-compiler'), 'batch-compile', '@OUTPUT@', '@CURRENT_SOURCE_DIR@', '@INPUT@']
)

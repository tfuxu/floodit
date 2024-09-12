blueprints = custom_target('blueprints',
  input: files(
    'color_keyboard.blp',
    'custom_mode_page.blp',
    'difficulty_page.blp',
    'game_page.blp',
    'game_rules_page.blp',
    'help_overlay.blp',
    'main_window.blp',
    'result_page.blp',
    'starting_view.blp'
    #'preferences_window.blp',
  ),
  output: '.',
  command: [find_program('blueprint-compiler'), 'batch-compile', '@OUTPUT@', '@CURRENT_SOURCE_DIR@', '@INPUT@']
)

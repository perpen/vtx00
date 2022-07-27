MECHANISM
=========

- reads file with escape sequences

IDEAS
=====

- http://doc.cat-v.org/bell_labs/squeak/
- On editor using the windows: https://youtu.be/hB05UFqOtFA?t=2549
- acme-style shell: some control sequences automatically move output to a new window,
  eg top.
  - Any switch to alt-screen creates new window.
  - Doing "exec" replaces the current shell, should be detected by vman so the current
    window is used.
- Some processes are known to handle some actions, eg if process is kak, pass the
  mouse events to it.
- Could be used for scrollback?
  - Navigation by command entered
    - Optionally deal with interpreters like python or gdb?
  - Folding
- Could selector only intercept print and LF? (except on alt screen)
- Support mouse, to support acme style
- Selected things shown in a section at the top. Option
- Clear selections
- Toggle between showing 1. all buffer, 2. matching lines, 3. matching tokens
- Space
- Selection via numbers/letters, a la surfingkeys? Typing again deselects. Or M-number deselects.
- All configuration via command line switches, and config file
- Screen resize. Reuse vterm code?
- For paths, read cwd of process to find files in current dir?
- Actions: copy to clipboard; stdout; open in Chrome or whatever default program.
- Fzf replacement?

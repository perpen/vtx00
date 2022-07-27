TODO
====

- Use interface from vterm

Operations
==========

DL:
- scr.ClearFromCoords(vpair(0, cursor.Y), count*size.X)
- scr.Shift(vpair(0, cursor.Y+count), linesToMove*size.X, -(count)*size.X)
ED:
- scr.ClearBetween(...) - delete from/to cursor on multiple lines
EL:
- scr.ClearBetween(...) - delete from/to cursor on 1 line
IL:
- scr.Shift(vpair(0, cursor.Y), e*size.X, n*size.X)
LF:
- term.Screen.Shift([0, ?], n, -size.X)
RI:
- scr.Shift([0, Y], n, size.X)

Resulting operations =>
- clear lines
- shift lines
- delete on line to/from cursor

Misc:
- Resize()
- CellAt(x, y)
- CursorTo(x, y)
- Print(string, Style)

- search for multiline regex
- fold = hide lines
- overlay lines

Overlays
========

- Folding of commands:
  - Previous commands are folded, the line(s) with the typed command in some
    style, and the output hidden. Show exit status on the summary line, maybe
    just as a color.
  - When pressing M-p/M-n, unfold previous/next command
  - Incremental search of command or output with M-r
  - Shortcuts to copy the command or output to clipboard
- Search for tokens, eg URLs:
  - Only show lines with URLs
  - URLs are styled
  - Optionally show the command summaries
- Scrolling:
  - Page up/down to scroll
  - M-u switches to URL selections on the viewport

Either:
1. Keep screen intact, create another one
2. Or edit the current screen, with line types

Operations on overlay:
- Fold/unfold
- Change style
- Insert lines
- Hide lines
- Clear up when underlying cleared
- Cache?

type Line struct {
    s *string
    styles []Style
    OR
    cells []Cell
    type LineType
}

type BySection struct {
    sections []Section
}

type Section struct {
    struct {
        cmd string //
        output string
    }
}

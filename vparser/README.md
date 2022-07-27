FIXME
=====
- Use 1 as coords base instead of 0?
- To test things like BEL: Store in state the number of BEL events.
- Supplement to ECMA-48 https://vt100.net/emu/ctrlseq_dec.html
- Great info in the hterm/doc/ControlSequences.md
- http://www.real-world-systems.com/docs/ANSIcode.html
- parser.ignoreFlagged - use? notify listeners when stuff has been ignored so they can log it?
- The state machine only recognises ascii 7 bit and utf-8, what about 8bit ascii?
- Charsets: russian, greek, hebrew, turkish mentioned on
  https://www.vt100.net/docs/vt510-rm/chapter7.html
- FIXME - 7/8 bit
  The PDF lists clashing key sequences for SI/SO and LS0/LS1, we ignore the former as they are
  for 7-bit environments only.

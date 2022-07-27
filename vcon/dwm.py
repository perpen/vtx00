#!/usr/bin/env python3

# - Compute borders in controller, for more flexibility.
#   If more performance needed, have something like a queue to limit recompute
#   on resize.
# - To preserve positions when resizing a lot, preserve original position
#   and have a scaling param.
# - To get rid of copy() calls, have a property always returning a copy.
# - Unify Api and Server, or Server and Controller
# - Could the controller methods return updates list instead of pushing?
#   Would simplify tests a bit.
# - Move tests into this file
# - Use SimpleQueue?

import os
import socket
import sys
import json
import threading
import time
import logging
from queue import Queue
from pprint import pprint


class Pair:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def __str__(self):
        return "[%i, %i]" % (self.x, self.y)


class Panel:
    def __init__(self, panel_id):
        self.panel_id = panel_id
        self.mapped = True
        self.top = None
        self.size = None
        self.border = None
        self.meta = None
        self.z = 0

    def unserial(serial):
        panel_id = serial["ID"]
        p = Panel(panel_id)
        p.panel_id = panel_id
        # p.mapped = serial["Mapped"]
        p.z = serial["Z"]
        p.meta = json.loads(serial["Meta"] or "{}")
        position = serial["Pos"]
        if position:
            p.top = Pair(position[0], position[1])
            p.size = Pair(position[2], position[3])
        return p

    def serial(self):
        serial = {
            "ID": self.panel_id,
            "Mapped": self.mapped,
            "Meta": json.dumps(self.meta or {}),
            "Z": self.z,
        }
        if self.top and self.size:
            serial["Pos"] = [
                self.top.x, self.top.y,
                self.size.x, self.size.y
            ]
        return serial


class Event:
    def __init__(self):
        pass

    def unseria(serial):
        evt = Event()
        evt.type = serial["Type"]
        evt.target = serial["Target"]
        evt.details = serial["Details"]
        return evt


class Api:
    def __init__(self, api_sock, notif_sock):
        self.api_sock = api_sock
        self.notif_sock = notif_sock
        self.events_queue = Queue()
        thread = threading.Thread(target=self.listen, args=())
        thread.daemon = True
        thread.start()

    def listen(self):
        log.info('listen')
        if os.path.exists(self.notif_sock):
            os.remove(self.notif_sock)
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        sock.bind(self.notif_sock)
        sock.listen(1)
        while True:
            log.info('accept')
            conn, client_address = sock.accept()
            try:
                log.info('conn')
                while True:
                    data_json = conn.recv(9999)
                    if data_json:
                        data = data_json.decode('utf-8')
                        log.info('data: %s', data)
                        evt = Event.unseria(json.loads(data))
                        self.events_queue.put_nowait(evt)
                    else:
                        log.info('data end')
                        break
                # conn.sendall('got it'.encode('utf-8'))
            finally:
                conn.close()

    def send(self, req):
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        try:
            sock.connect(self.api_sock)
            msg = json.dumps(req)
            log.info('sending:')
            log.info('sending: %s', msg)
            print(msg, file=sys.stderr)
            pprint(msg)
            sock.sendall(msg.encode('utf-8'))
            resp_json = sock.recv(9999)
            resp = json.loads(resp_json.decode('utf-8'))
            log.info('response:')
            pprint(resp)
            return resp
        except socket.error as e:
            log.info('socket error: %s', e)
            sys.exit(1)

    def state(self):
        req = {
            "GetStateCmd": {}
        }
        resp = self.send(req)
        return resp

    def create_panel(self, argv, meta=None):
        req = {
            "CreatePanelCmd": {
                "Argv": argv,
                "Cwd": ".",
                "Meta": json.dumps(meta)
            }
        }
        return self.send(req)


class Server:
    def __init__(self, api):
        self.api = api

    def state(self):
        state_serial = self.api.state()
        if not state_serial:
            return {}
        focus_id = state_serial["FocusID"]
        if focus_id < 0:
            focus_id = None
        log.info('Server.state: focus_id=%s', focus_id)
        size = state_serial["Size"]
        return {
            "focus_id": focus_id,
            "size": Pair(size[0], size[1]),
            "panels": [
                Panel.unserial(serial) for serial in state_serial["Panels"]
            ]
        }

    def create_panel(self, argv, meta=None):
        ''' Returns panel ID '''
        return Panel.unserial(self.api.create_panel(argv, meta=meta))

    def push_updates(self, focus_id, panels):
        cmd = {
            'LayoutCmd': {
                'FocusID': focus_id,
                'Panels': [panel.serial() for panel in panels],
            }
        }
        self.api.send(cmd)


class Controller:
    def __init__(self, srv):
        self.srv = srv
        self.focus_id = None
        self.size = None
        self.panels_by_id = {}
        self.panel_ids = []
        self.layout = TileLayout()

    @property
    def focus(self):
        if self.focus_id is None:
            log.info('no focus')
            return None
        return self.panels_by_id[self.focus_id]

    def event_loop(self):
        self.refresh_state(None)
        while True:
            log.info('event loop')
            evt = self.srv.api.events_queue.get()
            self.handle(evt)
            self.srv.api.events_queue.task_done()

    def handle(self, evt):
        evt_type = evt.type
        log.info('handle: evt: %s', evt)
        method = {
            'death': self.handle_death,
            'key': self.handle_key,
            'resize': self.handle_resize,
        }.get(evt_type)
        if not method:
            log.info('no handler for event type: %s', evt_type)
            return
        method(evt)

    def handle_death(self, evt):
        del self.panels_by_id[evt.target]
        self.focus_id = list(self.panels_by_id.values())[-1].panel_id
        self.handle_resize(None)
        # self.refresh_state(None)

    def handle_key(self, evt):
        key = str(chr(evt.details))
        log.info('handle_key: %s', key)
        handler = {
            'x': self.exit,
            'r': self.refresh_state,
            'c': self.new_panel,
            ',': self.layout.master_inc,
            '.': self.layout.master_dec,
            'p': self.popup_test,
        }.get(key)
        if not handler:
            log.error("handle_key: no binding for %s", key)
            return
        handler(evt)
        self.handle_resize(None)

    def exit(self, evt):
        log.info('exiting')
        sys.exit(0)

    def handle_resize(self, evt):
        if evt:
            w, h = evt.details
            self.size = Pair(w, h)
        updates = self.layout.layout(self.size, list(self.panels_by_id.values()))
        self.srv.push_updates(self.focus_id, updates)

    def refresh_state(self, evt):
        state = self.srv.state()
        self.panels_by_id = {}
        for panel in state["panels"]:
            self.panels_by_id[panel.panel_id] = panel
        self.focus_id = state["focus_id"]
        self.size = state["size"]

        if len(self.panels_by_id) == 0:
            log.info("refresh_state: creating initial panel and bar")
            # self.new_bar(None)
            self.new_panel(None)

    def new_panel(self, evt, argv=["bash"]):
        p = self.srv.create_panel(argv)
        self.focus_id = p.panel_id
        self.panels_by_id[p.panel_id] = p
        updates = self.layout.layout(self.size, list(self.panels_by_id.values()))
        self.srv.push_updates(p.panel_id, updates)

    def new_bar(self, evt):
        argv = ["../vcon/dwm-bar.py"]
        p = self.srv.create_panel(argv, meta={"bar": True})
        self.panels_by_id[p.panel_id] = p
        self.bar_id = p.panel_id
        updates = self.layout.layout(self.size, list(self.panels_by_id.values()))
        self.srv.push_updates(p.panel_id, updates)

    def popup_test(self, evt):
        argv = ["bash"]
        p = self.srv.create_panel(argv)
        p.z = 1
        self.panels_by_id[p.panel_id] = p
        self.focus_id = p.panel_id
        p.top = Pair(5, 5)
        p.size = Pair(50, 15)
        self.srv.push_updates(p.panel_id, [p])


class TileLayout:
    def __init__(self):
        self.nmaster = 1
        self.mfact = .5

    def _panels(self, all_panels):
        panels = []
        popups = []
        bar = None
        for p in all_panels:
            if p.meta and p.meta.get("bar", False):
                bar = p
            elif p.z > 0:
                popups.append(p)
            else:
                panels.append(p)
        return bar, panels, popups

    def layout(self, size, all_panels):
        print('all_panels:', all_panels)
        bar, panels, popups = self._panels(all_panels)
        if bar:
            bar_height = 2
            bar.top = Pair(0, size.y - bar_height)
            bar.size = Pair(size.x, bar_height)
            size = Pair(size.x, size.y - bar_height)
        n = len(panels)
        log.info("tile_layout: n=%i", n)
        if n == 0:
            return []
        nmaster = self.nmaster
        if n > nmaster:
            master_width = round(size.x * self.mfact)
            # master_width, slave_width = partition(size.x, [self.mfact])
            # slave_heights = partition(size.y, slave_factors(panels[nmaster:]))
        else:
            master_width = size.x
        master_top = 0
        slave_top = 0
        num_slaves = n - nmaster
        for i, p in enumerate(panels):
            if i < nmaster:
                h = size.y // min(n, nmaster)
                p.top = Pair(0, master_top)
                if i < nmaster - 1:
                    master_top += h
                else:
                    h = size.y - master_top
                p.size = Pair(master_width, h)
            else:
                if i < n - 1:
                    h = size.y // num_slaves
                else:
                    h = size.y - slave_top
                p.top = Pair(master_width, slave_top)
                p.size = Pair(size.x - master_width, h)
                slave_top += h
        all_panels.extend(popups)
        return all_panels

    def master_inc(self, evt=None):
        self.nmaster += 1

    def master_dec(self, evt=None):
        self.nmaster -= 1


log = logging.getLogger('vcon')
if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    log.info('################################################### starting')
    _, api_sock, notif_sock = sys.argv
    api = Api(api_sock, notif_sock)
    threading.Thread(target=api.listen)
    # api.listen()
    time.sleep(1)
    srv = Server(api)
    ctr = Controller(srv)
    ctr.event_loop()

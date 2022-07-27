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
from copy import copy


class Pair:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def shifted(self, dx, dy):
        return Pair(self.x + dx, self.y + dy)

    def __str__(self):
        return "[%i, %i]" % (self.x, self.y)


class Panel:
    def __init__(self, panel_id):
        self.panel_id = panel_id
        self.mapped = True
        self.top = None
        self.size = None
        self.border = None

    def unserial(serial):
        panel_id = serial["ID"]
        p = Panel(panel_id)
        p.panel_id = panel_id
        # p.mapped = serial["mapped"]
        # p.meta = serial["meta"]
        position = serial["Pos"]
        if position:
            p.top = Pair(position[0], position[1])
            p.size = Pair(position[2], position[3])
        return p

    def serial(self):
        serial = {
            "ID": self.panel_id,
            "Mapped": self.mapped,
            # "Meta": self.meta,
        }
        if self.top and self.size:
            serial["Pos"] = [
                self.top.x, self.top.y,
                self.size.x, self.size.y
            ]
        return serial


class Border:
    def __init__(self):
        pass

    def unserial(serial):
        b = Border()


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
                        evt = json.loads(data)
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
            log.info('sending: %s', msg)
            sock.sendall(msg.encode('utf-8'))
            resp_json = sock.recv(9999)
            resp = json.loads(resp_json.decode('utf-8'))
            log.info('response: %s', resp)
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

    def create_panel(self, argv):
        req = {
            "CreatePanelCmd": {
                "Argv": argv,
                "Cwd": ".",
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

    def create_panel(self, argv):
        ''' Returns panel ID '''
        return Panel.unserial(self.api.create_panel(argv))

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
        evt_type = evt['EvtType']
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
        pass

    def handle_key(self, evt):
        key = str(chr(evt['EvtDetails']))
        log.info('handle_key: %s', key)
        handler = {
            'x': self.exit,
            'r': self.refresh_state,
            'c': self.split_vertical,
            '"': self.split_vertical,
            '%': self.split_horizontal,
        }.get(key)
        if not handler:
            log.error("handle_key: no binding for %s", key)
            return
        handler(evt)

    def exit(self, evt):
        log.info('exiting')
        sys.exit(0)

    def handle_resize(self, evt):
        evt_details = evt["EvtDetails"]
        old_size = self.size
        new_size = Pair(evt_details[0], evt_details[1])
        rx = float(new_size.x) / old_size.x
        ry = float(new_size.y) / old_size.y
        log.info("handle_resize: %s -> %s", old_size, new_size)
        log.info('rx, ry: %f, %f', rx, ry)
        for p in self.panels_by_id.values():
            log.info(p)
            old_top = copy(p.top)
            p.top.x = round(old_top.x * rx)
            p.top.y = round(old_top.y * ry)
            p.size.x = round(p.size.x * rx)
            p.size.y = round(p.size.y * ry)
        self.size = new_size
        self.srv.push_updates(self.focus_id, self.panels_by_id.values())

    def refresh_state(self, evt):
        state = self.srv.state()
        self.panels_by_id = {}
        for panel in state["panels"]:
            self.panels_by_id[panel.panel_id] = panel
        self.focus_id = state["focus_id"]
        self.size = state["size"]

        if len(self.panels_by_id) == 0:
            log.info("refresh_state: creating an initial panel")
            panel = self.srv.create_panel(["bash"])
            log.info('refresh_state: created panel %s', panel.panel_id)
            panel.top = Pair(0, 0)
            panel.size = copy(self.size)
            self.focus_id = panel.panel_id
            self.panels_by_id[panel.panel_id] = panel
            self.srv.push_updates(panel.panel_id, [panel])
        log.info("refresh_state: focus_id=%s", self.focus_id)

    def split_horizontal(self, evt, argv=["bash"]):
        target_id = self.focus_id
        p1 = self.panels_by_id[target_id]
        p1 = self.focus
        p2 = self.srv.create_panel(argv)
        orig_width = p1.size.x
        p1.size.x = p1.size.x // 2
        p2.top = p1.top.shifted(p1.size.x, 0)
        p2.size = Pair(orig_width - p1.size.x, p1.size.y)
        self.focus_id = p2.panel_id
        self.panels_by_id[p2.panel_id] = p2
        self.srv.push_updates(p2.panel_id, [p1, p2])

    def split_vertical(self, evt, argv=["bash"]):
        log.info('split_vertical: panels_by_id=%s', self.panels_by_id)
        target_id = self.focus_id
        p1 = self.panels_by_id[target_id]
        p1 = self.focus
        p2 = self.srv.create_panel(argv)
        log.info("split_vertical: p1=%s", p1)
        log.info("split_vertical: p1.size=%s", p1.size)
        orig_height = p1.size.y
        p1.size.y = p1.size.y // 2
        p2.top = p1.top.shifted(0, p1.size.y)
        p2.size = Pair(p1.size.x, orig_height - p1.size.y)
        self.focus_id = p2.panel_id
        self.panels_by_id[p2.panel_id] = p2
        self.srv.push_updates(p2.panel_id, [p1, p2])


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

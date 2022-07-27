#!/usr/bin/env python3

import unittest
from unittest.mock import MagicMock
from simple import Controller, Server


class TestController(unittest.TestCase):

    def setUp(self):
        self.mock_api = MagicMock()
        self.p1 = self.panel_serial(1, [0, 0, 4, 5])
        self.p2 = self.panel_serial(2, [4, 0, 4, 5])
        self.p3 = self.panel_serial(3)
        self.mock_api.state = MagicMock(return_value={
            "FocusID": 1,
            "Size": [10, 10],
            "Panels": [self.p1, self.p2],
        })
        self.mock_api.create_panel = MagicMock(return_value=self.p3)

        self.srv = Server(self.mock_api)
        self.ctr = Controller(self.srv)
        self.ctr.focus_id = 1
        self.ctr.refresh_state(None)

    def test_handle_resize(self):
        '''
           [0, 0, 10,  5], [0,  5, 10,  5]
        => [0, 0, 10, 10], [0, 10, 10, 10]

        '''
        self.ctr.handle_resize({"EvtDetails": [20, 20]})

        self.mock_api.send.assert_called_with({
            'LayoutCmd': {
                'FocusID': 1,
                'Panels': [
                    self.panel_serial(1, [0, 0, 8, 10]),
                    self.panel_serial(2, [8, 0, 8, 10]),
                ]
            }
        })

    def test_split_horizontal(self):
        self.ctr.focus_id = 1
        self.ctr.split_horizontal({"EvtTarget": 1}, ["bash"])

        self.mock_api.send.assert_called_with({
            'LayoutCmd': {
                'FocusID': 3,
                'Panels': [
                    self.panel_serial(1, [0, 0, 2, 5]),
                    self.panel_serial(3, [2, 0, 2, 5]),
                ]
            }
        })

    def test_split_vertical(self):
        self.ctr.split_vertical({"EvtTarget": 1}, ["bash"])

        self.mock_api.send.assert_called_with({
            'LayoutCmd': {
                'FocusID': 3,
                'Panels': [
                    self.panel_serial(1, [0, 0, 4, 2]),
                    self.panel_serial(3, [0, 2, 4, 3]),
                ]
            }
        })

    def panel_serial(self, panel_id, position=[], focused=False):
        serial = {
            "ID": panel_id,
            "Mapped": True,
            # "Meta": None,
            "Pos": position,
        }
        return serial


if __name__ == '__main__':
    unittest.main()

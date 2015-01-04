#!/usr/bin/env python

import socket
import select
from utils import *
import json
from tun import TunDev
import logging
import os


class TunClient(object):

    def __init__(self, serverpwd, host, port, user, passwd):
        self.server_addr = (host, port)
        self.user = (user, passwd)
        self.enc = Enc(serverpwd)

    def run(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.sock.connect(self.server_addr)
        self.login()
        self.setup()
        self.proxy()

    def login(self):
        obj = {'user': self.user[0], 'passwd': self.user[1]}
        self.sock.send(self.enc.encrypt(
            PKT_LOGIN, json.dumps(obj).encode('utf8')))
        raw_data = self.sock.recv(65535)
        pkt_type, data = self.enc.decrypt(raw_data)
        if data is None:
            raise Exception('invalid data: %r' % raw_data)
        rep = json.loads(data.decode('utf8'))
        if rep['ok']:
            self.gateway = rep['gateway']
            self.ip = rep['ip']
        else:
            raise Exception('login fail: %r' % rep['msg'])

    def setup(self):
        self.tun = TunDev(self.ip, self.gateway)
        self.poll = select.epoll()
        self.poll.register(
            self.sock.fileno(), select.EPOLLERR | select.EPOLLIN)
        self.poll.register(
            self.tun.fileno(), select.EPOLLERR | select.EPOLLIN)

    def proxy(self):
        while True:
            for fd, evs in self.poll.poll():
                if fd == self.sock.fileno():
                    data = self.sock.recv(65535)
                    pkt_type, data = self.enc.decrypt(data)
                    if pkt_type == PKT_PROXY:
                        self.tun.write(data)
                    elif pkt_type == PKT_SHUTDOWN:
                        self.tun.close()
                        self.sock.close()
                        logging.info('shutdown')
                        return
                else:
                    data = self.tun.read(65535)
                    self.sock.send(self.enc.encrypt(PKT_PROXY, data))


if __name__ == '__main__':
    cli = TunClient('123', '127.0.0.1', 7780, 'lo', '123')
    cli.run()

#!/usr/bin/env python


import socket
import select
import struct
import re
import os
import json
from utils import *
from tun import TunDev
import logging


class TunManager(object):
    r_ipnm = re.compile(
        r'^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})/(\d{1,2})$')

    def __init__(self, selfip, iprange):
        self.gateway = selfip
        self.iprange = iprange
        self.make_ips()

    def make_ips(self):
        r = self.r_ipnm.findall(self.iprange)
        if not r or any([int(x) >= 255 for x in r[0][:4]]):
            raise Exception('invalid ip range: %s' % self.iprange)
        mask = int(r[0][-1])
        if mask < 24 or mask >= 30:
            raise Exception('ip range too large/small: %s' % self.iprange)
        net = struct.unpack('>I', bytes([int(x) for x in r[0][:4]]))[0]
        mask = (1 << (32 - mask)) - 1
        self.ips = set()
        for i in range(1, mask):
            ds = struct.pack('>I', net | i)
            ip = '.'.join(['%d' % x for x in ds])
            self.ips.add(ip)
        if self.gateway in self.ips:
            self.ips.remove(self.gateway)

    def new_tun(self):
        ip = self.ips.pop()
        tun = TunDev(self.gateway, ip)
        return ip, tun


class Client(object):

    def __init__(self, addr):
        self.addr = addr
        self.ip = None
        self.user = None
        self.tun = None


class TunServer(object):

    def __init__(self, passwd, ip, port, users, selfip, iprange):
        self.enc = Enc(passwd)
        self.addr = (ip, port)
        self.users = users
        self.tun_mgr = TunManager(selfip, iprange)
        self.addr2clis = {}
        self.tun2clis = {}

    def run(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.sock.bind(self.addr)
        self.poll = select.epoll()
        self.poll.register(self.sock.fileno(),
                           select.EPOLLIN | select.EPOLLERR)

        while True:
            for fd, evs in self.poll.poll():
                if fd == self.sock.fileno():
                    if evs & select.EPOLLIN:
                        self.on_server()
                    elif evs & select.EPOLLERR:
                        pass
                else:
                    self.on_tuns(fd, evs)

    def on_server(self):
        raw_data, addr = self.sock.recvfrom(65535)
        pkt_type, data = self.enc.decrypt(raw_data)
        if data is None:
            logging.debug('ignored invalid packet: %r' % raw_data)
            return
        if pkt_type == PKT_PROXY:
            cli = self.addr2clis.get('%s:%d' % addr, None)
            if cli is not None:
                cli.tun.write(data)
        elif pkt_type == PKT_LOGIN:
            self.on_new_cli(addr, data)
        elif pkt_type == PKT_SHUTDOWN:
            pass

    def on_new_cli(self, addr, data):
        try:
            data = json.loads(data.decode('utf8'))
        except Exception as e:
            logging.warn('decode new cli config fail: %r' % e)
            return
        username = data.get('user', None)
        passwd = self.users.get(username, None)
        if username is None or passwd != data.get('passwd', None):
            obj = {'ok': False, 'msg': 'invalid user/passwd'}
        else:
            ip, tun = self.tun_mgr.new_tun()
            cli = Client(addr)
            cli.ip = ip
            cli.user = username
            cli.tun = tun
            self.addr2clis['%s:%d' % addr] = cli
            self.tun2clis[tun.fileno()] = cli
            self.poll.register(tun.fileno(), select.EPOLLIN | select.EPOLLERR)
            obj = {'ok': True, 'ip': ip, 'gateway': self.tun_mgr.gateway}
        data = self.enc.encrypt(
            PKT_LOGIN_FAIL, json.dumps(obj).encode('utf8'))
        self.sock.sendto(data, addr)

    def on_tuns(self, fd, evs):
        cli = self.tun2clis[fd]
        if evs & select.EPOLLERR:
            logging.warn('tun on err')
            cli.tun.close()
            del self.tun2clis[fd]
            del self.addr2clis['%s:%d' % cli.addr]
            return
        data = os.read(fd, 2048)
        if data:
            self.sock.sendto(self.enc.encrypt(PKT_PROXY, data), cli.addr)


if __name__ == '__main__':
    l = logging.getLogger()
    l.setLevel(logging.DEBUG)
    ser = TunServer('123', '0.0.0.0', 7780, {'lo': '123'},
                    '192.168.11.1', '192.168.11.0/24')
    ser.run()

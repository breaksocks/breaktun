import fcntl
import struct
import subprocess
import os


TUNSETIFF = 0x400454ca
IFF_TUN = 0x0001
IFF_TAP = 0x0002
IFF_NO_PI = 0x1000


class TunDev(object):

    def __init__(self, selfip, remoteip):
        tun = os.open('/dev/net/tun', os.O_RDWR)
        ifr = struct.pack('16sH', b'', IFF_TUN | IFF_NO_PI)
        ifr = fcntl.ioctl(tun, TUNSETIFF, ifr)
        name = ifr[:ifr.index(b'\x00')].decode()
        subprocess.check_call('ifconfig %s %s pointopoint %s up' % (
            name, selfip, remoteip), shell=True)
        self.name = name
        self.tun = tun

    def fileno(self):
        return self.tun

    def read(self, n):
        return os.read(self.tun, n)

    def write(self, data):
        return os.write(self.tun, data)

    def close(self):
        os.close(self.tun)

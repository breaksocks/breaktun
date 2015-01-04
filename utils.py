from Crypto.Cipher import DES3
import struct
from hashlib import md5

PKT_LOGIN = 0
PKT_LOGIN_FAIL = 1
PKT_PROXY = 2
PKT_SHUTDOWN = 3


class Enc(object):

    def __init__(self, passwd):
        self.passwd = md5(passwd.encode('utf8')).digest()
        self.enc = DES3.new(self.passwd, DES3.MODE_ECB)

    def decrypt(self, data):
        if len(data) % 8 != 0:
            return 0, None
        data = self.enc.decrypt(data)
        if data[0:1] != b'M':
            return 0, None
        size = struct.unpack('>H', data[2:4])[0]
        if size > len(data) - 4:
            return 0, None
        return data[1], data[4:4 + size]

    def encrypt(self, pkt_type, data):
        size = len(data)
        bsize = (size + 4) % 8
        if bsize != 0:
            bsize = 8 - bsize
        data = b'M' + struct.pack('B', pkt_type) + \
            struct.pack('>H', len(data)) + data + b'\x00' * bsize
        return self.enc.encrypt(data)

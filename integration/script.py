#coding=utf8

from httplib2 import Http
from urllib import urlencode

h = Http()

r, c = h.request("http://127.0.0.1:4000/trade.do", "POST",
          "{\"symbol\": \"WSCN\", \"type\": \"sell\", \"amount\": 100, \"price\": 100.01}",
          {"Content-Type": "text/json"})
print c

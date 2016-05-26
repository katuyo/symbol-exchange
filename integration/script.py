#coding=utf8

import random
import logging
import json

from httplib2 import Http
from logging import handlers

LOG_FILE = '../logs/WSCN_client.log'
handler = logging.handlers.RotatingFileHandler(LOG_FILE, maxBytes=1024 * 1024, backupCount=5)  # 实例化handler
handler.setFormatter(logging.Formatter('%(asctime)s - %(filename)s:%(lineno)s - %(name)s - %(message)s'))

logger = logging.getLogger('client')
logger.addHandler(handler)
logger.setLevel(logging.DEBUG)

h = Http()


def send():
    exchange_type = random_type()
    # r, c = h.request("http://127.0.0.1:4000/trade.do", "POST",
    #                  "{\"symbol\": \"WSCN\", \"type\": \"sell\", \"amount\": 10, \"price\": 100.00}", {"Content-Type": "text/json"})
    r, c = h.request("http://127.0.0.1:4000/trade.do", "POST",
                     "{\"symbol\": \"WSCN\", \"type\": \"" + exchange_type + "\", \"amount\": " + random_amount() +
                     ", \"price\": " + random_price() + "}", {"Content-Type": "text/json"})
    if exchange_type == "buy" or exchange_type == "sell":
        obj = json.loads(c)
        logger.info("%s, %s", obj['order_id'], exchange_type)


def random_type():
    return str(random.choice(["buy", "sell", "buy_market", "sell_market"]))


def random_amount():
    return str(random.randrange(1, 100, 1))


def random_price():
    return str(round(random.uniform(90.00, 110.00), 2))

if __name__ == '__main__':
    send()

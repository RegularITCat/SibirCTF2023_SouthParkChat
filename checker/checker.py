#!/usr/bin/python3
import sys
import errno
import requests
import string
import random
import base64
import pickle
import time
import json

TIMEOUT = 5

# put-get flag to service success
def service_up():
    print("[service is worked] - 101")
    exit(101)

# service is available (available tcp connect) but protocol wrong could not put/get flag
def service_corrupt():
    print("[service is corrupt] - 102")
    exit(102)

# waited time (for example: 5 sec) but service did not have time to reply
def service_mumble():
    print("[service is mumble] - 103")
    exit(103)

# service is not available (maybe blocked port or service is down)
def service_down():
    print("[service is down] - 104")
    exit(104)

if len(sys.argv) != 5:
    print("\nUsage:\n\t" + sys.argv[0] + " <host> (put|check) <flag_id> <flag>\n")
    print("Example:\n\t" + sys.argv[0] + " \"127.0.0.1\" put \"abcdifghr\" \"123e4567-e89b-12d3-a456-426655440000\" \n")
    print("\n")
    exit(0)

def debug(err):
    pass
    # if isinstance(err, str):
    #     err = Exception(err)
    # traceback.print_exc()
    # raise err

host = sys.argv[1]
port = 8888
command = sys.argv[2]
f_id = sys.argv[3]
flag = sys.argv[4]

url = "http://" + host + ":" + str(port)
names = pickle.loads(base64.b64decode("gASVngYAAAAAAABdlCh9lCiMCmZpcnN0X25hbWWUjARTdGFulIwLc2Vjb25kX25hbWWUjAVNYXJzaJR1fZQoaAKMBEt5bGWUaASMCkJyb2Zsb3Zza2mUdX2UKGgCjARFcmljlGgEjAdDYXJ0bWFulHV9lChoAowFS2VubnmUaASMCU1jQ29ybWlja5R1fZQoaAKMB0J1dHRlcnOUaASMBlN0b3RjaJR1fZQoaAKMBVJhbmR5lGgEjAVNYXJzaJR1fZQoaAKMB0hlcmJlcnSUaASMCEdhcnJpc29ulHV9lChoAowDTXIulGgEjAZNYWNrZXmUdX2UKGgCjAZHZXJhbGSUaASMCkJyb2Zsb3Zza2mUdX2UKGgCjAZTaGVpbGGUaASMCkJyb2Zsb3Zza2mUdX2UKGgCjAVMaWFuZZRoBIwHQ2FydG1hbpR1fZQoaAKMBUppbW15lGgEjAZWYWxtZXKUdX2UKGgCjAdUb2xraWVulGgEjAVCbGFja5R1fZQoaAKMBVdlbmR5lGgEjAtUZXN0YWJ1cmdlcpR1fZQoaAKMBUNseWRllGgEjAdEb25vdmFulHV9lChoAowFQ3JhaWeUaASMBlR1Y2tlcpR1fZQoaAKMBEJlYmWUaASMB1N0ZXZlbnOUdX2UKGgCjAVIZWlkaZRoBIwGVHVybmVylHV9lChoAowFU2NvdHSUaASMCU1hbGtpbnNvbpR1fZQoaAKMBVRpbW15lGgEjAVCdXJjaJR1fZQoaAKMBVR3ZWVrlGgEjAVUd2Vha5R1fZQoaAKMAlBDlGgEjAlQcmluY2lwYWyUdX2UKGgCjAZTdHJvbmeUaASMBVdvbWFulHV9lChoAowGU2hhcm9ulGgEjAVNYXJzaJR1fZQoaAKMB1NoZWxsZXmUaASMBU1hcnNolHV9lChoAowGTWFydmlulGgEjAVNYXJzaJR1fZQoaAKMBUppbWJvlGgEjARLZXJulHV9lChoAowDSWtllGgEjApCcm9mbG92c2tplHV9lChoAowGU3R1YXJ0lGgEjAlNY0Nvcm1pY2uUdX2UKGgCjAVDYXJvbJRoBIwJTWNDb3JtaWNrlHV9lChoAowFS2FyZW6UaASMCU1jQ29ybWlja5R1fZQoaAKMB1N0ZXBoZW6UaASMBlN0b3RjaJR1fZQoaAKMBUxpbmRhlGgEjAZTdG90Y2iUdX2UKGgCjAdPZmZpY2VylGgEjAhCYXJicmFkeZR1fZQoaAKMBkJpZ0dheZRoBIwCQWyUdX2UKGgCjANOZWSUaASMCkdlcmJsYW5za3mUdX2UKGgCjAVKZXN1c5RoBIwKbm90LUNocmlzdJR1fZQoaAKMBVR1b25nlGgEjAVMdUtpbZR1fZQoaAKMBkZhdGhlcpRoBIwETWF4aZR1fZQoaAKMBE1hcnmUaASMCU1jRGFuaWVsc5R1fZQoaAKMC0RyLkFscGhvbnNllGgEjAhNZXBoZXN0b5R1fZQoaAKMBVNhbnRhlGgEjAVDbGF1c5R1fZQoaAKMA01yLpRoBIwFU2xhdmWUdX2UKGgCjAdUb3dlbGlllGgEjAlXYXNoY2xvdGiUdX2UKGgCjAhIYXJyaXNvbpRoBIwFWWF0ZXOUdX2UKGgCjANQaXCUaASMBlBpcnJ1cJR1fZQoaAKMBUJldHN5lGgEjAdEb25vdmFulHV9lChoAowDTXMulGgEjApDaG9rc29uZGlrlHV9lChoAowDTXMulGgEjAhDcmFidHJlZZR1fZQoaAKMBlNhZGRhbZRoBIwHSHVzc2VpbpR1fZQoaAKMCVByaW5jaXBhbJRoBIwIVmljdG9yaWGUdX2UKGgCjANNci6UaASMBkhhbmtleZR1fZQoaAKMBVNhdGFulGgEjAVEZW1vbpR1fZQoaAKMBUtldmlulGgEjAZTdG9sZXmUdX2UKGgCjAVKYXNvbpRoBIwFV2hpdGWUdX2UKGgCjARNcnMulGgEjAZOZWxzb26UdX2UKGgCjAZEYXJyeWyUaASMCFdlYXRoZXJzlHV9lChoAowHQnJhZGxleZRoBIwGQmlnZ2xllHV9lChoAowER290aJRoBIwES2lkc5R1fZQoaAKMBkRvdWdpZZRoBIwITydDb25uZWyUdX2UKGgCjARtb21zlGgEjAZmdWNrZXKUdX2UKGgCjAluYWdpYmF0b3KUaASMBDEzMzeUdX2UKGgCjARNRUdBlGgEjAZIT1JPU0iUdX2UKGgCjARUUlVFlGgEjAd3YXJyMTBylHVlLg=="))

# will be mumble (2) - for test jury
# while True: time.sleep(10);

def id_generator(size=6, chars=string.ascii_uppercase + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))

def put_flag():
    global host, port, f_id, flag, url, names
    # try put
    login = id_generator(size=32)
    password = id_generator(size=32)
    name = random.choice(names)
    message = ' '.join([id_generator(size=random.randint(4,12)) for _ in range(random.randint(2,5))])
    
    try:
        #TODO
        r = requests.get(url + "/api/v1/health", timeout=TIMEOUT)
        sess = requests.Session()
        r = sess.post(url + "/api/v1/register", json={"login":login,"password":password, "first_name":name["first_name"], "second_name":name["second_name"]}, timeout=TIMEOUT)
        #WTF FUCKING PYTHON WHY
        cookies = requests.utils.cookiejar_from_dict(sess.cookies.get_dict())
        sess.cookies.update(cookies)
        check_random_func = random.randint(1,3)
        check_random_func = 3
        if check_random_func == 1:
            fileobj = sess.post(url + "/api/v1/file", files={"file":("test.file", message, "multipart/form-data")}, timeout=TIMEOUT).json()
            r = sess.get(url + "/api/v1/file/" + str(fileobj["id"]), timeout=TIMEOUT)
            r = sess.get(url + "/api/v1/file/" + str(fileobj["id"]) + "/download", timeout=TIMEOUT)
            if r.text != message:
                service_corrupt()
        elif check_random_func == 2:
            r = sess.get(url + "/api/v1/chat", timeout=TIMEOUT)
            resp = sess.post(url + "/api/v1/chat/0/message", json={"message":message}, timeout=TIMEOUT).json()
            res = sess.get(url + "/api/v1/chat/0/message/" + str(resp), timeout=TIMEOUT).json()
            if res["message"] != message:
                service_corrupt()
        elif check_random_func == 3:
            r = sess.post(url + "/api/v1/post", json={"name":"test", "content":message}, timeout=TIMEOUT).json()
            res = sess.get(url + "/api/v1/post/" + str(r), timeout=TIMEOUT).json()
            if res["content"] != message:
                service_corrupt()
        r = sess.post(url + "/api/v1/card", json={"comment":' '.join([id_generator(size=random.randint(4,12)) for _ in range(random.randint(2,5))])}, timeout=TIMEOUT)
        r = sess.get(url + "/api/v1/card", timeout=TIMEOUT)
        r = r.json()
        card_id = r[0]["id"]
        r = sess.post(url + "/api/v1/card/" + str(r[0]["id"]) + "/transaction", json={"to_card":r[1]["id"],"amount": r[0]["balance"] - 1, "comment":flag}, timeout=TIMEOUT)
        with open("./flags/" + flag + ".json", "w") as f:
            json.dump({"login":login,"password":password,"first_name":name["first_name"],"second_name":name["second_name"], "message":message,"flag":flag,"card_id":card_id,"transaction_id":r.json()},f)
    except requests.exceptions.ConnectionError:
        service_down()
    except requests.exceptions.RequestException as e:
        service_mumble()
    except Exception as e:
        print(e)
        debug(e)
        service_corrupt()

def check_flag():
    global host, port, f_id, flag, url
    # try get
    flag2 = ""
    try:
        with open("./flags/" + flag + ".json", "r") as f:
            data = json.load(f)
        sess = requests.Session()
        r = sess.post(url + "/api/v1/login", json={"login":data["login"],"password":data["password"]}, timeout=TIMEOUT)
        #WTF FUCKING PYTHON WHY
        cookies = requests.utils.cookiejar_from_dict(r.cookies.get_dict())
        sess.cookies.update(cookies)
        r = sess.get(url + "/api/v1/card/" + str(data["card_id"]) + "/transaction/" + str(data["transaction_id"]), timeout=TIMEOUT)
        flag2 = r.json()["comment"]
    except Exception as e:
        debug(e)
        service_corrupt()

    if flag != flag2:
        debug('flag: [' + flag +  '] flag2: [' + str(flag2) + ']')
        service_corrupt()


if command == "put":
    put_flag()
    check_flag()
    service_up()

if command == "check":
    check_flag()
    service_up()

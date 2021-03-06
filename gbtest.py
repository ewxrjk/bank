#!/usr/bin/env python3
import requests
import os
import subprocess
import time
import socket

try:
    os.remove("_test.db")
except FileNotFoundError:
    pass

# Initialize database
subprocess.check_call(["./bank", "-d", "_test.db", "init"])

# Create a user
subprocess.check_call(["./bank", "-d", "_test.db",
                       "user", "add", "-p", "pass1", "fred"])

# List users
cp = subprocess.run(["./bank", "-d", "_test.db", "user", "list"], check=True, stdout=subprocess.PIPE)
users = cp.stdout[:-1].split(b"\n")
assert(users == [b"fred"])

# Change password
subprocess.check_call(["./bank", "-d", "_test.db",
                       "user", "pw", "-p", "pass2", "fred"])

host = "localhost"
port = 9823
address = "%s:%d" % (host, port)
server = subprocess.Popen(
    ["./bank", "-d", "_test.db", "server", "-a", address])
try:
    # Wait for server
    while True:
        try:
            s = socket.create_connection((host, port))
        except ConnectionRefusedError:
            time.sleep(0.1)
            pass
        break

    # Login
    r = requests.post("http://%s/v1/login" %
                      address, json={"User": "fred", "Password": "pass1"})
    assert r.status_code == 403
    r = requests.post("http://%s/v1/login" %
                      address, json={"User": "fred", "Password": "pass2"})
    r.raise_for_status()
    cookies = r.cookies
    token = r.json()["Token"]

    # Users
    r = requests.get("http://%s/v1/user/" % address, cookies=cookies)
    r.raise_for_status()
    assert r.json() == ['fred']
    r = requests.post("http://%s/v1/user/" % address,
                      json={"User": "bob", "Password": "pass3", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.get("http://%s/v1/user/" % address, cookies=cookies)
    r.raise_for_status()
    assert r.json() == ['bob', 'fred']
    r = requests.put("http://%s/v1/user/bob/password" % address,
                     json={"Password": "pass4", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.post("http://%s/v1/user/" % address,
                      json={"User": "joe", "Password": "pass3", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.delete("http://%s/v1/user/joe" % address, cookies=cookies)
    r.raise_for_status()
    r = requests.delete("http://%s/v1/user/joe" % address, cookies=cookies)
    assert r.status_code == 404

    # Enforce login credentials
    r = requests.get("http://%s/v1/account/" % address)
    assert r.status_code == 403
    r = requests.get("http://%s/v1/account/" %
                     address, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "fred", "Token": "whatever"}, cookies=cookies)
    assert r.status_code == 403
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "fred", "Token": token}, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.delete("http://%s/v1/account/fred" %
                        address, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.get("http://%s/v1/user/" % address)
    assert r.status_code == 403
    r = requests.get("http://%s/v1/user/" %
                     address, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.post("http://%s/v1/user/" % address,
                      json={"User": "fred", "Token": "whatever"}, cookies=cookies)
    assert r.status_code == 403
    r = requests.post("http://%s/v1/user/" % address,
                      json={"User": "fred", "Token": token}, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.delete("http://%s/v1/user/joe" %
                        address, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.put("http://%s/v1/user/bob/password" % address,
                     json={"Password": "pass4", "Token": "whatever"}, cookies=cookies)
    assert r.status_code == 403
    r = requests.put("http://%s/v1/user/bob/password" % address,
                     json={"Password": "pass4", "Token": token}, cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.post("http://%s/v1/transaction/" % address,
                      json={"Origin": "house", "Destination": "fred",
                            "Description": "buy a bus", "Amount": 2000, "Token": 'whatever'},
                      cookies=cookies)
    assert r.status_code == 403
    r = requests.post("http://%s/v1/transaction/" % address,
                      json={"Origin": "house", "Destination": "fred",
                            "Description": "buy a bus", "Amount": 2000, "Token": token},
                      cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.post("http://%s/v1/distribute/" % address,
                      json={"Origin": "house", "Destinations": [
                          "bob", "fred"], "Description": "distribution", "Token": 'whatever'},
                      cookies=cookies)
    assert r.status_code == 403
    r = requests.post("http://%s/v1/distribute/" % address,
                      json={"Origin": "house", "Destinations": [
                          "bob", "fred"], "Description": "distribution", "Token": token},
                      cookies={'bank': 'whatever'})
    assert r.status_code == 403
    r = requests.get("http://%s/v1/config" % address)
    assert r.status_code == 403
    r = requests.get("http://%s/v1/config/title" % address)
    assert r.status_code == 403
    r = requests.put("http://%s/v1/config/" % address,
                     json={"Value": "gbtest.py", "Token": 'whatever'},
                     cookies=cookies)
    assert r.status_code == 405
    r = requests.put("http://%s/v1/config/" % address,
                     json={"Value": "gbtest.py", "Token": token},
                     cookies={'bank': 'whatever'})

    # Accounts
    r = requests.get("http://%s/v1/account/" % address, cookies=cookies)
    r.raise_for_status()
    assert r.json() == []
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "fred", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "bob", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "house", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.get("http://%s/v1/account/" % address, cookies=cookies)
    r.raise_for_status()
    accounts = r.json()
    assert accounts == ["bob", "fred", "house"]
    r = requests.post("http://%s/v1/account/" % address,
                      json={"Account": "joe", "Token": token}, cookies=cookies)
    r.raise_for_status()
    r = requests.delete("http://%s/v1/account/joe" % address, cookies=cookies)
    r.raise_for_status()
    r = requests.delete("http://%s/v1/account/joe" % address, cookies=cookies)
    assert r.status_code == 404

    # Transactions
    r = requests.get("http://%s/v1/transaction/?offset=0&limit=10" %
                     address, cookies=cookies)
    r.raise_for_status()
    assert r.json() == []
    r = requests.post("http://%s/v1/transaction/" % address,
                      json={"Origin": "house", "Destination": "fred",
                            "Description": "buy a bus", "Amount": 2000, "Token": token},
                      cookies=cookies)
    r.raise_for_status()
    r = requests.post("http://%s/v1/transaction/" % address,
                      json={"Origin": "house", "Destination": "bob",
                            "Description": "buy a car", "Amount": 3001, "Token": token},
                      cookies=cookies)
    r.raise_for_status()
    r = requests.post("http://%s/v1/distribute/" % address,
                      json={"Origin": "house", "Destinations": [
                          "bob", "fred"], "Description": "distribution", "Token": token},
                      cookies=cookies)
    r.raise_for_status()
    r = requests.get("http://%s/v1/transaction/?offset=0&limit=10" %
                     address, cookies=cookies)
    r.raise_for_status()
    transactions = r.json()
    assert len(transactions) == 4
    assert transactions[3]["ID"] == 1
    assert transactions[3]["User"] == "fred"
    assert transactions[3]["Origin"] == "house"
    assert transactions[3]["Destination"] == "fred"
    assert transactions[3]["Description"] == "buy a bus"
    assert transactions[3]["Amount"] == 2000
    assert transactions[3]["OriginBalanceAfter"] == -2000
    assert transactions[3]["DestinationBalanceAfter"] == 2000
    assert transactions[2]["ID"] == 2
    assert transactions[2]["User"] == "fred"
    assert transactions[2]["Origin"] == "house"
    assert transactions[2]["Destination"] == "bob"
    assert transactions[2]["Description"] == "buy a car"
    assert transactions[2]["Amount"] == 3001
    assert transactions[2]["OriginBalanceAfter"] == -5001
    assert transactions[2]["DestinationBalanceAfter"] == 3001
    assert transactions[1]["ID"] == 3
    assert transactions[1]["User"] == "fred"
    assert transactions[1]["Origin"] == "bob"
    assert transactions[1]["Destination"] == "house"
    assert transactions[1]["Description"] == "distribution"
    assert transactions[1]["Amount"] == 2500
    assert transactions[1]["OriginBalanceAfter"] == 501
    assert transactions[1]["DestinationBalanceAfter"] == -2501
    assert transactions[0]["ID"] == 4
    assert transactions[0]["User"] == "fred"
    assert transactions[0]["Origin"] == "fred"
    assert transactions[0]["Destination"] == "house"
    assert transactions[0]["Description"] == "distribution"
    assert transactions[0]["Amount"] == 2500
    assert transactions[0]["OriginBalanceAfter"] == -500
    assert transactions[0]["DestinationBalanceAfter"] == -1

    # Cannot delete accounts with a balance
    r = requests.delete("http://%s/v1/account/bob" % address, cookies=cookies)
    assert r.status_code == 400

    # Config
    r = requests.get("http://%s/v1/config/title" % address, cookies=cookies)
    assert r.status_code == 200
    assert r.text == "Test Bank"
    r = requests.put("http://%s/v1/config/title" % address,
                     json={"Value": "gbtest.py", "Token": token}, cookies=cookies)
    assert r.status_code == 200
    r = requests.get("http://%s/v1/config/title" % address, cookies=cookies)
    assert r.status_code == 200
    assert r.text == "gbtest.py"

    # HEAD support
    r = requests.head("http://%s/v1/user/" % address, cookies=cookies)
    r.raise_for_status()
    assert r.text == ''

    # Logout
    r = requests.post("http://%s/v1/logout" % address, cookies=cookies)
    assert r.status_code == 200
    r = requests.get("http://%s/v1/account/" % address)
    assert r.status_code == 403

finally:
    server.kill()

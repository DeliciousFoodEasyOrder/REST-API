import urllib
import urllib3
import httplib2
import requests

url = "https://www.sysu-easyorder.top"

# Create a merchant
headers = {'Content-Type':'application/json'}
post_data = '{"name":"zhang", "password":"123456", "email":"930994408@qq.com", "phone":"12345678901", "address":"hh"}'
pPost = requests.post(url+"/merchants", headers=headers, data=post_data)
if pPost.json()['code'] == 201 and pPost.json()['email'] == "930994408@qq.com" and pPost.json()['phone'] == "12345678901":
    print("Create a merchant success")
else:
    print("Create a merchant failed, program shut down")
    exit(0)

# Password Grant
parameters={"grant_type":"password", "username":"12345678901", "password":"123456", "scope":"111"}
pPost = requests.post(url+"/oauth/token", data=parameters)
# print(pPost.json()["code"])
token = pPost.json()["access_token"]
print("Password Grant Success")

# Get a merchant
headers={"Authorization":"Bearer"+" "+token}
parameters={"email":"930994408@qq.com", "phone":"12345678901"}
pGet = requests.get(url+"/merchants", params=parameters, headers=headers)
if pGet.json()['data']['email']=="930994408@qq.com" and pGet.json()['data']['phone']=="12345678901":
    print("Get a merchant success")
else:
    print("Get a merchant failed")
    exit(0)
# print(pGet.text)

# Get a merchant by id
pGet = requests.get(url+"/merchants/5")
if pGet.json()['data']['merchant_id'] == 5:
    print("Get a merchant by id success")
else:
    print("Get a merchant by id failed")
    exit(0)
#print(pGet.json()['data']['merchant_id'])

# Update a merchant partially
headers={'Content-Type':'application/json',
         'Authorization':'Bearer'+' '+token}

post_data = '{"on": 1}'
pPatch = requests.patch(url+"/merchants/1", headers=headers, data=post_data)
if pPatch.json()['data']['on'] == 1:
    print("Update a merchant partially success")
else:
    print("Update a merchant partially failed")
    exit(0)


# Create an icon of a merchant

# Create a customer
headers = {'Content-Type':'application/json'}
post_data = '{"wechat_id":"aaaaa", "balance":0}'
pPost = requests.post(url+"/customers", headers=headers, data=post_data)

if pPost.json()['data']['wechat_id'] == "aaaaa" and pPost.json()['data']['balance'] == 0:
    print("Create a customer success")
else:
    print("Create a customer failed")

# Get a customer
pGet = requests.get(url+"/customers/"+str(pPost.json()['data']['customer_id']))
if pGet.json()['data']['customer_id'] == pPost.json()['data']['customer_id']:
    print("Get a customer success")
else:
    print("Get a customer failed")
    exit(0)

# Update a customer
headers = {'Content-Type':'application/json'}
post_data = '{"customer_id": 20, "wechat_id":"aaaaa", "balance":10}'
pPut = requests.put(url+"/customers/21", headers=headers, data=post_data)
if pPut.json()['data']['customer_id'] == 20 and pPut.json()['data']['balance'] == 10:
    print("Update a customer success")
else:
    print("Update a customer failed")
    exit(0)

# Create a seat
headers={'Content-Type':'application/json',
         'Authorization':'Bearer'+' '+token}

post_data = '{"number":"28A", "qr_code_url":"https://www.example.com", "merchant_id": 5}'
pPost = requests.post(url+"/seats", headers=headers, data=post_data)
# print(pPost.text)
if pPost.json()['code'] == 201 and pPost.json()['data']['number'] == "28A" and pPost.json()['data']['merchant_id'] == 5:
    print("Create a seat success")
else:
    print("Create a seat failed, program shut down")
    exit(0)

# List seats
parameters={"merchant_id":5}
headers={'Authorization':'Bearer'+' '+token}
pGet = requests.get(url+"/seats", params=parameters, headers=headers)
# print(pGet.text)
if pGet.json()['code'] == 200:
    print("List seats success")
else:
    print("List seats failed")
    exit(0)

# Delete a seat
headers={'Authorization':'Bearer'+' '+token}
pDelete = requests.delete(url+"/seats/47", headers=headers)
if pDelete.content == b'':
    print("Delete a seat success")
else:
    print("Delete a seat failed, program shut down")
    exit(0)

# Create an icon of a food

# Create a food
headers={'Content-Type':'application/json',
         'Authorization':'Bearer'+' '+token}
post_data = '{"name":"11", "description":"22", "price": 12.00, "merchant_id": 5, "icon_url":""}'
pPost = requests.post(url+"/foods", headers=headers, data=post_data)
if pPost.json()['code'] == 200 and pPost.json()['data']['merchant_id'] == 5:
    print("Create a food success")
else:
    print("Create a food failed, program shut down")
    exit(0)


# List foods
parameters = {"merchant_id": "5"}
pGet = requests.get(url+"/foods", params=parameters)
# print(pGet.text)
if pPost.json()['code'] == 200:
    print("List foods success")
else:
    print("List foods failed, program shut down")
    exit(0)

# Get a food
headers={'Authorization':'Bearer'+' '+token}
pGet = requests.get(url+"/foods/21", headers=headers)
if pGet.json()['code'] == 200 and pGet.json()['data']['food_id'] == 21:
    print("Get a food success")
else:
    print("Get a food failed, program shut down")
    exit(0)

# Delete a food
headers={'Authorization':'Bearer'+' '+token}
pDelete = requests.delete(url+"/foods/43", headers=headers)
if pDelete.content == b'':
    print("Delete a food success")
    exit(0)
else:
    print("Delete a food failed, program shut down")

# Create an icon of a food


# Create an order
headers={'Content-Type':'application/json'}
post_data='{"status":0, "seat_id":39, "customer_id":14, "merchant_id":5, "order_time": "2018-01-01T12:00:00+08:00", "complete_time": "2018-01-02T12:00:00+08:00", "foods":[{"food_id":47, "name":"11", "description":"22", "price": 12.00, "merchant_id":5, "amount":5}]}'
pPost = requests.post(url+"/orders", headers=headers, data=post_data)
if pPost.json()['code'] == 201 and pPost.json()['data']['seat_id'] == 39 and pPost.json()['data']['merchant_id'] == 5 and pPost.json()['data']['customer_id'] == 14:
    print("Create an order success")
else:
    print("Create an order failed, program shut down")
    exit(0)

# List orders by merchant_id
headers={'Authorization':'Bearer'+' '+token}
parameters={"merchant_id":5, "status":0}
pGet = requests.get(url+"/orders", headers=headers, params=parameters)
if pGet.json()['code'] == 200:
    print("List orders by merchant_id success")
else:
    print("List orders by merchant_id failed, program shut down")
    exit(0)

# List orders by customer_id
parameters={"customer_id":14, "status":0}
pGet = requests.get(url+"/orders", params=parameters)
if pGet.json()['code'] == 200:
    print("List orders by customer_id success")
else:
    print("List orders by customer_id failed, program shut down")
    exit(0)

# Get an order
headers={'Authorization':'Bearer'+' '+token}
pGet = requests.get(url+"/orders/14", headers=headers)
if pGet.json()['code'] == 200 and pGet.json()['data']['order_id'] == 14:
    print("Get an order success")
else:
    print("List orders by customer_id failed, program shut down")
    exit(0)

# Update an order partially
headers={'Content-Type':'application/json',
         'Authorization':'Bearer'+' '+token}
post_data='{"status":1}'
pPatch = requests.patch(url+"/orders/17", headers=headers, data=post_data)
if pPatch.json()['code'] == 201 and pPatch.json()['data']['order_id'] == 17 and pPatch.json()['data']['status']==1:
    print("Update an order partially success")
else:
    print("Update an order partially failed, program shut down")
    exit(0)
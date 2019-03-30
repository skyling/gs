import json, time
import requests
headers = {
    'Cookie': ''
}

s = requests.session()
url = "http://post.aliexpress.com/offer/category/ajaxGetProductCat.do?catIds={}&locale=zh_CN"
def getItem(id):
    id = str(id)
    surl = url.format(id)
    req = s.get(surl, headers=headers)
    data = json.loads(req.content)
    time.sleep(1)
    print(surl)
    return data['data'][id]

def getList(data):
    for item in data:
        if item['isLeaf']==False:
            subData = getItem(item['catId'])
            item['sub'] = subData
            getList(subData)
        else:
            item['hasPromission'] = checkPermission(item['catId'])

def checkPermission(catId):
    catId = str(catId)
    url = "http://post.aliexpress.com/offer/re_select_category.htm?catId={}&locale=zh_CN"
    resp = s.get(url.format(catId), headers=headers)
    data = resp.json()
    ret = False
    try:
        for item in data['data']['hasPromissionBrandAttrjson'][0]['data']:
            if item['id'] == "203470182":
                ret = True
    except:
        pass
    print("catId", catId, "hasPromission", ret)
    return ret 

def getall():
    data = getItem(15)
    getList(data)
    with open("./ok.json", '+w') as fp:
        json.dump(data, fp)


def outcsv():
    with open("./ok.json", "r") as fp:
        data = json.load(fp)
    with open("./ok.csv", "+a") as f:
        for item in data:
            if item['isLeaf'] == False:
                for item1 in item['sub']:
                    if item1['isLeaf'] == False:
                        for item2 in item1['sub']:
                            if item2['hasPromission']:
                                f.write('"{}({})","{}({})","{}({})"\n'.format(item['cnName'], item['enName'], 
                                    item1['cnName'], item1['enName'], item2['cnName'], item2['enName']))
                    elif item1['hasPromission']:
                        f.write('"{}({})","{}({})"\n'.format(item['cnName'], item['enName'], item1['cnName'], item1['enName']))

outcsv()
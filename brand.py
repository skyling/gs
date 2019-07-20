import json, time
import requests
headers = {
    'Cookie': 'ali_apache_id=11.251.145.137.1546750689596.182102.2; _ga=GA1.2.1530005911.1546750693; ali_beacon_id=11.251.145.137.1546750689596.182102.2; _uab_collina=154687032039323000760812; UM_distinctid=16828ab313617e-04f1262cf6c587-10306653-fa000-16828ab313a456; cna=MzcrD24we0ACAduGMASn+E+W; _fbp=fb.1.1549371120601.1120273472; _umdata=GF438946D3CC68A9D3245A4985B99EFE9523B9A; aefeMsite=amp-IuSZMnOtD8eu54WFlyFh4g; AMP_ECID_GOOGLE=amp-G_wOaT5in9Hn2s9-D8v6gw; amp-user-notification=amp-ebdXxsPEhRR_F7xyProibg; aep_common_f=FRnk9n1sL2Lums19wRQXCt3YRbusfyNRHKyT1EcYiri4fqUzC2BvkA==; CNZZDATA1254571644=116242106-1546866338-%7C1556505246; _gid=GA1.2.1879898646.1557146072; _m_h5_tk=346536f4986c53af0e09d5f9760652d6_1557330718209; _m_h5_tk_enc=07e70f3663214408b8a6cbd7459845ea; aep_history=keywords%5E%0Akeywords%09%0A%0Aproduct_selloffer%5E%0Aproduct_selloffer%0932899376888%0932922619977%0932982457355%091000008319275%0932987307288%0932948283704%091000008327085%0932983719235; _gat=1; xman_us_t=x_lid=cn1528211297apvd&sign=y&x_user=4elW7sYjkXDjYzCR1oiMeYxyvAS2jFiIqwUEA52Tffc=&ctoken=106o9pivgxpw3&need_popup=y&l_source=aliexpress; aep_usuc_t=ber_l=A0; xman_f=B1LGVi9fUM282OuivhQKhAKQw3Dk0hGxomSrD22F3pgURCxYF7HcBc5Kft/+9eElP7I/GsK8SdAPLd2KF2y/Sp+tAQIU3o2KMQdiAZliK/Ml6mCAWu8T0QXCoYxU29FX9+20zBAWPvFoNXbxTYF+Dg3wqtY4ABHzzak9piIuu0V/BOBo12o4c7oyWCuDBkeA5ehDtuSgEwS13Xb3bb3m6B4PAT59XC4FFb5YqhNZLwegziBBpJzgGBEPGgsUKXgHL+6ylzY9orG40ujA1TEpgX85E6obt8MkIM7y01BZFS4DNvaozOSptNWdQPhaM1g7q+jnMnkck/AACPNfgYuOoigBFXoe4xxk4GHxv1XW663DqUYIDZ0OqPRZ7JRimrM7rtwmOij1CIVi5Kn4S4qxioKhYjDG0unK; JSESSIONID=FC0A95A02AEE3F6CA0F0E9D3ED944115; ali_apache_track=ms=|mt=2|mid=cn1528211297apvd; ali_apache_tracktmp=W_signed=Y; xman_us_f=zero_order=y&x_locale=zh_CN&x_l=1&last_popup_time=1546868396401&x_user=CN|xiuqing|su|cnfm|238485195&no_popup_today=n; acs_usuc_t=x_csrf=ibiwwrl5wnkj&acs_rt=e6f2f1b54d8f40ffb9e8fcedc4c6d1ee; intl_locale=zh_CN; aep_usuc_f=isfm=y&site=glo&c_tp=USD&x_alimid=238485195&iss=y&s_locale=zh_CN&region=US&b_locale=en_US; isg=BBkZNMfIoVq9Zn3Mgp7Of3LSKAMzDg1RHTrXijvOlcC_QjnUg_YdKIfQQEaRYaWQ; intl_common_forever=jMWmFvSNQDufKpTYdWYrZActVKsRnDER3WRouhsIeJtFBKaUCrOGDg==; l=bBO4gpagvDa_i4gDBOCanurza77OSIRYYuPzaNbMi_5aN6TsA6_OlKlOcF96Vj5RsPTB4T0o1199-etbZ; xman_t=G49n7tq5+UyugWS5slg2GKkqGPKWYQfUu8wEkGAKJF9A60v7zVlit1MUvl8LM3deX4UM9w2FqTBgg6rUO8hV5N2KEEhpLD05zwjvyrlRhrvqJeMp1B5/VQ2IZY5yKCv0DQKwFJGmXwIs1WLVeolJnARXbkzVAlsTFPPLXgKXlwka8HfT0ObhHsPnukZi31v5wI1mBCf3nqj3C4+W8U8QhJ02YHkLvbgoevlli3AqDKUlJqxh+Njj1IblBZK79H4N9revP8v+Zmk06B4Ha81OS7fPTZmdW2nul07fKSfTIItsuyRrzw6K7w3GneLryrS3wLYfJdQ6JJVv8rvrYuAGreGM5rJBc/aFq2nIqaAiqyRXa5tosX9Z7R7AaRdLUDUjLXWQ60lYuQYSxRcTjDNt6tpYNKPKJ9xCcEMvBvVzKuAV694BKz2KA1Un2rnxHD4yiNXCwITDIdBHd3tBVdm/yxvs1W6Ve9WvMuZXRtOAou851QWgpeRrX7doLMP4AsN11PVtYJFGvPS/5cpEmcyi40A5Uev+zEhj600C50TBY+YMRNP5x5YhmmwrA7gNS4xbVNLiPltN0JcnX9tsdsmr4Te5TOcJs2XGPeQ/LCIEMJC5AoKGB3v4vEkpO400a41w0+GPKNRB79ty4Ib6Ug/p0hTgbBejHFcX'
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
            if item['id'] == "203019826":
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
# getall()
outcsv()
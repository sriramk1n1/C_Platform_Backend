from bs4 import BeautifulSoup
import requests
import re
import json 
from concurrent.futures import ThreadPoolExecutor
import threading

BACKEND_URL = "https://database.skapi.online"
CSES_SESSION_ID = "e79b248270acb86865b1f9fdc2838368ec05fcb6"


cookies = {
    "PHPSESSID": CSES_SESSION_ID
}

def parsequestion(url):
    global notyet
    id = url[-4:]
    id = int(id)
    print("Working on ",id,"-", url, "...")
    flag = 0
    page = requests.get(url, cookies=cookies)
    # Parse the HTML
    soup = BeautifulSoup(page.content, "html.parser")
    # print(soup.prettify())
    s = soup.find("div", class_="title-block")
    p = s.find("h1")
    name = p.text
    # Find the div with class "with-sidebar"
    with_sidebar_div = soup.find("body", class_="with-sidebar")

    # Find the div with class "content-wrapper" inside the "with-sidebar" div
    content_wrapper_div = with_sidebar_div.find("div", class_="content-wrapper")

    # Find the div with class "content" inside the "content-wrapper" div
    content_div = content_wrapper_div.find("div", class_="content")

    # Find all divs with class "md" inside the "content" div
    md_divs = content_div.find_all("div", class_="md")

    # Extract text from each div with class "md"
    str = ""
    for md_div in md_divs:
        str += md_div.get_text(strip=False)
    # print(str)
    t1 = str.split("Input")
    desc = t1[0]
    str = t1[1]+t1[2]
    # print(str)

    t1 = str.split("Output")
    constraint = t1[0]
    str = t1[1]+t1[2]
    # print(str)

    t1 = str.split("Constraints")
    constraint += t1[0]
    str = t1[1]
    # print(str)

    t1 = str.split("Example")
    constraint += t1[0]
    str = t1[1]
    # print(str)

    t1 = str[2:].split(":")
    # print(t1)
    input1 = t1[0]
    str = t1[1]
    t1 = str.split("Explanation")
    output1 = t1[0]

    dict = {
        "id": id,
        "name": name,
        "desc": desc,
        "constraints": constraint,
        "input1": input1.strip('\n '),
        "output1": output1.strip('\n ')
    }

    # print(dict)

    nav_sidebar_div = content_wrapper_div.find("div", class_="nav sidebar")

    # Find the 9th anchor tag inside the "nav sidebar" div
    anchors = nav_sidebar_div.find_all("a")
    ninth_anchor_href = ""
    if len(anchors) >= 9:
        ninth_anchor_href = anchors[8].get("href")
    else:
        print("NO SOLUTION")
        flag = 1


    if flag==0:
        print("Uploading...")
        # json_data = json.dumps(dict)
        # print(json_data)
        res = requests.post(BACKEND_URL+"/question", data=json.dumps(dict), headers={"Content-Type": "application/json"})
        if res.status_code == 200:
            print("uploaded ",dict["name"])
            lock.acquire()
            notyet+=1  
            lock.release()
        else:
            print("ERRRRRRRRRRROR", dict['name'], res.text)
        ninth_anchor_href = "https://cses.fi"+ninth_anchor_href
        page = requests.get(ninth_anchor_href, cookies=cookies)

        soup = BeautifulSoup(page.content, 'html.parser')

        with_sidebar_div = soup.find("body", class_="with-sidebar")

        # Find the div with class "content-wrapper" inside the "with-sidebar" div
        content_wrapper_div = with_sidebar_div.find("div", class_="content-wrapper")

        # Find the div with class "content" inside the "content-wrapper" div
        content_div = content_wrapper_div.find("div", class_="content")

        # print(content_div)
        l = content_div.find_all("a", class_="view")
        # print(l)
        c=0
        count = 0
        inp = ""
        out = ""
        print("Uploading test cases... ", id)
        for i in l:
            if c==2:
                c=0
                dict = {
                    "id": id,
                    "input": inp.strip('\n '),
                    "output": out.strip('\n '),
                }

                count+=1
                # print(json_data)
                res = requests.post(BACKEND_URL+"/testcase",data=json.dumps(dict), headers={"Content-Type": "application/json"})
                if res.status_code==200:
                    print(id , '-' , count, "successful")
                else:
                    print("Error uploading testcase...",count, res.text)
        
            k = i.get("href")
            if re.match(r"(/view/\d+/!.*)",k):
                pass
            else:
                k="https://cses.fi"+k
                z = requests.get(k,cookies=cookies)
                if c==0:
                    inp = z.content.decode()
                if c==1:
                    out = z.content.decode()
                c+=1
        if c==2:
                c=0
                dict = {
                    "id": id,
                    "input": inp.strip('\n '),
                    "output": out.strip('\n '),
                }
                count+=1
                res = requests.post(BACKEND_URL+"/testcase",data=json.dumps(dict), headers={"Content-Type": "application/json"})
                if res.status_code==200:
                    print(count, "successful")
                else:
                    print("Error uploading testcase...",count, res.text)


urls = [
    "https://cses.fi/problemset/task/1164"
]
# print(len(urls))

lock = threading.Lock()
notyet = 0
executor = ThreadPoolExecutor(max_workers=10)
for j in urls:
    executor.submit(parsequestion, j)
executor.shutdown(wait=True)
print("Uploaded ",notyet)
import os

def save(resp):
    status = resp["status"]
    content = resp["content"]
    f = open("./test.html", "w+")
    f.write(content)
    f.close()

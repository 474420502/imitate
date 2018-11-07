import os

def save(params):


    status = params["status"]
    content = params["content"]
    f = open("./test.html", "w+")
    f.write(content)
    f.close()

import os

def save(params):
    #status = params["status"]
    content = params["content"]
    f = open("/tmp/test.html", "w+")
    f.write(content)
    f.close()


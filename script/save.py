import os

def save(params):
    #status = params["status"]
    content = params["content"]
    with open("/tmp/test.html", "w+") as f:
        f.write(content)
        f.flush()


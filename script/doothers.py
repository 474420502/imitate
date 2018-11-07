

def doothers(params):
    content = params["content"]
    params["content"] = content + "\ndoothers"
    return "save", params
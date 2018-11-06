import imp
import os
import sys
import re
import inspect


def load_script():
    func_book = {}
    for t in os.walk("script"):
        src_path, _ , files = t
        for f in files:
            if f.endswith(".py"):
                script_path = src_path + "/" + f
                m = imp.load_source("script" , script_path)
                for t in inspect.getmembers(m, predicate=inspect.isfunction):
                    # t[1]({"content": 13, "status": 123})
                    func_book[t[0]] = t[1]
    return func_book
                
if __name__ == "__main__":
    load_script()
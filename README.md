# aoc-azure-aap-installer

### Mock server
As of now to run this UI locally we also require python mock server which contains the mock data and API's.

Here is UI folder you can see run.py file.Then just run this run.py and you'll have a server on port 9090.  It will return steps on http://127.0.0.1:9090/step" and will accept restart requests at "http://127.0.0.1:9090/execution/<id>/restart".
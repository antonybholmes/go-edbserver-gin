import json
from datetime import datetime

from dotenv import dotenv_values

INC = 2

# current time and date
# datetime object
time = datetime.now()

config = json.load(open("version.json"))

# print(config)

build = config.get("build", 0)
major, minor, patch = [int(x) for x in config["version"].split(".")]

build += INC

if build % 10 == 0:
    patch += INC

if patch % 10 == 0:
    minor += INC
    patch = 0

if minor > 9:
    major += INC
    minor = 0

config["version"] = f"{major}.{minor}.{patch}"
config["build"] = build
config["updated"] = time.strftime("%b %d, %Y")

with open("version.json", "w") as f:
    json.dump(config, f, indent=4)


# with open("version.env", "w") as f:
#     for key in sorted(config):
#         print(f'{key}="{config[key]}"', file=f)

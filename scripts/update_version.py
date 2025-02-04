from dotenv import dotenv_values
from datetime import datetime
 
# current time and date
# datetime object
time = datetime.now()

config = dotenv_values("version.env") 

#print(config)

major, minor, patch, build = [int(x) for x in config["VERSION"].split(".")]

build += 1
patch += 1

if patch > 9:
    minor += 1
    patch = 0

if minor > 9:
    major += 1
    minor = 0

config["VERSION"] = f"{major}.{minor}.{patch}.{build}"
config["UPDATED"] = time.strftime("%b %d, %Y")

with open("version.env", "w") as f:
    for key in sorted(config):
        print(f'{key}="{config[key]}"', file=f)
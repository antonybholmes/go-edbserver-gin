import base64
import uuid

id = uuid.uuid4()

print(id)
print(base64.urlsafe_b64encode(id.bytes).decode("utf8"))
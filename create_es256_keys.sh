openssl ecparam -name prime256v1 -genkey -noout -out jwt.es256.private.pem
openssl ec -in jwt.es256.private.pem -pubout -out jwt.es256.public.pem
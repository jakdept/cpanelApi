# cpapi
Golang cPanel API bindings

```bash
cookiejar="$(mktemp)"

url="$(ssh cpanel.jakdept.dev 'whmapi1 create_user_session --output=json user=root service=whostmgrd locale=en' |jq -r '.data.url')"
echo "${url}"

curl \
--location \
--output /dev/null \
--cookie "${cookiejar}" \
--cookie-jar "${cookiejar}" \
"${url}"

cat "${cookiejar}"

url="$(echo "${url}" | sed 's/login.*$/json-api\/listaccts\?want=user/g')"
echo "${url}"

echo

curl \
--cookie "${cookiejar}" \
--cookie-jar "${cookiejar}" \
"${url}" | jq
```

```bash
cookiejar="$(mktemp)"

url="$(ssh cpanel.jakdept.dev 'whmapi1 create_user_session --output=json user=root service=whostmgrd locale=en' |jq -r '.data.url')"
echo "${url}"

curl \
--location \
--output /dev/null \
--cookie "${cookiejar}" \
--cookie-jar "${cookiejar}" \
"${url}"

cat "${cookiejar}"

url="$(echo "${url}" | sed 's/login.*$/json-api\/listresellers/g')"
echo "${url}"

echo

curl \
--cookie "${cookiejar}" \
--cookie-jar "${cookiejar}" \
"${url}" | jq
```
# cpapi
Golang cPanel API bindings

A few demo examples:

> `demo/listaccts`

```bash
jhayhurst@jack:~/Downloads|⇒  GOOS=linux GOARCH=amd64 go build github.com/jakdept/cpapi/demo/listresellers
jhayhurst@jack:~/Downloads|⇒  GOOS=linux GOARCH=amd64 go build github.com/jakdept/cpapi/demo/listaccts
jhayhurst@jack:~/Downloads|⇒  scp listaccts listresellers wopr:~/
listaccts                                                                         100%   11MB  63.4MB/s   00:00
listresellers                                                                     100%   11MB  55.8MB/s   00:00
```

```bash
{20-02-11 19:57}wopr:~ jack% ./listaccts --host cpanel.jakdept.dev --port 12345 --keyfile /home/jack/.ssh/id_rsa
All accounts on server:
wordpressy
singledb
whousescpanel
{20-02-11 19:59}wopr:~ jack% ./listresellers --host cpanel.jakdept.dev --port 12345 --keyfile /home/jack/.ssh/id_rsa
All resellers on server:
whousescpanel
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

url="$(echo "${url}" | sed 's/login.*$/json-api\/listaccts\?want=user/g')"
echo "${url}"

echo

curl \
--cookie "${cookiejar}" \
--cookie-jar "${cookiejar}" \
"${url}" | jq
```


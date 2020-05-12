# GoThreaded
### Asynchronous DB Query Client built for use with systems like ShardMatrix for PHP

* Runs over a TCP Client
* Supports Docker
* Debug Mode
* Username & Password protected access
* Change Host and Ports that runs on

### Docker Compose example
```
version: "3.5"

services:
  threaded:
    image: jrsaunders/gothreaded
    environment:
      GOTHREADED_DEBUG: "true"
    ports:
      - 1534:1534
```


Default is to run on ::1534

### Environment Variables can be used to change these features:

* Debug mode
```
GOTHREADED_DEBUG=true
```
* Docker mode
```
GOTHREADED_DOCKER=true
```
* Set UserName
```
GOTHREADED_USER=myusernameforGoThreaded
```
* Set Password
```
GOTHREADED_PASS=mypasswordforGoThreaded
```
* Set Host
if left blank will run from ip of machine
```
GOTHREADED_HOST=mydomain.com
or
GOTHREADED_HOST=1.45.123.1
```
* Set Port
```
GOTHREADED_PORT=1534
```


### Example Query
```json
{
   "auth":{
      "username":"gothreaded",
      "password":"password"
   },
   "node_queries":[
      {
         "node":{
            "name":"DB0001",
            "dsn":{
               "driver":"mysql",
               "host":"127.0.0.1",
               "port":"3301",
               "dbname":"shard",
               "user":"root",
               "password":"password",
               "charset":"utf8mb4"
            },
            "geo":"UK",
            "docker_network":{
               "port":"3306",
               "host":"DB0001"
            }
         },
         "sql":"select * from `users` where `uuid` in (?, ?) order by `created` desc limit 15",
         "binds":[
            {
               "value":"06a00233-1ea8af83-d2dd-6e9c-8430-444230303031",
               "key":"0"
            },
            {
               "value":"06a00233-1ea8af83-cdcb-6fbc-83c6-444230303031",
               "key":"1"
            }
         ]
      },
      {
         "node":{
            "name":"DB0002",
            "dsn":{
               "driver":"mysql",
               "host":"127.0.0.1",
               "port":"3302",
               "dbname":"shard",
               "user":"root",
               "password":"password",
               "charset":"utf8mb4"
            },
            "geo":"UK",
            "docker_network":{
               "port":"3306",
               "host":"DB0002"
            }
         },
         "sql":"select * from `users` where `uuid` in (?, ?, ?, ?) order by `created` desc limit 15",
         "binds":[
            {
               "value":"06a00233-1ea8af83-d484-6002-8fee-444230303032",
               "key":"0"
            },
            {
               "value":"06a00233-1ea8af83-d359-6204-9134-444230303032",
               "key":"1"
            },
            {
               "value":"06a00233-1ea8af83-d11d-6a76-9622-444230303032",
               "key":"2"
            },
            {
               "value":"06a00233-1ea8af83-cc87-600c-a3ea-444230303032",
               "key":"3"
            }
         ]
      },
      {
         "node":{
            "name":"DB0003",
            "dsn":{
               "driver":"mysql",
               "host":"127.0.0.1",
               "port":"3303",
               "dbname":"shard",
               "user":"root",
               "password":"password",
               "charset":"utf8mb4"
            },
            "geo":"UK",
            "docker_network":{
               "port":"3306",
               "host":"DB0003"
            }
         },
         "sql":"select * from `users` where `uuid` in (?, ?) order by `created` desc limit 15",
         "binds":[
            {
               "value":"06a00233-1ea8af83-d3ec-6fa4-a355-444230303033",
               "key":"0"
            },
            {
               "value":"06a00233-1ea8af83-d25f-66d2-a29c-444230303033",
               "key":"1"
            }
         ]
      },
      {
         "node":{
            "name":"DB0007",
            "dsn":{
               "driver":"pgsql",
               "host":"127.0.0.1",
               "port":"5407",
               "dbname":"shard",
               "user":"postgres",
               "password":"password",
               "charset":"utf8"
            },
            "geo":"UK",
            "docker_network":{
               "port":"5432",
               "host":"DB0007"
            }
         },
         "sql":"select * from \"users\" where \"uuid\" in (?, ?, ?, ?, ?, ?, ?) order by \"created\" desc limit 15",
         "binds":[
            {
               "value":"06a00233-1ea8af83-d514-6a76-83ae-444230303037",
               "key":"0"
            },
            {
               "value":"06a00233-1ea8af83-d1b3-6116-8286-444230303037",
               "key":"1"
            },
            {
               "value":"06a00233-1ea8af83-d06a-65a2-b050-444230303037",
               "key":"2"
            },
            {
               "value":"06a00233-1ea8af83-cfb8-62c6-877a-444230303037",
               "key":"3"
            },
            {
               "value":"06a00233-1ea8af83-cf0a-62ac-beb7-444230303037",
               "key":"4"
            },
            {
               "value":"06a00233-1ea8af83-ce4e-66c4-b236-444230303037",
               "key":"5"
            },
            {
               "value":"06a00233-1ea8af83-cd22-6930-b3d3-444230303037",
               "key":"6"
            }
         ]
      }
   ]
}

```

Example Output
```json
{
   "nodes":[
      {
         "node_name":"DB0001",
         "data":[
            {
               "created":"2020-04-30 15:35:37",
               "email":"timmy19881345eaaf0491c4f2@google.com",
               "password":"cool!!81036",
               "something":"4",
               "username":"randy45339455eaaf0491c4ee",
               "uuid":"06a00233-1ea8af83-d2dd-6e9c-8430-444230303031"
            },
            {
               "created":"2020-04-30 15:35:36",
               "email":"timmy64913925eaaf0488ea83@google.com",
               "password":"cool!!34659",
               "something":"4",
               "username":"randy38638385eaaf0488ea80",
               "uuid":"06a00233-1ea8af83-cdcb-6fbc-83c6-444230303031"
            }
         ],
         "error":""
      },
      {
         "node_name":"DB0002",
         "data":[
            {
               "created":"2020-04-30 15:35:37",
               "email":"timmy59230775eaaf04928a18@google.com",
               "password":"cool!!53253",
               "something":"4",
               "username":"randy78956665eaaf04928a15",
               "uuid":"06a00233-1ea8af83-d359-6204-9134-444230303032"
            },
            {
               "created":"2020-04-30 15:35:37",
               "email":"timmy60420605eaaf04946846@google.com",
               "password":"cool!!86667",
               "something":"4",
               "username":"randy29941385eaaf04946842",
               "uuid":"06a00233-1ea8af83-d484-6002-8fee-444230303032"
            },
            {
               "created":"2020-04-30 15:35:36",
               "email":"timmy36838625eaaf0486e275@google.com",
               "password":"cool!!50751",
               "something":"4",
               "username":"randy33203695eaaf0486e270",
               "uuid":"06a00233-1ea8af83-cc87-600c-a3ea-444230303032"
            },
            {
               "created":"2020-04-30 15:35:36",
               "email":"timmy77103925eaaf048e39fc@google.com",
               "password":"cool!!68870",
               "something":"4",
               "username":"randy24805065eaaf048e39f9",
               "uuid":"06a00233-1ea8af83-d11d-6a76-9622-444230303032"
            }
         ],
         "error":""
      },
      {
         "node_name":"DB0003",
         "data":[
            {
               "created":"2020-04-30 15:35:37",
               "email":"timmy44145775eaaf0490fa94@google.com",
               "password":"cool!!77126",
               "something":"4",
               "username":"randy36538695eaaf0490fa91",
               "uuid":"06a00233-1ea8af83-d25f-66d2-a29c-444230303033"
            },
            {
               "created":"2020-04-30 15:35:37",
               "email":"timmy65560345eaaf049376a3@google.com",
               "password":"cool!!88802",
               "something":"4",
               "username":"randy24991755eaaf0493769f",
               "uuid":"06a00233-1ea8af83-d3ec-6fa4-a355-444230303033"
            }
         ],
         "error":""
      },
      {
         "node_name":"DB0007",
         "data":[
            {
               "created":"2020-04-30T15:35:37Z",
               "email":"timmy730485eaaf04954fb9@google.com",
               "password":"cool!!20823",
               "something":4,
               "username":"randy40609185eaaf04954fb5",
               "uuid":"06a00233-1ea8af83-d514-6a76-83ae-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy61932595eaaf0489bb33@google.com",
               "password":"cool!!13331",
               "something":4,
               "username":"randy34827445eaaf0489bb2f",
               "uuid":"06a00233-1ea8af83-ce4e-66c4-b236-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy16846205eaaf048ae795@google.com",
               "password":"cool!!9538",
               "something":4,
               "username":"randy90611535eaaf048ae792",
               "uuid":"06a00233-1ea8af83-cf0a-62ac-beb7-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy86084035eaaf0487db76@google.com",
               "password":"cool!!72042",
               "something":4,
               "username":"randy69330755eaaf0487db72",
               "uuid":"06a00233-1ea8af83-cd22-6930-b3d3-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy90589025eaaf048d1b1e@google.com",
               "password":"cool!!55238",
               "something":4,
               "username":"randy58136935eaaf048d1b1b",
               "uuid":"06a00233-1ea8af83-d06a-65a2-b050-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy46611915eaaf048f2909@google.com",
               "password":"cool!!58419",
               "something":4,
               "username":"randy77256935eaaf048f2905",
               "uuid":"06a00233-1ea8af83-d1b3-6116-8286-444230303037"
            },
            {
               "created":"2020-04-30T15:35:36Z",
               "email":"timmy71923665eaaf048bfdfe@google.com",
               "password":"cool!!77496",
               "something":4,
               "username":"randy23415065eaaf048bfdfb",
               "uuid":"06a00233-1ea8af83-cfb8-62c6-877a-444230303037"
            }
         ],
         "error":""
      }
   ]
}
```
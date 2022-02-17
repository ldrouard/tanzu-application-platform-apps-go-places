## go-places

Forked from [pvdbleek/go-places](https://github.com/pvdbleek/go-places)

This is probably one of the most useless API's there is, but it's not intended to be useful. 
It's mere purpose is to demonstrate how to create a workload written in `golang` in Tanzu Application Platform.

### Deploy the back-end database

We'll use an image from the VMware Application Catalog:
[MariaDB packaged by VMware](https://tac.bitnami.com/apps/158b6299c514cae9f673ba80ff08542b4d87b8bf)


```
kubectl create secret generic mariadb-secret \
  --from-literal=MARIADB_USER=dbuser\
  --from-literal=MARIADB_PASSWORD=secretpass
```

```
kubectl apply -f https://raw.githubusercontent.com/pvdbleek/go-places/main/mariadb/pvc.yaml
kubectl apply -f https://raw.githubusercontent.com/pvdbleek/go-places/main/mariadb/configmap.yaml
kubectl apply -f https://raw.githubusercontent.com/pvdbleek/go-places/main/mariadb/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/pvdbleek/go-places/main/mariadb/service.yaml
```
### Create a workload

Create the workload in TAP:

```
tanzu apps workload create go-places --git-repo https://github.com/pvdbleek/go-places --git-branch main --type web --label app.kubernetes.io/part-of=go-places --env "MARIADB_USER=dbuser" --env "MARIADB_PASS=secretpass" --env "MARIADB_HOST=mariadb" --yes
```

Watch it build and deploy:

```
tanzu apps workload tail go-places
```

Once done, fetch the URL:

```
tanzu apps workload get go-places
```

Which should result in something like this:

```
# go-places: Ready
---
lastTransitionTime: "2022-01-14T20:00:35Z"
message: ""
reason: Ready
status: "True"
type: Ready

Workload pods
NAME                                STATE       AGE
go-places-build-1-build-pod         Succeeded   10m
go-places-config-writer-dhpwc-pod   Succeeded   9m12s

Workload Knative Services
NAME        READY   URL
go-places   Ready   http://go-places.default.bleekie.tk
```

The repo also has a `catalog.yaml` which can be registered in your TAP GUI Catalogs.
When TAP GUI asks for the Repository URL, simply paste this and hit analyze:

```
https://github.com/pvdbleek/go-places/blob/main/catalog-info.yaml
```
### Using the API

The API has a the following endpoints:

| Endpoint     | Method      | Description                                |
| ------------ | ----------- | ------------------------------------------ |
| /places      | GET         | Fetches all places                         |
| /place/{id}  | GET         | Fetches a specific place                   |
| /places      | POST        | Add/replace a place                        |
| /url/{id}    | GET         | Generates a Google Maps URL for this place |

### Example POST

```
curl http://go-places.default.bleekie.tk/api/v1/place \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "Santa Clause Village","country": "Finland","description": "This is where Santa really lives.","latitude": "66.543733961877220","longitude": "25.847250943083463"}'
```

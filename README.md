# Keycloak

Setup on Rancher Desktop.

## Tasks

### deploy-keycloak

As per https://www.keycloak.org/getting-started/getting-started-kube

You will likely need to use the port forwarding feature to be able to access the API and web app. Forward port 41555 to 80.

```bash
kubectl create -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/latest/kubernetes/keycloak.yaml
```

### create-client

To create a client that can access an API, you create clients within a Realm.

* [Create the Realm](./01-create-realm.png)
* [Create a client for your API, remember the ID](./02-create-client.png)
* [Enable Authorization, which also enables the client credential flow](./03-create-client.png)
* [Set an optional root URL](./04-create-client.png)
* [Get the client secret](./05-get-secret.png)

```bash
echo No automated tasks
```

### get-token

The client must call Keycloak, passing the client ID and secret we created in the previous step.

This returns a JWT that the client can then pass to the API server.

Inputs: CLIENT_SECRET
Env: CLIENT_ID=file-api-sync
Env: KEYCLOAK_URL=http://keycloak:41555/realms/file-api/protocol/openid-connect/token

```bash
curl "$KEYCLOAK_URL" \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d 'grant_type=client_credentials' \
 -d "client_id=$CLIENT_ID" \
 -d "client_secret=$CLIENT_SECRET"
```

### start-server

Dir: server

```bash
go run .
```

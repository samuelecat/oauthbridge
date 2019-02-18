# oauthbridge
A reverse proxy that provides transparent authentication to resources protected by OAuth

---

## License
See the LICENSE file for details.

---

## Usage
1. From your Bitbucket account: go to OAuth, generate new client id and secret, set the scope on Repository
2. Edit the file ./conf/configuration.yml and add your client_id and client_secret 
3. docker-compose up

from now on you can access to your Bitbucket repository replacing the base URL from "https://bitbucket.org/" to "http://localhost:8081/bitbucket/"
the reverse proxy will authenticate.

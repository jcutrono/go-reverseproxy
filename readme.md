# Welcome to GWebProxy
- This is a reverse proxy built to enable authorization via a url key. All key maintenence is done via the admin api. It utilizes a basic authorizaton protocol.

# To initialize make an http post request as described

	https://{url}/api/key

	Headers:
	Referer: *url to proxy* (Ex. https://google.com)
	Authorization: Basic *authkey*

- The response will contain the key to your website
- You can configure the code to return the full URL
- You can configure the authentication for the API in the auth.go file

Other REST verbs perform the similar behaviors. The keys are associated to a URL through the Referer header. This will be used as the proxy value. It should contain the Scheme, Domain, and Path if necessary. It should not reference a page (like index.html)

# Including GWebProxy auth in your app (Node Example)


```
#!javascript

var express = require('express');
var app = express();

var authMiddleware = function(req, res, next) {
    if (req._parsedUrl.pathname === "/") {
        var options = {
            host: 'gwebproxy.ultilabs.xyz',
            path: '/api/key/' + req.query.key
        };

        var callback = function(response) {
            if(response.statusCode === 200) {
                next();
            } else {
                res.send(401)
            }
        };

        https.request(options, callback).end();
    } else {
        next();
    }
};

// Auth for EVERY request
app.use(authMiddleware);
app.use(express.static(__dirname + "/client"));

// On a route
app.post('/keywords', authMiddleware, function(req, res) {
    res.send(JSON.stringify({
        cool: "beans"
    }));
});
```
# Dependencies

[Gorilla](http://www.gorillatoolkit.org/)

[UUID](https://github.com/nu7hatch/gouuid)
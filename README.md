# http-echo

HTTP Echo is a small go web server that serves the contents it was started with
as an HTML page.

The default port is 5678, but this is configurable via the `-listen` flag:

```
http-echo -listen=:8080 -text="hello world"
```

Then visit http://localhost:8080/ in your browser.


# Flags

Flag  | Description | Default
----- | ----------- | -------
-listen | Listen port | 5678
-text | response text |
-version | display version information | false
-lagSecs | response lag in seconds | 5
-enableLag | enable response lag | false
-statusCode | custom return status code | 200

# Signals

Sig USR1 toggles the enableLag setting.

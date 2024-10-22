## Errors
- It is adviced to define sentinel errors as variable and not constants.
- As this makes comparison and wrapping of errors easier.


## Cobra
- Help won't be printed for a cmd, unless it or one of its sub-command have a Run defined.

## net/http
- The resp from the http request should be read using io.ReadAll, instead of using resp.Body.ReadAll([] byte).
- The later only fills the data as per the size of the buffer, if the buffer is empty nothing would be printed.
- The later should be used only when we need to read in a loop, with size restrictions.
- When mocking server with httptest servers. We need to have functions making the call to either accept the mock URL or a `http.Client` as a parameter.
- A client from a test server is configured to make request to the particular server
- Usually `http.Client` are set as package level variables.
- If getting certificate issue with a httptest TLS server, following can help.
```go
    client := ts.Client()
    client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
        InsecureSkipVerify: true,
    }

```

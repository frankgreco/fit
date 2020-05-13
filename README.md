# fit
> *WARNING:* The API should be considered unstable at this time. 

## overview
This project provides an encoder/decoder for the _Flexible and Interoperable Data Transfer_ (FIT) Protocol.

## usage
```
$ go build
$ ./fit decode -f test_data.fit
```

## todo
- [x] Profile agnostic API.
- [ ] Integrate with _Global FIT Profile_.
- [ ] Provide profile aware API.
- [ ] 100% unit test coverage.
- [ ] Finalize/document the API.
- [ ] Finalize/document the CLI.
- [ ] Implement gRPC and HTTP web servers.
- [ ] Implement encoder.

## references
[FIT SDK 21.30.00](https://www.thisisant.com/resources/fit-sdk/)
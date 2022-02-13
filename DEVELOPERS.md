# Developer Notes

## Testing

### Unit Tests

To run unit tests:

    make check

### Application Tests

After `make check` you can run a local server
using the bank database left over from the tests:

    ./bank -d _test.db server -a 127.0.0.1:8080

The test script sets up two users:

- `fred` with password `pass2`
- `bob` with password `pass4`

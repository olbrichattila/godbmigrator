# Test package usage

From root directory use make run-test or:

```
cd ./test
go test
```

This test is separated to a go module itself, it was done to test the base package but leave the dev dependencies out form the main migrator package and keep it tiny.

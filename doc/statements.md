## switch statement

```go
switch expr {
case expr:
    // no fallthrough
case expr, expr, expr:
    // no fallthrough
}
```

Within each case, a pre-declared name __case is the cased value.

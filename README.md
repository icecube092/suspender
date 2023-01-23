# Suspending execution based on value

Suspender is a value-based mutex.

## Example

### Simple

```go
s := suspender.New[uint64](suspender.Config{Count: 1})

userID := uint64(1)

err := s.Inc(userID) // nil

err = s.Inc(userID) // error: already locked for this user

userID2 := uint64(2)

err = s.Inc(userID2) // nil

err = s.Dec(userID) // nil
err = s.Inc(userID) // nil
```

### Auto-unlock on context done

```go
s := suspender.New[uint64](suspender.Config{Count: 1})

userID := uint64(1)

ctx, cancel := context.WithCancel(context.Background())
err := s.IncWithCtx(ctx, userID) // nil
err = s.IncWithCtx(ctx, userID) // error: already locked

cancel()

// notice: unlock may take some time
err = s.IncWithCtx(ctx, userID) // nil
```
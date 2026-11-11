context can carry values, but most importantly it can be used to signal a cancellation, for example, in case of a http request being cancelled by the client.
Passing a context with timeout when requesting a http page, then, http client's .Do(req) method will return error when the context is read.
You can also manually cancel a context, by context.WithCancel(ctx).

### r.Context().Done() channel

the following pattern does not make sense or useful since a go routine cannot kill the process.
```go
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    go func() {
        <-ctx.Done()
        // What do I do here?
    }()

    result := expensiveCalculation(r.URL.Query().Get("text"))

    json.NewEncoder(w).Encode(result)
}
```

correct implementation is something like this. periodically checking for the ctx.Err() in the CPU heavy function.
```go
func expensiveCalculation(ctx context.Context, input string) (string, error) {
    for i := 0; i < 1_000_000_000; i++ {
        if i%10000 == 0 {
        
	        // Periodically check for cancellation.
            if err := ctx.Err(); err != nil {
                return "", err
            }
        }

        // Expensive calculation
    }

    return "result", nil
}
```


> [!NOTE]
> pass ctx as function argument and watch for `ctx.Err()` every once in a while, for doing time consuming work, to exit the execution, saving un-needed work
> for http requests, always pass the ctx, then, once the user concel the request, the third party API call get cancelled since we have passed the ctx.

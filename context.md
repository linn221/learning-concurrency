context can carry values, but most importantly it can be used to signal a cancellation, for example, in case of a http request being cancelled by the client.
Passing a context with timeout when requesting a http page, then, http client's .Do(req) method will return error when the context is read.
You can also manually cancel a context, by context.WithCancel(ctx).


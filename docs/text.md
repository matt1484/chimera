# Binary/Text
`chimera` provides basic types that support both binary and plain/text requests/responses
- `BinaryRequest[Params any]`
- `BinaryResponse[Params any]`
- `Binary[Params any]` (can be used as a request or response)
- `PlainTextRequest[Params any]`
- `PlainTextResponse[Params any]`
- `PlainText[Params any]` (can be used as a request or response)
In general, these types are used to handle raw `string` or `[]byte` values in their body
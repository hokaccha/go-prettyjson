# prettyjson

JSON pretty print for Golang.

## Example

```go
v := map[string]interface{}{
    "str": "foo",
    "num": 100,
    "bool": false,
    "null": nil,
    "array": []string{"foo", "bar", "baz"},
    "map": map[string]interface{}{
        "foo": "bar",
    },
}
s, _ := prettyjson.Marshal(v)
fmt.Println(string(s))
```

![Output](http://i.imgur.com/cUFj5os.png)

## CLI(command-line interface)

Install:

```cmd
go get -u github.com/hokaccha/go-prettyjson
cd $GOPATH github.com/hokaccha/go-prettyjson/prettyjson
go install
```

Using:
```cmd
prettyjson -h
```
```
Usage: prettyjson [flags] file1.json file2.json ...

  -h	print help information
```

## License

MIT

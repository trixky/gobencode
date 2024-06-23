# gobencode

A simple library for parse and unmarshall the [bencode](https://en.wikipedia.org/wiki/Bencode) format with some utilities.

## Usage

### Get the package
```bash
go get -u github.com/trixky/gobencode
```

### Import the package
```go
import bencode "github.com/trixky/gobencode"
```

### Unmarshall from reader and get a bencode object

```golang
bc, err := gobencode.UnmarshallFromReader(reader)
```
### Or parse and unmarshall manually only what you want

```golang
data, err := gobencode.ParseFromReader(reader)

if err != nil {
    return err
}

bc := gobencode.bencode.Bencode{
    Data: data
}

if err := bc.UnmarshallAnnounce(); err != nil {
    return err
}
if err := bc.UnmarshallAnnounceList(); err != nil {
    return err
}
if err := bc.RandomizeAnnounceList(); err != nil {
    return err
}
if err := bc.UnmarshallInfo(); err != nil {
    return err
}
if err := bc.GetInfoHash(); err != nil {
    return err
}
```

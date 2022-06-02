# gobencode

A simple library for parse and unmarshall the [bencode](https://en.wikipedia.org/wiki/Bencode) format with some utilies.

## Usage

```golang
// quick start method

// parse, unmarshall and compute automatically the input
bc, err := gobencode.UnmarshallFromReader(reader)

/*
    available:

    bc.Announce
    bc.AnnounceList
    bc.RandomizedAnnounceList
    bc.Comment
    bc.CreatedBy
    bc.CreationDate
    bc.Info
        bc.Info.DirectoryName
        bc.Info.Files
            bc.Info.Files[...].Length
            bc.Info.Files[...].Path
            bc.Info.Files[...].Path
            bc.Info.Files[...].DecomposedPath
        bc.Info.PieceLength
        bc.Info.Pieces
    bc.InfoHash
    bc.UrlList
*/

```

```golang
// step by step method

// parse the input
data, err := gobencode.ParseFromReader(reader)

if err != nil {
    return err
}

bc := gobencode.bencode.Bencode{
    Data: data
}

bc.Data = data

// unmarshall and compute only what you need
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
// ...

/*
    available:

    bc.Announce
    bc.AnnounceList
    bc.RandomizedAnnounceList
    bc.Info
        bc.Info.DirectoryName
        bc.Info.Files
            bc.Info.Files[...].Length
            bc.Info.Files[...].Path
            bc.Info.Files[...].Path
            bc.Info.Files[...].DecomposedPath
        bc.Info.PieceLength
        bc.Info.Pieces
    bc.InfoHash
*/

```

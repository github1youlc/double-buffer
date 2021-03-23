## double buffer

A flexible library to manage resources which can be updated.


## Usage

Simple double buffer
```
// initialize
buffer := NewFileDoubleBuffer(&testLoader{}, func() interface{} { return 1 },)
buffer.Start()

// get latest data
buffer.Data() 
```


Local file double buffer
```
fl := NewFileLoader(
    LocalSingleFileLoadSystem,
    filePath,
    func(reader io.Reader, i interface{}) error {
        content, _ := ioutil.ReadAll(reader)
        s := i.(*string)
        *s = string(content)
        return nil
    })

buffer := NewDoubleBuffer(fl, func() interface{} {
    var s string
    return &s
})
```

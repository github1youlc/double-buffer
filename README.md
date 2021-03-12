## double buffer

A flexible library to manage resources which can be upgraded


## Usage

```
// initialize
buffer := NewFileDoubleBuffer(&testLoader{})
buffer.Start()

// get latest data
buffer.Data() 
```

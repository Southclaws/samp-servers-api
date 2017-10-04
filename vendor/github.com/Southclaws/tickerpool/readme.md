# TickerPool

[![Travis](https://img.shields.io/travis/Southclaws/tickerpool.svg)](https://travis-ci.org/Southclaws/tickerpool)[![Coverage](http://gocover.io/_badge/github.com/Southclaws/tickerpool)](http://gocover.io/github.com/Southclaws/tickerpool)

A worker pool of timed tasks, balanced equally to prevent cpu spikes.

This package is for creating a pool of workers which fire on a constant interval. The workers are balanced to avoid cpu spikes by simply dividing the interval by the amount of workers and offsetting their execution by this fraction.

For example, you want to scrape 3 pages one every second you could create a `time.NewTicker(time.Second)` and scrape all 3 at once but it would be better to create a `time.NewTicker(time.Second / 3)` and scrape each page 1/3 of a second apart from each other.

## Usage

Create a new TickerPool:
```go
func main() {
    tp := tickerpool.NewTickerPool(time.Second)
}
```

Add some workers:
```go
func addSomePages() {
    tp.Add("page1", func(){
        // scrape
    })
    tp.Add("page2", func(){
        // scrape
    })
    tp.Add("page3", func(){
        // scrape
    })
}
```

Maybe remove one or two when you don't need them any more:
```go
func imDone() {
    tp.Remove("page1")
}
```

All while the TickerPool balances the timing between each call automatically.

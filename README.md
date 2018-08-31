# ScaleImages-golang
small go app that scales images or prints them to stdout, will be containerized soon

# Usage

scale an image down by a factor of N

```go run shring.go scale input.png scale-factor(int)```

```go run shrink.go scale drake.png 5```

print image to stdout with optional width

```go run shrink.go print rc.png```

```go run shrink.go print rc.png 160```

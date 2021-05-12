# Yong Simple Web Framework
This project benchmarked gin-gonic.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://img.shields.io/github/license/rladyd818/yong) ![issue badge](https://img.shields.io/github/issues/rladyd818/yong) ![forks badge](https://img.shields.io/github/forks/rladyd818/yong) ![stars badge](https://img.shields.io/github/stars/rladyd818/yong)

### Installation
1. Go command to install Yong.
    ```sh
    $ go get -u github.com/rladyd818/yong
    ```
2. Import it in your code
    ```go
    import "github.com/rladyd818/yong"
    ```

### Simple example
```go
package main
import (
    "github.com/rladyd818/yong"
    "fmt"
)
func main() {
    r := yong.Default()
    r.GET("/", func(c *yong.Context) {
        fmt.Fprint(c.Writer, "Hello World")
    })
    r.Run(":8000")
}
```

### Method type list
GET, POST, PUT, PATCH, DELETE, OPTIONS

### Using Middleware example
```go
package main
import (
    "github.com/rladyd818/yong"
    "fmt"
)

func middleware() yong.HandlerFunc {
    return func(c *yong.Context) {
        c.Writer.Header().Set("m1", "Execute middleware")
    }
}

func middleware2() yong.HandlerFunc {
    return func(c *yong.Context) {
        c.Writer.Header().Set("m2", "Execute middleware2")
    }
}

func main() {
    r := yong.Default()
    r.USE("/", middleware(), middleware2()) // The "USE" method defines middleware.
    r.GET("/", func(c *yong.Context) {
        m1 := c.Writer.Header().Get("m1")
        m2 := c.Writer.Header().Get("m2")
        fmt.Println(m1)
	fmt.Println(m2)
        
        fmt.Fprint(c.Writer, "middleware Test") // response "middleware Test"
    })
    r.Run(":8000")
}
```

Result:
```sh
Execute middleware
Execute middleware2
```

### CORS Policy processing using middleware example
```go
func CORSMiddleware() yong.HandlerFunc {
	return func(c *yong.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Your client origin
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(204)
		}
	}
}

func main() {
    r := yong.Default()
    r.USE("/users", CORSMiddleware()) // The "USE" method defines middleware.
    r.GET("/users", func(c *yong.Context) {
        fmt.Fprint(c.Writer, "CORS Policy processing") // response "CORS Policy processing"
    })
    r.Run(":8000")
}
```

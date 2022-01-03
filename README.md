# Tistory SDK for Go

## How to use

```go
package main

import (
  "github.com/boltlessengineer/tistorysdk"
)

const (
  id string = "{{App ID}}"
  sk string = "{{Secret Key}}"
)

func main() {
  tistory := tistorysdk.NewClient(id, sk)
  //...
}
```

# 💬 Inspired By

- [jomei/notionapi](https://github.com/jomei/notionapi)
- [tistory/document-tistory-apis](https://github.com/tistory/document-tistory-apis)

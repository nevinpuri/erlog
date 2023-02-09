## zerolog

To send logs from zerolog, start by creating a new file and appending the following text to it

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NetLogger struct {
	queue []string
}

func New() NetLogger {
	return NetLogger{}
}

func (logger NetLogger) Write(p []byte) (n int, err error) {
	if !json.Valid(p) {
		fmt.Printf("%s is not valid, skipping\n", p)
		return
	}

	_, err = http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(p))

	if err != nil {
		return len(p), err
	}

	return len(p), nil
}
```

After you've created the NetLogger struct, then override the global zerolog log object.

```go
log.Logger = zerolog.New(netlogger.New()).With().Logger()
```

Then, any calls to `log.Print`, `Log.Error`, etc will be forwarded to erlog.

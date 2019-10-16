# logging

## `lib/common/log`
```go
func AddLogger(name string, l *log.Logger) error
func AddLoggerByPackageName() (*log.Entry, string)
func GetLogger(name string) *log.Entry

func GetDefaultLogger() *log.Entry
func GetSecurityLogger() *log.Entry
```

- `logrus` wrapper for common codes
- Use this package in codes that needs logging
- Allows creating and retriving logger by name, but cannot change logger settings
- `GetDefaultLogger()` and `GetSecurityLogger()` provide easy access to predefined loggers


## `lib/common/log/setup`
```go
func AddLogger(name, field string, l *log.Logger) error
func GetLogger(name string) *log.Entry
func SetLogger(name string, lv log.Level, fmt log.Formatter, out io.Writer, rc bool)
```

- `logrus` wrapper for application and testing codes
- This package holds the map structure for all loggers
- **Only** use it in application or test codes
- package named `internal` cannot be import from another repository


### In packages

- Get its own `logrus.Entry` instance with a call to `lib/common/log.AddLoggerByPackageName`
- The returned `logrus.Entry` logs the package name with every message
- These packages should not import `logrus`
- Any special loggers, e.g. the logger with name specified as *security*, can be
  declared prior to setup called in application
- `logrus.Entry` supports all logging functionality of a logger

```go
package example

import (
    commLog "lib/common/log"
)

var log, _ = commLog.AddLoggerByPackageName()
var slog, _ = commLog.GetSecurityLogger()

func DoSomething() {
    // Preferable to have it only in functions that use it
    // seclog := commLog.GetLogger("security")

    // trace log
    log.Trace()

    // debug log
    log.Debug()

    // write to security log
    slog.Trace()
}
```

#### Unit test code
- Use specified logger name for unit tests
- Redirect output to some buffer for auto checking
- Allows multiple tests logging at the same time

```go
package example_test

import (
    "testing"
    commLogInt "lib/common/log/setup"
)

var b bytes.Buffer
var testLog = commLogInt.GetLogger("example_test")

func testLogSetup() {
    // safely set named logger without interfering other tests
    commLogInt.SetLogger("example_test", log.Trace, nil, io.Writer(&b), false)
}

func TestDoSomething(t *testing.T) {

    testLogSetup()

    // log test process
    testLog.Trace("testing....")

    // this write to a buffer, automatically inspect the output with some code
    testLog.Trace("testing....")
}
```

### Inspecting Common Log in Tests
- Alter package global logger handle in test code
- In this case, test code should reside in a same package

```go
package example

// example_test.go

import (
    "testing"
    commLogInt "lib/common/log/setup"
)

var b1 bytes.Buffer
var b2 bytes.Buffer

func doFirst() {
    // change package global handle of these loggers
    log = commLogInt.GetLogger("example_test_default")
    slog = commLogInt.GetLogger("example_test_security")

    commLogInt.SetLogger("example_test_default", log.Trace, nil, io.Writer(&b1), false)
    commLogInt.SetLogger("example_test_security", log.Trace, nil, io.Writer(&b2), false)
}
```


### In application
- `lib/common/log` provides exported constants to default logger names
- import `lib/common/log/setup` for setting up loggers

```go
import (
    "github.com/sirupsen/logrus"

    commLog "lib/common/log"
    commLogInt "lib/common/log/setup"
)

func SetupLog() {
    secLogFile := os.OpenFile(consts.SecurityLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

    commLogInt.SetLogger(commLog.DefaultLoggerName, a.config().LogLevel, nil, os.Stdout, false)
    commLogInt.SetLogger(commLog.SecurityLoggerName, a.config().LogLevel, nil, secLogFile, false)
}
```


## Formatter

Implemented in `lib/common/log`, add another package `lib/common/log/formatter` if needed

```go
type IsecLFormatter struct {
	FormatDelimiter rune
	LineFormat      string
	TimeFormat      string
	LevelLength     int
}
```

### LevelLength
- How many alphabets are used to print level
- Should be 4 to 7 alphabets
- Truncate or padding space on right most part

### Format
- `FormatDelimiter`
    - The delimiter used in format string
    - Default is `$`
- `LineFormat`
    - Format string that will be used
    - A `\n` is automatically appended to the formatting string, no need explicitly add it
    - Arrange `tokens` for desired message
        - token `lv`    Print log level
        - token `pid`   Prints the process ID of logged item
        - token `t`     Print time
        - token `pkg`   Print pacakge name according to `package` field in log instance
        - token `msg`   Print log message
        - token `flds`  Print all fields that is not `package`
    - `tokens` should be prceeded and followed by `FormatDelimiter`
        - Example (Default format): `"$lv$[$pid$] $t$ $pkg$: $msg$; $flds$"`
        - Gets:
    ```log
    INFO[27628] 2019-10-14T14:32:22-07:00 intel/isecl/lib/common/log_test: Hello; field1=test field2=test
    ````
    - Token can be missing, but the total number of delimiters should equal to 10
        - Good: `"$lv$[$pid$] $t$: $msg$; $$$$"`, this will not show `pkg` and `flds`
            - Gets: 
            ```
            ERRO[10872] 2019-10-14T16:57:22-07:00: Hello;
            ```
        - Bad: `"[$lv$] t $pkg$: msg; $flds$"`, this will lead to an error and fall back to default format
- `TimeFormat`
    - Default is `time.RFC3339`
    - https://golang.org/pkg/time/#Time.Format

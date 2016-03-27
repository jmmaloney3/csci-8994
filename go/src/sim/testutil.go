package sim

import "testing"
import "runtime"
import "fmt"
import "path/filepath"

// Log an error with the source file and line number where the
// error occured
func LogErr(t *testing.T, msg string) {
  _, file, line, ok := runtime.Caller(2)
  file = filepath.Base(file)
  var msg2 string
  if (ok) {
    msg2 = fmt.Sprintf("[%s:%d]: %s", file, line, msg)
  } else {
    msg2 = fmt.Sprintf("[?file?:?line?]: %s", msg)
  }
  // log the error
  t.Error(msg2)
}
// assert that the two pointers point to the same object
func AssertAgentEqual(t *testing.T, r1 *Agent, r2 *Agent) {
  if (r1 != r2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", r1, r2))
  }
}
// assert that the two action modules are equal
func AssertActModEqual(t *testing.T, r1 *ActionModule, r2 *ActionModule) {
  if (!r1.SameBits(r2)) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", r1, r2))
  }
}
// assert that the two reputations are equal
func AssertRepEqual(t *testing.T, r1 Rep, r2 Rep) {
  if (r1 != r2) {
    var rep1 string
    var rep2 string
    if (r1 == 0) { rep1 = "GOOD" } else { rep1 = "BAD" }
    if (r2 == 0) { rep2 = "GOOD" } else { rep2 = "BAD" }
    LogErr(t, fmt.Sprintf("%v does not equal %v", rep1, rep2))
  }
}
// assert that the two reputations are NOT equal
func AssertRepNotEqual(t *testing.T, r1 Rep, r2 Rep) {
  if (r1 == r2) {
    var rep1 string
    var rep2 string
    if (r1 == 0) { rep1 = "GOOD" } else { rep1 = "BAD" }
    if (r2 == 0) { rep2 = "GOOD" } else { rep2 = "BAD" }
    LogErr(t, fmt.Sprintf("%v does equals %v", rep1, rep2))
  }
}
// assert that the two Int8s aee equal
func AssertIntEqual(t *testing.T, v1 int, v2 int) {
  if (v1 != v2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", v1, v2))
  }
}
// assert that the two Int8s aee equal
func AssertInt8Equal(t *testing.T, v1 int8, v2 int8) {
  if (v1 != v2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", v1, v2))
  }
}
// assert that the two Int32s aee equal
func AssertInt32Equal(t *testing.T, v1 int32, v2 int32) {
  if (v1 != v2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", v1, v2))
  }
}
// assert that the value is true
func AssertTrue(t *testing.T, b bool) {
  if (!b) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", b, true))
  }
}
// assert that the value is false
func AssertFalse(t *testing.T, b bool) {
  if (b) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", b, false))
  }
}

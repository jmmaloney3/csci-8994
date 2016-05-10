package testutil

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
// assert that v1 is greater than v2
func AssertInt32GT(t *testing.T, v1 int32, v2 int32) {
  if (v1 <= v2) {
    LogErr(t, fmt.Sprintf("%v is less than %v", v1, v2))
  }
}
// assert that v1 is equal to v2
func AssertFloat32Equal(t *testing.T, v1 float32, v2 float32) {
  if (v1 != v2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", v1, v2))
  }
}
// assert that v1 is equal to v2
func AssertFloat64Equal(t *testing.T, v1 float64, v2 float64) {
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

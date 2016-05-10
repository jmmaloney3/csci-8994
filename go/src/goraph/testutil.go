package goraph

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
// assert that the two vertices are equal
func AssertVertexEqual(t *testing.T, v1 Vertex, v2 Vertex) {
  if (v1 != v2) {
    LogErr(t, fmt.Sprintf("%v does not equal %v", v1, v2))
  }
}

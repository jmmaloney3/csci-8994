package sim

import "testing"

// assert that the two pointers point to the same object
func AssertAgentEqual(t *testing.T, r1 *Agent, r2 *Agent) {
  if (r1 != r2) {
    t.Error(r1, " does not equal ", r2)
  }
}
// assert that the two action modules are equal
func AssertActModEqual(t *testing.T, r1 *ActionModule, r2 *ActionModule) {
  if (r1 != r2) {
    t.Error(r1, " does not equal ", r2)
  }
}
// assert that the two reputations are equal
func AssertRepEqual(t *testing.T, r1 Rep, r2 Rep) {
  if (r1 != r2) {
    t.Error(r1, " does not equal ", r2)
  }
}
// assert that the two Int8s aee equal
func AssertIntEqual(t *testing.T, v1 int, v2 int) {
  if (v1 != v2) {
    t.Error(v1, " does not equal ", v2)
  }
}
// assert that the two Int8s aee equal
func AssertInt8Equal(t *testing.T, v1 int8, v2 int8) {
  if (v1 != v2) {
    t.Error(v1, " does not equal ", v2)
  }
}
// assert that the two Int32s aee equal
func AssertInt32Equal(t *testing.T, v1 int32, v2 int32) {
  if (v1 != v2) {
    t.Error(v1, " does not equal ", v2)
  }
}
// assert that the value is true
func AssertTrue(t *testing.T, b bool) {
  if (!b) {
    t.Error(b, " is not true")
  }
}
// assert that the value is false
func AssertFalse(t *testing.T, b bool) {
  if (b) {
    t.Error(b, " is not true")
  }
}

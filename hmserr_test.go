package base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////
//
// RFC 7807-compliant Problem Details struct
//
////////////////////////////////////////////////////////////////////////////

var probDetailsGood = &ProblemDetails{
	Type:     "https://example.com/probs/CustomError",
	Title:    "A custom problem type occurred",
	Detail:   "Your operation generated a custom error type during GET",
	Instance: "/object/1/errors/custom-error-on-get/err_instance=1",
	Status:   http.StatusMethodNotAllowed,
}

// Same as above, make sure not modfied unexpectedly.
var probDetailsCopy = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   probDetailsGood.Detail,
	Instance: probDetailsGood.Instance,
	Status:   probDetailsGood.Status,
}

// Instance and Detail change, rest inherited from parent
var probDetailsGoodChild1 = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   "Your operation generated a custom error during GET of object/2",
	Instance: "/object/2/errors/custom-error-on-get/err_instance=1",
	Status:   probDetailsGood.Status,
}

// Only Detail changed, rest inherited from parent.
var probDetailsGoodChild2 = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   "Your operation generated a custom error during GET of object/2",
	Instance: probDetailsGood.Instance,
	Status:   probDetailsGood.Status,
}

var probDetailsNoStatus = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   probDetailsGood.Detail,
	Instance: probDetailsGood.Instance,
}

var probDetailsStatus = &ProblemDetails{
	Type:   ProblemDetailsHTTPStatusType,
	Title:  "Not Found",
	Detail: "Your resource could not be found",
	Status: http.StatusNotFound,
}

var probDetailsStatusDefault = &ProblemDetails{
	Type:   probDetailsStatus.Type,
	Title:  "Bad Request",
	Detail: probDetailsStatus.Detail,
	Status: http.StatusBadRequest,
}

func testCheckProblemDetails(expected, test *ProblemDetails) error {
	if expected.Type != test.Type {
		return fmt.Errorf("MISMATCH (Type) expected: %s, test: %s",
			expected.Type, test.Type)
	}
	if expected.Title != test.Title {
		return fmt.Errorf("MISMATCH (Title) expected: %s, test: %s",
			expected.Title, test.Title)
	}
	if expected.Detail != test.Detail {
		return fmt.Errorf("MISMATCH (Detail) expected: %s, test: %s",
			expected.Detail, test.Detail)
	}
	if expected.Instance != test.Instance {
		return fmt.Errorf("MISMATCH (Instance) expected: %s, test: %s",
			expected.Instance, test.Instance)
	}
	if expected.Status != test.Status {
		return fmt.Errorf("MISMATCH (Status) expected: %d, test: %d",
			expected.Status, test.Status)
	}
	return nil
}

func TestNewProblemDetails(t *testing.T) {
	p := NewProblemDetails(
		"https://example.com/probs/CustomError",                   // ptype
		"A custom problem type occurred",                          // title
		"Your operation generated a custom error type during GET", // detail
		"/object/1/errors/custom-error-on-get/err_instance=1",     // instance
		http.StatusMethodNotAllowed,                               // status
	)

	if err := testCheckProblemDetails(probDetailsGood, p); err != nil {
		t.Errorf("Testcase 1: FAIL: Comparison failed: %s", err)
	} else {
		t.Logf("Testcase 1: Pass: Comparison succeeded")
	}
}

func TestNewProblemDetailsStatus(t *testing.T) {
	// Normal usage
	p := NewProblemDetailsStatus(
		"Your resource could not be found", // detail
		http.StatusNotFound,                // status
	)
	if err := testCheckProblemDetails(probDetailsStatus, p); err != nil {
		t.Errorf("Testcase 1: FAIL: Comparison failed: %s", err)
	} else {
		t.Logf("Testcase 1: Pass: Comparison succeeded (normal usage)")
	}
	// Abnormal usage, status is not valid (out-of-range)
	p = NewProblemDetailsStatus(
		"Your resource could not be found", // detail
		1234,                               // status (invalid)
	)
	if err := testCheckProblemDetails(probDetailsStatusDefault, p); err != nil {
		t.Errorf("Testcase 2a: FAIL: Comparison failed (bad status 1234): %s",
			err)
	} else {
		t.Logf("Testcase 2a: Pass: Comparison succeeded (bad status 1234)")
	}
	// Abnormal usage, status is not valid (zero)
	p = NewProblemDetailsStatus(
		"Your resource could not be found", // detail
		0,                                  // status (invalid)
	)
	if err := testCheckProblemDetails(probDetailsStatusDefault, p); err != nil {
		t.Errorf("Testcase 2b: FAIL: Comparison failed (bad status 0): %s",
			err)
	} else {
		t.Logf("Testcase 2b: Pass: Comparison succeeded (bad status 0)")
	}
	// Abnormal usage, status is not valid (negative)
	p = NewProblemDetailsStatus(
		"Your resource could not be found", // detail
		-1,                                 // status (invalid)
	)
	if err := testCheckProblemDetails(probDetailsStatusDefault, p); err != nil {
		t.Errorf("Testcase 2a: FAIL: Comparison failed (bad status -1): %s",
			err)
	} else {
		t.Logf("Testcase 2a: Pass: Comparison succeeded (bad status -1)")
	}
}

func TestNewChild(t *testing.T) {
	// Case 1, both blank, child should match parent
	p := probDetailsGood.NewChild("", "")
	if err := testCheckProblemDetails(probDetailsGood, p); err != nil {
		t.Errorf("Testcase 1: FAIL: Blank args, but child != parent: %s",
			err)
	} else {
		t.Logf("Testcase 1: Pass: Both arg blank, child matched parent")
	}
	// Case 2: Should override both Instance and Detail but keep other fields
	p = probDetailsGood.NewChild(
		probDetailsGoodChild1.Detail,
		probDetailsGoodChild1.Instance,
	)
	if err := testCheckProblemDetails(probDetailsGoodChild1, p); err != nil {
		t.Errorf("Testcase 2: FAIL: Child didn't match expected changes: %s",
			err)
	} else {
		t.Logf("Testcase 2: Pass: Child matched but with new Detail/Instance")
	}
	if err := testCheckProblemDetails(probDetailsCopy,
		probDetailsGood); err != nil {

		t.Errorf("Testcase 2a: FAIL: Original WAS modified: %s", err)
	} else {
		t.Logf("Testcase 2a: Pass: Original Problem not modified.")
	}
	// Case 3: Only Detail should be overridden
	p = probDetailsGood.NewChild(
		probDetailsGoodChild2.Detail,
		"",
	)
	if err := testCheckProblemDetails(probDetailsGoodChild2, p); err != nil {
		t.Errorf("Testcase 3: FAIL: Child didn't match expected changes: %s",
			err)
	} else {
		t.Logf("Testcase 3: Pass: Child matched but with new Detail only")
	}

}

///////////////////////////////////////////////////////////////////////////
// http.ResponseWriter response formatting for ProblemDetails
///////////////////////////////////////////////////////////////////////////

func TestSendProblemDetails(t *testing.T) {
	// Testcase 1 - Fully populated stucture, use Status field in
	// ProblemDetails since status arg is 0
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetails(w, probDetailsGood, 0)
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler1(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != probDetailsGood.Status {
		t.Errorf("Testcase 1a: FAIL: Status was: %d, not %d",
			resp.StatusCode, probDetailsGood.Status)
	} else {
		t.Logf("Testcase 1a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 1b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 1b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	barr, err := json.Marshal(probDetailsGood)
	if err != nil {
		t.Errorf("Testcase 1c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(string(barr)) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 1c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(string(barr)),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 1c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}

	// Testcase 2 - Force status field override, use as response StatusCode
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetails(w, probDetailsGood, http.StatusNotFound)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handler2(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Testcase 2a: FAIL: Status was: %d, not %d",
			resp.StatusCode, http.StatusNotFound)
	} else {
		t.Logf("Testcase 2a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 2b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 2b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	// Body should be the same, only StatusCode for HTTP response should be
	// different
	barr, err = json.Marshal(probDetailsGood)
	if err != nil {
		t.Errorf("Testcase 2c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(string(barr)) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 2c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(string(barr)),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 2c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}

	// Testcase 3 - 0 status in ProblemDetails, and 0 used as status arg
	// Should default to HTTP StatusCode http.StatusBadRequest
	handler3 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetails(w, probDetailsNoStatus, 0)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handler3(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Testcase 3a: FAIL: Status was: %d, not %d",
			resp.StatusCode, http.StatusBadRequest)
	} else {
		t.Logf("Testcase 3a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 3b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 3b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	// Body should be the same, only StatusCode for HTTP response should be
	// different.  The unset(i.e. 0) ProblemDetails should not be encoded.
	barr, err = json.Marshal(probDetailsNoStatus)
	if err != nil {
		t.Errorf("Testcase 3c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(string(barr)) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 3c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(string(barr)),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 3c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}

	// Testcase 4 - No Status in ProblemDetails, left blank, but
	// HTTP StatusCode should be the status arg, http.StatusNotFound
	handler4 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetails(w, probDetailsNoStatus, http.StatusNotFound)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handler4(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Testcase 4a: FAIL: Status was: %d, not %d",
			resp.StatusCode, http.StatusNotFound)
	} else {
		t.Logf("Testcase 4a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 4b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 4b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	// Body should be the same, only StatusCode for HTTP response should be
	// different.  The unset(i.e. 0) ProblemDetails should not be encoded.
	barr, err = json.Marshal(probDetailsNoStatus)
	if err != nil {
		t.Errorf("Testcase 4c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(strings.TrimSpace(string(barr))) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 4c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(strings.TrimSpace(string(barr))),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 4c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}
}

func TestSendProblemDetailsGeneric(t *testing.T) {
	// Testcase 1 - msg and status set properly.
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetailsGeneric(w,
			probDetailsStatus.Status, // status
			probDetailsStatus.Detail, // msg
		)
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler1(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != probDetailsStatus.Status {
		t.Errorf("Testcase 1a: FAIL: Status was: %d, not %d",
			resp.StatusCode, probDetailsStatus.Status)
	} else {
		t.Logf("Testcase 1a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 1b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 1b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	barr, err := json.Marshal(probDetailsStatus)
	if err != nil {
		t.Errorf("Testcase 1c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(string(barr)) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 1c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(string(barr)),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 1c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}

	// Testcase 2 - Status set improperly
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		SendProblemDetailsGeneric(w,
			1234,                            // status arg - bad
			probDetailsStatusDefault.Detail, // msg
		)
	}
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	handler2(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != probDetailsStatusDefault.Status {
		t.Errorf("Testcase 2a: FAIL: Status was: %d, not %d",
			resp.StatusCode, probDetailsStatusDefault.Status)
	} else {
		t.Logf("Testcase 2a: Pass: Status was %d as expected.",
			resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != ProblemDetailContentType {
		t.Errorf("Testcase 2b: FAIL: Content-Type was %s not %s",
			resp.Header.Get("Content-Type"), ProblemDetailContentType)
	} else {
		t.Logf("Testcase 2b: Pass: Content-Type was %s as expected.",
			ProblemDetailContentType)
	}
	// Body should match good example given above, with default status
	// code and matching title given.
	barr, err = json.Marshal(probDetailsStatusDefault)
	if err != nil {
		t.Errorf("Testcase 2c: INTERNAL ERROR: can't encode good response")
	} else {
		if strings.TrimSpace(string(barr)) !=
			strings.TrimSpace(strings.TrimSpace(string(body))) {

			t.Errorf("Testcase 2c: FAIL: wanted body '%s', got '%s'",
				strings.TrimSpace(string(barr)),
				strings.TrimSpace(strings.TrimSpace(string(body))))
		} else {
			t.Logf("Testcase 2c: Pass: got body '%s' as expected",
				strings.TrimSpace(string(barr)))
		}
	}
}

////////////////////////////////////////////////////////////////////////////
//
// HMSError - Works like a standard Go 'error', but can be distinguished
//            as an HMS-specific error with an optional class for better
//            handling upstream.  An RFC7807 error can optionally be added
//            in case we need to pass those through multiple layers- but
//            without forcing us to (since they still look like regular
//            Go errors)
//
//            Part of the motivation is so we can determine which errors
//            are safe to return to users, i.e. we don't want to give
//            them things that expose details about the database structure.
//
////////////////////////////////////////////////////////////////////////////

var hmsError = &HMSError{
	Class:   "HMSErrorClass1",
	Message: "This is an HMSError message",
}

// Same as above, make sure it isn't modified unintentionally.
var hmsErrorCopy = &HMSError{
	Class:   hmsError.Class,
	Message: hmsError.Message,
}

// Same as above, make sure it isn't modified unintentionally.
var hmsErrorNewMessage = &HMSError{
	Class:   hmsError.Class,
	Message: "This is a new Message",
}

var hmsErrorNoClass = &HMSError{
	Message: hmsError.Message,
}

var hmsErrorNoMessage = &HMSError{
	Class: hmsError.Class,
}

var hmsErrorEmpty = &HMSError{}

var hmsErrorWProb = &HMSError{
	Class:   hmsError.Class,
	Message: hmsError.Message,
	Problem: probDetailsGood,
}

var hmsErrorWProbCopy = &HMSError{
	Class:   hmsErrorWProb.Class,
	Message: hmsErrorWProb.Message,
	Problem: probDetailsCopy,
}

func testCheckHMSError(expected, test *HMSError) error {
	if expected.Class != test.Class {
		return fmt.Errorf("MISMATCH (Class) expected: %s, test: %s",
			expected.Class, test.Class)
	}
	if expected.Message != test.Message {
		return fmt.Errorf("MISMATCH (Message) expected: %s, test: %s",
			expected.Message, test.Message)
	}
	if expected.Problem == nil && test.Problem != nil {
		return fmt.Errorf("MISMATCH (Problem) expected: nil, test: non-nil")
	}
	if expected.Problem != nil && test.Problem == nil {
		return fmt.Errorf("MISMATCH (Problem) expected: non-nil, test: nil")
	}
	if expected.Problem != nil && test.Problem != nil {
		return testCheckProblemDetails(expected.Problem, test.Problem)
	}
	return nil
}

func TestNewHMSError(t *testing.T) {
	p := NewHMSError(hmsError.Class, hmsError.Message)
	if p.Class != hmsError.Class || p.Message != hmsError.Message {
		t.Errorf("Testcase 1: FAIL: got '%s' and '%s'", p.Class, p.Message)
	} else {
		t.Logf("Testcase 1: Pass: got '%s' and '%s'", p.Class, p.Message)
	}
	if p.Problem != nil {
		t.Logf("Testcase 2: FAIL: Problem pointer was not nil")
	} else {
		t.Logf("Testcase 2: Pass: No default problem")
	}
}

// Test the message that is printed when HMSError is printed like an
// ordinary Go error.
func TestError(t *testing.T) {
	// Both class and message set, error string should be message
	if hmsError.Error() != hmsError.Message {
		t.Errorf("Testcase 1: FAIL: Error() got '%s', expected '%s'",
			hmsError.Error(), hmsError.Message)
	} else {
		t.Logf("Testcase 1: Pass: Error() got '%s' as expected",
			hmsError.Message)
	}
	// Only message set, error string should again be message
	if hmsErrorNoClass.Error() != hmsErrorNoClass.Message {
		t.Errorf("Testcase 2: FAIL: Error() got '%s', expected '%s'",
			hmsErrorNoClass.Error(), hmsErrorNoClass.Message)
	} else {
		t.Logf("Testcase 2: Pass: Error() got '%s' as expected",
			hmsErrorNoClass.Message)
	}
	// Message is not set, only Class, use that as it is better than
	// a default message.
	if hmsErrorNoMessage.Error() != hmsErrorNoMessage.Class {
		t.Errorf("Testcase 3: FAIL: Error() got '%s', expected '%s'",
			hmsErrorNoMessage.Error(), hmsErrorNoMessage.Class)
	} else {
		t.Logf("Testcase 3: Pass: Error() got '%s' as expected",
			hmsErrorNoMessage.Class)
	}
	// Neither Message or Class is set, FAILreturn a generic string instead of
	// an empty string so we at least can see what went wrong.
	if hmsErrorEmpty.Error() != HMSErrorUnsetDefault {
		t.Errorf("Testcase 4: FAIL: Error() got '%s', expected '%s'",
			hmsErrorEmpty.Error(), HMSErrorUnsetDefault)
	} else {
		t.Logf("Testcase 4: Pass: Error() got '%s' as expected",
			HMSErrorUnsetDefault)
	}

	// This just needs to compile.  Make sure we can use HMSError wherever
	// a standard Go 'error' is returned.
	testFunc := func() error {
		return hmsError
	}
	t.Logf("Testcase 5: Pass: Check 'error' returning function: %s",
		testFunc().Error())
}

// Test that we can distinguish an HMSError from other structs that implement
// Error()
func TestIsHMSError(t *testing.T) {
	var someErr error

	if IsHMSError(hmsErrorWProb) {
		t.Logf("Testcase 1: Pass: Real HMSError returned 'true' as expected.")
	} else {
		t.Errorf("Testcase 1: FAIL: Real HMSError returned 'false'.")
	}
	if !IsHMSError(someErr) {
		t.Logf("Testcase 2: Pass: vanilla error returned 'false' as expected.")
	} else {
		t.Errorf("Testcase 1: FAIL: vanilla error returned 'true'.")
	}
}

// If a generic Go error is also an HMSError, return it and set true, else
// return nil and set false.
// This checks that we can convert an error into the full HMSError struct
// underneath.
func TestGetHMSError(t *testing.T) {
	// Testcase 1: returned 'error' is HMSError
	testFuncHMSErr := func() error {
		return hmsErrorWProb
	}
	errTest := testFuncHMSErr()
	if herr, valid := GetHMSError(errTest); valid == true {
		t.Logf("Testcase 1a: Pass: Real HMSError returned 'true' as expected.")
		if err := testCheckHMSError(hmsErrorWProb, herr); err != nil {
			t.Errorf("Testcase 1b: FAIL: Returned HMSError != orig: %s", err)
		} else {
			t.Logf("Testcase 1b: Pass: Returned HMSError matches input.")
		}
	} else {
		t.Errorf("Testcase 1a: FAIL: Real HMSError returned 'false'.")
	}

	// Testcase 2: returned 'error' is a vanilla, non-HMSError error.
	testFuncNormErr := func() error {
		return fmt.Errorf("I'm a normal error, not an HMSError")
	}
	errTest2 := testFuncNormErr()
	if herrNil, valid := GetHMSError(errTest2); valid == false {
		t.Logf("Testcase 2a: Pass: Non-HMSError returned 'false' as expected.")
		if herrNil != nil {
			t.Errorf("Testcase 2b: FAIL: Returned HMSError was not nil.")
		} else {
			t.Logf("Testcase 2b: Pass: Returned HMSError was nil as expected.")
		}
	} else {
		t.Errorf("Testcase 2a: FAIL: Non-HMSError returned 'true'.")
	}
}

// Return true if 'class' exactly matches Class field for HMSError
func TestIsClass(t *testing.T) {
	if !hmsError.IsClass(hmsError.Class) {
		t.Errorf("Testcase 1: FAIL: Should have matched class")
	} else {
		t.Logf("Testcase 1: Pass: Class matched")
	}
	if hmsError.IsClass(strings.ToUpper(hmsError.Class)) {
		t.Errorf("Testcase 2: FAIL: False true, Should've detected mixed case")
	} else {
		t.Logf("Testcase 1: Pass: false as expected - case was mismatched")
	}
	if hmsError.IsClass("sadasdasdas") {
		t.Errorf("Testcase 3: FAIL: False true, Should've detected bad class")
	} else {
		t.Logf("Testcase 3: Pass: Bad class did not match")
	}

}

// Return true if 'class' matches Class field for HMSError (case insensitive)
func TestIsClassIgnCase(t *testing.T) {
	if !hmsError.IsClassIgnCase(hmsError.Class) {
		t.Errorf("Testcase 1: FAIL: Should have matched class")
	} else {
		t.Logf("Testcase 1: Pass: Class matched")
	}
	if !hmsError.IsClassIgnCase(strings.ToUpper(hmsError.Class)) {
		t.Errorf("Testcase 2: FAIL: Got false, should've ignored mixed case")
	} else {
		t.Logf("Testcase 2: Pass: true as expected, case was mismatched")
	}
	if hmsError.IsClassIgnCase("sadasdasdas") {
		t.Errorf("Testcase 3: FAIL: False true, Should've detected bad class")
	} else {
		t.Logf("Testcase 3: Pass: Bad class did not match")
	}
}

// Returns false if 'err' is not an HMSError, or, if it is, if 'class' doesn't
// match the HMSError's Class field.
func TestIsHMSErrorClass(t *testing.T) {
	var regErr error
	if !IsHMSErrorClass(hmsError, hmsError.Class) {
		t.Errorf("Testcase 1: FAIL: Should have matched class")
	} else {
		t.Logf("Testcase 1: Pass: Class matched")
	}
	if IsHMSErrorClass(hmsError, strings.ToUpper(hmsError.Class)) {
		t.Errorf("Testcase 2: FAIL: Got true, should've caught mixtrueed case")
	} else {
		t.Logf("Testcase 2: Pass: false as expected, case was mismatched")
	}
	if IsHMSErrorClass(hmsError, "sadasdasdas") {
		t.Errorf("Testcase 3: FAIL: False true, Should've detected bad class")
	} else {
		t.Logf("Testcase 3: Pass: Bad class did not match")
	}
	if IsHMSErrorClass(regErr, "") {
		t.Errorf("Testcase 4: FAIL: Not HMS error - should always be false")
	} else {
		t.Logf("Testcase 4: Pass: Was false for non-HMSError")
	}
}

// Returns false if 'err' is not an HMSError, or, if it is, if 'class' doesn't
// match the HMSError's Class field (case insensitive).
func TestIsHMSErrorClassIgnCase(t *testing.T) {
	var regErr error
	if !IsHMSErrorClass(hmsError, hmsError.Class) {
		t.Errorf("Testcase 1: FAIL: Should have matched class")
	} else {
		t.Logf("Testcase 1: Pass: Class matched")
	}
	if !IsHMSErrorClassIgnCase(hmsError, strings.ToUpper(hmsError.Class)) {
		t.Errorf("Testcase 2: FAIL: Got false, should've ignored mixed case")
	} else {
		t.Logf("Testcase 2: Pass: true as expected, mismatched case ok")
	}
	if IsHMSErrorClassIgnCase(hmsError, "sadasdasdas") {
		t.Errorf("Testcase 3: FAIL: False true, Should've detected bad class")
	} else {
		t.Logf("Testcase 3: Pass: Bad class did not match")
	}
	if IsHMSErrorClassIgnCase(regErr, "") {
		t.Errorf("Testcase 4: FAIL: Not HMS error - should always be false")
	} else {
		t.Logf("Testcase 4: Pass: Was false for non-HMSError")
	}
}

// Test adding a ProblemDetails struct to an HMSError
func TestAddProblem(t *testing.T) {
	// Don't modify comparison copy, this makes a copy.
	copiedErr := hmsError.NewChild("")

	// Associate given problem with HMSError with no existing problem
	copiedErr.AddProblem(probDetailsGood)
	if err := testCheckHMSError(hmsErrorWProb, copiedErr); err != nil {
		t.Errorf("Testcase 1: FAIL: Modified HMSError != expected: %s", err)
	} else {
		t.Logf("Testcase 1: Pass: Modified HMSError == expected")
	}
}

// Test retrieving a ProblemDetails struct to an HMSError
func TestGetProblem(t *testing.T) {
	// HMSError has problem associated with it, should return it.
	p := hmsErrorWProb.GetProblem()
	if err := testCheckProblemDetails(probDetailsGood, p); err != nil {
		t.Errorf("Testcase 1: FAIL: ProblemDetails != expected: %s", err)
	} else {
		t.Logf("Testcase 1: Pass: Got expected ProblemDetails")
	}

	// No ProblemDetails associated wit hmsError, should return nil
	p = hmsError.GetProblem()
	if p != nil {
		t.Errorf("Testcase 1: FAIL: Problem != nil, should be nil.")
	} else {
		t.Logf("Testcase 1: Pass: ProblemDetails == nil as expected.")
	}
}

// Verify creating a new HMSError with same class as e, but a new message
// field if msg is non-empty.  If msg is "", it basically acts like a (deep)
// copy.
//
// DOES NOT COPY ProblemDetails if set, use NewChildWithProblem for that.
func TestNewChild2(t *testing.T) {
	// Test case 1:  New message for child.
	herr := hmsError.NewChild(hmsErrorNewMessage.Message)
	if err := testCheckHMSError(hmsErrorCopy, hmsError); err != nil {
		t.Errorf("Testcase 1a: FAIL: NewChild modified parent: %s", err)
	} else {
		t.Logf("Testcase 1a: Pass: NewChild did NOT modify parent.")
	}
	if err := testCheckHMSError(hmsErrorNewMessage, herr); err != nil {
		t.Errorf("Testcase 1b: FAIL: NewChild not as expected: %s", err)
	} else {
		t.Logf("Testcase 1b: Pass: NewChild matched expected result.")
	}
	// Testcase 2: No new message for child, should act like copy.
	herr = hmsError.NewChild("")
	if err := testCheckHMSError(hmsError, herr); err != nil {
		t.Errorf("Testcase 2: FAIL: NewChild not as expected: %s", err)
	} else {
		t.Logf("Testcase 2: Pass: NewChild matched expected result.")
	}
	// Testcase 3: Modify field, but should not carry over ProblemDetails
	herr = hmsErrorWProb.NewChild(hmsErrorNewMessage.Message)
	if err := testCheckHMSError(hmsErrorNewMessage, herr); err != nil {
		t.Errorf("Testcase 1: FAIL: NewChild not as expected: %s", err)
	} else {
		t.Logf("Testcase 1: Pass: NewChild matched expected result.")
	}
}

var hmsErrorWProbNewMsg = &HMSError{
	Class:   hmsErrorWProb.Class,
	Message: "New message used as details",
	Problem: probDetailsNewMsg,
}

var hmsErrorWProbNewInst = &HMSError{
	Class:   hmsErrorWProb.Class,
	Message: hmsErrorWProb.Message,
	Problem: probDetailsNewInst,
}

var hmsErrorWProbNewMsgInst = &HMSError{
	Class:   hmsErrorWProb.Class,
	Message: hmsErrorWProbNewMsg.Message,
	Problem: probDetailsNewMsgInst,
}

// Same as above, make sure not modfied unexpectedly.
var probDetailsNewMsg = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   "New message used as details",
	Instance: probDetailsGood.Instance,
	Status:   probDetailsGood.Status,
}

// Same as above, make sure not modfied unexpectedly.
var probDetailsNewInst = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   probDetailsGood.Detail,
	Instance: "New Instance",
	Status:   probDetailsGood.Status,
}

// Same as above, make sure not modfied unexpectedly.
var probDetailsNewMsgInst = &ProblemDetails{
	Type:     probDetailsGood.Type,
	Title:    probDetailsGood.Title,
	Detail:   "New message used as details",
	Instance: "New Instance",
	Status:   probDetailsGood.Status,
}

// This is the ProblemDetails aware-version of the above.  It allows the
// ProblemDetails Instance field to be updated if ProblemDetails are
// associated and the instance arg is non-"".  The msg field can also,
// if non-"" update the Detail field in ProblemDetails if they are present,
// in addition to working like msg in NewChild.
func TestNewChildWithProblem(t *testing.T) {
	// Test case 1:  New msg field set for child, no ProblemDetails are set.
	//     should expect same result as NewChild.
	herr := hmsError.NewChildWithProblem(
		hmsErrorNewMessage.Message, // msg
		"",                         // instance
	)
	if err := testCheckHMSError(hmsErrorCopy, hmsError); err != nil {
		t.Errorf("Testcase 1a: FAIL: Creating child modified parent: %s", err)
	} else {
		t.Logf("Testcase 1a: Pass: Creating child did NOT modify parent.")
	}
	if err := testCheckHMSError(hmsErrorNewMessage, herr); err != nil {
		t.Errorf("Testcase 1b: FAIL: Child != expected: %s", err)
	} else {
		t.Logf("Testcase 1b: Pass: Child matched parent but with new Message.")
	}
	// instance field should be ignored if no ProblemDetails
	herr = hmsError.NewChildWithProblem(
		hmsErrorNewMessage.Message,
		"instance- should be ignored",
	)
	if err := testCheckHMSError(hmsErrorCopy, hmsError); err != nil {
		t.Errorf("Testcase 1c: FAIL: Creating child modified parent: %s", err)
	} else {
		t.Logf("Testcase 1c: Pass: Creating child did NOT modify parent.")
	}
	if err := testCheckHMSError(hmsErrorNewMessage, herr); err != nil {
		t.Errorf("Testcase 1d: FAIL: instance w.o. prob caused change: %s", err)
	} else {
		t.Logf("Testcase 1d: Pass: instance was ignored since no Problem set.")
	}

	// Testcase 2: HMSError has ProblemDetails, but no msg or instance fields
	//     set for child, should act like (deep) copy.
	herr = hmsErrorWProb.NewChildWithProblem("", "")
	if err := testCheckHMSError(hmsErrorWProb, herr); err != nil {
		t.Errorf("Testcase 2: FAIL: Didn't act like copy: %s", err)
	} else {
		t.Logf("Testcase 2: Pass: Blank msg and instance did a copy.")
	}

	// Testcase 3: HMSError has ProblemDetails, msg field is set,
	//     and should modify both HMSError.Message and Problem's Details
	//     However, instance is unset and should remain the same.
	herr = hmsErrorWProb.NewChildWithProblem(
		hmsErrorWProbNewMsg.Message, //msg field
		"",                          // instance field
	)
	if err := testCheckHMSError(hmsErrorWProbNewMsg, herr); err != nil {
		t.Errorf("Testcase 3: FAIL: NewChildWithProblem != expected: %s", err)
	} else {
		t.Logf("Testcase 3: Pass: NewChild matched expected result.")
	}

	// Testcase 4: msg is blank, instance arg is given.  Copy should be
	//     identical except for Instance field in HMSError's ProblemDetails.
	herr = hmsErrorWProb.NewChildWithProblem(
		hmsErrorWProbNewMsg.Message, //msg field
		"",                          // instance field
	)
	if err := testCheckHMSError(hmsErrorWProbCopy, hmsErrorWProb); err != nil {
		t.Errorf("Testcase 4a: FAIL: Creating child modified parent: %s", err)
	} else {
		t.Logf("Testcase 4a: Pass: Creating child did NOT modify parent.")
	}
	if err := testCheckHMSError(hmsErrorWProbNewMsg, herr); err != nil {
		t.Errorf("Testcase 4b: FAIL: NewChildWithProblem != expected: %s", err)
	} else {
		t.Logf("Testcase 4b: Pass: NewChild matched expected result.")
	}

	// Testcase 5: msg is given, instance arg is given.  Copy should have
	//    updated Message field, and Detail (==msg) and Instance field in
	//    HMSError's ProblemDetails should also be updated.
	herr = hmsErrorWProb.NewChildWithProblem(
		hmsErrorWProbNewMsgInst.Message,          //msg field
		hmsErrorWProbNewMsgInst.Problem.Instance, // instance field
	)
	if err := testCheckHMSError(hmsErrorWProbCopy, hmsErrorWProb); err != nil {
		t.Errorf("Testcase 5a: FAIL: Creating child modified parent: %s", err)
	} else {
		t.Logf("Testcase 5c: Pass: Creating child did NOT modify parent.")
	}
	if err := testCheckHMSError(hmsErrorWProbNewMsgInst, herr); err != nil {
		t.Errorf("Testcase 5b: FAIL: NewChildWithProblem != expected: %s", err)
	} else {
		t.Logf("Testcase 5b: Pass: Both Message and Problem Detail/Instance ok.")
	}
}

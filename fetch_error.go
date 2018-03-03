package main

import "fmt"

// We define a custom error type so that we can provide friendlier error messages
type FetchError struct {
	errorCode int    // error code
	details   string // underlying error message, if any
	err       error  // underlying golang error, if any
}

// Implement the golang Error interface
func (e *FetchError) Error() string {
	var str string
	if str = getErrorMessage(e.errorCode, e.details) ; str == "" {
		str = fmt.Sprintf("%d - %s", e.errorCode, e.details)
	}

	return str
}

func newError(errorCode int, details string) *FetchError {
	return &FetchError{
		errorCode: errorCode,
		details:   details,
		err:       nil,
	}
}

func wrapError(err error) *FetchError {
	if err == nil {
		return nil
	}
	return &FetchError{
		errorCode: -1,
		details:   err.Error(),
		err:       err,
	}
}

func getErrorMessage(errorCode int, errorDetails string) string {
	switch errorCode {
	case INVALID_TAG_CONSTRAINT_EXPRESSION:
		return fmt.Sprintf(`
The --tag value you entered is not a valid constraint expression.
See https://github.com/opsgang/fetch#version-constraint-operators for examples.

Underlying error message:
%s
`, errorDetails)
	case INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED:
		return fmt.Sprintf(`
Received an HTTP 401 Response when attempting to query the repo for its tags.

Either your GitHub OAuth Token is invalid, or that you don't have access to
the repo with that token. Is the repo private?

Underlying error message:
%s
`, errorDetails)
	case REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED:
		return fmt.Sprintf(`
Received an HTTP 404 Response when attempting to query the repo for its tags.

Either the URL does not exist, or you don't have permission to access it.
If the repo is private, you will need to set GITHUB_TOKEN (or GITHUB_OAUTH_TOKEN)
in the env before invoking fetch.

Underlying error message:
%s
`, errorDetails)
	}

	return ""
}

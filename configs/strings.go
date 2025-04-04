// some constants you can return
package configs

import "fmt"

const (
	Oops              = "Oops"
	SomethingHappened = "Something happened, please try again later."
	PleaseTryLater    = "Please try again later."
	WelcomeBack       = "Welcome back"
	Null              = "null"
)

// /api/<params>
func API(s string) string {
	return fmt.Sprintf("/api/%s", s)
}

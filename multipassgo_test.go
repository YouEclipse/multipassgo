package multipassgo_test

import (
	"fmt"
	"github.com/YouEclipse/multipassgo"
	"testing"
	"time"
)

const secret = "4c191800e584533ab18fc533397fce1d"

// The customer information is represented as a hash which must contain at least the email address
// of the customer and a current timestamp (in ISO8601 encoding).
// You can also include the customer's first name, last name or several shipping addresses.
// Optionally, you can include an IP address of the customer's current browser session,
// that makes the token valid only for requests originating from this IP address.
// You can attribute tags to your customer by setting "tag_string" to a list of comma separated one-word values.
// These tags will override any tags that you may have already attributed to this customer.
// At Shopify, we use email addresses as unique identifiers for customers of a shop.

// If your site uses other identifiers (such as usernames),
// or if it is possible that two different users of your site registered with the same email address,
// you must set the unique identifier in the "identifier" field to avoid security problems.

// If the email address is always unique, you don't need to set the "identifier" field.

// If you want your users to see a specific page of your Shopify store, you can use the "return_to" field for that.
type ShopifyUserInfo struct {
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func TestTokenGenerate(t *testing.T) {
	var m = multipassgo.NewMultipass(secret)
	m.UserInfo = &ShopifyUserInfo{
		Email:     "bob@shopify.com",
		CreatedAt: time.Now().Format("2006-01-02T15:04:05-07:00"), //ISO8601
	}
	token, err := m.GenerateToken()
	if err != nil {
		t.Error("token generate failed: %v", err)
		return
	}
	fmt.Println(token)
}

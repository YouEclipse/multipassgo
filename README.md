multipassgo
============

 Golang implementation for  [Shopify Multipass](https://help.shopify.com/api/reference/multipass) 



## Installation
<pre>
    go get github.com/YouEclipse/multipassgo
</pre>

## Usage

The Multipass login feature is only available to Shopify Plus Customers. The secret can be found in your shop Admin (Settings > Checkout > Customer Accounts).
Make sure "Accounts are required" or "Accounts are optional" is selected and Multipass is enabled.

``` go

	const secret = "4c191800e584533ab18fc533397fce1d"

	type ShopifyUserInfo struct {
		Email      string `json:"email"`
		CreatedAt  string `json:"created_at"`
	}

	func main() { 
		var m = multipassgo.NewMultipass(secret)
		m.UserInfo = &ShopifyUserInfo{
			Email:      "chuiyouwu@gmail.com",
			CreatedAt:  time.Now().Format("2006-01-02T15:04:05-07:00"),
		}
		token, err := m.GenerateToken()
		if err != nil {
			fmt.Println("token generate failed: %v", err)
			return
		}
		fmt.Println(token)
	}

```
Once you have the token, you should trigger a HTTP GET request to your Shopify store.

```
http://yourstorename.myshopify.com/account/login/multipass/insert_token_here

```

The multipass token is only valid within a very short timeframe and each token can only be used once. For those reasons, you should not generate tokens in advance for rendering them into your HTML sites. You should create a redirect URL which generates tokens on-the-fly when needed and then automatically redirects the browser.
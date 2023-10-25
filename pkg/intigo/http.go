package intitools

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

const (
	ApiURL   = "https://app.intigriti.com/api"
	AppURL   = "https://app.intigriti.com"
	LoginURL = "https://login.intigriti.com"
)

type Client struct {
	ApiURL        string
	AppURL        string
	LoginURL      string
	apiKey        string
	Authenticated bool
	username      string
	password      string
	secret        string
	LastViewed    int64
	WebhookURL    string
	Ratelimiter   *rate.Limiter
	HTTPClient    *http.Client
	HttpCtx       context.Context
}

type ResponseState struct {
	Status              int    `json:"status"`
	Closereason         int    `json:"closeReason"`
	Duplicatesubmission string `json:"duplicateSubmission"`
}

type ResponsePayout struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type ResponseUser struct {
	Role     string `json:"role"`
	Email    string `json:"email"`
	Userid   string `json:"userId"`
	Avatarid string `json:"avatarId"`
	Username string `json:"userName"`
}

func NewClient(username string, password string, secret string, rl *rate.Limiter) *Client {

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	lastVisited := time.Now().UTC().Unix()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		ApiURL:     ApiURL,
		LoginURL:   LoginURL,
		AppURL:     AppURL,
		apiKey:     "",
		username:   username,
		password:   password,
		secret:     secret,
		LastViewed: lastVisited,
		HTTPClient: &http.Client{
			Timeout:   time.Minute,
			Jar:       jar,
			Transport: tr,
		},
		Authenticated: false,
		Ratelimiter:   rl,
	}
}

func (c *Client) Authenticate() error {

	// First request to get login page (and CSRF token / cookies)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/dashboard", c.AppURL), nil)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	finalPath := res.Request.URL.Path

	// If last redirect was to /researcher we are already logged in
	if finalPath != "/researcher" {
		// Parse HTML and find CSRF token and Return URL
		root, err := html.Parse(res.Body)
		if err != nil {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		csrfToken, err := c.getElementValue("__RequestVerificationToken", root)
		if err != nil {
			log.Fatal(err.Error())
		}

		returnURL, err := c.getElementValue("Input.ReturnUrl", root)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Prepare form for POST request
		form := url.Values{}
		form.Add("__RequestVerificationToken", csrfToken)
		form.Add("Input.ReturnUrl", returnURL)
		form.Add("Input.Email", c.username)
		form.Add("Input.LocalLogin", "True")
		form.Add("Input.Password", c.password)

		// Second request to submit username and password
		// We do not expect response body. Cookie is all we need (handled by CookieJar)
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/Account/Login", c.LoginURL), strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		// Check status
		if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		finalURL := res.Request.URL.String()

		resBody := res.Body
		// If last redirect was to /account/loginwith2fa we need a 2FA token
		if strings.Contains(finalURL, "/account/loginwith2fa") {
			if c.secret == "" {
				return fmt.Errorf("2FA is enabled but no secret is provided.")
			}

			// Parse HTML and find CSRF token and Return URL
			root, err := html.Parse(res.Body)
			if err != nil {
				return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
			}

			csrfToken, err := c.getElementValue("__RequestVerificationToken", root)
			if err != nil {
				log.Fatal(err.Error())
			}

			otpKey, err := totp.GenerateCode(strings.ToUpper(strings.Replace(c.secret, " ", "", -1)), time.Now())
			if err != nil {
				return err
			}

			// Prepare OTP form for POST request
			otpForm := url.Values{}
			otpForm.Add("__RequestVerificationToken", csrfToken)
			otpForm.Add("Input.TwoFactorAuthentication.VerificationCode", otpKey)

			req, err = http.NewRequest("POST", finalURL, strings.NewReader(otpForm.Encode()))
			if err != nil {
				return err
			}

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			res, err := c.HTTPClient.Do(req)
			if err != nil {
				return err
			}

			defer res.Body.Close()

			// Check status
			if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
				return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
			}

			finalURL := res.Request.URL.Path

			// If last redirect was not to /researcher/ the 2FA secret failed to authenticate
			if finalURL[len(finalURL)-10:] != "/authorize" {
				return fmt.Errorf("Failed to authenticate with 2FA")
			}
			resBody = res.Body
		}

		// Parse HTML and find code, state etc.
		root, err = html.Parse(resBody)
		if err != nil {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		code, err := c.getElementValue("code", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		scope, err := c.getElementValue("scope", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		state, err := c.getElementValue("state", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		session_state, err := c.getElementValue("session_state", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		iss, err := c.getElementValue("iss", root)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Prepare form for POST request
		form = url.Values{}
		form.Add("code", code)
		form.Add("scope", scope)
		form.Add("state", state)
		form.Add("session_state", session_state)
		form.Add("iss", iss)

		// Another request to /signin-oidc
		// We do not expect response body. Cookie is all we need (handled by CookieJar)
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/signin-oidc", c.AppURL), strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err = c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		// Check status
		if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		// Right now we should be redirected to second /connect/authorize and have another set of code, state etc.

		// Parse HTML and find code, state etc.
		root, err = html.Parse(res.Body)
		if err != nil {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		code, err = c.getElementValue("code", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		scope, err = c.getElementValue("scope", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		state, err = c.getElementValue("state", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		session_state, err = c.getElementValue("session_state", root)
		if err != nil {
			log.Fatal(err.Error())
		}
		iss, err = c.getElementValue("iss", root)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Prepare form for POST request
		form = url.Values{}
		form.Add("code", code)
		form.Add("scope", scope)
		form.Add("state", state)
		form.Add("session_state", session_state)
		form.Add("iss", iss)

		// Another request to /signin-oidc-researcher
		// We do not expect response body. Cookie is all we need (handled by CookieJar)
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/signin-oidc-researcher", c.AppURL), strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err = c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		// Check status
		if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		}

		log.Println("Client authenticated")

	}

	// Third request to get API token
	// req4, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/token", c.AppURL), nil)
	// if err != nil {
	// 	return err
	// }

	// res4, err := c.HTTPClient.Do(req4)
	// if err != nil {
	// 	return err
	// }

	// defer res4.Body.Close()

	// if res4.StatusCode < http.StatusOK || res4.StatusCode >= http.StatusBadRequest {
	// 	return fmt.Errorf("unknown error, status code: %d", res4.StatusCode)
	// }

	// // Parse response to get API Token
	// apiToken, err := io.ReadAll(res4.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //c.apiKey = string(apiToken[1 : len(apiToken)-1])
	c.Authenticated = true

	return nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {

	if !c.Authenticated {
		c.Authenticate()
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

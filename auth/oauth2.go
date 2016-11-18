package auth

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func Login() *http.Client {
	c := make(chan string)
	state := getRandomState(18)
	go startRedirectServer(c, state)
	port := <-c

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "frsyvri19q4rtqvkpamgyjexu8zlkaas",
		ClientSecret: "bPzkgJwOb4JgtaJ35gqfuxgCvm387GqT",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://account.box.com/api/oauth2/authorize",
			TokenURL: "https://app.box.com/api/oauth2/token",
		},
		RedirectURL: "http://localhost:" + port,
	}

	//url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	url := getDirectShibAuthCodeURL(conf, state)

	err := openURLInBrowser(url)
	if err != nil {
		fmt.Println("Visit this URL to log in:")
		fmt.Println(url)
	}

	code := <-c

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login successful")

	return conf.Client(ctx, tok)
}

func getDirectShibAuthCodeURL(conf *oauth2.Config, state string) string {
	idpURN := "urn:mace:incommon:uiuc.edu"
	targetResource := fmt.Sprintf("https://www.box.com/api/oauth2/authorize"+
		"?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		conf.ClientID, url.QueryEscape(conf.RedirectURL), state)
	return fmt.Sprintf("https://sso.services.box.net/sp/startSSO.ping?PartnerIdpId=%s&TargetResource=%s",
		url.QueryEscape(idpURN), url.QueryEscape(targetResource))
}

func openURLInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows", "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	fmt.Println("Complete the login process in your browser.")
	fmt.Println("If the page did not open automatically, visit this URL to log in:")
	fmt.Println(url)
	return nil
}

func startRedirectServer(c chan<- string, expectedState string) {
	listener, _ := net.Listen("tcp", ":0")
	defer listener.Close()

	_, port, _ := net.SplitHostPort(listener.Addr().String())
	c <- port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		query := r.URL.Query()
		actualState, ok := query["state"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing param \"state\"")
			return
		} else if actualState[0] != expectedState {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid state")
			return
		}

		if _, ok := query["error"]; ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Access denied: %v", query["error_description"][0])
			return
		}

		code, ok := query["code"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing param \"code\"")
			return
		}

		fmt.Fprintf(w, "Login successful -- you may now close this tab.")
		c <- code[0]
		listener.Close()
	})

	http.Serve(listener, nil)
}

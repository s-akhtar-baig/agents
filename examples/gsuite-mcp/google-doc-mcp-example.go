package main

import (
	"crypto/rand"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

type AuthCodeReceiver struct{
	server *httptest.Server
	code   chan string
	path   string
	state  string
}

func newAuthCodeReceiver() *AuthCodeReceiver{
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed generating oauth state: %v", err)
	}
	state := base64.RawURLEncoding.EncodeToString(b)
	return &AuthCodeReceiver{
		code: make(chan string),
		path: "/submit",
		state: state,
	}
}

func (o *AuthCodeReceiver) start() {
	mux := http.NewServeMux()
	mux.HandleFunc(o.path, func(w http.ResponseWriter, r *http.Request) {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		if e := params.Get("error"); e != "" {
			http.Error(w, "Authorization failed: "+e, http.StatusBadRequest)
			return
		}
		if s := params.Get("state"); s == "" || s != o.state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}
		if code := params.Get("code"); code != "" {
			o.code <- code
			fmt.Fprint(w, "Return to example!")
			return
		}
		http.Error(w, "Missing code", http.StatusBadRequest)
	})	
	o.server = httptest.NewServer(mux)
}

func (r *AuthCodeReceiver) url() string {
	return r.server.URL + r.path
}

func (r *AuthCodeReceiver) getCode() string {
	code := <-r.code
	return code
}

func (r *AuthCodeReceiver) stop() {
	r.server.Close()
}

func main() {
	codeReceiver := newAuthCodeReceiver()
	codeReceiver.start()
	defer codeReceiver.stop()
	clientId := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{
			"https://www.googleapis.com/auth/documents",
			"https://www.googleapis.com/auth/documents.readonly",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: codeReceiver.url(),
	}

	verifier := oauth2.GenerateVerifier()
	authURL := conf.AuthCodeURL(
		codeReceiver.state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	fmt.Printf("Visit the following URL to enable access:\n%s\n", authURL)

	code := codeReceiver.getCode()
	token, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}
	model := os.Getenv("INFERENCE_MODEL")
	if model == "" {
		model = openai.ChatModelGPT4o
	}
	log.Printf("Using model %s", model)
	client := openai.NewClient()
	params := responses.ResponseNewParams{
		Model:           model,
		Tools:           []responses.ToolUnionParam{
			{
				OfMcp: &responses.ToolMcpParam{
					ServerLabel: "google docs",
					ServerURL:   "http://localhost:9876/mcp",
					Headers:     map[string]string{
						"Authorization": "Bearer " + token.AccessToken,
					},
					AllowedTools: responses.ToolMcpAllowedToolsUnionParam{
						OfMcpAllowedTools: []string{
							"create_doc",
						},
					},
				},
			},
		},
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Using the tool provided, create a document containing a poem and return the link to it"),
		},
	}
	resp, err := client.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(resp.OutputText())
	
}

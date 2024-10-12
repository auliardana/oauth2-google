package main

import (
	"context"
	"encoding/json"
	// "html/template"
	// "log"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuth handler create an unique OAuth URL using the client ID and client secret.
// then it redirect the user to the OAuth provider website to complete the login.
func (a *App) oAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.config.AuthCodeURL("hello world", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuth callback handler handle the redirect request from the OAuth provider.
// It read the code query parameter and exchange it to get the access token.
// Then this handler call user info endpoint to get the user public detail eg.,
// Name, Email, Profile picture etc.
func (a *App) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")

    // Exchanging the code for an access token
    t, err := a.config.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Creating an HTTP client to make authenticated request using the access key.
    // This client method also regenerates the access key using the refresh key.
    client := a.config.Client(context.Background(), t)

    // Getting the user public details from Google API endpoint
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Closing the request body when this function returns.
    // This is a good practice to avoid memory leak
    defer resp.Body.Close()

    var v any

    // Reading the JSON body using JSON decoder
    err = json.NewDecoder(resp.Body).Decode(&v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the user data to JSON format
    responseData, err := json.Marshal(v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type to application/json and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseData)
}

package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "golang.org/x/oauth2/google"
    "google.golang.org/api/idtoken"
    vault "github.com/hashicorp/vault/api"
)

var (
    vaultAddr = os.Getenv("VAULT_ADDR")
    role      = os.Getenv("VAULT_ROLE")
)

type VaultLoginRequest struct {
    Role string `json:"role"`
    JWT  string `json:"jwt"`
}

func getVaultClient() (*vault.Client, error) {
    config := vault.DefaultConfig()
    config.Address = vaultAddr

    client, err := vault.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create Vault client: %w", err)
    }
    return client, nil
}

func getVaultToken() (string, error) {
    ctx := context.Background()
    tokenSource, err := google.DefaultTokenSource(ctx, idtoken.EndpointAudience(vaultAddr))
    if err != nil {
        return "", fmt.Errorf("failed to get default token source: %w", err)
    }

    token, err := tokenSource.Token()
    if err != nil {
        return "", fmt.Errorf("failed to get token: %w", err)
    }

    payload := VaultLoginRequest{
        Role: role,
        JWT:  token.AccessToken,
    }

    client, err := getVaultClient()
    if err != nil {
        return "", err
    }

    path := "auth/gcp/login"
    response, err := client.Logical().Write(path, map[string]interface{}{
        "role": role,
        "jwt":  payload.JWT,
    })
    if err != nil {
        return "", fmt.Errorf("failed to login to Vault: %w", err)
    }

    if response.Auth == nil {
        return "", fmt.Errorf("no authentication information in response")
    }

    return response.Auth.ClientToken, nil
}

func YourCloudFunction(w http.ResponseWriter, r *http.Request) {
    vaultToken, err := getVaultToken()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
    }

    client, err := getVaultClient()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating Vault client: %v", err), http.StatusInternalServerError)
        return
    }

    client.SetToken(vaultToken)
    secret, err := client.Logical().Read("secret/data/access")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error reading secret: %v", err), http.StatusInternalServerError)
        return
    }

    if secret == nil {
        http.Error(w, "No secret found at the specified path", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(secret.Data["data"])
}

func main() {
    http.HandleFunc("/", YourCloudFunction)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

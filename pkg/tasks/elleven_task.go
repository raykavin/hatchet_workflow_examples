package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
)

type PersonDataInput struct {
	DocumentNumber string `json:"document_number"`
}

type PersonDataOutput struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	CPFCNPJ     string `json:"cpf_cnpj,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	BirthDate   string `json:"birth_date,omitempty"`
}

type ssoTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// EllevenPersonInfo is a Hatchet task that retrieves person information from the Elleven ERP using PRest
func EllevenPersonInfo(ctx hatchet.Context, input PersonDataInput) (PersonDataOutput, error) {
	log.Printf("Fetching person information from ERP using document number: %s", input.DocumentNumber)

	// Validate input
	if input.DocumentNumber == "" {
		return PersonDataOutput{}, fmt.Errorf("document number is required")
	}

	// Get configuration
	personInfoURL := getEnvOrError("ELLEVEN_PREST_GET_PERSON_INFO_URL")

	// Retrieve SSO token
	token, err := getSSOToken(ctx)
	if err != nil {
		return PersonDataOutput{}, fmt.Errorf("failed to obtain SSO token: %w", err)
	}

	// Prepare headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	// Prest query parameters
	queryParams := map[string]string{"tx_id": input.DocumentNumber}

	// Call Prest API
	response, statusCode, err := makeRequest(
		ctx,
		personInfoURL,
		http.MethodGet,
		queryParams,
		headers,
		nil,
	)
	if err != nil {
		return PersonDataOutput{}, fmt.Errorf("request to Prest failed: %w", err)
	}

	if statusCode != http.StatusOK {
		return PersonDataOutput{}, fmt.Errorf("Prest returned status code %d", statusCode)
	}

	// Parse response
	var outputs []PersonDataOutput
	if err := json.Unmarshal(response, &outputs); err != nil {
		return PersonDataOutput{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(outputs) == 0 {
		return PersonDataOutput{}, fmt.Errorf("no person information found for document number %s", input.DocumentNumber)
	}

	return outputs[0], nil
}

func getSSOToken(ctx context.Context) (string, error) {
	log.Print("Requesting SSO authorization token...")

	// Load SSO configuration
	ssoURL := getEnvOrError("SSO_AUTHORIZATION_URL")
	clientID := getEnvOrError("SSO_CLIENT_ID")
	clientSecret := getEnvOrError("SSO_CLIENT_SECRET")
	grantType := getEnvOrError("SSO_GRANT_TYPE")
	scope := getEnvOrError("SSO_SCOPE")

	// Build form content data
	formContent := url.Values{}
	formContent.Set("grant_type", grantType)
	formContent.Set("scope", scope)
	formContent.Set("client_id", clientID)
	formContent.Set("client_secret", clientSecret)

	// Create quest headers
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	// Execute SSO request
	response, statusCode, err := makeRequest(
		ctx,
		ssoURL,
		http.MethodPost,
		nil,
		headers,
		[]byte(formContent.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("sso request failed: %w", err)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("sso returned status code %d", statusCode)
	}

	// Decode response
	var tokenResp ssoTokenResponse
	if err := json.Unmarshal(response, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode SSO response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("sso response does not contain an access token")
	}

	return tokenResp.AccessToken, nil
}

func makeRequest(
	ctx context.Context,
	urlStr string,
	method string,
	queryParams map[string]string,
	headers map[string]string,
	body []byte,
) ([]byte, int, error) {
	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid URL: %w", err)
	}

	// Append query parameters
	if len(queryParams) > 0 {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Create request
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	fmt.Println(u.String())

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

func getEnvOrError(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is not set", key))
	}
	return value
}

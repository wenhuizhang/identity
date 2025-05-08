// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/agntcy/identity/pkg/log"
)

// ------------------------ GLOBAL -------------------- //

// Timeout : API timeout time
const Timeout = 5

// Authorization : Authorization header
const Authorization = "Authorization"

// ------------------------ GLOBAL -------------------- //

// PostData : Post data without authenticated
func PostData(uri string, payload interface{}, headers map[string]string) interface{} {
	return postDataWithHeaders(uri, payload, headers)
}

// GetDataAuth : Get data authenticated
func GetDataAuth(ctx context.Context, uri, accessToken string, result interface{}) {
	getDataWithHeaders(ctx, uri, authHeaders(accessToken), result)
}

func GetRawDataWithHeaders(
	ctx context.Context,
	uri string,
	headers map[string]string,
	result *string,
) error {
	// Create context
	ctx, cancel := context.WithTimeout(ctx, Timeout*time.Second)
	defer cancel()

	// Create a new request using http
	log.Debug("Getting uri ", uri)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)

	// Add headers
	log.Debug("Getting with headers ", headers)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Got error", err)
		return err
	}

	defer resp.Body.Close()

	// Read and parse data
	body, _ := io.ReadAll(resp.Body)

	log.Debug("Server response status code is ", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code = %d", resp.StatusCode)
	}

	*result = string(body)

	return nil
}

func getDataWithHeaders(
	ctx context.Context,
	uri string,
	headers map[string]string,
	result interface{},
) {
	// Create context
	ctx, cancel := context.WithTimeout(ctx, Timeout*time.Second)
	defer cancel()

	// Create a new request using http
	log.Debug("Getting uri ", uri)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)

	// Add headers
	log.Debug("Getting with headers ", headers)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Got error", err)
		return
	}

	// Read and parse data
	body, _ := io.ReadAll(resp.Body)

	log.Debug("Server response status code is ", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return
	}

	jsonErr := json.Unmarshal(body, &result)

	if jsonErr != nil {
		log.Debug("Got error ", jsonErr)
		return
	}

	_ = resp.Body.Close()
}

func postDataWithHeaders(uri string, payload interface{}, headers map[string]string) interface{} {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	// Create a new request using http
	log.Debug("Posting to uri ", uri)

	rawPayload, payloadErr := json.Marshal(payload)

	log.Debug("Payload is ", string(rawPayload))

	if payloadErr != nil {
		return nil
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(rawPayload))

	// Add headers
	log.Debug("Posting with headers ", headers)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send req using http Client
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Got error", err)
		return nil
	}

	// Read and parse data
	body, _ := io.ReadAll(resp.Body)

	var data interface{}

	log.Debug("Server response body is ", string(body))
	log.Debug("Server response status code is ", resp.StatusCode)

	_ = json.Unmarshal(body, &data)
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	return data
}

func authHeaders(accessToken string) map[string]string {
	headers := make(map[string]string)
	headers[Authorization] = "Bearer " + accessToken
	headers["Content-Type"] = "application/json"

	return headers
}

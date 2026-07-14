package documentcallbackattempts

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
	docrepo "go-document-generator/internal/repository/documents"
	cbrepo "go-document-generator/internal/repository/documentcallbackattempts"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/pagination"
)

type TestCallbackInput struct {
	CallbackURL   string
	SamplePayload map[string]any
}

type TestCallbackResult struct {
	Success            bool
	ResponseStatusCode int
	ErrorMessage       string
}

type Service interface {
	ListByDocumentID(ctx context.Context, documentID int64, page pagination.Params) ([]cbEntity.CallbackAttempt, pagination.Meta, error)
	TestCallback(ctx context.Context, in TestCallbackInput) (TestCallbackResult, error)
}

type service struct {
	attempts   cbrepo.DocumentCallbackAttemptsRepository
	docs       docrepo.DocumentsRepository
	client     *http.Client
	hmacSecret string
}

func NewService(attempts cbrepo.DocumentCallbackAttemptsRepository, docs docrepo.DocumentsRepository, hmacSecret string) Service {
	return &service{
		attempts:   attempts,
		docs:       docs,
		client:     &http.Client{Timeout: 15 * time.Second},
		hmacSecret: hmacSecret,
	}
}

func (s *service) ListByDocumentID(ctx context.Context, documentID int64, page pagination.Params) ([]cbEntity.CallbackAttempt, pagination.Meta, error) {
	if _, err := s.docs.GetByID(ctx, nil, documentID, nil); err != nil {
		return nil, pagination.Meta{}, mapRepoErr(err)
	}
	page = pagination.Normalize(page.Page, page.Limit)
	items, total, err := s.attempts.ListByDocumentID(ctx, nil, documentID, page)
	if err != nil {
		return nil, pagination.Meta{}, err
	}
	return items, pagination.Meta{Page: page.Page, Limit: page.Limit, Total: total}, nil
}

func (s *service) TestCallback(ctx context.Context, in TestCallbackInput) (TestCallbackResult, error) {
	if in.CallbackURL == "" {
		return TestCallbackResult{}, apperror.ErrInvalidInput
	}
	body, _ := json.Marshal(in.SamplePayload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, in.CallbackURL, bytes.NewReader(body))
	if err != nil {
		return TestCallbackResult{Success: false, ErrorMessage: err.Error()}, nil
	}
	req.Header.Set("Content-Type", "application/json")
	if s.hmacSecret != "" {
		req.Header.Set("X-Document-Signature", computeHMAC(body, s.hmacSecret))
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return TestCallbackResult{Success: false, ErrorMessage: err.Error()}, nil
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300
	result := TestCallbackResult{
		Success:            ok,
		ResponseStatusCode: resp.StatusCode,
	}
	if !ok {
		result.ErrorMessage = resp.Status
	}
	return result, nil
}

func computeHMAC(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "hmac-sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	return err
}

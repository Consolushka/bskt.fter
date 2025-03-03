package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/log"
	"net/url"
)

type SearchPlayerByFullNameRequest struct {
	custom_request.BaseRequest
	fullName string
}

func (s *SearchPlayerByFullNameRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{
		s.parseFullName,
	}
}

func (s *SearchPlayerByFullNameRequest) FullName() string {
	return s.fullName
}

func (s *SearchPlayerByFullNameRequest) parseFullName(storage custom_request.CustomRequestStorage) error {
	encodedName := storage.GetQueryParam("fullName")
	decodedName, err := url.QueryUnescape(encodedName)
	if err != nil {
		return err
	}
	s.fullName = decodedName
	log.Info(s.fullName)
	return nil
}

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type AuthService struct {
	serviceUrl string
}

func NewAuthService(serviceUrl string) *AuthService {
	return &AuthService{
		serviceUrl: serviceUrl,
	}
}

func (s *AuthService) Login(email string, password string) (string, error) {
	values := map[string]string{"email": email, "password": password}
	json_data, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(fmt.Sprintf("%s/login", s.serviceUrl), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error: %s", res["error"].(string))
	}

	return res["token"].(string), nil
}

func (s *AuthService) Validate(token string) error {
	req, err := http.NewRequest("POST", s.serviceUrl+"/validate", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		if value, ok := res["error"]; ok {
			return fmt.Errorf("error: %s", value.(string))
		} else {
			return errors.New("invalid token")
		}
	}

	return nil
}

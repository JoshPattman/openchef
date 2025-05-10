package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func PostAndReciveJson(url string, send any, recv any) error {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(send)
	if err != nil {
		return errors.Join(errors.New("failed to encode body"), err)
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return errors.Join(errors.New("failed to create request"), err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Join(errors.New("failed to do request"), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&recv)
		if err != nil {
			return errors.Join(errors.New("failed to decode body"), err)
		}
		return nil
	} else {
		bs, err := io.ReadAll(body)
		if err != nil {
			return errors.Join(errors.New("failed to decode body after non0ok status"), err)
		}
		return fmt.Errorf("api returned non-ok status code %v with data %v", resp.Status, string(bs))
	}
}

package worldlogo

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	api "ext-data-domain/internal/server/webapi/api/openapi"
)

func WriteDataFromCSVToAPI(filename, apiAddr string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	// Name Key Src
	var record []string
	i := 0
	for {
		record, err = reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		// download image
		img, err := downloadImage(record[2])
		if err != nil {
			return fmt.Errorf("downloadImage: %w", err)
		}
		imgStr := base64.StdEncoding.EncodeToString(img)
		// send to api
		err = sendToAPI(apiAddr, record[0], record[1], imgStr)
		if err != nil {
			return fmt.Errorf("sendToAPI: %w", err)
		}
		if i == 3 {
			break
		}
	}
	return nil
}

func downloadImage(src string) ([]byte, error) {
	// download image
	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code err %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func sendToAPI(apiAddr, name, key, img string) error {
	rec := api.WorldLogoInput{
		Name:          name,
		SrcKey:        key,
		LogoBase64Str: &img,
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	// send to api
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/world-logo/", apiAddr), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "myapikey")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return fmt.Errorf("response code err %d", resp.StatusCode)
	}

	return nil
}

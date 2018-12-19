package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type (
	Request struct {
		ID           int           `json:"id"`
		Hash         string        `json:"hash"`
		Packed       string        `json:"packed"`
		Iv           string        `json:"iv"`
		P2           string        `json:"p2"`
		P1           string        `json:"p1"`
		EncResp      bool          `json:"enc_resp"`
		Method       string        `json:"method"`
		NonEncParams string        `json:"non_enc_params"`
		Params       []interface{} `json:"params"`
	}

	Response struct {
		Result struct {
			Cp        []string      `json:"cp"`
			Df        []interface{} `json:"df"`
			H         string        `json:"h"`
			IPAddr    string        `json:"ip_addr"`
			MultiConn bool          `json:"multi_conn"`
			S         string        `json:"s"`
			T         string        `json:"t"`
			Vc        int           `json:"vc"`
		} `json:"result"`
		Error interface{} `json:"error"`
		ID    int         `json:"id"`
	}
)

var (
	url       = "http://ad.cafebazaar.ir/json/getAppDownloadInfo"
	userAgent = "Dalvik/1.6.0 (Linux; U; Android 4.4.2; SAMSUNG-SM-N900A Build/KOT49H)"
)

func main() {
	// >> put package name here
	pkgName := ""

	resp, err := makeRequest(pkgName)
	if err != nil {
		log.Fatalln(err)
	}

	link := getLink(resp)
	fmt.Println(link)
}

// parse response to generate link
func getLink(resp *Response) string {
	return fmt.Sprintf("%sapks/%s.apk?rand=%d", resp.Result.Cp[0], resp.Result.T, currentMillis())
}

func makeRequest(pkgName string) (*Response, error) {
	body := &Request{
		ID:           1,
		Hash:         hashed(pkgName),
		Packed:       "xzrBQdWmJqg/BQN+4Ll+XCuNIhYwIpWmFRH+I1wjEKfb2NwtXaU4OO6LmDY+dcNKPh6v1a2GdLYcCdZ6NliD0nbYjcglOT7OYB9fefCL5Ec=",
		Iv:           "UFDpSQCua3LwOKb8QWW4dS2PNSfMQ3ua1eWAuJY1G8xcaTS+Md+gbGMCSG3C5QJLmoiSFyOv/QRFv6hWYsrA31ji0fGhWNGiqY9sWltqBst7YKoCqPLG0fCjoPKWPhvVhxKhjO8yT3RPalmDuPKpqGwW2fdHH+xPnuCDU51uUaE=",
		P2:           "r7oshN8AYo64PZDDlJg8TmiEiXrrBjKlwPQITF94s/3tKsyB1PJRJM5cD/JZBEHK/wWvGb/jyj0GrOgbEMONHBoLCMR/X6RWeC59LaItQaDk/uY3+2cEisuBw3VCAkKL887SebW0xmB/16rNl3LxLL5/vgCZ4jaUvIb1dj0JEH4=",
		P1:           "Kvn/n9BLGkFAcpAWBQsAVbcF8SVnS6f3XGulLM/J6a3SQOS5q8CagfCm2zbzQxHT0kRb9z90eCIBP9huKDth0Mu9JaAuNn9SiV7pBTs6C3hVlolY41W93hKPwhBfNyWCATymDnSjqcX/KKNcKn3fvMU7zR0w9h/WM/sUkccX8pg=",
		EncResp:      false,
		Method:       "getAppDownloadInfo",
		NonEncParams: "{\"device\":{\"mn\":260,\"abi\":\"x86\",\"sd\":19,\"bv\":\"7.12.2\",\"us\":{},\"cid\":0,\"lac\":0,\"ct\":\"\",\"id\":\"YGrrXv9TQkGyRwo6GaU0kw\",\"dd\":\"hlteatt\",\"co\":\"\",\"mc\":310,\"dm\":\"samsung\",\"do\":\"SAMSUNG-SM-N900A\",\"dpi\":160,\"abi2\":\"armeabi-v7a\",\"sz\":\"l\",\"dp\":\"hlteuc\",\"bc\":701202,\"pr\":\"\"},\"referer\":{\"name\":\"page_home|!EX!PaidRowTest|!VA!empty_key|referrer_slug=home|row-2-Best New Updates|3|test_group=A|not_initiated\"}}",
		Params:       []interface{}{},
	}

	// marshal body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// prepare request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Connection", "Keep-Alive")

	// send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// close body
	defer resp.Body.Close()

	var response Response
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}

// generate hash for package
func hashed(pkg string) string {
	hash := fmt.Sprintf(`{"7cc78271-e338-4edc-849c-b105c5d51ba5":["getAppDownloadInfo","%s",19]}`, pkg)
	h := sha1.New()
	io.WriteString(h, hash)
	return hex.EncodeToString(h.Sum(nil))
}

func currentMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

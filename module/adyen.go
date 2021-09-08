package module

import (
	"encoding/json"
	"fmt"
	"github.com/CrimsonAIO/adyen"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"time"
)

type AdyenCardNumber struct {
	Activate            string `json:"activate"`
	DFValue             string `json:"dfValue"`
	Generationtime      string `json:"generationtime"`
	InitializeCount     string `json:"initializeCount"`
	LuhnCount           string `json:"luhnCount"`
	LuhnOkCount         string `json:"luhnOkCount"`
	LuhnSameLengthCount string `json:"luhnSameLengthCount"`
	Number              string `json:"number"`
}

type AdyenCardExpirationMonth struct {
	Activate        string `json:"activate"`
	DFValue         string `json:"dfValue"`
	ExpiryMonth     string `json:"expiryMonth"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}

type AdyenCardExpirationYear struct {
	Activate        string `json:"activate"`
	DFValue         string `json:"dfValue"`
	ExpiryYear      string `json:"expiryYear"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}

type AdyenCardCVV struct {
	Activate        string `json:"activate"`
	Cvc             string `json:"cvc"`
	DFValue         string `json:"dfValue"`
	Generationtime  string `json:"generationtime"`
	InitializeCount string `json:"initializeCount"`
}

type PaymentBody struct {
	SID                   string       `json:"sid,omitempty"`
	PreferredLanguage     string       `json:"preferredLanguage"`
	TermsAndCondition     bool         `json:"termsAndCondition"`
	DeviceID              string       `json:"deviceId"`
	CartID                string       `json:"cartId"`
	EncryptedCardNumber   string       `json:"encryptedCardNumber,omitempty"`
	EncryptedExpiryMonth  string       `json:"encryptedExpiryMonth,omitempty"`
	EncryptedExpiryYear   string       `json:"encryptedExpiryYear,omitempty"`
	EncryptedSecurityCode string       `json:"encryptedSecurityCode,omitempty"`
	PaymentMethod         string       `json:"paymentMethod"`
	ReturnURL             string       `json:"returnUrl,omitempty"`
	BrowserInfo           *BrowserInfo `json:"browserInfo"`
}

type BrowserInfo struct {
	ScreenWidth    int    `json:"screenWidth"`
	ScreenHeight   int    `json:"screenHeight"`
	ColorDepth     int    `json:"colorDepth"`
	UserAgent      string `json:"userAgent"`
	TimeZoneOffset int    `json:"timeZoneOffset"`
	Language       string `json:"language"`
	JavaEnabled    bool   `json:"javaEnabled"`
}

func (tk *Task) AdyenConfirm() {
	adyenClient, err := adyen.NewClient(sites[tk.Site].adyenKey)
	if err != nil {
		panic(err)
	}

	card, err := adyenClient.Encrypt(adyen.Version118, map[string]interface{}{
		"activate":"1",
		"dfValue":"ryEGX8eZpJ0030000000000000LOziC3ZM670050271576cVB94iKzBGzk6emGsPvH5S16Goh5Mk0045zgp4q8JSa00000qZkTE00000PRbZ1HbvOQ1B2M2Y8Asg:40",
		"generationtime":time.Now().Format("2006-01-02T15:04:05.0000Z"),
		"initializeCount":"1",
		"luhnCount":"1",
		"luhnOkCount":"1",
		"luhnSameLengthCount":"1",
		"number":tk.Data.Profile.Billing.Number[:4]+" "+tk.Data.Profile.Billing.Number[4:8]+" "+tk.Data.Profile.Billing.Number[8:12]+" "+tk.Data.Profile.Billing.Number[12:],
	}, adyen.GenerationTimeNow)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "problem with ayden" + err.Error())
		tk.Stop()
		return
	}

	month, err := adyenClient.Encrypt(adyen.Version118, map[string]interface{}{
		"activate":"1",
		"dfValue":"ryEGX8eZpJ0030000000000000LOziC3ZM670050271576cVB94iKzBGzk6emGsPvH5S16Goh5Mk0045zgp4q8JSa00000qZkTE00000PRbZ1HbvOQ1B2M2Y8Asg:40",
		"expiryMonth":tk.Data.Profile.Billing.ExpirationMonth,
		"generationtime":time.Now().Format("2006-01-02T15:04:05.0000Z"),
		"initializeCount":"1",
	}, adyen.GenerationTimeNow)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "problem with ayden" + err.Error())
		tk.Stop()
		return
	}

	year, err := adyenClient.Encrypt(adyen.Version118, map[string]interface{}{
		"activate":"1",
		"dfValue":"ryEGX8eZpJ0030000000000000LOziC3ZM670050271576cVB94iKzBGzk6emGsPvH5S16Goh5Mk0045zgp4q8JSa00000qZkTE00000PRbZ1HbvOQ1B2M2Y8Asg:40",
		"expiryYear":"20"+tk.Data.Profile.Billing.ExpirationYear,
		"generationtime":time.Now().Format("2006-01-02T15:04:05.0000Z"),
		"initializeCount":"1",
	}, adyen.GenerationTimeNow)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "problem with ayden" + err.Error())
		tk.Stop()
		return
	}

	cvv, err := adyenClient.Encrypt(adyen.Version118, map[string]interface{}{
		"activate":"1",
		"cvc":tk.Data.Profile.Billing.CVV,
		"dfValue":"ryEGX8eZpJ0030000000000000LOziC3ZM670050271576cVB94iKzBGzk6emGsPvH5S16Goh5Mk0045zgp4q8JSa00000qZkTE00000PRbZ1HbvOQ1B2M2Y8Asg:40",
		"generationtime":time.Now().Format("2006-01-02T15:04:05.0000Z"),
		"initializeCount":"1",
	}, adyen.GenerationTimeNow)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "problem with ayden" + err.Error())
		tk.Stop()
		return
	}


	body, err := json.Marshal(PaymentBody{
		// payment information
		EncryptedCardNumber:   card,
		EncryptedExpiryMonth:  month,
		EncryptedExpiryYear:   year,
		EncryptedSecurityCode: cvv,

		// device information and other stuff
		PreferredLanguage: "en",
		PaymentMethod:     "CREDITCARD",
		ReturnURL:         fmt.Sprintf("https://www.%s.com/adyen/checkout", sites[tk.Site].name),
		CartID:            string(tk.FastClient.Jar.PeekValue("cart-guid")),
		DeviceID:          `0400JapG4txqVP4Nf94lis1ztioT9A1DShgAnrp/XmcfWoVVgr+Rt2dAZPhMS97Z4yfjSLOS3mruQCzk1eXuO7gGCUfgUZuLE2xCJiDbCfVZTGBk19tyNs7g87Mc/hl2WkFqp/uhGlbxNVvph/T+lMlRTBNagwJDR5g5DJ0qCc4gMpergecSo06izVzqeMmBCH/i9cmKrLQxcxA5OE2KkOZe/0jXzk77ILZ/eUsQ7RNrLro1kTKIs1496YkpIh3A707lm2e25SQbo1MiBCDKvxfoGxWoFrUkg1NT6ApRBNlgrj7mY1XpDGBtyepPcth3j49FESBMMp4euY7lsCgiXA0+TGUWZdjbdCFwwm7dm7drU4yQZFAy942E22nEXGg6OlP+O9JJbvAtIb3b8RzWdzZK4RpNGfQbU0JlfhyrSbxA8f904hiAcL/dBjV9nBJW1Y3JLACqjtcRnAvzJ1F0W5Ivre9RJsn4u3PdHf7WUtiodkSvIDLNlC9c+1Y5qvpO1Sdq72B5V9rdDiApQs/EGzmZLz+XMK/iUlCMEat/PY0idym2EMlMtbRnNwQxSxzPuRtycb0H5IjHkCcKSs4KVYKB4tDsBYG3MWKtcD59LyzqFb1HDo0vYYDC32sBOC7Bb3k7bazR87feXj5NK8155+SfnP10F4hyuT+ZFiKry47CadKyPr0t8ztVFisjUV4dJsOym9ceHDKRCiK4xI1RTIYC8ouD71qCKcmZqa+c5UMfdLNXqLz+1vlqUAr9dE2jcfl0wgroQBfpyuLfpNFn2jbkcriBT5GH5dY9VTw2C4oQ2p6vWc20/w4QKST/riUqiozfAOitx40UDzaLaxNWMM2S8UTjbKzZpUNBxKb7FG+fia+fFCEvMT9cc6XakoCa7XCW5+Cltm6/m0VPMQF00uJew0LT2BH9Dx8Z6yFodg/w6rqT5xVcmJXbIoCZ40cfr7DqbD/ELZW7CFWzgCek8R+gfsDqKneOJc4zwzUn99uB52i0SRC0uAB1tEqb+owuFjP6d5x1R59JmV9m3MnvGpKCFLOTJfWQCzXyIo+te04C2rl8VfiNqDItBJs+tc4cxQdXEkjojeY23vPI8I7OQmtv8AXYYpvbyv1Lva3jH3nbRHFsVWEIaCzMHcO8ldW2m1EPRNmBDAZAU2gytgKeNyoUj2KNbKSCekZHsIWPV6z2c94=`,
		BrowserInfo: &BrowserInfo{
			ScreenWidth: 3148,
			ScreenHeight: 886,
			ColorDepth: 24,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
			TimeZoneOffset: 240,
			Language: "en-US",
			JavaEnabled: false,
		},
	})

	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt json body")
		tk.Stop()
		return
	}

	req, err := tk.NewRequest("POST", fmt.Sprintf(`https://www.%s.com/api/v2/users/orders?timestamp=%s`, sites[tk.Site].name, time.Now().Unix()), body)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "create order confirmation request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.%s.com", sites[tk.Site].name))

	//todo check valid response
	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "create make confirmation request")
		tk.Stop()
		return
	}
}

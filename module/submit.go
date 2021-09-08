package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"time"
)

func (tk *Task) SubmitShipping(){
	var addrline2, billingIsShipping string
	if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil{
		addrline2 = *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2
	}

	if tk.Data.Profile.Shipping.BillingIsShipping{
		billingIsShipping = "true"
	}else
	{
		billingIsShipping = "false"
	}

	req, err := tk.NewRequest("POST", fmt.Sprintf(`https://www.%s.com/api/users/carts/current/addresses/shipping?timestamp=%d`, tk.Site, time.Now().Unix()),
		[]byte(fmt.Sprintf(`{"shippingAddress":{"setAsDefaultBilling":false,"setAsDefaultShipping":false,"firstName":"%s","lastName":"%s","email":false,"phone":"%s","country":{"isocode":"US","name":"United States"},"id":null,"setAsBilling":%s,"saveInAddressBook":false,"region":{"countryIso":"US","isocode":"US-%s","isocodeShort":"%s","name":"%s"},"type":"default","LoqateSearch":"","line1":"%s","line2":"%s","postalCode":"%s","town":"%s","regionFPO":null,"shippingAddress":true,"recordType":" "}}`,
			tk.Data.Profile.Shipping.FirstName,
			tk.Data.Profile.Shipping.LastName,
			tk.Data.Profile.Shipping.PhoneNumber,
			billingIsShipping,
			tk.Data.Profile.Shipping.ShippingAddress.StateCode,
			tk.Data.Profile.Shipping.ShippingAddress.StateCode,
			tk.Data.Profile.Shipping.ShippingAddress.State,
			tk.Data.Profile.Shipping.ShippingAddress.AddressLine,
			addrline2,
			tk.Data.Profile.Shipping.ShippingAddress.ZIP,
			tk.Data.Profile.Shipping.ShippingAddress.City,
		)))

	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create shipping req")
		tk.Stop()
		return
	}

	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt make shipping req")
		tk.Stop()
		return
	}
}

func (tk *Task) SubmitBilling(){
	var addrline2 string
	if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil{
		addrline2 = *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2
	}
	req, err := tk.NewRequest("POST", fmt.Sprintf(`https://www.%s.com/api/users/carts/current/set-billing?timestamp=%d`, tk.Site, time.Now().Unix()), []byte(fmt.Sprintf(
		`{"setAsDefaultBilling":false,"setAsDefaultShipping":false,"firstName":"%s","lastName":"%s","phone":"%s","country":{"isocode":"US","name":"United States"},"id":null,"saveInAddressBook":false,"region":{"countryIso":"US","isocode":"US-%s","isocodeShort":"%s","name":"%s"},"type":"default","LoqateSearch":"","line1":"%s","line2":"%s","postalCode":"%s","town":"%s","regionFPO":null,"recordType":" ","shippingAddress":true,"setAsBilling":false}`,
		tk.Data.Profile.Shipping.FirstName,
		tk.Data.Profile.Shipping.LastName,
		tk.Data.Profile.Shipping.PhoneNumber,
		tk.Data.Profile.Shipping.ShippingAddress.StateCode,
		tk.Data.Profile.Shipping.ShippingAddress.StateCode,
		tk.Data.Profile.Shipping.ShippingAddress.State,
		tk.Data.Profile.Shipping.ShippingAddress.AddressLine,
		addrline2,
		tk.Data.Profile.Shipping.ShippingAddress.ZIP,
		tk.Data.Profile.Shipping.ShippingAddress.City,
		)))
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt create billing req")
		tk.Stop()
		return
	}
	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "couldnt post billing req")
		tk.Stop()
		return
	}
}
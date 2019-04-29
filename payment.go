package f3api

type Payment struct {
	Type           string            `json:"type"`
	ID             string            `json:"id"`
	Version        int               `json:"version"`
	OrganisationID string            `json:"organisation_id"`
	Attributes     PaymentAttributes `json:"attributes"`
}

type PaymentAttributes struct {
	Amount               fractionalAmount  `json:"amount"`
	BeneficiaryParty     TypedParty        `json:"beneficiary_party"`
	ChargesInformation   ChargeInformation `json:"charges_information"`
	Currency             string            `json:"currency"`
	DebtorParty          Party             `json:"debtor_party"`
	EndToEndReference    string            `json:"end_to_end_reference"`
	Fx                   FX                `json:"fx"`
	NumericReference     stringedInt       `json:"numeric_reference"`
	PaymentID            stringedInt       `json:"payment_id"`
	PaymentPurpose       string            `json:"payment_purpose"`
	PaymentScheme        string            `json:"payment_scheme"`
	PaymentType          string            `json:"payment_type"`
	ProcessingDate       date              `json:"processing_date"`
	Reference            string            `json:"reference"`
	SchemePaymentSubType string            `json:"scheme_payment_sub_type"`
	SchemePaymentType    string            `json:"scheme_payment_type"`
	SponsorParty         MinimalParty      `json:"sponsor_party"`
}

type MinimalParty struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}

type Party struct {
	MinimalParty
	AccountName       string `json:"account_name"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	Name              string `json:"name"`
}

type TypedParty struct {
	Party
	AccountType int `json:"account_type"`
}

type ChargeInformation struct {
	BearerCode              string           `json:"bearer_code"`
	SenderCharges           []SenderCharge   `json:"sender_charges"`
	ReceiverChargesAmount   fractionalAmount `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency"`
}

type SenderCharge struct {
	Amount   fractionalAmount `json:"amount"`
	Currency string           `json:"currency"`
}

type FX struct {
	ContractReference string           `json:"contract_reference"`
	ExchangeRate      exchangeRate     `json:"exchange_rate"`
	OriginalAmount    fractionalAmount `json:"original_amount"`
	OriginalCurrency  string           `json:"original_currency"`
}

package f3api

// The Payment Type, our main resource payload.
// Details not disclosed, same goes for other types in this file.
type Payment struct {
	Type           string            `json:"type"`
	ID             string            `json:"id"`
	Version        int               `json:"version"`
	OrganisationID string            `json:"organisation_id"`
	Attributes     PaymentAttributes `json:"attributes"`
}

// No documentation exists for this type.
type PaymentAttributes struct {
	Amount               FractionalAmount  `json:"amount"`
	BeneficiaryParty     TypedParty        `json:"beneficiary_party"`
	ChargesInformation   ChargeInformation `json:"charges_information"`
	Currency             string            `json:"currency"`
	DebtorParty          Party             `json:"debtor_party"`
	EndToEndReference    string            `json:"end_to_end_reference"`
	Fx                   FX                `json:"fx"`
	NumericReference     StringedInt       `json:"numeric_reference"`
	PaymentID            StringedInt       `json:"payment_id"`
	PaymentPurpose       string            `json:"payment_purpose"`
	PaymentScheme        string            `json:"payment_scheme"`
	PaymentType          string            `json:"payment_type"`
	ProcessingDate       Date              `json:"processing_date"`
	Reference            string            `json:"reference"`
	SchemePaymentSubType string            `json:"scheme_payment_sub_type"`
	SchemePaymentType    string            `json:"scheme_payment_type"`
	SponsorParty         MinimalParty      `json:"sponsor_party"`
}

// No documentation exists for this type.
type MinimalParty struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}

// No documentation exists for this type.
type Party struct {
	MinimalParty
	AccountName       string `json:"account_name"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	Name              string `json:"name"`
}

// No documentation exists for this type.
type TypedParty struct {
	Party
	AccountType int `json:"account_type"`
}

// No documentation exists for this type.
type ChargeInformation struct {
	BearerCode              string           `json:"bearer_code"`
	SenderCharges           []SenderCharge   `json:"sender_charges"`
	ReceiverChargesAmount   FractionalAmount `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency"`
}

// No documentation exists for this type.
type SenderCharge struct {
	Amount   FractionalAmount `json:"amount"`
	Currency string           `json:"currency"`
}

// No documentation exists for this type.
type FX struct {
	ContractReference string           `json:"contract_reference"`
	ExchangeRate      ExchangeRate     `json:"exchange_rate"`
	OriginalAmount    FractionalAmount `json:"original_amount"`
	OriginalCurrency  string           `json:"original_currency"`
}

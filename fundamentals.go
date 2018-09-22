package tradier

//go:generate ffjson $GOFILE
import "encoding/json"

type CorporateEvent struct {
	BeginDateTime *string `json:"begin_date_time"`
	CompanyID     *string `json:"company_id"`
	EndDateTime   *string `json:"end_date_time"`
	Event         *string `json:"event"`
	EventType     *int64  `json:"event_type"`
	TimeZone      *string `json:"time_zone,omitempty"`
}

// If there is only a single event, then tradier sends back
// an object, but if there are multiple events, then it sends
// a list of objects...
type CorporateCalendar []CorporateEvent

func (cc *CorporateCalendar) UnmarshalJSON(data []byte) error {
	events := make([]CorporateEvent, 0)
	if err := json.Unmarshal(data, &events); err == nil {
		*cc = events
		return nil
	}

	event := CorporateEvent{}
	err := json.Unmarshal(data, &event)
	if err == nil {
		*cc = []CorporateEvent{event}
	}
	return err
}

type GetCorporateCalendarsResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			CorporateCalendars *CorporateCalendar `json:"corporate_calendars"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

type NAICS []int64

func (n *NAICS) UnmarshalJSON(data []byte) error {
	ids := make([]int64, 0)
	if err := json.Unmarshal(data, &ids); err == nil {
		*n = ids
		return nil
	}

	var id int64
	err := json.Unmarshal(data, &id)
	if err == nil {
		*n = []int64{id}
	}
	return err
}

type SIC []int64

func (n *SIC) UnmarshalJSON(data []byte) error {
	ids := make([]int64, 0)
	if err := json.Unmarshal(data, &ids); err == nil {
		*n = ids
		return nil
	}

	var id int64
	err := json.Unmarshal(data, &id)
	if err == nil {
		*n = []int64{id}
	}
	return err
}

type AssetClassification struct {
	FinancialHealthGradeAsOfDate *string  `json:"FinancialHealthGrade.asOfDate"`
	GrowthGradeAsOfDate          *string  `json:"GrowthGrade.asOfDate"`
	ProfitabilityGradeAsOfDate   *string  `json:"ProfitabilityGrade.asOfDate"`
	StockTypeAsOfDate            *string  `json:"StockType.asOfDate"`
	StyleBoxAsOfDate             *string  `json:"StyleBox.asOfDate"`
	CANNAICS                     *int64   `json:"c_a_n_n_a_i_c_s"`
	CompanyID                    *string  `json:"company_id"`
	FinancialHealthGrade         *string  `json:"financial_health_grade"`
	GrowthGrade                  *string  `json:"growth_grade"`
	GrowthScore                  *float64 `json:"growth_score"`
	MorningstarEconomySphereCode *int64   `json:"morningstar_economy_sphere_code"`
	MorningstarIndustryCode      *int64   `json:"morningstar_industry_code"`
	MorningstarIndustryGroupCode *int64   `json:"morningstar_industry_group_code"`
	MorningstarSectorCode        *int64   `json:"morningstar_sector_code"`
	NACE                         *float64 `json:"n_a_c_e"`
	NAICS                        NAICS    `json:"n_a_i_c_s"`
	ProfitabilityGrade           *string  `json:"profitability_grade"`
	SIC                          SIC      `json:"s_i_c"`
	SizeScore                    *float64 `json:"size_score"`
	StockType                    *int64   `json:"stock_type"`
	StyleBox                     *int64   `json:"style_box"`
	StyleScore                   *float64 `json:"style_score"`
	ValueScore                   *float64 `json:"value_score"`
}

type CompanyHeadquarter struct {
	AddressLine1 *string `json:"address_line1"`
	City         *string `json:"city"`
	Country      *string `json:"country"`
	Fax          *string `json:"fax"`
	Homepage     *string `json:"homepage"`
	Phone        *string `json:"phone"`
	PostalCode   *string `json:"postal_code"`
	Province     *string `json:"province"`
}

type CompanyProfile struct {
	TotalEmployeeNumberAsOfDate *string             `json:"TotalEmployeeNumber.asOfDate"`
	CompanyID                   *string             `json:"company_id"`
	ContactEmail                *string             `json:"contact_email"`
	Headquarter                 *CompanyHeadquarter `json:"headquarter"`
	ShortDescription            *string             `json:"short_description"`
	TotalEmployeeNumber         *int64              `json:"total_employee_number"`
}

type HistoricalAssetClassification struct {
	AsOfDate                     *string  `json:"as_of_date"`
	CompanyID                    *string  `json:"company_id"`
	FinancialHealthGrade         *string  `json:"financial_health_grade"`
	GrowthScore                  *float64 `json:"growth_score"`
	MorningstarEconomySphereCode *int64   `json:"morningstar_economy_sphere_code"`
	MorningstarIndustryCode      *int64   `json:"morningstar_industry_code"`
	MorningstarIndustryGroupCode *int64   `json:"morningstar_industry_group_code"`
	MorningstarSectorCode        *int64   `json:"morningstar_sector_code"`
	ProfitabilityGrade           *string  `json:"profitability_grade"`
	SizeScore                    *float64 `json:"size_score"`
	StockType                    *int64   `json:"stock_type"`
	StyleBox                     *int64   `json:"style_box"`
	StyleScore                   *float64 `json:"style_score"`
	ValueScore                   *float64 `json:"value_score"`
}

type ShareClass struct {
	CompanyID           *string `json:"company_id"`
	CUSIP               *string `json:"c_u_s_i_p"`
	CurrencyID          *string `json:"currency_id"`
	DelistingDate       *string `json:"delisting_date"`
	ExchangeID          *string `json:"exchange_id"`
	IPODate             *string `json:"i_p_o_date"`
	ISIN                *string `json:"i_s_i_n"`
	InvestmentID        *string `json:"investment_id"`
	IsDepositaryReceipt *bool   `json:"is_depositary_receipt"`
	IsDirectInvest      *bool   `json:"is_direct_invest"`
	IsDividendReinvest  *bool   `json:"is_dividend_reinvest"`
	IsPrimaryShare      *bool   `json:"is_primary_share"`
	MIC                 *string `json:"m_i_c"`
	SEDOL               *string `json:"s_e_d_o_l"`
	SecurityType        *string `json:"security_type"`
	ShareClassID        *string `json:"share_class_id"`
	ShareClassStatus    *string `json:"share_class_status"`
	Symbol              *string `json:"symbol"`
	TradingStatus       *bool   `json:"trading_status"`
	Valoren             *string `json:"valoren"`
}

type ShareClassProfile struct {
	EnterpriseValueAsOfDate                     *string `json:"EnterpriseValue.asOfDate"`
	MarketCapAsOfDate                           *string `json:"MarketCap.asOfDate"`
	SharesOutstandingAsOfDate                   *string `json:"SharesOutstanding.asOfDate"`
	EnterpriseValue                             *int64  `json:"enterprise_value"`
	MarketCap                                   *int64  `json:"market_cap"`
	ShareClassID                                *string `json:"share_class_id"`
	ShareClassLevelSharesOutstanding            *int64  `json:"share_class_level_shares_outstanding"`
	SharesOutstanding                           *int64  `json:"shares_outstanding"`
	SharesOutstandingWithBalanceSheetEndingDate *string `json:"shares_outstanding_with_balance_sheet_ending_date"`
}

type OwnershipDetail struct {
	AsOfDate              *string  `json:"as_of_date"`
	CurrencyOfMarketValue *string  `json:"currencyof_market_value"`
	MarketValue           *int64   `json:"market_value"`
	NumberOfShares        *float64 `json:"number_of_shares"`
	OwnerCIK              *int64   `json:"owner_c_i_k"`
	OwnerID               *string  `json:"owner_id"`
	OwnerName             *string  `json:"owner_name"`
	OwnerType             *int64   `json:"owner_type,string"`
	PercentageInPortfolio *float64 `json:"percentage_in_portfolio"`
	PercentageOwnership   *float64 `json:"percentage_ownership"`
	ShareChange           *int64   `json:"share_change"`
	ShareClassID          *string  `json:"share_class_id"`
}

type OwnershipDetails []OwnershipDetail

func (ods *OwnershipDetails) UnmarshalJSON(data []byte) error {
	details := make([]OwnershipDetail, 0)
	if err := json.Unmarshal(data, &details); err == nil {
		*ods = details
		return nil
	}

	detail := OwnershipDetail{}
	err := json.Unmarshal(data, &detail)
	if err == nil {
		*ods = []OwnershipDetail{detail}
	}
	return err
}

type OwnershipSummary struct {
	AsOfDate                                     *string            `json:"as_of_date"`
	DaysToCoverShort                             map[string]float64 `json:"days_to_cover_short"`
	Float                                        *int64             `json:"float"`
	InsiderPercentOwned                          *float64           `json:"insider_percent_owned"`
	InsiderSharesBought                          *int64             `json:"insider_shares_bought"`
	InsiderSharesOwned                           *int64             `json:"insider_shares_owned"`
	InsiderSharesSold                            *int64             `json:"insider_shares_sold"`
	InstitutionHolderNumber                      *int64             `json:"institution_holder_number"`
	InstitutionSharesBought                      *int64             `json:"institution_shares_bought"`
	InstitutionSharesHeld                        *int64             `json:"institution_shares_held"`
	InstitutionSharesSold                        *int64             `json:"institution_shares_sold"`
	NumberOfInsiderBuys                          *int64             `json:"number_of_insider_buys"`
	NumberOfInsiderSellers                       *int64             `json:"number_of_insider_sellers"`
	ShareClassID                                 *string            `json:"share_class_id"`
	ShareClassLevelSharesOutstanding             *int64             `json:"share_class_level_shares_outstanding"`
	ShareClassLevelSharesOutstandingBalanceSheet *int64             `json:"share_class_level_shares_outstanding_balance_sheet"`
	ShareClassLevelSharesOutstandingInterim      *int64             `json:"share_class_level_shares_outstanding_interim"`
	ShareClassLevelTreasuryShareOutstanding      *int64             `json:"share_class_level_treasury_share_outstanding"`
	SharesOutstanding                            *int64             `json:"shares_outstanding"`
	SharesOutstandingWithBalanceSheetEndingDate  *string            `json:"shares_outstanding_with_balance_sheet_ending_date"`
	ShortInterest                                *int64             `json:"short_interest"`
	ShortInterestsPercentageChange               map[string]float64 `json:"short_interests_percentage_change"`
	ShortPercentageOfFloat                       *float64           `json:"short_percentage_of_float"`
	ShortPercentageOfSharesOutstanding           *float64           `json:"short_percentage_of_shares_outstanding"`
}

type CompanyInfoResult struct {
	ID     string `json:"id"`
	Tables struct {
		AssetClassification           *AssetClassification           `json:"asset_classification"`
		CompanyProfile                *CompanyProfile                `json:"company_profile"`
		HistoricalAssetClassification *HistoricalAssetClassification `json:"historical_asset_classification"`
		LongDescriptions              *string                        `json:"long_descriptions"`
		OwnershipDetails              OwnershipDetails               `json:"ownership_details"`
		OwnershipSummary              *OwnershipSummary              `json:"ownership_summary"`
		ShareClass                    *ShareClass                    `json:"share_class"`
		ShareClassProfile             *ShareClassProfile             `json:"share_class_profile"`
	} `json:"tables"`
	Type string `json:"type"`
}

type GetCompanyInfoResponse []struct {
	Error   string
	Request string              `json:"request"`
	Results []CompanyInfoResult `json:"results"`
	Type    string              `json:"type"`
}

type MergerAndAcquisition struct {
	AcquiredCompanyID *string  `json:"acquired_company_id"`
	CashAmount        *float64 `json:"cash_amount"`
	CurrencyID        *string  `json:"currency_id"`
	EffectiveDate     *string  `json:"effective_date"`
	Notes             *string  `json:"notes"`
	ParentCompanyID   *string  `json:"parent_company_id"`
}

type MergersAndAcquisitions []MergerAndAcquisition

func (maq *MergersAndAcquisitions) UnmarshalJSON(data []byte) error {
	events := make([]MergerAndAcquisition, 0)
	if err := json.Unmarshal(data, &events); err == nil {
		*maq = events
		return nil
	}

	event := MergerAndAcquisition{}
	err := json.Unmarshal(data, &event)
	if err == nil {
		*maq = []MergerAndAcquisition{event}
	}
	return err
}

type StockSplit struct {
	AdjustmentFactor *float64 `json:"adjustment_factor"`
	ExDate           *string  `json:"ex_date"`
	ShareClassID     *string  `json:"share_class_id"`
	SplitFrom        *float64 `json:"split_from"`
	SplitTo          *float64 `json:"split_to"`
	SplitType        *string  `json:"split_type"`
}

type StockSplits map[string]StockSplit

type GetCorporateActionsResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			MergersAndAcquisitions MergersAndAcquisitions `json:"mergers_and_acquisitions"`
			StockSplits            StockSplits            `json:"stock_splits"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

type CashDividend struct {
	CashAmount      *float64 `json:"cash_amount"`
	CurrencyID      *string  `json:"currency_i_d"`
	DeclarationDate *string  `json:"declaration_date"`
	DividendType    *string  `json:"dividend_type"`
	ExDate          *string  `json:"ex_date"`
	Frequency       *int64   `json:"frequency"`
	PayDate         *string  `json:"pay_date"`
	RecordDate      *string  `json:"record_date"`
	ShareClassID    *string  `json:"share_class_id"`
}

type CashDividends []CashDividend

func (cds *CashDividends) UnmarshalJSON(data []byte) error {
	dividends := make([]CashDividend, 0)
	if err := json.Unmarshal(data, &dividends); err == nil {
		*cds = dividends
		return nil
	}

	d := CashDividend{}
	err := json.Unmarshal(data, &d)
	if err == nil {
		*cds = []CashDividend{d}
	}
	return err
}

type GetDividendsResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			CashDividends CashDividends `json:"cash_dividends"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

type BalanceSheet struct {
	AccountsPayable                               *float64 `json:"accounts_payable"`
	AccountsReceivable                            *float64 `json:"accounts_receivable"`
	AccumulatedDepreciation                       *float64 `json:"accumulated_depreciation"`
	CapitalStock                                  *float64 `json:"capital_stock"`
	CashAndCashEquivalents                        *float64 `json:"cash_and_cash_equivalents"`
	CashCashEquivalentsAndMarketableSecurities    *float64 `json:"cash_cash_equivalents_and_marketable_securities"`
	CommercialPaper                               *float64 `json:"commercial_paper"`
	CommonStock                                   *float64 `json:"common_stock"`
	CommonStockEquity                             *float64 `json:"common_stock_equity"`
	CurrencyID                                    *string  `json:"currency_id"`
	CurrentAccruedExpenses                        *float64 `json:"current_accrued_expenses"`
	CurrentAssets                                 *float64 `json:"current_assets"`
	CurrentDebt                                   *float64 `json:"current_debt"`
	CurrentDebtAndCapitalLeaseObligation          *float64 `json:"current_debt_and_capital_lease_obligation"`
	CurrentDeferredLiabilities                    *float64 `json:"current_deferred_liabilities"`
	CurrentDeferredRevenue                        *float64 `json:"current_deferred_revenue"`
	CurrentLiabilities                            *float64 `json:"current_liabilities"`
	FileDate                                      *string  `json:"file_date"`
	FiscalYearEnd                                 *string  `json:"fiscal_year_end"`
	GainsLossesNotAffectingRetainedEarnings       *float64 `json:"gains_losses_not_affecting_retained_earnings"`
	Goodwill                                      *float64 `json:"goodwill"`
	GoodwillAndOtherIntangibleAssets              *float64 `json:"goodwill_and_other_int64angible_assets"`
	GrossPPE                                      *float64 `json:"gross_p_p_e"`
	Inventory                                     *float64 `json:"inventory"`
	InvestedCapital                               *float64 `json:"invested_capital"`
	InvestmentsAndAdvances                        *float64 `json:"investments_and_advances"`
	LandAndImprovements                           *float64 `json:"land_and_improvements"`
	Leases                                        *float64 `json:"leases"`
	LongTermDebt                                  *float64 `json:"long_term_debt"`
	LongTermDebtAndCapitalLeaseObligation         *float64 `json:"long_term_debt_and_capital_lease_obligation"`
	MachineryFurnitureEquipment                   *float64 `json:"machinery_furniture_equipment"`
	NetDebt                                       *float64 `json:"net_debt"`
	NetPPE                                        *float64 `json:"net_p_p_e"`
	NetTangibleAssets                             *float64 `json:"net_tangible_assets"`
	NonCurrentDeferredLiabilities                 *float64 `json:"non_current_deferred_liabilities"`
	NonCurrentDeferredRevenue                     *float64 `json:"non_current_deferred_revenue"`
	NonCurrentDeferredTaxesLiabilities            *float64 `json:"non_current_deferred_taxes_liabilities"`
	NumberOfShareHolders                          *int64   `json:"number_of_share_holders"`
	OrdinarySharesNumber                          *float64 `json:"ordinary_shares_number"`
	OtherCurrentAssets                            *float64 `json:"other_current_assets"`
	OtherCurrentBorrowings                        *float64 `json:"other_current_borrowings"`
	OtherIntangibleAssets                         *float64 `json:"other_int64angible_assets"`
	OtherNonCurrentAssets                         *float64 `json:"other_non_current_assets"`
	OtherNonCurrentLiabilities                    *float64 `json:"other_non_current_liabilities"`
	OtherReceivables                              *float64 `json:"other_receivables"`
	OtherShortTermInvestments                     *float64 `json:"other_short_term_investments"`
	Payables                                      *float64 `json:"payables"`
	PayablesAndAccruedExpenses                    *float64 `json:"payables_and_accrued_expenses"`
	Period                                        *string  `json:"period"`
	PeriodEndingDate                              *string  `json:"period_ending_date"`
	Receivables                                   *float64 `json:"receivables"`
	ReportType                                    *string  `json:"report_type"`
	RetainedEarnings                              *float64 `json:"retained_earnings"`
	ShareIssued                                   *float64 `json:"share_issued"`
	StockholdersEquity                            *float64 `json:"stockholders_equity"`
	TangibleBookValue                             *float64 `json:"tangible_book_value"`
	TotalAssets                                   *float64 `json:"total_assets"`
	TotalCapitalization                           *float64 `json:"total_capitalization"`
	TotalDebt                                     *float64 `json:"total_debt"`
	TotalEquity                                   *float64 `json:"total_equity"`
	TotalEquityGrossMinorityInterest              *float64 `json:"total_equity_gross_minority_int64erest"`
	TotalLiabilities                              *float64 `json:"total_liabilities"`
	TotalLiabilitiesNetMinorityInterest           *float64 `json:"total_liabilities_net_minority_int64erest"`
	TotalNonCurrentAssets                         *float64 `json:"total_non_current_assets"`
	TotalNonCurrentLiabilities                    *float64 `json:"total_non_current_liabilities"`
	TotalNonCurrentLiabilitiesNetMinorityInterest *float64 `json:"total_non_current_liabilities_net_minority_int64erest"`
	WorkingCapital                                *float64 `json:"working_capital"`
}

type CashFlowStatement struct {
	BeginningCashPosition             *float64 `json:"beginning_cash_position"`
	CapitalExpenditure                *float64 `json:"capital_expenditure"`
	CashDividendsPaid                 *float64 `json:"cash_dividends_paid"`
	ChangeInAccountPayable            *float64 `json:"change_in_account_payable"`
	ChangeInInventory                 *float64 `json:"change_in_inventory"`
	ChangeInOtherWorkingCapital       *float64 `json:"change_in_other_working_capital"`
	ChangeInPayable                   *float64 `json:"change_in_payable"`
	ChangeInPayablesAndAccruedExpense *float64 `json:"change_in_payables_and_accrued_expense"`
	ChangeInReceivables               *float64 `json:"change_in_receivables"`
	ChangeInWorkingCapital            *float64 `json:"change_in_working_capital"`
	ChangesInAccountReceivables       *float64 `json:"changes_in_account_receivables"`
	ChangesInCash                     *float64 `json:"changes_in_cash"`
	CommonStockIssuance               *float64 `json:"common_stock_issuance"`
	CommonStockPayments               *float64 `json:"common_stock_payments"`
	CurrencyID                        *string  `json:"currency_id"`
	DeferredIncomeTax                 *float64 `json:"deferred_income_tax"`
	DeferredTax                       *float64 `json:"deferred_tax"`
	DepreciationAmortizationDepletion *float64 `json:"depreciation_amortization_depletion"`
	DepreciationAndAmortization       *float64 `json:"depreciation_and_amortization"`
	DomesticSales                     *float64 `json:"domestic_sales"`
	EndCashPosition                   *float64 `json:"end_cash_position"`
	FileDate                          *string  `json:"file_date"`
	FinancingCashFlow                 *float64 `json:"financing_cash_flow"`
	FiscalYearEnd                     *string  `json:"fiscal_year_end"`
	ForeignSales                      *float64 `json:"foreign_sales"`
	FreeCashFlow                      *float64 `json:"free_cash_flow"`
	IncomeTaxPaidSupplementalData     *float64 `json:"income_tax_paid_supplemental_data"`
	InterestPaidSupplementalData      *float64 `json:"int64erest_paid_supplemental_data"`
	InvestingCashFlow                 *float64 `json:"investing_cash_flow"`
	IssuanceOfCapitalStock            *float64 `json:"issuance_of_capital_stock"`
	NetBusinessPurchaseAndSale        *float64 `json:"net_business_purchase_and_sale"`
	NetCommonStockIssuance            *float64 `json:"net_common_stock_issuance"`
	NetIncome                         *float64 `json:"net_income"`
	NetIncomeFromContinuingOperations *float64 `json:"net_income_from_continuing_operations"`
	NetIntangiblesPurchaseAndSale     *float64 `json:"net_int64angibles_purchase_and_sale"`
	NetInvestmentPurchaseAndSale      *float64 `json:"net_investment_purchase_and_sale"`
	NetIssuancePaymentsOfDebt         *float64 `json:"net_issuance_payments_of_debt"`
	NetOtherFinancingCharges          *float64 `json:"net_other_financing_charges"`
	NetOtherInvestingChanges          *float64 `json:"net_other_investing_changes"`
	NetPPEPurchaseAndSale             *float64 `json:"net_p_p_e_purchase_and_sale"`
	NetShortTermDebtIssuance          *float64 `json:"net_short_term_debt_issuance"`
	NumberOfShareHolders              *int64   `json:"number_of_share_holders"`
	OperatingCashFlow                 *float64 `json:"operating_cash_flow"`
	OtherNonCashItems                 *float64 `json:"other_non_cash_items"`
	Period                            *string  `json:"period"`
	PeriodEndingDate                  *string  `json:"period_ending_date"`
	PurchaseOfBusiness                *float64 `json:"purchase_of_business"`
	PurchaseOfIntangibles             *float64 `json:"purchase_of_int64angibles"`
	PurchaseOfInvestment              *float64 `json:"purchase_of_investment"`
	PurchaseOfPPE                     *float64 `json:"purchase_of_p_p_e"`
	ReportType                        *string  `json:"report_type"`
	RepurchaseOfCapitalStock          *float64 `json:"repurchase_of_capital_stock"`
	SaleOfInvestment                  *float64 `json:"sale_of_investment"`
	StockBasedCompensation            *float64 `json:"stock_based_compensation"`
}

type IncomeStatement struct {
	AccessionNumber                                     *string  `json:"accession_number"`
	CostOfRevenue                                       *float64 `json:"cost_of_revenue"`
	CurrencyID                                          *string  `json:"currency_id"`
	EBIT                                                *float64 `json:"e_b_i_t"`
	EBITDA                                              *float64 `json:"e_b_i_t_d_a"`
	FileDate                                            *string  `json:"file_date"`
	FiscalYearEnd                                       *string  `json:"fiscal_year_end"`
	FormType                                            *string  `json:"form_type"`
	GrossProfit                                         *float64 `json:"gross_profit"`
	InterestExpense                                     *float64 `json:"int64erest_expense"`
	InterestExpenseNonOperating                         *float64 `json:"int64erest_expense"`
	InterestIncome                                      *float64 `json:"int64erest_income"`
	InterestIncomeNonOperating                          *float64 `json:"int64erest_income_non_operating"`
	InterestAndSimilarIncome                            *float64 `json:"int64erestand_similar_income"`
	NetIncome                                           *float64 `json:"net_income"`
	NetIncomeCommonStockholders                         *float64 `json:"net_income_common_stockholders"`
	NetIncomeContinuousOperations                       *float64 `json:"net_income_continuous_operations"`
	NetIncomeFromContinuingAndDiscontinuedOperation     *float64 `json:"net_income_from_continuing_and_discontinued_operation"`
	NetIncomeFromContinuingOperationNetMinorityInterest *float64 `json:"net_income_from_continuing_operation_net_minority_int64erest"`
	NetIncomeIncludingNoncontrollingInterests           *float64 `json:"net_income_including_noncontrolling_int64erests"`
	NetInterestIncome                                   *float64 `json:"net_int64erest_income"`
	NetNonOperatingInterestIncomeExpense                *float64 `json:"net_non_operating_int64erest_income_expense"`
	NonOperatingExpenses                                *float64 `json:"non_operating_expenses"`
	NonOperatingIncome                                  *float64 `json:"non_operating_income"`
	NormalizedEBITDA                                    *float64 `json:"normalized_e_b_i_t_d_a"`
	NormalizedIncome                                    *float64 `json:"normalized_income"`
	NumberOfShareHolders                                *int64   `json:"number_of_share_holders"`
	OperatingExpense                                    *float64 `json:"operating_expense"`
	OperatingIncome                                     *float64 `json:"operating_income"`
	OperatingRevenue                                    *float64 `json:"operating_revenue"`
	OtherIncomeExpense                                  *float64 `json:"other_income_expense"`
	Period                                              *string  `json:"period"`
	PeriodEndingDate                                    *string  `json:"period_ending_date"`
	PretaxIncome                                        *float64 `json:"pretax_income"`
	ReconciledCostOfRevenue                             *float64 `json:"reconciled_cost_of_revenue"`
	ReconciledDepreciation                              *float64 `json:"reconciled_depreciation"`
	ReportType                                          *string  `json:"report_type"`
	ResearchAndDevelopment                              *float64 `json:"research_and_development"`
	SellingGeneralAndAdministration                     *float64 `json:"selling_general_and_administration"`
	TaxEffectOfUnusualItems                             *float64 `json:"tax_effect_of_unusual_items"`
	TaxProvision                                        *float64 `json:"tax_provision"`
	TaxRateForCalcs                                     *float64 `json:"tax_rate_for_calcs"`
	TotalExpenses                                       *float64 `json:"total_expenses"`
	TotalRevenue                                        *float64 `json:"total_revenue"`
}

type BalanceSheetResults []map[string]BalanceSheet

func (bsr *BalanceSheetResults) UnmarshalJSON(data []byte) error {
	results := make([]map[string]BalanceSheet, 0)
	if err := json.Unmarshal(data, &results); err == nil {
		*bsr = results
		return nil
	}

	r := make(map[string]BalanceSheet)
	err := json.Unmarshal(data, &r)
	if err == nil {
		*bsr = []map[string]BalanceSheet{r}
	}
	return err
}

type CashFlowStatements []map[string]CashFlowStatement

func (cfs *CashFlowStatements) UnmarshalJSON(data []byte) error {
	results := make([]map[string]CashFlowStatement, 0)
	if err := json.Unmarshal(data, &results); err == nil {
		*cfs = results
		return nil
	}

	r := make(map[string]CashFlowStatement)
	err := json.Unmarshal(data, &r)
	if err == nil {
		*cfs = []map[string]CashFlowStatement{r}
	}
	return err
}

type IncomeStatements []map[string]IncomeStatement

func (is *IncomeStatements) UnmarshalJSON(data []byte) error {
	results := make([]map[string]IncomeStatement, 0)
	if err := json.Unmarshal(data, &results); err == nil {
		*is = results
		return nil
	}

	r := make(map[string]IncomeStatement)
	err := json.Unmarshal(data, &r)
	if err == nil {
		*is = []map[string]IncomeStatement{r}
	}
	return err
}

type FinancialStatementsRestate struct {
	AsOfDate          *string             `json:"as_of_date"`
	BalanceSheet      BalanceSheetResults `json:"balance_sheet"`
	CashFlowStatement CashFlowStatements  `json:"cash_flow_statement"`
	CompanyID         *string             `json:"company_id"`
	IncomeStatement   IncomeStatements    `json:"income_statement"`
}

type Segmentation struct {
	AsOfDate                    *string  `json:"as_of_date"`
	CompanyID                   *string  `json:"company_id"`
	DepreciationAndAmortization *float64 `json:"depreciation_and_amortization"`
	OperatingIncome             *float64 `json:"operating_income"`
	OperatingRevenue            *float64 `json:"operating_revenue"`
	Period                      *string  `json:"period"`
	TotalAssets                 *float64 `json:"total_assets"`
}

type EarningReport struct {
	AccessionNumber                     *string  `json:"accession_number"`
	AsOfDate                            *string  `json:"as_of_date"`
	BasicAverageShares                  *float64 `json:"basic_average_shares"`
	BasicContinuousOperations           *float64 `json:"basic_continuous_operations"`
	BasicEPS                            *float64 `json:"basic_e_p_s"`
	ContinuingAndDiscontinuedBasicEPS   *float64 `json:"continuing_and_discontinued_basic_e_p_s"`
	ContinuingAndDiscontinuedDilutedEPS *float64 `json:"continuing_and_discontinued_diluted_e_p_s"`
	CurrencyID                          *string  `json:"currency_id"`
	DilutedAverageShares                *float64 `json:"diluted_average_shares"`
	DilutedContinuousOperations         *float64 `json:"diluted_continuous_operations"`
	DilutedEPS                          *float64 `json:"diluted_e_p_s"`
	DividendPerShare                    *float64 `json:"dividend_per_share"`
	FileDate                            *string  `json:"file_date"`
	FiscalYearEnd                       *string  `json:"fiscal_year_end"`
	FormType                            *string  `json:"form_type"`
	NormalizedBasicEPS                  *float64 `json:"normalized_basic_e_p_s"`
	NormalizedDilutedEPS                *float64 `json:"normalized_diluted_e_p_s"`
	Period                              *string  `json:"period"`
	PeriodEndingDate                    *string  `json:"period_ending_date"`
	ReportType                          *string  `json:"report_type"`
	ShareClassID                        *string  `json:"share_class_id"`
}

type HistoricalReturns struct {
	AsOfDate     *string  `json:"as_of_date"`
	Period       *string  `json:"period"`
	ShareClassID *string  `json:"share_class_id"`
	TotalReturn  *float64 `json:"total_return"`
}

type EarningReports []map[string]EarningReport

func (ers *EarningReports) UnmarshalJSON(data []byte) error {
	results := make([]map[string]EarningReport, 0)
	if err := json.Unmarshal(data, &results); err == nil {
		*ers = results
		return nil
	}

	r := make(map[string]EarningReport)
	err := json.Unmarshal(data, &r)
	if err == nil {
		*ers = []map[string]EarningReport{r}
	}
	return err
}

type GetFinancialsResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			FinancialStatementsRestate *FinancialStatementsRestate  `json:"financial_statements_restate"`
			Segmentation               map[string]Segmentation      `json:"segmentation"`
			EarningReportsAOR          EarningReports               `json:"earning_reports_a_o_r"`
			EarningReportsRestate      EarningReports               `json:"earning_reports_restate"`
			HistoricalReturns          map[string]HistoricalReturns `json:"historical_returns"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

type OperationRatio struct {
	AsOfDate                      *string  `json:"as_of_date"`
	AssetsTurnover                *float64 `json:"assets_turnover"`
	CapExSalesRatio               *float64 `json:"cap_ex_sales_ratio"`
	CashConversionCycle           *float64 `json:"cash_conversion_cycle"`
	CommonEquityToAssets          *float64 `json:"common_equity_to_assets"`
	CompanyID                     *string  `json:"company_id"`
	CurrentRatio                  *float64 `json:"current_ratio"`
	DaysInInventory               *float64 `json:"days_in_inventory"`
	DaysInPayment                 *float64 `json:"days_in_payment"`
	DaysInSales                   *float64 `json:"days_in_sales"`
	DebtToAssets                  *float64 `json:"debt_to_assets"`
	EBITDAMargin                  *float64 `json:"e_b_i_t_d_a_margin"`
	EBITMargin                    *float64 `json:"e_b_i_t_margin"`
	FCFNetIncomeRatio             *float64 `json:"f_c_f_net_income_ratio"`
	FCFSalesRatio                 *float64 `json:"f_c_f_sales_ratio"`
	FinancialLeverage             *float64 `json:"financial_leverage"`
	FiscalYearEnd                 *string  `json:"fiscal_year_end"`
	FixAssetsTurnover             *float64 `json:"fix_assets_turonver"`
	GrossMargin                   *float64 `json:"gross_margin"`
	InterestCoverage              *float64 `json:"int64erest_coverage"`
	InventoryTurnover             *float64 `json:"inventory_turnover"`
	LongTermDebtEquityRatio       *float64 `json:"long_term_debt_equity_ratio"`
	LongTermDebtTotalCapitalRatio *float64 `json:"long_term_debt_total_capital_ratio"`
	NetIncomeGrowth               *float64 `json:"net_income_growth"`
	NetIncomeContOpsGrowth        *float64 `json:"net_income_cont_ops_growth"`
	NetMargin                     *float64 `json:"net_margin"`
	NormalizedNetProfitMargin     *float64 `json:"normalized_net_profit_margin"`
	NormalizedROIC                *float64 `json:"normalized_r_o_i_c"`
	OperationIncomeGrowth         *float64 `json:"operation_income_growth"`
	OperationMargin               *float64 `json:"operation_margin"`
	PaymentTurnover               *float64 `json:"payment_turnover"`
	Period                        *string  `json:"period"`
	PretaxMargin                  *float64 `json:"pretax_margin"`
	QuickRatio                    *float64 `json:"quick_ratio"`
	ROA                           *float64 `json:"r_o_a"`
	ROE                           *float64 `json:"r_o_e"`
	ROIC                          *float64 `json:"r_o_i_c"`
	ReceivableTurnover            *float64 `json:"receivable_turnover"`
	ReportType                    *string  `json:"report_type"`
	SalesPerEmployee              *float64 `json:"sales_per_employee"`
	TaxRate                       *float64 `json:"tax_rate"`
	TotalDebtEquityRatio          *float64 `json:"total_debt_equity_ratio"`
}

type AlphaBeta struct {
	Alpha        *float64 `json:"alpha"`
	AsOfDate     *string  `json:"as_of_date"`
	Beta         *float64 `json:"beta"`
	NonDivAlpha  *float64 `json:"non_div_alpha"`
	NonDivBeta   *float64 `json:"non_div_beta"`
	Period       *string  `json:"period"`
	ShareClassID *string  `json:"share_class_id"`
}

type EarningsRatiosRestate struct {
	AsOfDate             *string  `json:"as_of_date"`
	DPSGrowth            *float64 `json:"d_p_s_growth"`
	DilutedContEPSGrowth *float64 `json:"diluted_cont_e_p_s_growth"`
	DilutedEPSGrowth     *float64 `json:"diluted_e_p_s_growth"`
	FiscalYearEnd        *string  `json:"fiscal_year_end"`
	Period               *string  `json:"period"`
	ReportType           *string  `json:"report_type"`
	ShareClassID         *string  `json:"share_class_id"`
}

type ValuationRatios struct {
	AsOfDate                       *string  `json:"as_of_date"`
	BookValuePerShare              *float64 `json:"book_value_per_share"`
	BookValueYield                 *float64 `json:"book_value_yield"`
	BuyBackYield                   *float64 `json:"buy_back_yield"`
	CFOPerShare                    *float64 `json:"c_f_o_per_share"`
	CFYield                        *float64 `json:"c_f_yield"`
	CashReturn                     *float64 `json:"cash_return"`
	DividendRate                   *float64 `json:"dividend_rate"`
	DividendYield                  *float64 `json:"dividend_yield"`
	EVToEBITDA                     *float64 `json:"e_v_to_e_b_i_t_d_a"`
	EarningYield                   *float64 `json:"earning_yield"`
	FCFPerShare                    *float64 `json:"f_c_f_per_share"`
	FCFRatio                       *float64 `json:"f_c_f_ratio"`
	FCFYield                       *float64 `json:"f_c_f_yield"`
	ForwardDividendYield           *float64 `json:"forward_dividend_yield"`
	ForwardEarningYield            *float64 `json:"forward_earning_yield"`
	ForwardPERatio                 *float64 `json:"forward_p_e_ratio"`
	NormalizedPERatio              *float64 `json:"normalized_p_e_ratio"`
	PBRatio                        *float64 `json:"p_b_ratio"`
	PCFRatio                       *float64 `json:"p_c_f_ratio"`
	PEGPayback                     *float64 `json:"p_e_g_payback"`
	PEGRatio                       *float64 `json:"p_e_g_ratio"`
	PERatio                        *float64 `json:"p_e_ratio"`
	PSRatio                        *float64 `json:"p_s_ratio"`
	PayoutRatio                    *float64 `json:"payout_ratio"`
	PriceChange1M                  *float64 `json:"price_change1_m"`
	PriceToEBITDA                  *float64 `json:"priceto_e_b_i_t_d_a"`
	RatioPE5YearAverage            *float64 `json:"ratio_p_e5_year_average"`
	SalesPerShare                  *float64 `json:"sales_per_share"`
	SalesYield                     *float64 `json:"sales_yield"`
	ShareClassID                   *string  `json:"share_class_id"`
	SustainableGrowthRate          *float64 `json:"sustainable_growth_rate"`
	TangibleBVPerShare3YearAvg     *float64 `json:"tangible_b_v_per_share3_yr_avg"`
	TangibleBVPerShare5YearAvg     *float64 `json:"tangible_b_v_per_share5_yr_avg"`
	TangibleBookValuePerShare      *float64 `json:"tangible_book_value_per_share"`
	TotalYield                     *float64 `json:"total_yield"`
	WorkingCapitalPerShare         *float64 `json:"working_capital_per_share"`
	WorkingCapitalPerShare3YearAvg *float64 `json:"working_capital_per_share3_yr_avg"`
	WorkingCapitalPerShare5YearAvg *float64 `json:"working_capital_per_share5_yr_avg"`
}

type OperationRatios []map[string]OperationRatio

func (ors *OperationRatios) UnmarshalJSON(data []byte) error {
	results := make([]map[string]OperationRatio, 0)
	if err := json.Unmarshal(data, &results); err == nil {
		*ors = results
		return nil
	}

	r := make(map[string]OperationRatio)
	err := json.Unmarshal(data, &r)
	if err == nil {
		*ors = []map[string]OperationRatio{r}
	}
	return err
}

type GetRatiosResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			OperationRatiosAOR     OperationRatios                  `json:"operation_ratios_a_o_r"`
			OperationRatiosRestate OperationRatios                  `json:"operation_ratios_restate"`
			AlphaBeta              map[string]AlphaBeta             `json:"alpha_beta"`
			EarningsRatiosRestate  map[string]EarningsRatiosRestate `json:"earnings_ratios_restate"`
			ValuationRatios        *ValuationRatios                 `json:"valuation_ratios"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

type PriceStatistics struct {
	ShareClassID              *string  `json:"share_class_id"`
	AsOfDate                  *string  `json:"as_of_date"`
	Period                    *string  `json:"period"`
	ArithmeticMean            *float64 `json:"arithmetic_mean"`
	AverageVolume             *float64 `json:"average_volume"`
	Best3MonthTotalReturn     *float64 `json:"best3_month_return_total"`
	ClosePriceToMovingAverage *float64 `json:"close_price_to_moving_average"`
	HighPrice                 *float64 `json:"high_price"`
	LowPrice                  *float64 `json:"low_price"`
	MovingAveragePrice        *float64 `json:"moving_average_price"`
	PercentageBelowHighPrice  *float64 `json:"percentage_below_high_price"`
	StandardDeviation         *float64 `json:"standard_deviation"`
	TotalVolume               *float64 `json:"total_volume"`
	Worst3MonthTotalReturn    *float64 `json:"worst3_month_total_return"`
}

type TrailingReturns struct {
	ShareClassID *string  `json:"share_class_id"`
	AsOfDate     *string  `json:"as_of_date"`
	Period       *string  `json:"period"`
	TotalReturn  *float64 `json:"total_return"`
}

type GetPriceStatisticsResponse []struct {
	Error   string
	Request string `json:"request"`
	Results []struct {
		ID     string `json:"id"`
		Tables struct {
			PriceStatistics map[string]PriceStatistics `json:"price_statistics"`
			TrailingReturns map[string]TrailingReturns `json:"trailing_returns"`
		} `json:"tables"`
		Type string `json:"type"`
	} `json:"results"`
	Type string `json:"type"`
}

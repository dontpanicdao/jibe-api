package data

type TypedSubject struct {
	Types struct {
		StarkNetDomain []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"StarkNetDomain"`
		Exam []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"Exam"`
	} `json:"types"`
	PrimaryType string `json:"primaryType"`
	Domain      struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		ChainID int    `json:"chainId"`
	} `json:"domain"`
	Message struct {
		Name         string `json:"name"`
		AssetAddress string `json:"assetAddress"`
		NPhases      string `json:"nPhases"`
		DaoScheme    string `json:"daoScheme"`
		SignerScheme string `json:"signerScheme"`
		Provider     string `json:"provider"`
	} `json:"message"`
}

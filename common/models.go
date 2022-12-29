package common

type Conta struct {
	Agencia uint   `json:"agencia"`
	Conta   string `json:"conta"`
}

type Transferencia struct {
	ContaOrigem  Conta   `json:"conta_origem"`
	ContaDestino Conta   `json:"conta_destino"`
	Valor        float64 `json:"valor"`
	Moeda        string  `json:"moeda"`
}

func NovaTransferencia(contaOrigem Conta, contaDestino Conta, valor float64, moeda string) *Transferencia {
	return &Transferencia{
		ContaOrigem:  contaOrigem,
		ContaDestino: contaDestino,
		Valor:        valor,
		Moeda:        moeda,
	}
}

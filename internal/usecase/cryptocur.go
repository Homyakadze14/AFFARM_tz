package usecase

type CryptocurrencyStorage interface {
}

type CryptocurrencyService struct {
	st CryptocurrencyStorage
}

func NewCryptocurrencyService(st CryptocurrencyStorage) *CryptocurrencyService {
	return &CryptocurrencyService{
		st: st,
	}
}

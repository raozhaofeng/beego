package blockchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/raozhaofeng/beego/components/blockchain/tokens"
	"math"
	"math/big"
	"regexp"
)

const (
	TokenPocketRPC = "https://web3.mytokenpocket.vip"
	BinanceRPC     = "https://bsc.nodereal.io"
)

type Ethereum struct {
	ethClient        *ethclient.Client //	RPC连接对象
	privateKey       *ecdsa.PrivateKey //	私钥
	contractInstance *tokens.Tokens    //	合约实例
}

// TransactionByHash 查询交易hex状态
func (_Ethereum *Ethereum) TransactionByHash(hashTxStr string) (*types.Transaction, bool) {
	hashTx := common.HexToHash(hashTxStr)
	tx, isPending, err := _Ethereum.ethClient.TransactionByHash(context.Background(), hashTx)
	if err != nil {
		return nil, true
	}
	return tx, isPending
}

// TransactionAsMessage 获取哈希消息
func (_Ethereum *Ethereum) TransactionAsMessage(tx *types.Transaction) (types.Message, error) {
	return tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
}

// TokenTransferFrom 合约授权转账
func (_Ethereum *Ethereum) TokenTransferFrom(fromAddress, toAddress common.Address, value *big.Float) (*types.Transaction, error) {
	//	获取私钥权限
	accountAuth, err := _Ethereum.GetAccountAuth()
	if err != nil {
		return nil, err
	}

	//	转换数量
	decimalValue, err := _Ethereum.GetTokenDecimalsAmount(value)
	if err != nil {
		return nil, err
	}

	transaction, err := _Ethereum.contractInstance.TransferFrom(accountAuth, fromAddress, toAddress, decimalValue)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// TokenTransfer 合约转账
func (_Ethereum *Ethereum) TokenTransfer(address common.Address, value *big.Float) (*types.Transaction, error) {
	//	获取私钥权限
	accountAuth, err := _Ethereum.GetAccountAuth()
	if err != nil {
		return nil, err
	}

	//	转换数量
	decimalValue, err := _Ethereum.GetTokenDecimalsAmount(value)
	if err != nil {
		return nil, err
	}

	transaction, err := _Ethereum.contractInstance.Transfer(accountAuth, address, decimalValue)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// TokenBalance 合约余额
func (_Ethereum *Ethereum) TokenBalance(address common.Address) (*big.Float, error) {
	balance, err := _Ethereum.contractInstance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return nil, err
	}

	decimals, err := _Ethereum.contractInstance.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return _Ethereum.FormatEther(balance, int(decimals.Int64())), nil
}

// Balance 获取余额
func (_Ethereum *Ethereum) Balance(address common.Address) (*big.Float, error) {
	//	RPC获取余额
	balance, err := _Ethereum.ethClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return _Ethereum.FormatEther(balance, 18), nil
}

// SetClient 设置PRC对象
func (_Ethereum *Ethereum) SetClient(rpcURL string) *Ethereum {
	_Ethereum.ethClient, _ = ethclient.Dial(rpcURL)
	return _Ethereum
}

// CloseClient 关闭连接
func (_Ethereum *Ethereum) CloseClient() {
	_Ethereum.ethClient.Close()
}

// SetPrivateKey 设置私钥
func (_Ethereum *Ethereum) SetPrivateKey(privateStr string) *Ethereum {
	_Ethereum.privateKey, _ = crypto.HexToECDSA(privateStr)
	return _Ethereum
}

// SetContract 设置合约
func (_Ethereum *Ethereum) SetContract(address string) *Ethereum {
	_Ethereum.contractInstance, _ = tokens.NewTokens(common.HexToAddress(address), _Ethereum.ethClient)
	return _Ethereum
}

// GetAccountAuth 获取账号权限
func (_Ethereum *Ethereum) GetAccountAuth() (*bind.TransactOpts, error) {
	chainID, err := _Ethereum.ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	accountAuth, err := bind.NewKeyedTransactorWithChainID(_Ethereum.privateKey, chainID)
	if err != nil {
		return nil, err
	}
	return accountAuth, nil
}

// GetTokenDecimalsAmount 转换Token的数量
func (_Ethereum *Ethereum) GetTokenDecimalsAmount(value *big.Float) (*big.Int, error) {
	decimals, err := _Ethereum.contractInstance.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	tenDecimal := big.NewFloat(math.Pow(10, float64(decimals.Int64())))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, value).Int(&big.Int{})
	return convertAmount, nil
}

// FormatEther 获取以太单位
func (_Ethereum *Ethereum) FormatEther(wei *big.Int, decimals int) *big.Float {
	bigBalance := new(big.Float)
	bigBalance.SetString(wei.String())
	return new(big.Float).Quo(bigBalance, big.NewFloat(math.Pow10(decimals)))
}

// FormatWei 浮点数转
func (_Ethereum *Ethereum) FormatWei(value *big.Float, decimals int) *big.Int {
	m := new(big.Float)
	convertAmount, _ := m.Mul(value, big.NewFloat(math.Pow10(decimals))).Int(&big.Int{})
	return convertAmount
}

// TokenDecimals 合约精度
func (_Ethereum *Ethereum) TokenDecimals() (*big.Int, error) {
	return _Ethereum.contractInstance.Decimals(&bind.CallOpts{})
}

// GenerateKey 生成 私钥｜地址
func (_Ethereum *Ethereum) GenerateKey() (string, string, error) {
	//	生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	//	转成字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateHex := hexutil.Encode(privateKeyBytes)[2:]

	//	私钥派生公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return privateHex, address, nil
}

// IsAddress 验证地址是否有效
func (_Ethereum *Ethereum) IsAddress(address string) bool {
	compile := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return compile.MatchString(address)
}

// NewEthereum 创建ETH对象
func NewEthereum() *Ethereum {
	return new(Ethereum).SetClient(TokenPocketRPC)
}

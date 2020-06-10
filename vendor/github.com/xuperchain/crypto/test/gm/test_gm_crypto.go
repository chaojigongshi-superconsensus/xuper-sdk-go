package main

import (
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/xuperchain/crypto/client/service/gm"
	"github.com/xuperchain/crypto/gm/account"
	"github.com/xuperchain/crypto/gm/hdwallet/rand"

	hdapi "github.com/xuperchain/crypto/gm/hdwallet/api"
)

func main() {
	gcc := new(gm.GmCryptoClient)

	// --- 哈希算法相关 start ---
	hashResult := gcc.HashUsingSM3([]byte("This is xchain crypto"))
	log.Printf("Hash result for [This is xchain crypto] is: %s", hashResult)
	// --- 哈希算法相关 end ---

	// --- 账户生成相关 start ---
	ecdsaAccount, err := gcc.CreateNewAccountWithMnemonic(rand.SimplifiedChinese, account.StrengthHard)
	if err != nil {
		log.Printf("CreateNewAccountWithMnemonic failed and err is: %v", err)
		return
	}
	log.Printf("mnemonic is %v, jsonPrivateKey is %v, jsonPublicKey is %v and address is %v", ecdsaAccount.Mnemonic, ecdsaAccount.JsonPrivateKey, ecdsaAccount.JsonPublicKey, ecdsaAccount.Address)
	// --- 账户生成相关 end ---

	// --- 账户恢复相关 start ---
	// 从助记词恢复账户
	// 测试错误助记词
	test_mnemonic := "This is a test"
	wrongEcdsaAccount, err := gcc.RetrieveAccountByMnemonic(test_mnemonic, rand.SimplifiedChinese)
	log.Printf("retrieve account from test mnemonic: [%v], ecdsaAccount is %v and err is %v", test_mnemonic, wrongEcdsaAccount, err)

	// 测试正确助记词
	ecdsaAccount, err = gcc.RetrieveAccountByMnemonic(ecdsaAccount.Mnemonic, rand.SimplifiedChinese)
	if err != nil {
		log.Printf("RetrieveAccountByMnemonic failed and err is: %v", err)
		return
	}
	log.Printf("retrieve account from mnemonic %v, ecdsaAccount is %v and err is %v", ecdsaAccount.Mnemonic, ecdsaAccount, err)
	// --- 账户恢复相关 end ---

	// --- ECDSA签名算法相关 start ---
	msg := []byte("Welcome to the world of super chain using GM.")
	strJsonPrivateKey := ecdsaAccount.JsonPrivateKey
	privateKey, err := gcc.GetEcdsaPrivateKeyFromJsonStr(strJsonPrivateKey)
	sig, err := gcc.SignECDSA(privateKey, msg)
	log.Printf("sig is %v and err is %v", sig, err)

	isSignatureMatch, err := gcc.VerifyECDSA(&privateKey.PublicKey, sig, msg)
	log.Printf("Verifying & Unmashalling GM ecdsa signature by VerifyECDSA, isSignatureMatch is %v and err is %v", isSignatureMatch, err)
	// --- ECDSA签名算法相关 end ---

	// --- 非对称加密算法相关 start ---
	msg = []byte("Hello encryption!")
	ct, err := gcc.EncryptByEcdsaKey(&privateKey.PublicKey, msg)
	if err != nil {
		log.Printf("Encrypt failed and err is: %v", err)
		return
	}

	pt, err := gcc.DecryptByEcdsaKey(privateKey, ct)
	if err != nil {
		log.Printf("Decrypt failed and err is: %v", err)
		return
	}
	log.Printf("pt msg after decryption is: %s", pt)
	// --- 非对称加密算法相关 end ---

	// --- 多重签名相关 start ---
	privateKey1 := privateKey

	ecdsaAccount2, _ := gcc.CreateNewAccountWithMnemonic(rand.SimplifiedChinese, account.StrengthHard)
	strJsonPrivateKey2 := ecdsaAccount2.JsonPrivateKey
	privateKey2, _ := gcc.GetEcdsaPrivateKeyFromJsonStr(strJsonPrivateKey2)

	// 开始算多重签名sig1
	var multiSignKeys []*ecdsa.PrivateKey
	multiSignKeys = append(multiSignKeys, privateKey1)
	multiSignKeys = append(multiSignKeys, privateKey2)

	multiSig, err := gcc.MultiSign(multiSignKeys, msg)
	log.Printf("generate XuperSignature of multiSig is: %s and err is %v", multiSig, err)

	// 开始验证多重签名
	var multiSignVerifyKeys []*ecdsa.PublicKey
	multiSignVerifyKeys = append(multiSignVerifyKeys, &privateKey1.PublicKey)
	multiSignVerifyKeys = append(multiSignVerifyKeys, &privateKey2.PublicKey)

	isSignatureMatch, err = gcc.VerifyXuperSignature(multiSignVerifyKeys, multiSig, msg)
	log.Printf("Verifying & Unmashalling multiSign signature by VerifyXuperSignature, isSignatureMatch is: %v and err is %v", isSignatureMatch, err)
	// -- 多重签名相关 end ---

	// --- Schnorr签名算法相关 start ---
	schnorrSig, err := gcc.SignSchnorr(privateKey, msg)
	log.Printf("Schnorr signature is %s and err is %v", schnorrSig, err)

	var schnorrKeys []*ecdsa.PublicKey
	schnorrKeys = append(schnorrKeys, &privateKey.PublicKey)

	isSignatureMatch, err = gcc.VerifyXuperSignature(schnorrKeys, schnorrSig, msg)
	log.Printf("Verifying & Unmashalling Schnorr signature by VerifyXuperSignature, isSignatureMatch is %v and err is %v", isSignatureMatch, err)
	// --- Schnorr签名算法相关 end ---

	// --- Schnorr环签名算法相关 start ---
	ecdsaAccount3, _ := gcc.CreateNewAccountWithMnemonic(rand.SimplifiedChinese, account.StrengthHard)
	strJsonPrivateKey3 := ecdsaAccount3.JsonPrivateKey
	privateKey3, _ := gcc.GetEcdsaPrivateKeyFromJsonStr(strJsonPrivateKey3)

	ringSig, err := gcc.SignSchnorrRing(multiSignVerifyKeys, privateKey3, msg)
	log.Printf("Schnorr ring signature is %s and err is %v", ringSig, err)

	var schnorrRingSignVerifyKeys []*ecdsa.PublicKey
	schnorrRingSignVerifyKeys = append(schnorrRingSignVerifyKeys, &privateKey1.PublicKey)
	schnorrRingSignVerifyKeys = append(schnorrRingSignVerifyKeys, &privateKey2.PublicKey)
	schnorrRingSignVerifyKeys = append(schnorrRingSignVerifyKeys, &privateKey3.PublicKey)
	log.Printf("schnorrRingSignVerifyKeys is [%v]", schnorrRingSignVerifyKeys)

	isSignatureMatch, err = gcc.VerifyXuperSignature(schnorrRingSignVerifyKeys, ringSig, msg)
	log.Printf("Verifying & Unmashalling Schnorr ring signature, isSignatureMatch is %v and err is %v", isSignatureMatch, err)

	// 生成环签名地址
	ringAddress, err := gcc.GetAddressFromPublicKeys(schnorrRingSignVerifyKeys)
	log.Printf("Schnorr ring signature address is %s and err is %v", ringAddress, err)
	isAddressValid, _ := gcc.VerifyAddressUsingPublicKeys(ringAddress, schnorrRingSignVerifyKeys)
	log.Printf("Schnorr ring signature address[%s] is %v", ringAddress, isAddressValid)
	// --- Schnorr环签名算法相关 end ---

	// --- hd crypto api ---
	log.Printf("hd crypto api ----------")

	hdMnemonic := "呈 仓 冯 滚 刚 伙 此 丈 锅 语 揭 弃 精 塘 界 戴 玩 爬 奶 滩 哀 极 样 费"
	// 中心化控制中心产生根密钥
	rootKey, _ := gcc.GenerateMasterKeyByMnemonic(hdMnemonic, rand.SimplifiedChinese)
	// 中心化控制中心产生父私钥
	parentPrivateKey, _ := gcc.GenerateChildKey(rootKey, hdapi.HardenedKeyStart+8)
	// 中心化控制中心产生父公钥，并分发给客户端
	parentPublicKey, _ := gcc.ConvertPrvKeyToPubKey(parentPrivateKey)

	hdMsg := "Hello hd msg!"

	// 客户端为每次加密产生子公钥
	newChildPublicKey, err := gcc.GenerateChildKey(parentPublicKey, 18)
	log.Printf("newChildPublicKey is %v and err is %v", newChildPublicKey, err)
	// 客户端使用子公钥加密，产生密文
	cryptoMsg, err := gcc.EncryptByHdKey(newChildPublicKey, hdMsg)
	log.Printf("cryptoMsg is %v and err is %v", []byte(cryptoMsg), err)

	// 中心化控制中心使用根密钥、子公钥、密文，解密出原文
	realMsg, err := gcc.DecryptByHdKey(newChildPublicKey, rootKey, cryptoMsg)
	log.Printf("realMsg decrypted by root key is: [%s] and err is %v", realMsg, err)

	// 全节点使用一级父私钥、二级子公钥、密文，解密出原文
	realMsg, err = gcc.DecryptByHdKey(newChildPublicKey, parentPrivateKey, cryptoMsg)
	log.Printf("realMsg decrypted by parent private key is: [%s] and err is %v", realMsg, err)
	log.Printf("hd crypto api end----------")
	// -- hd crypto api end ---

	// --- secret share start ---
	//	msg = []byte("Welcome to the world of secret share.")
	secretMsg := 2147483647
	//	log.Printf("max int is %d", int(^uint32(0)>>1))
	log.Printf("secret_share secret is %d", secretMsg)
	totalShareNumber := 7
	minimumShareNumber := 3

	// 不能太大，否则会由于超出有限域范围产生数据丢失
	complexSecretBigInt, _ := big.NewInt(0).SetString("469507068585669108987494430799457046190734249189690901954100429825889257211", 0)
	complexSecretMsg := complexSecretBigInt.Bytes()
	//	complexSecretMsg := []byte("a")
	//	log.Printf("secret_share complexSecretMsg is: %s", complexSecretMsg)
	log.Printf("secret_share complexSecretMsg is: %d", complexSecretBigInt)

	complexShares, err := gcc.SecretSplit(totalShareNumber, minimumShareNumber, complexSecretMsg)
	log.Printf("secret_share ComplexSecretSplit result is %v and err is %v", complexShares, err)

	retrieveComplexShares := make(map[int]*big.Int, minimumShareNumber)
	number := 0
	for k, v := range complexShares {
		if number >= minimumShareNumber {
			break
		}
		retrieveComplexShares[k] = v
		number++
	}

	secretBytes, _ := gcc.SecretRetrieve(retrieveComplexShares)
	//	log.Printf("secret_share ComplexSecretRetrieve result is: %s", secretBytes)
	log.Printf("secret_share ComplexSecretRetrieve result is: %d", big.NewInt(0).SetBytes(secretBytes))
	// --- secret share end ---

	//--------------密钥分存-----------------
	strPrivKeyShares, err := gcc.SplitPrivateKey(strJsonPrivateKey, totalShareNumber, minimumShareNumber)
	log.Printf("share_key SplitPrivateKey result is: %s, and err is: %v", strPrivKeyShares, err)

	jsonPrivKey, err := gcc.RetrievePrivateKeyByShares(strPrivKeyShares[0:minimumShareNumber])
	log.Printf("share_key RetrievedPrivateKey fragments are: %s, result is: %s, and err is: %v", strPrivKeyShares, jsonPrivKey, err)

	//--------------密钥分存结束-----------------
}
